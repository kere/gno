require.config({
	waitSeconds :30,
	paths: {
    'zepto' : MYENV+'/mylib/zepto'
	}
});
define('util', ['zepto'], function(){
	var util = {}

	util.language = function(){
	var lang = util.getCookie('lang');
		if(!lang){
			lang = navigator.language || navigator.userLanguage;
		}
		return lang.replace('_', '-').toLowerCase();
	}

  util.money = function(v){
    return util.numStr(v, 2);
  }

  util.numStr = function(v, deci){
    if(!v) return (0).toFixed(deci);

    if(typeof(v)=="string"){
      v = parseFloat(v);
    }
    if(deci==0){
      return v.toFixed(deci);
    }
    var deciV = Math.pow(10, deci);
    var tmp = Math.round(v * deciV)/100;
    return tmp.toFixed(deci);
  }

  util.weekCH = function(v){
		if(v==null || typeof v =='undefined') return '';
    var weeks = ['天', '一', '二', '三', '四', '五', '六' ];
    switch (typeof v) {
      case "number":
        return weeks[parseInt(v)];
      default:
        return weeks[this.str2date(v).getDay()];
    }
  }

	util.DATE_DAY = 86400000;
	util.DATE_HOUR = 3600000;

	util.timeAgoStr = function (b, e){
    var arr = this.timeAgo(b, e), str='';
    var l = arr.length, isskip = true;
    for (var i = l-1; i > -1; i--) {
      if(isskip && arr[i].value==0) continue;
      isskip = false;
      str += arr[i].value + arr[i].label + ' ';
    }
    return str;
  }
	util.timeAgo = function (b, e){
    var diff;
    if(typeof(b)=='number' && typeof(e) == 'undefined'){
      diff = b;
    }else{
      diff = Math.abs(e.getTime() - b.getTime());
    }
    var v, n, arr=[];

    for(var i=0;i<5;i++) {
      switch (i) {
        case 0: // second
          v = diff % 60000;
          arr.push({value:Math.floor(v/1000), 'label':'秒', 'ext': 's'});
          diff= Math.floor(diff/60000);
          break;
        case 1: // 分钟
          v = diff % 60;
          arr.push({value:v, 'label':'分', 'ext': 'm'});
          diff= Math.floor(diff/60);
          break;
        case 2: //小时
          v = diff % 24;
          arr.push({value:v, 'label':'小时', 'ext': 'h'});
          diff= Math.floor(diff/24);
          break;
        case 3: //天
          v = diff % 30;
          arr.push({value:v, 'label':'天', 'ext': 'd'});
          diff= Math.floor(diff/30);
          break;
        case 4: //月
          v = diff % 12;
          arr.push({value:v, 'label':'月', 'ext': 'n'});
          diff= Math.floor(diff/12);
          break;
        case 5:
          arr.push({value:v, 'label':'年', 'ext': 'y'});
          break;
      }
      if(Math.floor(diff) < 0) break;
    }

    return arr;
	}

	util.clone = function (obj){
    var newObj = {}
    for (let key in obj) {
        if (typeof obj[key] !== 'object') {
            newObj[key] = obj[key];
        } else {
            newObj[key] = this.clone(obj[key]);
        }
    }
    return newObj;
  }

	util.str2date = function (str){
    if(!str) return;
    if(typeof str != 'string') return str;

		str = str.replace(/[A-Za-z日]/g, ' ').substr(0,19);
    str = str.replace(/[年月]/g, '-');

    var d = new Date(Date.parse(str));
    if(!d || isNaN(d.getFullYear())){
      str = str.replace(/[-]/g, '/');
      return new Date(Date.parse(str));
    }
    return d;
	}

	util.date2str = function(time, ctype){
		if(!time){
			return '';
		}
    ctype = ctype ? ctype : 'date';

		switch(typeof(time)){
			case 'number':
				time = new Date(time);
				break;
			case 'string':
				time = this.str2date(time);
				break;
		}
		switch (ctype) {
			case 'date':
				return time.getFullYear()+'-'+this.lpad(time.getMonth()+1, '0', 2)+'-'+this.lpad(time.getDate(), '0', 2)
      case 'date2':
        return this.lpad(time.getMonth()+1, '0', 2)+'-'+this.lpad(time.getDate(), '0', 2)
      case 'dateCH':
        return time.getFullYear()+'年'+this.lpad(time.getMonth()+1, '0', 2)+'月'+this.lpad(time.getDate(), '0', 2) + '日'
      case 'date2CH':
        return this.lpad(time.getMonth()+1, '0', 2)+'月'+this.lpad(time.getDate(), '0', 2) + '日'
			case 'datetime':
				return time.getFullYear()+'-'+this.lpad(time.getMonth()+1, '0', 2)+'-'+this.lpad(time.getDate(), '0', 2)+' '+this.lpad(time.getHours(), '0', 2)+':'+this.lpad(time.getMinutes(), '0', 2)
			case 'time':
				return this.lpad(time.getHours(), '0', 2)+':'+this.lpad(time.getMinutes(), '0', 2)
		}
		return 'unknow'
	}

	util.getUrlParameter = function (sParam){
    var sPageURL = window.location.search.substring(1);
    var sURLVariables = sPageURL.split('&');
    for (var i = 0; i < sURLVariables.length; i++)
    {
      var sParameterName = sURLVariables[i].split('=');
      if (sParameterName[0] == sParam)
      {
          return decodeURI(sParameterName[1]);
      }
    }
    return '';
	};

	// util.cipherString = function(rsaData, nick, pwd){
	// 	var rsa = new RSAKey(),
	// 		ts = Server.getTime().toString(),
	// 		userkey = CryptoJS.MD5( ts + nick )
	// 	rsa.setPublic(rsaData.hex, '10001');
  //
	// 	var cipher = rsa.encrypt(util.lpad(ts, '0', 16)+userkey.toString(CryptoJS.enc.Base64)),
	// 		text = nick + "|" +pwd;
  //
	// 	var aesCipher = util.aesEncrypto(text, ts, userkey);
  //
	// 	var s = rsaData.keyid.toString()+"|"+
	// 			CryptoJS.enc.Hex.parse(cipher.toString()).toString(CryptoJS.enc.Base64)+"|"+
	// 			aesCipher.toString();
  //
	// 	return s;
	// };
  //
	// util.aesEncrypto = function(text, ts, key){
	// 	ts = ts.toString()
	//     var iv  = CryptoJS.MD5(util.lpad(ts, '0', 16)),
	//     	encrypted = CryptoJS.AES.encrypt(text, key, { iv: iv })
  //
	// 	return encrypted.ciphertext.toString(CryptoJS.enc.Base64)
	// };
  //
	// util.aesDecrypto = function(src, ts, key){
	// 	ts = ts.toString()
	//     var iv  = CryptoJS.MD5(util.lpad(ts, '0', 16)),
	//     	obj = {
	// 			ciphertext: CryptoJS.enc.Base64.parse(src),
	// 			salt: ""
	// 		}
	//     	decrypted = CryptoJS.AES.decrypt(obj, key, { iv: iv })
	// 	return decrypted.toString(CryptoJS.enc.Utf8)
	// };

	util.lpad = function(str, padString, l) {
    if(typeof(str)!='string')
      str = str.toString();
    while (str.toString().length < l)
      str = padString + str;
    return str;
	};

	//pads right
	util.rpad = function(str, padString, l) {
    if(typeof(str)!='string')
      str = str.toString();
    while (str.toString().length < l)
      str = str + padString;
    return str;
	};

  util.find = function(field, value, data){
		if(!data)
			return null;

		var i,len = data.length
		for(i=0;i<len;i++)
			if(data[i] && data[i][field] == value){
				return data[i];
			}
		return null;
	}

	util.inArray = function(val, arr){
		for(var i in arr){
			if(arr[i]==val)
				return true;
		}
		return false;
	}
	util.in2Array = function(val, arr){
		for(var i in arr){
			if(val.toString().indexOf(arr[i].toString())>-1)
				return true;
		}
		return false;
	}

	util.findIndex = function(field, value, data){
		if(!data)
			return -1;

		var i,len = data.length
		for(i=0;i<len;i++)
			if(data[i] && data[i][field] == value){
				return i;
			}
		return -1;
	}

  util.setCookie = function (name,value,days) {
		var expires = ""
      if (days) {
          var date = new Date();
          date.setTime(date.getTime()+(days*86400000));
          expires = "; expires="+date.toGMTString();
      }
      document.cookie = name+"="+value+expires+"; path=/";
  };

	util.getCookie = function (name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0){
          var s = c.substring(nameEQ.length,c.length);
          if(s[0] == '"'){
            return s.substring(1, s.length-1);
          }
          return c.substring(nameEQ.length,c.length);
        }
    }
    return '';
  };

  util.deleteCookie = function (name) {
      this.SetCookie(name,"",-1);
  };

  util.randomElement = function(arr){
    return arr[Math.floor(Math.random() * arr.length)]
  }

  util.shuffle = function(array) {
    var currentIndex = array.length, temporaryValue, randomIndex;
    // While there remain elements to shuffle...
    while (0 !== currentIndex) {

      // Pick a remaining element...
      randomIndex = Math.floor(Math.random() * currentIndex);
      currentIndex -= 1;

      // And swap it with the current element.
      temporaryValue = array[currentIndex];
      array[currentIndex] = array[randomIndex];
      array[randomIndex] = temporaryValue;
    }

    return array;
  }

  util.NumStr = function(v, deci){
    var s = v.toFixed(deci);
    var val = parseFloat(s);
    if (val==parseInt(val)){
      return val.toString();
    }
    return s;
  }

  util.showLoading = function(){
    var $el;
    if($('#toast').length==0){
      $('body').append(util.toast());
    }

    $el = $('#toast');
    $el.removeClass('hide');
    $el.find('.toast-success').hide();
    $el.find('.toast-loading').show();
    setTimeout(function(){$el.addClass('in');}, 50);
    setTimeout(function(){util.hideToast()}, 10000);
  }

  util.toast = function(){
    return $('<div id="toast" class="weui-toast fade hide">'+
      '<div class="weui-mask_transparent"></div>'+
      '<div class="weui-toast">'+
          '<i class="toast-success weui-icon-success-no-circle weui-icon_toast"></i>'+
          '<p class="toast-success weui-toast__content">已完成</p>'+
          '<i class="toast-loading weui-loading weui-icon_toast"></i>'+
          '<p class="toast-loading weui-toast__content">数据加载中</p>'+
      '</div>'+
    '</div>');
  }
  util.showSuccess = function(){
    var $el;
    if($('#toast').length==0){
      $('body').append(util.toast());
    }

    $el = $('#toast');
    $el.removeClass('hide');
    $el.find('.toast-success').show();
    $el.find('.toast-loading').hide();
    setTimeout(function(){$el.addClass('in');}, 50);
    setTimeout(function(){util.hideToast()}, 10000);
  }

  util.hideToast = function(){
    var $el = $('#toast');
    $el.removeClass('in');
    setTimeout(function(){$el.addClass('hide');}, 150);
  }

  util.taggle = function(e){
    var $t;
    if(e.currentTarget){
      $t = $(e.currentTarget);
      if($t.hasClass('parent-pp')){
        $t = $t.parent().parent();
      } else if($t.hasClass('parent-p')){
        $t = $t.parent();
      }
    }else{
      $t = e;
    }

    // index = index ? '-'+index : '';
    var $box = $t.next('.default-hide');
    if($t.hasClass('open')){
      if($box.hasClass('fade')){
        $box.removeClass('in');
        setTimeout(function(){
          $t.removeClass('open');
        }, 150)
      }else{
        $t.removeClass('open');
      }

      return 'close';
    }else{
      $t.addClass('open').siblings('.open').removeClass('open');
      setTimeout(function(){
        $box.addClass('in');
      }, 10)
      $t.parent().siblings().find('.open').removeClass('open');
      return 'open';
    }
  }

	return util;
});
