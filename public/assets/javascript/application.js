document.addEventListener("DOMContentLoaded", function(e){
  var date = new Date()
  var hours = date.getHours()
  if (hours > 12) {
    hours -= 12;
  } else if (hours === 0) {
    hours = 12;
  }
  var minutes = date.getMinutes()
  if(minutes < 12) {minutes = '0' + minutes}
  var xhr = new XMLHttpRequest();
  xhr.onload = function() {
    if (xhr.status === 200) {
        var data = JSON.parse(xhr.responseText)
        console.log("Current local time: "+hours+":"+minutes+"\n"+"Current Server Time: "+data.time)
    }
  }
  xhr.open("GET", "/api/time")
  xhr.send()
})



