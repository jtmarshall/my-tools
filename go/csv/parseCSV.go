package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Person struct {
	Firstname string //Annotation example: `json:"firstname"`
	Lastname  string
	Email     string
	Address   *Address
}

type Address struct {
	City  string
	State string
}

func main() {
	// Open up the file
	csvFile, _ := os.Open("test.csv")
	// Start the reader
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var peopleData []Person
	// Iterate through each line and append each data struct into the peopleData array
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		peopleData = append(peopleData, Person{
			Firstname: line[0],
			Lastname:  line[1],
			Email:     line[2],
			Address: &Address{
				City:  line[3],
				State: line[4],
			},
		})
	}
	// Convert to JSON, then print it out
	dataJSON, _ := json.Marshal(peopleData)
	fmt.Println(string(dataJSON))
}
