define(
  'ajax',
  ['util', 'accto'],

  function(util, accto) {

    // options.timeout,options.url, options.async
    function ajaxFunc(options) {
      // XHR
      var xhr = new XMLHttpRequest();
      xhr.open('post', options.url, options.async);
      xhr.timeout = options.timeout * 1000;

      //设置请求头
      //   xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8");
      xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");

      for (var key in options.headers) {
        xhr.setRequestHeader(key, options.headers[key]);
      }

      var promise = new Promise(function(resolve, reject){
        xhr.onload = (e) => {
          if (xhr.status == 200) {
            resolve(xhr.response);
          } else {
            reject(xhr.response, xhr.status);
          }
        }

        xhr.onerror = (e) => {
          reject(xhr);
        }
        xhr.ontimeout = () => {
          reject('timeout');
        }
      });

      //发送请求
      if(options.data){
        xhr.send(JSON.stringify(options.data));
      }else{
        xhr.send();
      }

      return promise;
    }

    var comp  = {
			url : "/home/openapi/data",
      NewClient : function(path, timeout){
        return new Client(path, timeout);
      },
      NewUpload : function(path){
        return new Upload(path);
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
      }
    };

    var Client = function(path){
      this.path = path || comp.url;
      this.isrun = false;
      this.timeout= 10;
      this.async = true;
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
        util.tool.showLoading(this.timeout);
      }
      if(opt.disable){
        util.$.each(opt.disable, (e) => {
          e.setAttribute("disabled", true);
        })
      }
      if(opt.busy){
        util.$.each(opt.busy, (e) => {
          util.tool.showBusy(e, this.timeout);
        })
      }

      var ts = comp.serverTime.utctime().toString(),
        ptoken = window[this.pfield] || '',
      	// method + ts + src + agent + ts + ptoken + window.location.hostname
        str = method+ts+(args?JSON.stringify(args): '')+navigator.userAgent+ts+ptoken + window.location.hostname;

      return ajaxFunc({
        async:    this.async,
        url:      this.path + '/' + method,
        data:     args,
        headers:  {"Accto":accto(str), "Accts": ts, "AccPage": ptoken, "Api": method},
        timeout:  this.timeout * 1000

      }).then(result =>{
        clas.isrun = false;
        if(opt.loading){
          util.tool.hideToast();
        }
        if(opt.disable){
          util.$.each(opt.disable, (e) => {
            e.removeAttribute("disabled");
          })
        }
        if(opt.busy){
          util.$.each(opt.busy, (e) => {
            util.tool.hideBusy(e);
          })
        }

        if(result=="") return null;

				if(typeof(result)=='string'){
        	result = JSON.parse(result);
				}

        clas.trySetDataVer(result);

				return result;
      }).catch(function (err, status){
        clas.isrun = false;
        if(opt.loading){
          util.tool.hideToast();
        }
        if(opt.disable){
          util.$.each(opt.disable, (e) => {
            e.removeAttribute("disabled");
          })
        }

        return Promise.reject(err);
      });

    }

    var Upload = function(path){
        this.path = path;
        this.pfield = 'accpt';
        this.verField = "_data_version";
    }

    Upload.prototype.upload = function(blob, filename){
      var xhr = new XMLHttpRequest();
      var ts = comp.serverTime.utctime().toString(), ptoken = window[this.pfield] || '';
      xhr.open('POST', this.path, true);

      var str = ts+blob.name +  blob.size + blob.lastModified + blob.type+navigator.userAgent+ts+ptoken + window.location.hostname;
      console.log(str);
      // xhr.setRequestHeader("Content-Type", "multipart/form-data");
      xhr.setRequestHeader('Accto', accto(str));
      xhr.setRequestHeader('Accts', ts);
      xhr.setRequestHeader('AccPage', ptoken);
      var promise = new Promise(function(resolve, reject){
        xhr.onload = function(e) {
          resolve(e);
        };
        xhr.onerror = function(e) {
          reject(e);
        };
      });

      var formData = new FormData();
			formData.append('filename', filename ? filename : "");
			formData.append('file', blob);
			formData.append('name', blob.name);
			formData.append('size', blob.size);
			formData.append('lastModified', blob.lastModified);
			formData.append('type', blob.type);
      xhr.send(formData);

      return promise;
    };

    var WS = function(path){
        this.path = path;
    }

    WS.prototype.onclose = function(r){
      if(!this.ws) return false;
      ws.close();
    };
    WS.prototype.onopen = function(e){
      console.log('open');
    };
    WS.prototype.receive = function(method, args, result){
    };
    WS.prototype.error = function(msg, e){
      console.log(msg);
    };

    WS.prototype.path = '';
    WS.prototype.Connect = function(){
      if(this.conn){
        this.conn.close();
        this.conn = null;
      }
      var ths = this;
      var sign = accto(navigator.userAgent + window.location.host),
          ws = new WebSocket('ws://'+window.location.host+this.path+"?url="+encodeURI(document.location.pathname)+"&sign="+sign);

      ws.onopen = this.onopen;
      ws.onclose = this.onclose;
      ws.onmessage = (e) => {
        var obj = JSON.parse(e.data);
        if (obj['iserror']){
          ths.error(obj.error, e);
          return;
        }
        if(obj.method){
          ths.receive(obj.method, obj.args, obj.result);
        }

      };

      this.conn = ws;
    };

    WS.prototype.Send = function(method, args){
      if(this.conn.readyState==1){
        this.conn.send(JSON.stringify({"method":method, "args": args}))
      }
    };

    comp.NewWS = function(path){
      return new WS(path);
    }

    return comp;
  }
)
