let url = "ws://" + window.location.host + "/ws";
let ws = new WebSocket(url);

let name = "Guest" + Math.floor(Math.random() * 1000);

let chat = document.getElementById("chat");
let text = document.getElementById("text");

let now = function () {
    let iso = new Date().toISOString();
    return iso.split("T")[1].split(".")[0];
};

ws.onmessage = function (msg) {
    let line =  now() + " " + msg.data + "\n";
    chat.innerText = line + chat.innerText;
};

text.onkeydown = function (e) {
    if (e.keyCode === 13 && text.value !== "") {
        ws.send("<" + name + "> " + text.value);
        text.value = "";
    }
};

let $table = $('#table');

$(document).ready(function() {
  $.getJSON("http://" + window.location.host + "/competition", function(r) {
    chat.innerText = now() + " " + JSON.stringify(r, null, 2) + "\n" + chat.innerText;
  });
});

// vim: set ts=2 sw=2 et:
