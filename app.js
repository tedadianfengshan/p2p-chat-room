window.addEventListener("load", function(evt) {
    var display_div = document.getElementById("div1")
    var input = document.getElementById("input1")
    var ws
    var print = function(message) {
        var d = document.createElement("div")
        d.innerHTML = message
        display_div.appendChild(d)
    }

    document.getElementById("button1").onclick = function(evt) {
        if (ws) {
            return false
        }

        var ws_url = "ws://" + location.host + "/ws"
        ws = new WebSocket(ws_url)
        
        ws.onopen = function(evt) {
            print("OPEN")
        }

        ws.onclose = function(evt) {
            print("CLOSE")
            ws = null
        }

        ws.onmessage = function(evt) {
            print("Server: " + evt.data)
        }

        ws.onerror = function(evt) {
            print("ERROR: " + evt.data)
        }

        return false
    }

    document.getElementById("button2").onclick = function(evt) {
        if (!ws) {
            return false
        }

        print("Send: " + input.value)
        ws.send(input.value)
        input.value = ""
    }

    document.getElementById("button3").onclick = function(evt) {
        if (!ws) {
            return false
        }

        ws.close()
        return false
    }

    document.onkeydown=function(event){
        var e = event || window.event || arguments.callee.caller.arguments[0];
        if(e && e.key=="Enter"){
            document.getElementById("button2").click();
        }
    };
})