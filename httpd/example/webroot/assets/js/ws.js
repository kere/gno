define(
  'websock',
  ['util','accto'],

  function(util, accto) {
    var Client = function(path){
        this.path = path;
        this.pfield = 'accpt';
    }

    Client.prototype.error = function(msg, e){
      console.log(msg);
    };
    Client.prototype.onclose = function(r){
      console.log(r);
    };

    Client.prototype.path = '';
    Client.prototype.connect = function(){
      var ws = new WebSocket(this.path+"?url="+encodeURI(document.location.pathname));
      var ths = this;
      ws.onclose = this.onclose;
      ws.onmessage = function(e){
        var obj = JSON.parse(e.data);
        if (obj['errmsg']){
          ths.onerror(obj.errmsg, e);
          return;
        }

        ths.onmessage(obj, e);
      };
      ws.onopen = this.onopen;

      this.conn = ws;
    };

    Client.prototype.onmessage = function(e){

    };

    Client.prototype.onopen = function(e){

    };

    Client.prototype.Send = function(args){
      if(this.conn.readyState==1){
        this.conn.send(JSON.stringify(args))
      }
    };

    var comp = {
      New : function(path){
        return new Client(path);
      }
    }

    return comp;
  }
)
