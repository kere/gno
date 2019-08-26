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
      methods : {
        _onClick : function(){
					util.$.addClass('#box2', 'm-t-lg');
        },
        _onClick2 : function(){
					util.$.removeClass('#box2', 'm-t-lg');
        },
        _onClick3 : function(){
					util.tool.showSuccess(5);
        },
        _onClick4 : function(){
					util.tool.hideToast();
        },
        _onClick5 : function(){
					util.$.show('#box2');
        },
        _onClick6 : function(){
					util.$.hide('#box2');
        },
        _onClick7 : function(){
					util.tool.viewImage('http://n.sinaimg.cn/blog/250/w640h410/20190826/f58a-icuacrz5631047.jpg');
        },
        _onClick8 : function(e){
					util.tool.taggle(e);
        },
			},
      mounted : function(){
    		var client = ajax.NewClient("/openapi/app");
				client.timeout = 3;
		    client.send("PageData", {name:'tom', age: 22}).then(function(result){
		      console.log("-----1-----", result)
		    })

		    client.send("PageData", null).then(function(result){
		      console.log("-----2-----", result)
		    })
      }
    })
});
