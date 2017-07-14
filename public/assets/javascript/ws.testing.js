document.addEventListener("DOMContentLoaded", function(e){
  var websocket = new WebSocket("ws://localhost:3000/api/websocket")
  websocket.onopen = function(e) {
    websocket.send("What am I?")
  }
  websocket.onmessage = function(e) {
    console.log(e.data)
  }
})