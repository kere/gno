define('util', ['zepto'], function(){
	var util = {}

  util.env = () => {
    var o = {}
    var agent = navigator.userAgent.toLowerCase();
    if(/android/.test(agent)){
      o.os = 'android';
    }else if(/iphone|ipod|ipad|ios/.test(agent)) {
      o.os = 'ios';
    }else if(/windows/.test(agent)){
      o.os = 'windows';
    }
    o.iswxwork = /wxwork/.test(agent);
    o.inwx = /micromessenger/.test(agent);

    return o;
  }

	util.language = () => {
	var lang = util.getCookie('lang');
		if(!lang){
			lang = navigator.language || navigator.userLanguage;
		}
		return lang.replace('_', '-').toLowerCase();
	}

  util.money = (v) => {
    return util.numStr(v, 2);
  }

  util.getElement = ($el, name) => {
    var arr = name.split('.')
    if(arr.length!=2) return
    var $t = $el.getElementsByTagName(arr[0]);
    for (var i in $t) {
      if($t[i].className == arr[1]){
        return $t[i];
      }
    }
  }

	util.arrayEq = (arr1, arr2) => {
		if(typeof(arr1)!='object' || !arr1.length){
			return arr1 === arr2;
		}
		if(arr1.length!= arr2.length) return false;

		for (var i = 0; i < arr1.length; i++) {
			if(!util.objectEq(arr1[i], arr2[i])){
				return false;
			}
		}
		return true;
	}
	util.objectEq = (o1, o2) => {
		if(o1 == null && o2 == null){
			return true;
		}
		if(!o1 || !o2){
			return false;
		}
		if(typeof(o1)!='object'){
			return o1 === o2;
		}
		if (o1.hasOwnProperty('length')) {
			return util.arrayEq(o1, o2);
		}
		// 检查key是否对应
		var k ;
		for (k in o2) {
			if(k.substr(0,2)=='__'){
				continue;
			}
			if (!o1.hasOwnProperty(k)) {
				return false;
			}
		}

		for (k in o1) {
			if(k.substr(0,2)=='__'){
				continue;
			}
			if(typeof(o1[k])=='object') {
				if(util.objectEq(o1[k], o2[k])){
					continue;
				}else{
					return false;
				}
			}
			if (o1[k] != o2[k]) {
				return false;
			}
		}
		return true;
	}

  util.numStr = (v, deci) => {
		if(!v) return '';
		if(typeof(v)== 'string') v= parseFloat(v);

		var isPad = false;
		if(deci < 0){
			var s = v.toString(), arr = v.toString().split('.')
			if(arr.length == 1) return s;
			deci = -deci;
			isPad = true;
		}

    if(deci==0){
      return v.toFixed(deci);
    }
    var deciV = Math.pow(10, deci);
    var s = (Math.round(v * deciV)/deciV).toFixed(deci);
		if(!isPad) return s;

		var arr = new Array(s.length);
		for (var i = 0; i < s.length; i++) {
			arr[i] = s[i];
		}

		for (var i = s.length-1; i > -1; i--) {
			if(arr[i]!= '0') break;
			arr.pop();
		}
		return arr.join('');
  }

  util.weekCH = (v) => {
		if(v==null || typeof v =='undefined') return '';
    var weeks = ['天', '一', '二', '三', '四', '五', '六' ];
    switch (typeof v) {
      case "number":
        return weeks[parseInt(v)];
      default:
        return weeks[util.str2date(v).getDay()];
    }
  }

	util.DATE_DAY = 86400000;
	util.DATE_HOUR = 3600000;

	// n 长度，年月日小时分钟
	util.timeAgoStr = (b, e, n) => {
		if(!n) n = 2;
    var arr = util.timeAgo(b, e), str='';
    var l = arr.length, isskip = true, i, k=0;
    for (i = l-1; i > -1; i--) {
      if(isskip && arr[i].value==0) continue;
      isskip = false;
      str += arr[i].value + arr[i].label + ' ';
			k++
			if(k == n) break;
    }
    return str;
  }

	util.timeAgo = (b, e) => {
    var diff;
    if(typeof(b)=='number' && typeof(e) == 'undefined'){
      diff = b;
    }else{
      diff = Math.abs(e.getTime() - b.getTime());
    }
    var v, n, arr=[];

    for(var i=0;i<6;i++) {
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

	util.clone = (obj) => {
    var newObj = {}
    for (let key in obj) {
        if (typeof obj[key] !== 'object') {
            newObj[key] = obj[key];
        } else {
            newObj[key] = util.clone(obj[key]);
        }
    }
    return newObj;
  }

	util.str2date = (str) => {
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

	util.date2str = (time, ctype) => {
		if(!time){
			return '';
		}
    ctype = ctype ? ctype : 'date';

		switch(typeof(time)){
			case 'number':
				time = new Date(time);
				break;
			case 'string':
				time = util.str2date(time);
				break;
		}
		switch (ctype) {
			case 'date':
				return time.getFullYear()+'-'+util.lpad(time.getMonth()+1, '0', 2)+'-'+util.lpad(time.getDate(), '0', 2)
      case 'date2':
        return util.lpad(time.getMonth()+1, '0', 2)+'-'+util.lpad(time.getDate(), '0', 2)
      case 'dateCH':
        return time.getFullYear()+'年'+util.lpad(time.getMonth()+1, '0', 2)+'月'+util.lpad(time.getDate(), '0', 2) + '日'
      case 'date2CH':
        return util.lpad(time.getMonth()+1, '0', 2)+'月'+util.lpad(time.getDate(), '0', 2) + '日'
			case 'datetime':
				return time.getFullYear()+'-'+util.lpad(time.getMonth()+1, '0', 2)+'-'+util.lpad(time.getDate(), '0', 2)+' '+util.lpad(time.getHours(), '0', 2)+':'+util.lpad(time.getMinutes(), '0', 2)
			case 'time':
				return util.lpad(time.getHours(), '0', 2)+':'+util.lpad(time.getMinutes(), '0', 2)
		}
		return 'unknow'
	}

	util.getUrlParameter = (sParam) => {
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
	};

	util.getUrlRouterParam = (index) => {
    var arr = window.location.pathname.split('/'),
      l = arr.length,
      i = l-1-index;
    if (i<0) return null;

    return arr[i];
  }

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

	util.lpad = (str, padString, l) => {
    if(typeof(str)!='string')
      str = str.toString();
    while (str.toString().length < l)
      str = padString + str;
    return str;
	};

	//pads right
	util.rpad = (str, padString, l) => {
    if(typeof(str)!='string')
      str = str.toString();
    while (str.toString().length < l)
      str = str + padString;
    return str;
	};

  util.find =(field, value, arr) => {
		if(!arr)
			return null;

		var i,len = arr.length
		for(i=0;i<len;i++)
			if(arr[i] && arr[i][field] == value){
				return arr[i];
			}
		return null;
	}

	function indexOfSortedI(val, arr, b, e, desc){
		if(b == e){
			return arr[b] == val ? b : -1;
		}else if(b>e){
			return -1
		}
		var l = e-b+1, i = b+Math.floor(l/2), v = arr[i];

		if(v == val){
			return i;
		}else if(v > val){
			if(desc){
				return indexOfSortedI(val, arr, i+1, e)
			}
			// small zone
			return indexOfSortedI(val, arr, b, i-1)
		}else{
			if(desc){
				return indexOfSortedI(val, arr, b, i-1)
			}
			return indexOfSortedI(val, arr, i+1, e)
		}
	}

  util.indexOfSortedI = (val, arr) => {
		if(!arr) return -1;
		var isdesc = false;
		if(arr.length> 1)
			isdesc = arr[0] > arr[1];
		return indexOfSortedI(val, arr, 0, arr.length-1, isdesc);
	}

	// b begin; e end index
	function getSortedI(field, val, arr, b, e, desc) {
		if(b == e){
			return arr[b][field] == val ? b : -1;
		}else if(b>e){
			return -1
		}
		var l = e-b+1, i = b+Math.floor(l/2), v;
		v = arr[i][field]

		if(v == val){
			return i;
		}else if(v > val){
			if(desc){
				return getSortedI(field, val, arr, i+1, e)
			}
			// small zone
			return getSortedI(field, val, arr, b, i-1)
		}else{
			if(desc){
				return getSortedI(field, val, arr, b, i-1)
			}
			return getSortedI(field, val, arr, i+1, e)
		}
	}

  util.findSortedI = (field, val, arr) => {
		if(!arr) return -1;
		var isdesc = false;
		if(arr.length> 1)
			isdesc = arr[0][field] > arr[1][field];
		return getSortedI(field, val, arr, 0, arr.length-1, isdesc);
	}

  util.findSorted = (field, val, arr) => {
		if(!arr) return null;
		var isdesc = false;
		if(arr.length> 1)
			isdesc = arr[0][field] > arr[1][field];
		var i = getSortedI(field, val, arr, 0, arr.length-1, isdesc);
		if(i < 0) return null;
		return arr[i];
	}

	util.inArray = (val, arr) => {
		for(var i in arr){
			if(arr[i]==val)
				return true;
		}
		return false;
	}
	util.in2Array = (val, arr) => {
		for(var i in arr){
			if(val.toString().indexOf(arr[i].toString())>-1)
				return true;
		}
		return false;
	}

	util.findIndex = (field, value, data) => {
		if(!data)
			return -1;

		var i,len = data.length
		for(i=0;i<len;i++)
			if(data[i] && data[i][field] == value){
				return i;
			}
		return -1;
	}

  util.setCookie =  (name,value,days) => {
		var expires = ""
      if (days) {
          var date = new Date();
          date.setTime(date.getTime()+(days*86400000));
          expires = "; expires="+date.toGMTString();
      }
      document.cookie = name+"="+value+expires+"; path=/";
  };

	util.getCookie =  (name) => {
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

  util.deleteCookie = (name) => {
      util.SetCookie(name,"",-1);
  };

  util.randomElement = (arr) => {
    return arr[Math.floor(Math.random() * arr.length)]
  }

	// 洗牌
  util.shuffle = (array)  => {
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

  util.showLoading = ()=> {
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

  util.viewImage = (url) => {
    if($('#viewImage').length==0){
      var $t = $('<div id="viewImage" class="fade hide" style="z-index:1000;position:fixed;width:100%;height:100%;top:0;left:0;text-align:center;top:0;left:0;"><div style="width:100%;height:100%;background: #000;opacity: 0.6;" class="view-image fade in"></div><img src="" style="z-index:1001;max-width:98%;max-height:98%;transform: translate(-50%, -50%);top:50%;position:absolute;"></div>');
      $('body').append($t);
      $t[0].addEventListener('click', function(e){
        var $t = $(e.currentTarget)
        $t.removeClass('in');
        setTimeout(function(){
          $t.addClass('hide');
        },150)
      })
    }
    var $t = $('#viewImage');
    $t.find('img').attr('src', url);
    $t.removeClass('hide').addClass('in');
  }

  util.toast = () => {
    return $('<div id="toast" class="weui-toast fade hide">'+
      '<div class="weui-mask_bk"></div>'+
      '<div class="weui-toast">'+
          '<i class="toast-success weui-icon-success-no-circle weui-icon_toast"></i>'+
          '<p class="toast-success weui-toast__content">已完成</p>'+
          '<i class="toast-loading weui-loading weui-icon_toast"></i>'+
          '<p class="toast-loading weui-toast__content">数据加载中</p>'+
      '</div>'+
    '</div>');
  }
  util.showSuccess = () => {
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

  util.hideToast = ()=>{
    var $el = $('#toast');
    $el.removeClass('in');
    setTimeout(function(){$el.addClass('hide');}, 150);
  }

  util.taggle = function(e, isCloseOther){
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
      $t.addClass('open')
      setTimeout(function(){
        $box.addClass('in');
      }, 10)
			if(isCloseOther){
				$t.siblings('.open').removeClass('open');
      	$t.parent().siblings().find('.open').removeClass('open');
			}
      return 'open';
    }
  }

  util.copy = (o) => {
    var dat = {};
		if(typeof(o) == 'object' && o.hasOwnProperty('length')){
			dat = new Array(o.length);
		}
    for (var k in o) {
      if (!o.hasOwnProperty(k) || k.substr(0,2)=="__") continue;
			if(o[k]== null){
				dat[k] = null
				continue;
			}
      if(typeof(o[k]) == 'object'){
        dat[k] = util.copy(o[k])
      }else{
        dat[k] = o[k];
      }
    }
    return dat;
  }

	return util;
});
