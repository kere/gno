require.config({
	waitSeconds :30,
	baseUrl : "/assets/js",
	paths: {
		'Compressor' : 'compressorjs',
		'imageUpload': 'vue-imageUpload'
	}
});

require(
	['ajax', 'util', 'imageUpload', 'Compressor'],
	function (ajax, util, imageUpload, Compressor){
    var client = ajax.NewClient("/openapi/app");

    if ("WebSocket" in window) {
       // 打开一个 web socket
			 var ws = ajax.NewWS('/ws');
			 ws.receive = (method, args, result) =>{
				 console.log(method, args, result);
			 }
			 ws.onclose = () => {
				 alert("closed");
			 }

       window.sendTo = function(){
				 ws.Send('SayHi', {name: "tome", msg: "hello"})
       }

			 ws.Connect();
    } else {
       // 浏览器不支持 WebSocket
       alert("您的浏览器不支持 WebSocket!");
    }


    var main = new Vue({
      el : '#main-div',
			components : {
				'image-upload' : imageUpload
			},
      data: {
      },
      methods : { },
      mounted : function(){
		    client.send("PageData", {name:'tom', age: 22}).then(function(result){
		      console.log(result)
		    })
      }
    })
});
