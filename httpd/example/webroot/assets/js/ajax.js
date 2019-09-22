define(
  'ajax',
  ['util', 'accto'],

  function(util, accto) {

    // options.timeout,options.url, options.async
    function ajaxFunc(options) {
      // XHR
      let xhr = new XMLHttpRequest();
      xhr.open('post', options.url, options.async);
      xhr.timeout = options.timeout * 1000;

      //设置请求头
      //   xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8");
      xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");

      for (let key in options.headers) {
        xhr.setRequestHeader(key, options.headers[key]);
      }

      let promise = new Promise(function(resolve, reject){
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

    let comp  = {
			url : "/home/openapi/data",
      NewClient : function(path, timeout){
        return new Client(path, timeout);
      },
      NewUpload : function(path){
        return new Upload(path);
      },

      DataSet: function(dat){
        this.Fields = dat.fields;
        this.Columns = dat.columns;
        this.Len = function(){
          if(this.Columns.length == 0 ) return 0;
          return this.Columns[0].length;
        }

        this.FieldI = function(name){
          return this.Fields.indexOf(name);
        }

        this.RowAt = function(index){
          let l = this.Len(), row = {};
          for (var k = 0; k < this.Fields.length; k++) {
            row[this.Fields[k]] = this.Columns[k][index]
          }
          return row;
        }
      },

      torows : (dat, callback) => {
        let n = dat.fields.length;
        if(!dat || !dat.columns || dat.columns.length==0 || n===0 || !dat.columns[0]){
          return [];
        }

        let cols = dat.columns;
        let l = cols[0].length;
        if(l === 0) return [];

        let rows = new Array(l);
        let i,k, key, obj, fields = dat.fields;

        for (i = 0; i < l; i++) {
          obj = {}
          for (k = 0; k < n; k++) {
            key = fields[k];
            if(callback){
              obj[key] = callback(k, i);
            }else{
              obj[key] = cols[k][i];
            }
          }
          rows[i] = obj;
        }

        return rows;
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
          let d = this.now();
          return d.getTime() - d.getTimezoneOffset();
        }
      },

      getUrlVar: function(sParam){
        let sPageURL = window.location.search.substring(1);
        let sURLVariables = sPageURL.split('&');
        for (let i = 0; i < sURLVariables.length; i++) {
          let sParameterName = sURLVariables[i].split('=');
          if (sParameterName[0] == sParam)
          {
              return decodeURI(sParameterName[1]);
          }
        }
        return null;
      }
    };

    let Client = function(path){
      this.path = path || comp.url;
      this.isrun = false;
      this.timeout= 10;
      this.async = true;
      this.pfield = 'accpt';
      this.verField = "_data_version";
    }

    Client.prototype.trySetDataVer = function(result) {
      if(!result || typeof(result) != "object" || !result[this.verField]){
        return;
      }
      let ver = result[this.verField];
      this[this.verField] = ver;
      window[this.verField] = ver;
    }

    // GetData 先从local缓存里查看，如果没有、或版本不匹配，则发送ajax抓取
    // return Promise()
    Client.prototype.getData = function(method, args, opt) {
      return new Promise((resolve, reject) => {
        let key = method + (args ? JSON.stringify(args): '');

        let doit = (resolve, reject) =>{
          this.send(method, args, opt).then(result =>{
            if(result && result[this.verField]){
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

        let src = window.localStorage.getItem(key);
        if(!src){
          doit(resolve, reject);
          return;
        }

        let dat = JSON.parse(src);
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

      let ts = opt.ts ? opt.ts : comp.serverTime.utctime().toString(),
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
        if(opt.busy){
          util.$.each(opt.busy, (e) => {
            util.tool.hideBusy(e);
          })
        }

        return Promise.reject(err);
      });

    }

    let Upload = function(path){
        this.path = path;
        this.pfield = 'accpt';
        this.verField = "_data_version";
    }

    Upload.prototype.upload = function(file, opt){
      let xhr = new XMLHttpRequest();
      let ts = comp.serverTime.utctime().toString(), ptoken = window[this.pfield] || '';
      xhr.open('POST', this.path, true);

      let str = ts+file.name +  file.size + file.lastModified + file.type+navigator.userAgent+ts+ptoken + window.location.hostname;
      // console.log(str);
      // xhr.setRequestHeader("Content-Type", "multipart/form-data");
      xhr.setRequestHeader('Accto', accto(str));
      xhr.setRequestHeader('Accts', ts);
      xhr.setRequestHeader('AccPage', ptoken);
      let promise = new Promise(function(resolve, reject){
        xhr.onload = function(e) {
          if (xhr.status < 200 || xhr.status >= 300) {
            return reject(e);
          }
          resolve(e.currentTarget.responseText);
          if(opt && opt.onSuccess) opt.onSuccess(e.currentTarget.responseText, file);
        };

        xhr.onerror = function(e) {
          reject(e);
        };
      });

      if(opt && opt.onProgress) {
        xhr.upload.onprogress = function(e) {
          if (e.total > 0) {
            e.percent = e.loaded / e.total * 100;
          }
          opt.onProgress(e, file);
        };
      }

      let formData = new FormData();
			formData.append('filename', (opt && opt.filename) ? opt.filename : "");
			formData.append('file', file);
			formData.append('name', file.name);
			formData.append('size', file.size);
			formData.append('lastModified', file.lastModified);
			formData.append('type', file.type);
      xhr.send(formData);

      return promise;
    };

    let WS = function(path){
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
      let ths = this;
      let sign = accto(navigator.userAgent + window.location.host),
          ws = new WebSocket('ws://'+window.location.host+this.path+"?url="+encodeURI(document.location.pathname)+"&sign="+sign);

      ws.onopen = this.onopen;
      ws.onclose = this.onclose;
      ws.onmessage = (e) => {
        let obj = JSON.parse(e.data);
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
