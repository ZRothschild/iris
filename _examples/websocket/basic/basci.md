# go iris websocket 基本示例

在本示例的最后，您将能够运行Websocket服务器和所有平台（Go，Browser和Nodejs）的客户端。

![](overview.png)

打开想要尝试的任意数量的客户端并开始输入。

此示例仅包含基础知识，但是该库支持rooms，本机Websocket消息，可以发送和接收任何数据（即protobufs，json）以及各种广播和连接集合。

## 怎么运行

### 服务端

打开终端窗口实例并执行：

```sh
$ go run server.go # start the websocket server.
```

### 客户端 （go-client）

启动一个新的终端实例并执行：

```sh
$ cd ./go-client
$ go run client.go # start the websocket client.
# start typing...
```

### 客户端 （browser）

导航到<http://localhost:8080>并开始输入。应该提供`./browser/index.html`，它包含客户端代码

### 客户端 (browserify)

首先安装[NPM](https://nodejs.org) ，然后启动一个新的终端实例并执行：

```sh
$ cd ./browserify
$ npm install
$ npm安装

# 构建现代的浏览器端客户端：
# 嵌入neffos.js节点模块和app.js
# 放入单个./browserify/bundle.js文件
# 导入./browserify/client.html。

# build the modern browser-side client:
# embed the neffos.js node-module and app.js
# into a single ./browserify/bundle.js file
# which ./browserify/client.html imports.
$ npm run-script build
```
导航到<http://localhost:8080/browserify/client.html>并开始输入

### 客户端 (Nodejs)

如果尚未安装[NPM](https://nodejs.org)，则启动一个新的终端实例并执行：

```sh
$ cd nodejs-client
$ npm install
$ node client.js # start the websocket client.
# start typing.
```
## 目录结构
> 主目录`boltdb`
```html
    —— browser
        —— index.html
    —— browserify
        —— app.js
        —— bundle.js
        —— client.html
        —— package.json
    —— go-client
        —— client.go
    —— nodejs-client
        —— client.js
        —— package.json
    —— server.go

```
## 代码示例
> `browser/index.html`
```html
<!-- 消息的输入 -->

<!-- the message's input -->
<input id="input" type="text" />

<!-- 当单击时，一个websocket事件将被发送到服务器，在此示例中，我们注册了'chat'-->

<!-- when clicked then a websocket event will be sent to the server, at this example we registered the 'chat' -->
<button id="sendBtn" disabled>Send</button>

<!-- 消息将显示在这里 -->
<!-- the messages will be shown here -->
<pre id="output"></pre>
<!-- 从CDN或本地导入iris客户端库以供浏览器使用。但是，`neffos.(min.)js`也是NPM软件包，
     因此也可以将其用作package.json和所有nodejs-npm的依赖项 工具可用：
     有关更多信息，请参见"browserify"示例-->

<!-- import the iris client-side library for browser from a CDN or locally.
     However, `neffos.(min.)js` is a NPM package too so alternatively,
     you can use it as dependency on your package.json and all nodejs-npm tooling become available:
     see the "browserify" example for more-->
<script src="https://cdn.jsdelivr.net/npm/neffos.js@latest/dist/neffos.min.js"></script>
<script>
    //`neffos`全局变量现在可用

    // `neffos` global variable is available now.
    var scheme = document.location.protocol == "https:" ? "wss" : "ws";
    var port = document.location.port ? ":" + document.location.port : "";
    var wsURL = scheme + "://" + document.location.hostname + port + "/echo";

    const enableJWT = true;
    if (enableJWT) {
        //这只是示例内容的签名和有效内容，请使用您的逻辑替换它

        // This is just a signature and a payload of an example content, 
        // please replace this with your logic.

        //在令牌前添加一个随机字母以使其无效，
        //并确保不允许该客户端链接Websocket服务器

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
        window.alert("error: see the dev console");
    }

    function handleNamespaceConnectedConn(nsConn) {
        nsConn.emit("Hello from browser client side!");

        let inputTxt = document.getElementById("input");
        let sendBtn = document.getElementById("sendBtn");

        sendBtn.disabled = false;
        sendBtn.onclick = function () {
            const input = inputTxt.value;
            inputTxt.value = "";
            nsConn.emit("chat", input);
            addMessage("Me: " + input);
        };
    }

    const username = window.prompt("Your username?");

    async function runExample() {
        //您可以省略"default"，而仅定义事件，名namespace将是一个空字符串“”，
        // 但是，如果您决定对此示例进行任何更改，请确保这些更改反映在../server.go内部 文件

        // You can omit the "default" and simply define only Events, the namespace will be an empty string"",
        // however if you decide to make any changes on this example make sure the changes are reflecting inside the ../server.go file as well.
        try {
            const conn = await neffos.dial(wsURL, {
                default: { // "default" namespace.
                    _OnNamespaceConnected: function (nsConn, msg) {
                        addMessage("connected to namespace: " + msg.Namespace);
                        handleNamespaceConnectedConn(nsConn)
                    },
                    _OnNamespaceDisconnect: function (nsConn, msg) {
                        addMessage("disconnected from namespace: " + msg.Namespace);
                    },
                    chat: function (nsConn, msg) { // "chat" event.
                        addMessage(msg.Body);
                    }
                }
            },{
                headers: {
                    "X-Username": username,
                }
            });
            //您可以等待连接，也可以只是conn.connect("connect")，
            //然后将`handleNamespaceConnectedConn`放在`_OnNamespaceConnected`回调中
            // const nsConn = await conn.connect("default");
            // nsConn.emit(...); handleNamespaceConnectedConn(nsConn);

            // You can either wait to conenct or just conn.connect("connect")
            // and put the `handleNamespaceConnectedConn` inside `_OnNamespaceConnected` callback instead.
            // const nsConn = await conn.connect("default");
            // nsConn.emit(...); handleNamespaceConnectedConn(nsConn);
            conn.connect("default");

        } catch (err) {
            handleError(err);
        }
    }

    runExample();

    //如果可用"await"和"async"，请改用它们^，所有现代浏览器都支持它们，
    //所有的JavaScript示例都将使用async/await方法而不是promise then/catch回调编写。
    // 然后/捕获如下：
    
    // If "await" and "async" are available, use them instead^, all modern browsers support those,
    // all of the javascript examples will be written using async/await method instead of promise then/catch callbacks.
    // A usage example of promise then/catch follows:
    
    // neffos.dial(wsURL, {
    //     default: { // "default" namespace.
    //         _OnNamespaceConnected: function (ns, msg) {
    //             addMessage("connected to namespace: " + msg.Namespace);
    //         },
    //         _OnNamespaceDisconnect: function (ns, msg) {
    //             addMessage("disconnected from namespace: " + msg.Namespace);
    //         },
    //         chat: function (ns, msg) { // "chat" event.
    //             addMessage(msg.Body);
    //         }
    //     }
    // }).then(function (conn) {
    //     conn.connect("default").then(handleNamespaceConnectedConn).catch(handleError);
    // }).catch(handleError);
</script>
```
> `browserify/app.js`
```js
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
```
> `browserify/bundle.js`
```js
(function(){function b(d,e,g){function a(j,i){if(!e[j]){if(!d[j]){var f="function"==typeof require&&require;if(!i&&f)return f(j,!0);if(h)return h(j,!0);var c=new Error("Cannot find module '"+j+"'");throw c.code="MODULE_NOT_FOUND",c}var k=e[j]={exports:{}};d[j][0].call(k.exports,function(b){var c=d[j][1][b];return a(c||b)},k,k.exports,b,d,e,g)}return e[j].exports}for(var h="function"==typeof require&&require,c=0;c<g.length;c++)a(g[c]);return a}return b})()({1:[function(a){function b(a){j.innerHTML+=a+"\n"}function c(a){console.log(a),window.alert(a)}function d(a){a.emit("chat","Hello from browser(ify) client-side!");const c=document.getElementById("input"),d=document.getElementById("sendBtn");d.disabled=!1,d.onclick=function(){const d=c.value;c.value="",a.emit("chat",d),b("Me: "+d)}}async function e(){try{const a=await f.dial(i,{default:{_OnNamespaceConnected:function(a,c){b("connected to namespace: "+c.Namespace),d(a)},_OnNamespaceDisconnect:function(a,c){b("disconnected from namespace: "+c.Namespace)},chat:function(a,c){b(c.Body)}}});a.connect("default")}catch(a){c(a)}}const f=a("neffos.js");var g="https:"==document.location.protocol?"wss":"ws",h=document.location.port?":"+document.location.port:"",i=g+"://"+document.location.hostname+h+"/echo";{i+="?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjozMjEzMjF9.8waEX7-vPKACa-Soi1pQvW3Rl8QY-SUFcHKTLZI4mvU"}var j=document.getElementById("output");e()},{"neffos.js":2}],2:[function(a,b,c){(function(d,e){function f(a){return!("_OnNamespaceConnect"!==a&&"_OnNamespaceConnected"!==a&&"_OnNamespaceDisconnect"!==a&&"_OnRoomJoin"!==a&&"_OnRoomJoined"!==a&&"_OnRoomLeave"!==a&&"_OnRoomLeft"!==a)}function g(a){return!(void 0!==a)||!(null!==a)||(""==a||"string"==typeof a||a instanceof String?0===a.length||""===a:!!(a instanceof Error)&&g(a.message))}function h(a){return JSON.stringify(a)}function i(a){return g(a)?"":a.replace(K,"@%!semicolon@%!")}function j(a){return g(a)?"":a.replace(L,";")}function k(a){if(a.IsNative&&g(a.wait))return a.Body;var b="0",c="0",d=a.Body||"";return a.isError&&(d=a.Err,b="1"),a.isNoOp&&(c="1"),[a.wait||"",i(a.Namespace),i(a.Room),i(a.Event),b,c,d].join(";")}function l(a,b,c){if(0==c)return[a];var d=a.split(b,c);if(d.length==c){var e=d.join(b)+b;return d.push(a.substr(e.length)),d}return[a]}function m(a,b){var c=new J;if(0==a.length)return c.isInvalid=!0,c;var d=l(a,";",6);if(7!=d.length)return b?(c.Event="_OnNativeMessage",c.Body=a):c.isInvalid=!0,c;c.wait=d[0],c.Namespace=j(d[1]),c.Room=j(d[2]),c.Event=j(d[3]),c.isError="1"==d[4]||!1,c.isNoOp="1"==d[5]||!1;var e=d[6];return g(e)?c.Body="":c.isError?c.Err=e:c.Body=e,c.isInvalid=!1,c.IsForced=!1,c.IsLocal=!1,c.IsNative=b&&"_OnNativeMessage"==c.Event||!1,c}function n(){if(!A){var a=d.hrtime();return"$"+1e9*a[0]+a[1]}var b=window.performance.now();return"$"+b.toString()}function o(a){return a+";".repeat(6)}function p(a,b){return a.events.has(b.Event)?a.events.get(b.Event)(a,b):a.events.has("_OnAnyEvent")?a.events.get("_OnAnyEvent")(a,b):null}function q(a){return null===a||a===void 0||"undefined"==typeof a}function r(a,b){if(q(a))return q(b)||b("connHandler is empty."),null;var c=new Map,d=new Map,e=0;if(Object.keys(a).forEach(function(b){e++;var f=a[b];if(f instanceof Function)d.set(b,f);else if(f instanceof Map)c.set(b,f);else{var g=new Map;Object.keys(f).forEach(function(a){g.set(a,f[a])}),c.set(b,g)}}),0<d.size){if(e!=d.size)return q(b)||b("all keys of connHandler should be events, mix of namespaces and event callbacks is not supported "+d.size+" vs total "+e),null;c.set("",d)}return c}function s(a,b){return a.has(b)?a.get(b):null}function t(a,b){if(q(a))return b;for(var c in a)if(a.hasOwnProperty(c)){var d=a[c];c=encodeURIComponent("X-Websocket-Header-"+c),d=encodeURIComponent(d);var e=c+"="+d;b=-1==b.indexOf("?")?-1==b.indexOf("#")?b+"?"+e:b.split("#")[0]+"?"+e+"#"+b.split("#")[1]:b.split("?")[0]+"?"+e+"&"+b.split("?")[1]}return b}function u(a,b,c){return v(a,b,0,c)}function v(a,b,c,d){return-1==a.indexOf("ws")&&(a="ws://"+a),new Promise(function(e,f){WebSocket||f("WebSocket is not accessible through this browser.");var h=r(b,f);if(!q(h)){q(d)&&(d={}),q(d.headers)&&(d.headers={});var i=d.reconnect?d.reconnect:0;0<c&&0<i?d.headers["X-Websocket-Reconnect"]=c.toString():!q(d.headers["X-Websocket-Reconnect"])&&delete d.headers["X-Websocket-Reconnect"];var j=w(a,d),k=new T(j,h);k.reconnectTries=c,j.binaryType="arraybuffer",j.onmessage=function(a){var b=k.handle(a);return g(b)?void(k.isAcknowledged()&&e(k)):void f(b)},j.onopen=function(){j.send("M")},j.onerror=function(a){k.close(),f(a)},j.onclose=function(){if(k.isClosed());else{if(j.onmessage=void 0,j.onopen=void 0,j.onerror=void 0,j.onclose=void 0,0>=i)return k.close(),null;var c=new Map;k.connectedNamespaces.forEach(function(a,b){var d=[];!q(a.rooms)&&0<a.rooms.size&&a.rooms.forEach(function(a,b){d.push(b)}),c.set(b,d)}),k.close(),x(a,i,function(g){v(a,b,g,d).then(function(a){return q(e)||"function () { [native code] }"==e.toString()?void c.forEach(function(b,c){a.connect(c).then(function(a){return function(b){a.forEach(function(a){b.joinRoom(a)})}}(b))}):void e(a)}).catch(f)})}return null}}})}function w(a,b){return A&&!q(b)?(b.headers&&(a=t(b.headers,a)),b.protocols?new WebSocket(a,b.protocols):new WebSocket(a)):new WebSocket(a,b)}function x(a,b,c){var d=a.replace(/(ws)(s)?\:\/\//,"http$2://"),e=1,f={method:"HEAD",mode:"no-cors"},g=function(){B(d,f).then(function(){c(e)}).catch(function(){e++,setTimeout(function(){g()},b)})};setTimeout(g,b)}var y=this&&this.__awaiter||function(a,b,c,d){return new(c||(c=Promise))(function(e,f){function g(a){try{i(d.next(a))}catch(a){f(a)}}function h(a){try{i(d["throw"](a))}catch(a){f(a)}}function i(a){a.done?e(a.value):new c(function(b){b(a.value)}).then(g,h)}i((d=d.apply(a,b||[])).next())})},z=this&&this.__generator||function(a,b){function c(a){return function(b){return d([a,b])}}function d(c){if(e)throw new TypeError("Generator is already executing.");for(;k;)try{if(e=1,h&&(i=2&c[0]?h["return"]:c[0]?h["throw"]||((i=h["return"])&&i.call(h),0):h.next)&&!(i=i.call(h,c[1])).done)return i;switch((h=0,i)&&(c=[2&c[0],i.value]),c[0]){case 0:case 1:i=c;break;case 4:return k.label++,{value:c[1],done:!1};case 5:k.label++,h=c[1],c=[0];continue;case 7:c=k.ops.pop(),k.trys.pop();continue;default:if((i=k.trys,!(i=0<i.length&&i[i.length-1]))&&(6===c[0]||2===c[0])){k=0;continue}if(3===c[0]&&(!i||c[1]>i[0]&&c[1]<i[3])){k.label=c[1];break}if(6===c[0]&&k.label<i[1]){k.label=i[1],i=c;break}if(i&&k.label<i[2]){k.label=i[2],k.ops.push(c);break}i[2]&&k.ops.pop(),k.trys.pop();continue;}c=b.call(a,k)}catch(a){c=[6,a],h=0}finally{e=i=0}if(5&c[0])throw c[1];return{value:c[0]?c[1]:void 0,done:!0}}var e,h,i,j,k={label:0,sent:function(){if(1&i[0])throw i[1];return i[1]},trys:[],ops:[]};return j={next:c(0),throw:c(1),return:c(2)},"function"==typeof Symbol&&(j[Symbol.iterator]=function(){return this}),j},A="undefined"!=typeof window,B="undefined"==typeof fetch?void 0:fetch;A?WebSocket=window.WebSocket:(WebSocket=a("ws"),B=a("node-fetch"));var C="_OnNamespaceConnected",D="_OnNamespaceDisconnect",E="_OnRoomJoin",F="_OnRoomJoined",G="_OnRoomLeave",H="_OnRoomLeft",I="_OnNativeMessage",J=function(){function a(){}return a.prototype.isConnect=function(){return"_OnNamespaceConnect"==this.Event||!1},a.prototype.isDisconnect=function(){return this.Event==D||!1},a.prototype.isRoomJoin=function(){return this.Event==E||!1},a.prototype.isRoomLeft=function(){return this.Event==H||!1},a.prototype.isWait=function(){return!g(this.wait)&&(this.wait[0]=="#"||this.wait[0]=="$"||!1)},a.prototype.unmarshal=function(){return JSON.parse(this.Body)},a}(),K=new RegExp(";","g"),L=new RegExp("@%!semicolon@%!","g"),M=function(){function a(a,b){this.nsConn=a,this.name=b}return a.prototype.emit=function(a,b){var c=new J;return c.Namespace=this.nsConn.namespace,c.Room=this.name,c.Event=a,c.Body=b,this.nsConn.conn.write(c)},a.prototype.leave=function(){var a=new J;return a.Namespace=this.nsConn.namespace,a.Room=this.name,a.Event=G,this.nsConn.askRoomLeave(a)},a}(),N=function(){function a(a,b,c){this.conn=a,this.namespace=b,this.events=c,this.rooms=new Map}return a.prototype.emit=function(a,b){var c=new J;return c.Namespace=this.namespace,c.Event=a,c.Body=b,this.conn.write(c)},a.prototype.ask=function(a,b){var c=new J;return c.Namespace=this.namespace,c.Event=a,c.Body=b,this.conn.ask(c)},a.prototype.joinRoom=function(a){return y(this,void 0,void 0,function(){return z(this,function(b){switch(b.label){case 0:return[4,this.askRoomJoin(a)];case 1:return[2,b.sent()];}})})},a.prototype.room=function(a){return this.rooms.get(a)},a.prototype.leaveAll=function(){return y(this,void 0,void 0,function(){var a,b=this;return z(this,function(){return a=new J,a.Namespace=this.namespace,a.Event=H,a.IsLocal=!0,this.rooms.forEach(function(c,d){return y(b,void 0,void 0,function(){var b;return z(this,function(c){switch(c.label){case 0:a.Room=d,c.label=1;case 1:return c.trys.push([1,3,,4]),[4,this.askRoomLeave(a)];case 2:return c.sent(),[3,4];case 3:return b=c.sent(),[2,b];case 4:return[2];}})})}),[2,null]})})},a.prototype.forceLeaveAll=function(a){var b=this,c=new J;c.Namespace=this.namespace,c.Event=G,c.IsForced=!0,c.IsLocal=a,this.rooms.forEach(function(a,d){c.Room=d,p(b,c),b.rooms.delete(d),c.Event=H,p(b,c),c.Event=G})},a.prototype.disconnect=function(){var a=new J;return a.Namespace=this.namespace,a.Event=D,this.conn.askDisconnect(a)},a.prototype.askRoomJoin=function(a){var b=this;return new Promise(function(c,d){return y(b,void 0,void 0,function(){var b,e,f,h;return z(this,function(i){switch(i.label){case 0:if(b=this.rooms.get(a),void 0!==b)return c(b),[2];e=new J,e.Namespace=this.namespace,e.Room=a,e.Event=E,e.IsLocal=!0,i.label=1;case 1:return i.trys.push([1,3,,4]),[4,this.conn.ask(e)];case 2:return i.sent(),[3,4];case 3:return f=i.sent(),d(f),[2];case 4:return(h=p(this,e),!g(h))?(d(h),[2]):(b=new M(this,a),this.rooms.set(a,b),e.Event=F,p(this,e),c(b),[2]);}})})})},a.prototype.askRoomLeave=function(a){return y(this,void 0,void 0,function(){var b,c;return z(this,function(d){switch(d.label){case 0:if(!this.rooms.has(a.Room))return[2,Q];d.label=1;case 1:return d.trys.push([1,3,,4]),[4,this.conn.ask(a)];case 2:return d.sent(),[3,4];case 3:return b=d.sent(),[2,b];case 4:return(c=p(this,a),!g(c))?[2,c]:(this.rooms.delete(a.Room),a.Event=H,p(this,a),[2,null]);}})})},a.prototype.replyRoomJoin=function(a){if(!(g(a.wait)||a.isNoOp)){if(!this.rooms.has(a.Room)){var b=p(this,a);if(!g(b))return a.Err=b.message,void this.conn.write(a);this.rooms.set(a.Room,new M(this,a.Room)),a.Event=F,p(this,a)}this.conn.writeEmptyReply(a.wait)}},a.prototype.replyRoomLeave=function(a){return g(a.wait)||a.isNoOp?void 0:this.rooms.has(a.Room)?void(p(this,a),this.rooms.delete(a.Room),this.conn.writeEmptyReply(a.wait),a.Event=H,p(this,a)):void this.conn.writeEmptyReply(a.wait)},a}(),O=new Error("invalid payload"),P=new Error("bad namespace"),Q=new Error("bad room"),R=new Error("use of closed connection"),S=new Error("write closed"),T=function(){function a(a,b){this.conn=a,this.reconnectTries=0,this._isAcknowledged=!1,this.namespaces=b;var c=b.has("");this.allowNativeMessages=c&&b.get("").has(I),this.queue=[],this.waitingMessages=new Map,this.connectedNamespaces=new Map,this.closed=!1}return a.prototype.wasReconnected=function(){return 0<this.reconnectTries},a.prototype.isAcknowledged=function(){return this._isAcknowledged},a.prototype.handle=function(a){if(!this._isAcknowledged){var b=this.handleAck(a.data);return null==b?(this._isAcknowledged=!0,this.handleQueue()):this.conn.close(),b}return this.handleMessage(a.data)},a.prototype.handleAck=function(a){var b=a[0];switch(b){case"A":var c=a.slice(1);this.ID=c;break;case"H":var d=a.slice(1);return new Error(d);default:return this.queue.push(a),null;}},a.prototype.handleQueue=function(){var a=this;null==this.queue||0==this.queue.length||this.queue.forEach(function(b,c){a.queue.splice(c,1),a.handleMessage(b)})},a.prototype.handleMessage=function(a){var b=m(a,this.allowNativeMessages);if(b.isInvalid)return O;if(b.IsNative&&this.allowNativeMessages){var c=this.namespace("");return p(c,b)}if(b.isWait()){var d=this.waitingMessages.get(b.wait);if(null!=d)return void d(b)}var e=this.namespace(b.Namespace);switch(b.Event){case"_OnNamespaceConnect":this.replyConnect(b);break;case D:this.replyDisconnect(b);break;case E:if(void 0!==e){e.replyRoomJoin(b);break}case G:if(void 0!==e){e.replyRoomLeave(b);break}default:if(void 0===e)return P;b.IsLocal=!1;var f=p(e,b);if(!g(f))return b.Err=f.message,this.write(b),f;}return null},a.prototype.connect=function(a){return this.askConnect(a)},a.prototype.waitServerConnect=function(a){var b=this;return q(this.waitServerConnectNotifiers)&&(this.waitServerConnectNotifiers=new Map),new Promise(function(c){return y(b,void 0,void 0,function(){var b=this;return z(this,function(){return this.waitServerConnectNotifiers.set(a,function(){b.waitServerConnectNotifiers.delete(a),c(b.namespace(a))}),[2]})})})},a.prototype.namespace=function(a){return this.connectedNamespaces.get(a)},a.prototype.replyConnect=function(a){if(!(g(a.wait)||a.isNoOp)){var b=this.namespace(a.Namespace);if(void 0!==b)return void this.writeEmptyReply(a.wait);var c=s(this.namespaces,a.Namespace);return q(c)?(a.Err=P.message,void this.write(a)):void(b=new N(this,a.Namespace,c),this.connectedNamespaces.set(a.Namespace,b),this.writeEmptyReply(a.wait),a.Event=C,p(b,a),!q(this.waitServerConnectNotifiers)&&0<this.waitServerConnectNotifiers.size&&this.waitServerConnectNotifiers.has(a.Namespace)&&this.waitServerConnectNotifiers.get(a.Namespace)())}},a.prototype.replyDisconnect=function(a){if(!(g(a.wait)||a.isNoOp)){var b=this.namespace(a.Namespace);return void 0===b?void this.writeEmptyReply(a.wait):void(b.forceLeaveAll(!0),this.connectedNamespaces.delete(a.Namespace),this.writeEmptyReply(a.wait),p(b,a))}},a.prototype.ask=function(a){var b=this;return new Promise(function(c,d){return b.isClosed()?void d(R):(a.wait=n(),b.waitingMessages.set(a.wait,function(a){return a.isError?void d(new Error(a.Err)):void c(a)}),!b.write(a))?void d(S):void 0})},a.prototype.askConnect=function(a){var b=this;return new Promise(function(c,d){return y(b,void 0,void 0,function(){var b,e,f,h,i;return z(this,function(j){switch(j.label){case 0:if(b=this.namespace(a),void 0!==b)return c(b),[2];if(e=s(this.namespaces,a),q(e))return d(P),[2];if(f=new J,f.Namespace=a,f.Event="_OnNamespaceConnect",f.IsLocal=!0,b=new N(this,a,e),h=p(b,f),!g(h))return d(h),[2];j.label=1;case 1:return j.trys.push([1,3,,4]),[4,this.ask(f)];case 2:return j.sent(),[3,4];case 3:return i=j.sent(),d(i),[2];case 4:return this.connectedNamespaces.set(a,b),f.Event=C,p(b,f),c(b),[2];}})})})},a.prototype.askDisconnect=function(a){return y(this,void 0,void 0,function(){var b,c;return z(this,function(d){switch(d.label){case 0:if(b=this.namespace(a.Namespace),void 0===b)return[2,P];d.label=1;case 1:return d.trys.push([1,3,,4]),[4,this.ask(a)];case 2:return d.sent(),[3,4];case 3:return c=d.sent(),[2,c];case 4:return b.forceLeaveAll(!0),this.connectedNamespaces.delete(a.Namespace),a.IsLocal=!0,[2,p(b,a)];}})})},a.prototype.isClosed=function(){return this.closed},a.prototype.write=function(a){if(this.isClosed())return!1;if(!a.isConnect()&&!a.isDisconnect()){var b=this.namespace(a.Namespace);if(void 0===b)return!1;if(!g(a.Room)&&!a.isRoomJoin()&&!a.isRoomLeft()&&!b.rooms.has(a.Room))return!1}return this.conn.send(k(a)),!0},a.prototype.writeEmptyReply=function(a){this.conn.send(o(a))},a.prototype.close=function(){var a=this;if(!this.closed){var b=new J;b.Event=D,b.IsForced=!0,b.IsLocal=!0,this.connectedNamespaces.forEach(function(c){c.forceLeaveAll(!0),b.Namespace=c.namespace,p(c,b),a.connectedNamespaces.delete(c.namespace)}),this.waitingMessages.clear(),this.closed=!0,this.conn.readyState===this.conn.OPEN&&this.conn.close()}},a}();(function(){var a={dial:u,isSystemEvent:f,OnNamespaceConnect:"_OnNamespaceConnect",OnNamespaceConnected:C,OnNamespaceDisconnect:D,OnRoomJoin:E,OnRoomJoined:F,OnRoomLeave:G,OnRoomLeft:H,OnAnyEvent:"_OnAnyEvent",OnNativeMessage:I,Message:J,Room:M,NSConn:N,Conn:T,ErrInvalidPayload:O,ErrBadNamespace:P,ErrBadRoom:Q,ErrClosed:R,ErrWrite:S,marshal:h};if("undefined"!=typeof c)c=a,b.exports=a;else{var d="object"==typeof self&&self.self===self&&self||"object"==typeof e&&e.global===e&&e;d.neffos=a}})()}).call(this,a("_process"),"undefined"==typeof global?"undefined"==typeof self?"undefined"==typeof window?{}:window:self:global)},{_process:4,"node-fetch":3,ws:5}],3:[function(a,b,c){(function(a){"use strict";var d=function(){if("undefined"!=typeof self)return self;if("undefined"!=typeof window)return window;if("undefined"!=typeof a)return a;throw new Error("unable to locate global object")},a=d();b.exports=c=a.fetch,c.default=a.fetch.bind(a),c.Headers=a.Headers,c.Request=a.Request,c.Response=a.Response}).call(this,"undefined"==typeof global?"undefined"==typeof self?"undefined"==typeof window?{}:window:self:global)},{}],4:[function(a,b){function c(){throw new Error("setTimeout has not been defined")}function d(){throw new Error("clearTimeout has not been defined")}function e(a){if(l===setTimeout)return setTimeout(a,0);if((l===c||!l)&&setTimeout)return l=setTimeout,setTimeout(a,0);try{return l(a,0)}catch(b){try{return l.call(null,a,0)}catch(b){return l.call(this,a,0)}}}function f(a){if(m===clearTimeout)return clearTimeout(a);if((m===d||!m)&&clearTimeout)return m=clearTimeout,clearTimeout(a);try{return m(a)}catch(b){try{return m.call(null,a)}catch(b){return m.call(this,a)}}}function g(){q&&o&&(q=!1,o.length?p=o.concat(p):r=-1,p.length&&h())}function h(){if(!q){var a=e(g);q=!0;for(var b=p.length;b;){for(o=p,p=[];++r<b;)o&&o[r].run();r=-1,b=p.length}o=null,q=!1,f(a)}}function j(a,b){this.fun=a,this.array=b}function k(){}var l,m,n=b.exports={};(function(){try{l="function"==typeof setTimeout?setTimeout:c}catch(a){l=c}try{m="function"==typeof clearTimeout?clearTimeout:d}catch(a){m=d}})();var o,p=[],q=!1,r=-1;n.nextTick=function(a){var b=Array(arguments.length-1);if(1<arguments.length)for(var c=1;c<arguments.length;c++)b[c-1]=arguments[c];p.push(new j(a,b)),1!==p.length||q||e(h)},j.prototype.run=function(){this.fun.apply(null,this.array)},n.title="browser",n.browser=!0,n.env={},n.argv=[],n.version="",n.versions={},n.on=k,n.addListener=k,n.once=k,n.off=k,n.removeListener=k,n.removeAllListeners=k,n.emit=k,n.prependListener=k,n.prependOnceListener=k,n.listeners=function(){return[]},n.binding=function(){throw new Error("process.binding is not supported")},n.cwd=function(){return"/"},n.chdir=function(){throw new Error("process.chdir is not supported")},n.umask=function(){return 0}},{}],5:[function(a,b){'use strict';b.exports=function(){throw new Error("ws does not work in the browser. Browser clients must use the native WebSocket object")}},{}]},{},[1]);
```
> `browserify/client.html`
```html
<!-- 消息的输入 -->

<!-- the message's input -->
<input id="input" type="text" />
<!-- 当单击时，一个websocket事件将被发送到服务器，在此示例中，我们注册了'chat'-->

<!-- when clicked then a websocket event will be sent to the server, at this example we registered the 'chat' -->
<button id="sendBtn" disabled>Send</button>
<!-- 消息将显示在这里 -->

<!-- the messages will be shown here -->
<pre id="output"></pre>

<script src="./bundle.js"></script>
```
> `browserify/package.json`
```json
{
    "name": "neffos.js.example.browserify",
    "version": "0.0.1",
    "scripts": {
        "browserify": "browserify ./app.js -o ./bundle.js",
        "minifyES6": "minify ./bundle.js --outFile ./bundle.js",
        "build": "npm run-script browserify && npm run-script minifyES6"
    },
    "dependencies": {
        "neffos.js": "latest"
    },
    "devDependencies": {
        "browserify": "^16.2.3",
        "babel-minify": "^0.5.0"
    }
}
```
> `go-client/client.go`
```golang
package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kataras/iris/v12/websocket"
)

const (
	endpoint              = "ws://localhost:8080/echo"
	namespace             = "default"
	dialAndConnectTimeout = 5 * time.Second
)

//可以与server.go共享
//NSConn.Conn具有`IsClient() bool`方法，
//可用于检查它是客户端还是服务器端回调

// this can be shared with the server.go's.
// `NSConn.Conn` has the `IsClient() bool` method which can be used to
// check if that's is a client or a server-side callback.
var clientEvents = websocket.Namespaces{
	namespace: websocket.Events{
		websocket.OnNamespaceConnected: func(c *websocket.NSConn, msg websocket.Message) error {
			log.Printf("connected to namespace: %s", msg.Namespace)
			return nil
		},
		websocket.OnNamespaceDisconnect: func(c *websocket.NSConn, msg websocket.Message) error {
			log.Printf("disconnected from namespace: %s", msg.Namespace)
			return nil
		},
		"chat": func(c *websocket.NSConn, msg websocket.Message) error {
			log.Printf("%s", string(msg.Body))
			return nil
		},
	},
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(dialAndConnectTimeout))
	defer cancel()

	// username := "my_username"
	// dialer := websocket.GobwasDialer(websocket.GobwasDialerOptions{Header: websocket.GobwasHeader{"X-Username": []string{username}}})
	dialer := websocket.DefaultGobwasDialer
	client, err := websocket.Dial(ctx, dialer, endpoint, clientEvents)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	c, err := client.Connect(ctx, namespace)
	if err != nil {
		panic(err)
	}

	c.Emit("chat", []byte("Hello from Go client side!"))

	fmt.Fprint(os.Stdout, ">> ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			log.Printf("ERROR: %v", scanner.Err())
			return
		}

		text := scanner.Bytes()

		if bytes.Equal(text, []byte("exit")) {
			if err := c.Disconnect(nil); err != nil {
				log.Printf("reply from server: %v", err)
			}
			break
		}

		ok := c.Emit("chat", text)
		if !ok {
			break
		}

		fmt.Fprint(os.Stdout, ">> ")
	}
	//尝试两次运行该程序，或者/并且运行服务器的http://localhost:8080来检查浏览器客户端
} // try running this program twice or/and run the server's http://localhost:8080 to check the browser client as well.
```
> `nodejs-client/client.js`
```javascript
const neffos = require('neffos.js');
const stdin = process.openStdin();

const wsURL = "ws://localhost:8080/echo";

async function runExample() {
  try {
    const conn = await neffos.dial(wsURL, {
      default: { // "default" namespace.
        _OnNamespaceConnected: function (nsConn, msg) {
          console.log("connected to namespace: " + msg.Namespace);
        },
        _OnNamespaceDisconnect: function (nsConn, msg) {
          console.log("disconnected from namespace: " + msg.Namespace);
        },
        chat: function (nsConn, msg) { // "chat" event.
          console.log(msg.Body);
        }
      }
    });

    const nsConn = await conn.connect("default");
    nsConn.emit("chat", "Hello from Nodejs client side!");

    stdin.addListener("data", function (data) {
      const text = data.toString().trim();
      nsConn.emit("chat", text);
    });

  } catch (err) {
    console.error(err);
  }
}

runExample();
```
> `nodejs-client/package.json`
```json
{
    "name": "neffos.js.example.nodejsclient",
    "version": "0.0.1",
    "main": "client.js",
    "dependencies": {
        "neffos.js": "latest"
    }
}
```
> `server.go`
```golang
package main

import (
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	//当"enableJWT"常量为true时使用：

	// Used when "enableJWT" constant is true:
	"github.com/iris-contrib/middleware/jwt"
)
//值也应与客户端匹配

// values should match with the client sides as well.
const enableJWT = true
const namespace = "default"

//如果名称空间为空，则可以简单地使用websocket.Events {...}

// if namespace is empty then simply websocket.Events{...} can be used instead.
var serverEvents = websocket.Namespaces{
	namespace: websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			//通过`websocket.GetContext`可以获取iris的`Context`

			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			ctx := websocket.GetContext(nsConn.Conn)

			log.Printf("[%s] connected to namespace [%s] with IP [%s]",
				nsConn, msg.Namespace,
				ctx.RemoteAddr())
			return nil
		},
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			return nil
		},
		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// room.String() 返回 -> NSConn.String() 返回 -> Conn.String() 返回 -> Conn.ID()

			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))
			// 使用以下命令将消息写回客户端消息所有者：
			// nsConn.Emit("chat", msg)
			// 使用以下命令向除此客户端之外的所有用户写消息：

			// Write message back to the client message owner with:
			// nsConn.Emit("chat", msg)
			// Write message to all except this client with:
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	},
}

func main() {
	app := iris.New()
	websocketServer := websocket.New(
		/* 也可以使用DefaultGobwasUpgrader */
		websocket.DefaultGorillaUpgrader, /* DefaultGobwasUpgrader can be used too. */
		serverEvents)

	j := jwt.New(jwt.Config{
		//通过“token” URL提取，
		//因此客户端应使用ws://localhost:8080/echo?token=$token请求

		// Extract by the "token" url,
		// so the client should dial with ws://localhost:8080/echo?token=$token
		Extractor: jwt.FromParameter("token"),

		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// 设置后，中间件将验证令牌是否已使用特定的签名算法进行签名
		// 如果签名方法不是恒定的，则可使用`Config.ValidationKeyGetter`回调字段来实施其他检查，
		// 这对于避免此处所述的安全问题很重要：
		// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/

		// When set, the middleware verifies that tokens are signed
		// with the specific signing algorithm
		// If the signing method is not constant the
		// `Config.ValidationKeyGetter` callback field can be used
		// to implement additional checks
		// Important to avoid security issues described here:
		// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	idGen := func(ctx iris.Context) string {
		if username := ctx.GetHeader("X-Username"); username != "" {
			return username
		}

		return websocket.DefaultIDGenerator(ctx)
	}
	//通过可选的自定义ID生成器为ws://localhost:8080/echo的端点提供服务

	// serves the endpoint of ws://localhost:8080/echo
	// with optional custom ID generator.
	websocketRoute := app.Get("/echo", websocket.Handler(websocketServer, idGen))

	if enableJWT {
		//注册jwt中间件（握手时）：

		// Register the jwt middleware (on handshake):
		websocketRoute.Use(j.Serve)

		// 或|OR
		//
		//通过websocket连接或任何事件通过jwt中间件检查令牌：
		//
		// Check for token through the jwt middleware
		// on websocket connection or on any event:

		/* websocketServer.OnConnect = func(c *websocket.Conn) error {
		ctx := websocket.GetContext(c)
		if err := j.CheckJWT(ctx); err != nil {
			//将在客户端上发送上述错误，并且完全不允许它连接到websocket服务器

			// will send the above error on the client
			// and will not allow it to connect to the websocket server at all.
			return err
		}

		user := ctx.Values().Get("jwt").(*jwt.Token)
		// or just: user := j.Get(ctx)

		log.Printf("This is an authenticated request\n")
		log.Printf("Claim content:")
		log.Printf("%#+v\n", user.Claims)

		log.Printf("[%s] connected to the server", c.ID())

		return nil
		} */
	}
	//提供基于浏览器的websocket客户端

	// serves the browser-based websocket client.
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("./browser/index.html", false)
	})
	//提供npm浏览器websocket客户端用法示例

	// serves the npm browser websocket client usage example.
	app.HandleDir("/browserify", "./browserify")

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
```