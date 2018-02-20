(function() {
    var count = 0,
        counter = document.getElementById("count"),
        url = new URL(location.href),
        addr = (url.port ? url.hostname + ':' + url.port : url.hostname),
        conn = new WebSocket('ws://'+addr+'/ws');

    conn.onclose = function(evt) {
        console.log("connection closed")
    }

    conn.onmessage = function(evt) {
        console.log(evt.data);

        if (evt.data == "ping") {
            conn.send("pong");
        }

        count++;
        counter.innerText = count;
    }
})();