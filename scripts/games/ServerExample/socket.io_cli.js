var io = require('socket.io-client');

var socket = io('ws://127.0.0.1:81', {transports: ['websocket'], rejectUnauthorized: false });

// listen for messages
socket.on('text', function(message) {
  message = Buffer.from(message, 'base64').toString('ascii');
  console.log('on message:', message);
});

socket.on('connect', function () {
  console.log('socket connected');
  var someData = {name:"hello"}
  var msg = Buffer.from(JSON.stringify(someData)).toString('base64')
  socket.emit('text', msg , function(data){
    data = Buffer.from(data, 'base64').toString('ascii');
    console.log('ACK from server wtih data: ', data);
  });
});