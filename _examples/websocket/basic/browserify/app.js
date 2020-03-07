const neffos = require('neffos.js');

var scheme = document.location.protocol == "https:" ? "wss" : "ws";
var port = document.location.port ? ":" + document.location.port : "";

var wsURL = scheme + "://" + document.location.hostname + port + "/echo";

const enableJWT = true;
if (enableJWT) {
  //这只是示例内容的签名和payload，
  //请将其替换为您的逻辑。

  // This is just a signature and a payload of an example content, 
  // please replace this with your logic.

  //在令牌前添加一个随机字母以使其无效，并确保不允许该客户端拨打Websocket服务器

  // Add a random letter in front of the token to make it
  // invalid and see that this client is not allowed to dial the websocket server.
  const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjozMjEzMjF9.8waEX7-vPKACa-Soi1pQvW3Rl8QY-SUFcHKTLZI4mvU";
  wsURL += "?token=" + token;
}

var outputTxt = document.getElementById("output");

function addMessage(msg) {
  outputTxt.innerHTML += msg + "\n";
}

function handleError(reason) {
  console.log(reason);
  window.alert(reason);
}

function handleNamespaceConnectedConn(nsConn) {
  nsConn.emit("chat", "Hello from browser(ify) client-side!");

  const inputTxt = document.getElementById("input");
  const sendBtn = document.getElementById("sendBtn");

  sendBtn.disabled = false;
  sendBtn.onclick = function () {
    const input = inputTxt.value;
    inputTxt.value = "";

    nsConn.emit("chat", input);
    addMessage("Me: " + input);
  };
}

async function runExample() {
  try {
    const conn = await neffos.dial(wsURL, {
      default: { // "default" namespace.
        _OnNamespaceConnected: function (nsConn, msg) {
          addMessage("connected to namespace: " + msg.Namespace);
          handleNamespaceConnectedConn(nsConn);
        },
        _OnNamespaceDisconnect: function (nsConn, msg) {
          addMessage("disconnected from namespace: " + msg.Namespace);
        },
        chat: function (nsConn, msg) { // "chat" event.
          addMessage(msg.Body);
        }
      }
    });

    //您可以等待连接，也可以只是conn.connect("connect")，
    //然后将`handleNamespaceConnectedConn`放在`_OnNamespaceConnected`回调中
    // const nsConn = await conn.connect("default");
    // nsConn.emit(...); handleNamespaceConnectedConn(nsConn);

    // You can either wait to conenct or just conn.connect("connect")
    // and put the `handleNamespaceConnectedConn` inside `_OnNamespaceConnected` callback instead.
    // const nsConn = await conn.connect("default");
    // handleNamespaceConnectedConn(nsConn);
    // nsConn.emit(...); handleNamespaceConnectedConn(nsConn);
    conn.connect("default");

  } catch (err) {
    handleError(err);
  }
}

runExample();