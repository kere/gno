require.config({
	waitSeconds :15,
	paths: {
		'util' : MYENV+'/mylib/util',
		'accto' : MYENV+'/mylib/accto',
    'zepto' : MYENV+'/mylib/zepto'
	}
});
define(
  'ajax',
  ['util', 'accto', 'zepto'],

  function(util, accto) {
    var ajax  = {
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
        }
      },

      getHostName : function(){
        var hostArr = window.location.host.split('.')
        return hostArr[hostArr.length-2]+'.'+hostArr[hostArr.length-1]
      },

      getUrlVars: function(){
        var vars = [], hash;
        var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
        for(var i = 0; i < hashes.length; i++) {
        hash = hashes[i].split('=');
        vars.push(hash[0]);
        vars[hash[0]] = hash[1];
        }
        return vars;
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
        this.path = path || "/api/web";
        this.isrun = false;
        this.timeout= 300;
        this.pfield = 'accpt';
    }

    Client.prototype.errorHandler = function(r){console.log(r);};

    // 排他性运行
    Client.prototype.sendEx = function(method, args, opt){
      if(this.deferred && this.deferred.state()=='pending'){
          var func = function(){}
          return {done:func, fail:func, always:func};
      }
      this.deferred = this.send(method, args, opt);
      return this.deferred;
    };

    Client.prototype.send = function(method, args, opt){
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

      var ts = (new Date()).getTime().toString(),
        jsonStr = args ? JSON.stringify(args) : '',
        token = window.atob(decodeURI(util.getCookie(this.pfield))),
        str = ts+method+jsonStr + token;

      return $.ajax({
        url:        this.path + '/' + method,
        type:       'POST',
        dataType:   'json',
        cache:      false,
        data:       {'_src': jsonStr, 'method': method},
        headers: {'Accto':accto(str), 'Accts': ts},
        timeout:    this.timeout * 1000

      }).always(function(){
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

      }).fail(function (jqXHR, textStatus, errorThrown){
        if(!jqXHR) return;

        switch(jqXHR.status){
          case 599:
            if(clas.errorHandler){
              try{
                  clas.errorHandler(JSON.parse(jqXHR.responseText))
              }catch(e){
                  console.log(e)
              }
            }
            break;
          case 404:
            console.log(textStatus+': api not found!')
            break;
          default:
          if(jqXHR.responseText)
            console.log(textStatus+': '+jqXHR.responseText);
          else
            console.log('请稍后再试');
            break;
        }
      });
    }

    return ajax;
  }
)
