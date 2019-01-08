package tcp_ip

import (
	"bufio"
	"encoding/gob"
	"flag"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Struct with mixed fields to use with gob
type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}

// Function to handle an incoming command
type HandleFunc func(*bufio.ReadWriter)

// Provides endpoint to other processes that they can send data to
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc
	m        sync.RWMutex // maps NOT thread safe, need mutex for access control
}

const (
	Port = ":61000"
)

// Starts as either a client or server depending on the connect flag set
// With flag start process as client and connect to the host specified by flag value
func main() {
	connect := flag.String("connect", "", "IP address of process to join. If empty, go into listen mode.")
	flag.Parse()

	// If connect flag is set we start as client
	if *connect != "" {
		err := client(*connect)
		if err != nil {
			log.Println("Err:", errors.WithStack(err))
		}
		log.Println("Client done.")
		return
	}

	// Otherwise we start as server
	err := server()
	if err != nil {
		log.Println("Err:", errors.WithStack(err))
	}

	log.Println("Server done.")
}

// Lshortfile flag includes file name & line number in logs
func init() {
	log.SetFlags(log.Lshortfile)
}

// Open: connects to TCP address and returns TCP connection with a timeout and wrapped into a buffered ReadWriter
func Open(addr string) (*bufio.ReadWriter, error) {

	// Dial the remote process
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}

	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// Create a new endpoint; listens on a fixed port number
func NewEndpoint() *Endpoint {

	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

// Add new function for handling incoming data
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// Start listening on the endpoint port on all interfaces.
// Need to have at least one handler func added through AddHandleFunc() before.
func (e *Endpoint) Listen() error {
	var err error

	e.listener, err = net.Listen("tcp", Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", Port)
	}

	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a conn request: ", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go e.handleMessages(conn)
	}
}

// Read connection up to newline, use command in this string to call appropriate HandleFunc
func (e *Endpoint) handleMessages(conn net.Conn) {
	// wrap in buffered reader for easy reading
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// read from connection until EOF. Expecting command name, then call handler registered for that command
	for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close the connection\n  ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}

		// trim request string
		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd + "'")

		// get appropriate handler func from "handler" map, then call it
		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw)
	}
}

// Handle "STRING" command from request
func handleStrings(rw *bufio.ReadWriter) {
	log.Print("Receive STRING message: ")

	// read message up to newline
	s, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.\n", err)
	}
	s = strings.Trim(s, "\n ")
	log.Println(s)

	_, err = rw.WriteString("Thanks for the STRING command!\n")
	if err != nil {
		log.Println("Can't write to connection.\n", err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println("Something clogged because the Flush failed.", err)
	}
}

// Handle "GOB" command from request
func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive GOB data: ")
	var data complexData

	// decoder to decode straight into struct var
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data: ", err)
		return
	}

	// print out complexData struct and nested
	log.Printf("Outer complexData struct: \n%#v\n", data)
	log.Printf("Inner complexData struct: \n%#v\n", data.C)
}

// Simulate client interactions
func client(ip string) error {
	testStruct := complexData{
		N: 23,
		S: "string data",
		M: map[string]int{
			"one": 1,
			"two": 2,
			"three": 3,
		},
		P: []byte("abc"),
		C: &complexData{
			N: 256,
			S: "Recursive struct, for funsies",
			M: map[string]int{
				"01": 1,
				"10": 2,
				"11": 3,
			},
		},
	}

	rw, err := Open(ip + Port)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+Port)
	}
	log.Println("Send the string request.")
	n, err := rw.WriteString("STRING\n")
	if err != nil {
		return errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)")
	}
	n , err = rw.WriteString("Additional data.\n")
	if err != nil {
		return errors.Wrap(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)")
	}
	log.Println("Flush buffer.")
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}

	log.Println("Read the reply.")
	response, err := rw.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "Client: Failed to read the reply: '"+response+"'")
	}

	log.Println("STRING request: got response:", response)

	log.Println("Send a struct as GOB:")
	log.Printf("Outer complexData struct: \n%#v\n", testStruct)
	log.Printf("Inner complexData struct: \n%#v\n", testStruct.C)
	enc := gob.NewEncoder(rw)
	n, err = rw.WriteString("GOB\n")
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}

	err = enc.Encode(testStruct)
	if err != nil {
		return errors.Wrapf(err, "Encode failed for struct: %#v", testStruct)
	}

	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed")
	}

	return nil
}

// Listens for incoming requests and dispatches to registered handler funcs
func server() error {
	endpoint := NewEndpoint()

	// add the handler funcs
	endpoint.AddHandleFunc("STRING", handleStrings)
	endpoint.AddHandleFunc("GOB", handleGob)

	// start listening
	return endpoint.Listen()
}
