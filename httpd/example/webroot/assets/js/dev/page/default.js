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

    if ("WebSocket" in window) {
       // 打开一个 web socket
			 var ws = ajax.NewWS('/ws');
			 ws.receive = (method, args, result) =>{
				 console.log(method, args, result);
			 }

       window.sendTo = function(){
				 ws.Send('SayHi', {name: "tome", msg: "hello"})
       }

			 ws.Connect();
    } else {
       // 浏览器不支持 WebSocket
       alert("您的浏览器不支持 WebSocket!");
    }
});
