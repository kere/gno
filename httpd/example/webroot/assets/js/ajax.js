define(
  'ajax',
  ['util', 'accto'],

  function(util, accto) {
    function ajaxFunc(options) {
        //编码数据
        function setData() {
          //设置对象的遍码
          function setObjData(data, parentName) {
            function encodeData(name, value, parentName) {
              var items = [];
              name = parentName === undefined ? name : parentName + "[" + name + "]";
              if (typeof value === "object" && value !== null) {
                items = items.concat(setObjData(value, name));
              } else {
                name = encodeURIComponent(name);
                value = encodeURIComponent(value);
                items.push(name + "=" + value);
              }
              return items;
            }
            var arr = [],value;
            if (Object.prototype.toString.call(data) == '[object Array]') {
              for (var i = 0, len = data.length; i < len; i++) {
                value = data[i];
                arr = arr.concat(encodeData( typeof value == "object"?i:"", value, parentName));
              }
            } else if (Object.prototype.toString.call(data) == '[object Object]') {
              for (var key in data) {
                value = data[key];
                arr = arr.concat(encodeData(key, value, parentName));
              }
            }
            return arr;
          };
          //设置字符串的遍码，字符串的格式为：a=1&b=2;
          function setStrData(data) {
              var arr = data.split("&");
              for (var i = 0, len = arr.length; i < len; i++) {
                  name = encodeURIComponent(arr[i].split("=")[0]);
                  value = encodeURIComponent(arr[i].split("=")[1]);
                  arr[i] = name + "=" + value;
              }
              return arr;
          }

          if (data) {
            if (typeof data === "string") {
                data = setStrData(data);
            } else if (typeof data === "object") {
                data = setObjData(data);
            }
            data = data.join("&").replace("/%20/g", "+");
            //若是使用get方法或JSONP，则手动添加到URL中
            if (type === "get") {
                url += url.indexOf("?") > -1 ? (url.indexOf("=") > -1 ? "&" + data : data) : "?" + data;
            }
          }
        }

        //设置请求超时
        function TimeOut(n){
          var p = new Promise(function(resolve, reject){
            setTimeout(function(){
              reject('timeout');
            }, n);
          });
          return p;
        }

        // XHR
        function createXHR() {
          //由于IE6的XMLHttpRequest对象是通过MSXML库中的一个ActiveX对象实现的。
          //所以创建XHR对象，需要在这里做兼容处理。
          function getXHR() {
            if (window.XMLHttpRequest) {
              return new XMLHttpRequest();
            } else {
              //遍历IE中不同版本的ActiveX对象
              var versions = ["Microsoft", "msxm3", "msxml2", "msxml1"];
              for (var i = 0; i < versions.length; i++) {
                try {
                  var version = versions[i] + ".XMLHTTP";
                  return new ActiveXObject(version);
                } catch (e) {

                }
              }
            }
          }
          //创建对象。
          xhr = getXHR();
          xhr.open(type, url, async);
          //设置请求头
          if (type === "post" && !contentType) {
            //若是post提交，则设置content-Type 为application/x-www-four-urlencoded
            xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8");
          } else if (contentType) {
            xhr.setRequestHeader("Content-Type", contentType);
          }

          for (var key in headers) {
            xhr.setRequestHeader(key, headers[key]);
          }

          var promise = new Promise(function(resolve, reject){
            //添加监听
            xhr.onreadystatechange = function() {
              if (xhr.readyState === 4) {
                if ((xhr.status >= 200 && xhr.status < 300) || xhr.status == 304) {
                  resolve(xhr.response);
                } else {
                  reject(xhr);
                }
              }
            };
            //发送请求
            xhr.send(type === "get" ? null : data);
          });
          return promise;
        }


        var url = options.url || "", //请求的链接
            type = (options.type || "get").toLowerCase(), //请求的方法,默认为get
            data = options.data || {}, //请求的数据
            contentType = options.contentType || "", //请求头
            dataType = options.dataType || "", //请求的类型
            headers = options.headers || [], //请求headers
            async = options.async === undefined ? true : options.async, //是否异步，默认为true.
            timeOut = options.timeout || 30000, //超时时间。
            // before = options.before || function() {}, //发送之前执行的函数
            // error = options.error || function() {}, //错误执行的函数
            success = options.success || function() {}; //请求成功的回调函数

        var xhr = null; //xhr对角

        setData();
        return Promise.race([createXHR(), TimeOut(timeOut)])
    }

    var comp  = {
			url : "/home/openapi/data",
      NewClient : function(path, timeout){
        return new Client(path, timeout);
      },

      serverTime : {
        diff : 0,
        set : function(unix){
          this.diff = unix>0 ? (new Date()).getTime() - unix*1000 : 0;
        },
        now : function(){
          return new Date(this.time());
        },
        time : function(){
          return (new Date()).getTime() - this.diff;
        },
        utctime : function(){
          var d = this.now();
          return d.getTime() - d.getTimezoneOffset();
        }
      },

      getHostName : function(){
        var hostArr = window.location.host.split('.')
        return hostArr[hostArr.length-2]+'.'+hostArr[hostArr.length-1]
      },

      getUrlVar: function(sParam){
        var sPageURL = window.location.search.substring(1);
        var sURLVariables = sPageURL.split('&');
        for (var i = 0; i < sURLVariables.length; i++) {
          var sParameterName = sURLVariables[i].split('=');
          if (sParameterName[0] == sParam)
          {
              return decodeURI(sParameterName[1]);
          }
        }
        return null;
      },

      datasetDecode : function(data){
        if(!data || data.length==0)
            return data

        var fields=data[0], l = data.length, obj, value, items=[]
        for(var i=1; i < l; i++) {
            obj = new Object()
            for(var k=0; k < fields.length; k++)
                obj[fields[k]] = data[i][k]
            items.push(obj)
        }
        return items;
      },
    };

    // ajax.setCookie('lang', navigator.language || navigator.userLanguage);
    var Client = function(path){
        this.path = path || comp.url;
        this.isrun = false;
        this.timeout= 10;
        this.pfield = 'accpt';
        this.verField = "_data_version";
    }

    Client.prototype.trySetDataVer = function(result) {
      if(typeof(result) != "object" || !result[this.verField]){
        return;
      }
      var ver = result[this.verField];
      this[this.verField] = ver;
      window[this.verField] = ver;
    }

    // GetData 先从local缓存里查看，如果没有、或版本不匹配，则发送ajax抓取
    // return Promise()
    Client.prototype.getData = function(method, args, opt) {
      return new Promise((resolve, reject) => {
        var key = method + (args ? JSON.stringify(args): '');
        var doit = (resolve, reject) =>{
          this.send(method, args, opt).then(result =>{
            if(result[this.verField]){
              window.localStorage.setItem(key, JSON.stringify(result));
            }
            resolve(result);
          }).catch(err => {
            reject(err)
          });
        }

        if(!window.localStorage){
          doit(resolve, reject);
          return;
        }

        var src = window.localStorage.getItem(key);
        if(!src){
          doit(resolve, reject);
          return;
        }

        var dat = JSON.parse(src);
        if(!dat){
          doit(resolve, reject);
          return;
        }
        opt = opt || {ver: window[this.verField]};
        if(dat[this.verField] != opt.ver){
          doit(resolve, reject);
          return;
        }

        resolve(dat);
      });
    };

    // 排他性运行
    Client.prototype.sendEx = function(method, args, opt){
      if(this.isrun) return;
      return this.send(method, args, opt);
    };

    Client.prototype.send = function(method, args, opt) {
      var clas = this;
      this.isrun = true;
      opt = opt || {}

      if(opt.loading){
        util.showLoading();
      }

      if(opt.target){
        var $t = opt.target;
        if(typeof $t == 'string'){
          $t = $($t);
        }
        $t.addClass('weui-btn_loading');
      }

      var ts = comp.serverTime.utctime().toString(),
        jsonStr = args ? JSON.stringify(args) : '',
        ptoken = window[this.pfield] || '',
        str = ts+method+ts+jsonStr + ptoken;

      var promise = ajaxFunc({
          url:        this.path + '/' + method,
          type:       'post',
          dataType:   'json',
          cache:      false,
          data:       {'_src': jsonStr, 'method': method},
          headers: {'Accto':accto(str), 'Accts': ts, 'AccPage': ptoken},
          timeout:    this.timeout * 1000

        }).then(result =>{
          clas.isrun = false;
          if(opt.loading){
            util.hideToast();
          }

          if(opt.target){
            var $t = opt.target;
            if(typeof $t == 'string'){
              $t = $($t);
            }
            $t.removeClass('weui-btn_loading');
          }
          if(result=="") return null;

					if(typeof(result)=='string'){
          	result = JSON.parse(result);
					}
          clas.trySetDataVer(result);

					return result;

        }).catch(function (xhr){
					var status = xhr.status;
          clas.isrun = false;
          if(opt.loading){
            util.hideToast();
          }

          if(opt.target){
            var $t = opt.target;
            if(typeof $t == 'string'){
              $t = $($t);
            }
            $t.removeClass('weui-btn_loading');
          }
					var error;
          switch(status){
          case 599:
          case 500:
      			error = xhr.response;
						break;
          case 404:
            error = 'api not found';
            break;
          case 502:
            error = '内部错误';
            break;
					default:
						error = status+': '+xhr.response
          }
          return Promise.reject(error);
        });

        return promise;
    }

    return comp;
  }
)
