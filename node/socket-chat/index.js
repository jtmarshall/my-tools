var express = require('express');
var app = express();
var http = require('http').Server(app);
var io = require('socket.io')(http); // Pass in http object into new instance of socket.io
var port = process.env.PORT || 3000;
var path = require("path");


// Start listening on port
http.listen(port, () => {
    console.log('listening on *:' + port);
});

// Serve html template when endpoint is hit
// app.get('/', function (req, res) {
//     res.sendFile(path.join(__dirname, 'app'));
// });

// Serve app on /chat route
app.use('/chat', express.static(path.join(__dirname, 'app')));

// Count ppl in the chat
let numActive = 0;

// Listen on connection event for incoming sockets
io.on('connection', (socket) => {
    // Increment active usrs
    numActive++;

    // broadcast new user connect event to others in chatroom
    socket.broadcast.emit('newMsg', {
        username: socket.username,
        message: "new usr joined chat! There are now " + numActive + " usrs in chat."
    });

    // send message to everyone
    socket.on('newMsg', (msg) => {
        // tell client to execute 'newMsg'
        socket.broadcast.emit('newMsg', {
            username: socket.username,
            message: msg
        });
    });

    // log disconnect event
    socket.on('disconnect', () => {
        // Decrement active usrs
        numActive--;

        socket.broadcast.emit('newMsg', {
            username: socket.username,
            message: "usr left chat! There are now " + numActive + " usrs in chat."
        });
    });
});