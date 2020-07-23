package hotzone

import (
    "fmt"
    "github.com/BabyEngine/Backend/logger"
    "github.com/BabyEngine/Backend/core"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "time"
)

var (
    HTTPData = `<!DOCTYPE html>
<html>
<head>
    <title>Websocket Debug Viewer</title>
    <meta charset="UTF-8">
    <link href='http://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700,400italic,700italic' rel='stylesheet' type='text/css'>
    <script src="https://code.jquery.com/jquery-3.1.1.slim.min.js"></script>
    <style type="text/css">
    body {
        font-family: 'Ubuntu Mono', monospace;
        font-size: 12px
    }

    #log_container div.file {
        border: 1px solid #CCC;
        padding: 1em;
        margin: 1em;
    }

    #log_container div.file div.line {
        font-size: 10px
    }

    /* Clearfix */
    .cf:before,
    .cf:after {
        content: " "; /* 1 */
        display: table; /* 2 */
    }
    .cf:after {
        clear: both;
    }
    .cf {
        *zoom: 1;
    }
    </style>

    <script type="text/javascript">
    var socket;
    var max_line = 50;

    function socketize(wsUri) {
        socket = new WebSocket('ws://' + wsUri);
        socket.onopen = function(evt) { onEvent(evt) };
        socket.onclose = function(evt) { onEvent(evt) };
        socket.onmessage = function(evt) { onMessage(evt) };
        socket.onerror = function(evt) { onEvent(evt) };
        console.log(wsUri);
    }

    function onEvent(evt) {
        switch (evt.type) {
            case 'open':
            $("#wsconnect").text('Connected');
            $("#wsaddr").prop("disabled", "true");
            $("#log_container").html('');
            break;

            default:
            $("#wsconnect").text('Disconnected');
            $("#wsaddr").prop("disabled", "false");
            break;
        }
    }

    function onMessage(evt) {
        var lines = evt.data.replace(/^\s+|\s+$/g, '').split("\n");
        var container = $("#log_container");
        var elemid = 'h_';
        var lastfile = "";

        if ($("#" + elemid, container).length) {
            var elem = $("#" + elemid + ' div.lines', container);
        } else {
            var elem = $('<div/>', {id: elemid}).addClass("file").append($('<h3/>')).append($('<div/>').addClass('lines')).appendTo(container).find('div.lines');
        }

        // Currently we receive only one line of log each call, but assumptions is bad.
        for (i = 0; i < lines.length; ++i) {
            $('<div/>').addClass('line').text(lines[i]).prependTo(elem);
        }

        // Show last 25 records
        var elemcount = $("div.line", elem).length;
        if (elemcount > max_line) {
            $("div.line", elem).slice(max_line-elemcount).remove();
        }

        // Move to top
        $('#' + elemid).prependTo(container);


        //console.log(evt.data);
    }

    $(function() {
        $("#wsconnect").click(function() {
            socketize($("#wsaddr").val());
        }).click();
    });
    $(function() {
        $("#hotreload").click(function() {
            var url = 'http://'+$("#wsaddr").val() + "/r";
            const Http = new XMLHttpRequest();
            Http.open("GET", url);
            Http.send();
console.log("xx");
            Http.onreadystatechange = (e) => {
                if (Http.readyState == 4) {
                    appendLog(Http.responseText);
                }
            }
        });
    });

    $(function() {
        $("#clean").click(function() {
            console.log('clean log');
            var container = $("#log_container");
            var elem = $("#h_" + ' div.lines', container);
            $("div.line", elem).slice(0).remove();
        });
    });
    $(function() {
        $("#log_num").change(function(val) {
            var n = $("#log_num").val();
            n = parseInt(n, 10);
            if (n < 5) {
                n = 5;
            }
            max_line = n;
            $("#log_num").val(max_line);
        }).change();
    });

    function appendLog(msg) {
        console.log(msg);
        var lines = msg.replace(/^\s+|\s+$/g, '').split("\n");
        var container = $("#log_container");
        var elemid = 'h_';
        var lastfile = "";

        if ($("#" + elemid, container).length) {
            var elem = $("#" + elemid + ' div.lines', container);
        } else {
            var elem = $('<div/>', {id: elemid}).addClass("file").append($('<h3/>')).append($('<div/>').addClass('lines')).appendTo(container).find('div.lines');
        }

        // Currently we receive only one line of log each call, but assumptions is bad.
        for (i = 0; i < lines.length; ++i) {
            $('<div/>').addClass('line').text(lines[i]).prependTo(elem);
        }

        // Show last 25 records
        var elemcount = $("div.line", elem).length;
        if (elemcount > max_line) {
            $("div.line", elem).slice(max_line-elemcount).remove();
        }

        // Move to top
        $('#' + elemid).prependTo(container);
    }

    function validateNumber(evt) {
        var e = evt || window.event;
        var key = e.keyCode || e.which;

        if (!e.shiftKey && !e.altKey && !e.ctrlKey &&
        // numbers
        key >= 48 && key <= 57 ||
        // Numeric keypad
        key >= 96 && key <= 105 ||
        // Backspace and Tab and Enter
        key == 8 || key == 9 || key == 13 ||
        // Home and End
        key == 35 || key == 36 ||
        // left and right arrows
        key == 37 || key == 39 ||
        // Del and Ins
        key == 46 || key == 45) {
            // input is VALID
        }
        else {
            // input is INVALID
            e.returnValue = false;
            if (e.preventDefault) e.preventDefault();
        }
    }

    </script>
</head>
<body>
    <div id="wssrc_container">
        <p><strong>Websocket 日志查看器</strong> <br>连接以 TextMessage 回复的 Websocket</p>
        ws://<input type="text" id="wsaddr" value="127.0.0.1:80" />
        <button id="wsconnect">Connect</button>
        <button id="clean">Clean log</button>
        <button id="hotreload">Hot Reload</button>
        <label for="Name">Max Line:</label>
        <input type="" id="log_num" placeholder="text" onkeydown="validateNumber(event);" value="50"/>
    </div>
    <hr>
    <div id="log_container">
        Your logs will show here.
    </div>
</body>
</html>
`
)
var upgrader = websocket.Upgrader{
    // 解决跨域问题
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


var (
    logChan = make(chan string, 100000)
    isRunning bool
)
func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func EnableHotRestart(app *core.Application, restartCb func()) {
    if isRunning { return }
    isRunning = true
    defer func() {
        isRunning = false
    }()
    logger.Debug("Enable Hot Reload")
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
        w.Write([]byte(HTTPData))
    })
    http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
        app.Stop()
        go restartCb()
        w.Write([]byte(fmt.Sprintf("restart success %v", time.Now().Format(time.RFC3339))))
    })
    http.HandleFunc("/log", showLog)
    if err := http.ListenAndServe(":80", nil); err != nil {
        logger.Debug("start hot restart fail:", err)
    }
}

func showLog(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()
    n := 0
    for {
        //time.Sleep(time.Second)
        msg := <- logChan
        fmt.Println(msg)
        err := c.WriteMessage(websocket.TextMessage, []byte(msg))
        if err != nil {
            break
        }
        n++
    }
}
