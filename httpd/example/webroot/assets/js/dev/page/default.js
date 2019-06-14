require.config({
	waitSeconds :30,
	baseUrl : "/assets/js/",
	paths: {

	}
});

require(
	['ajax', 'util'],
	function (ajax, util){
    var client = ajax.NewClient("/openapi/app");
    client.send("PageData", {name:'tom', age: 22}).then(function(result){
      console.log(result)
    })

    // if ("WebSocket" in window) {
    //    // 打开一个 web socket
    //    var ws = new WebSocket("ws://localhost:9000/ws");
		//
    //    ws.onopen = function() {
    //       // Web Socket 已连接上，使用 send() 方法发送数据
    //       var obj = {"method":"abc", "name":"tom"};
    //       ws.send(JSON.stringify(obj));
    //       // alert("数据发送中...");
    //    };
		//
    //    ws.onmessage = function (evt) {
    //       var received_msg = evt.data;
    //       alert("接收:" + evt.data);
    //    };
		//
    //    ws.onclose = function() {
    //       // 关闭 websocket
    //       alert("连接关闭...");
    //    };
		//
    //    $('#btnSendFrom').click(function(){
    //      client.send("ServerSend", {message: "this is from server."});
    //    })
    //    $('#btnSendTo').click(function(){
    //      var obj = {"method":"toServer", "content":"send to server message"};
    //      ws.send(JSON.stringify(obj));
    //    })
    // } else {
    //    // 浏览器不支持 WebSocket
    //    alert("您的浏览器不支持 WebSocket!");
    // }
});
