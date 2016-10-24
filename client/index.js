var msg = null;
var log = null;

function appendLog(item) {
    var doScroll = log.scrollTop === log.scrollHeight - log.clientHeight;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

document.addEventListener("DOMContentLoaded", function(event) {
    msg = document.getElementById("msg");
    log = document.getElementById("log");

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        const id = location.search.substring(1).replace(/room=/, '');
        conn.send(msg.value);
        msg.value = "";
        return false;
    };

    const id = location.search.substring(1).replace(/room=/, '');
    const conn = new WebSocket(`ws://127.0.0.1:8088/room/${id}`);
    conn.onclose = function (evt) {
        let item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        appendLog(item);
    };
    conn.onmessage = function (evt) {
        var messages = evt.data.split('\n');
        for (var i = 0; i < messages.length; i++) {
            var item = document.createElement("div");
            item.innerText = messages[i];
            appendLog(item);
        }
    };
});
