require.config({
	waitSeconds :30,
	baseUrl : "/assets/js/",
	paths: {

	}
});

require(
	['ajax'],
	function (ajax, util){
    var client = ajax.NewClient("/openapi/app");
    client.send("PageData", {name:'tom', age: 22}).done(function(result){

    })

});
