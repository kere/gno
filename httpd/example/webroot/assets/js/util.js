define('util', [], function(){
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

	var util = {
		DATE_DAY : 86400000,
		DATE_HOUR : 3600000,
	  env : () => {
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
	  },

		language : function() {
			var lang = this.getCookie('lang');
			if(!lang){
				lang = navigator.language || navigator.userLanguage;
			}
			return lang.replace('_', '-').toLowerCase();
		},

	  num2str : (v, deci) => {
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
	  },

		// n 长度，年月日小时分钟
		timeAgoStr : function(b, e, n) {
			if(!n) n = 2;
	    var arr = this.timeAgo(b, e), str='';
	    var l = arr.length, isskip = true, i, k=0;
	    for (i = l-1; i > -1; i--) {
	      if(isskip && arr[i].value==0) continue;
	      isskip = false;
	      str += arr[i].value + arr[i].label + ' ';
				k++
				if(k == n) break;
	    }
	    return str;
	  },

		timeAgo : (b, e) => {
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
		},

		str2date : (str) => {
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
		},

		date2str : function(time, ctype) {
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
		},

		lpad : (str, padString, l) => {
			if(typeof(str)!='string')
				str = str.toString();
			while (str.toString().length < l)
				str = padString + str;
			return str;
		},

		//pads right
		rpad : (str, padString, l) => {
			if(typeof(str)!='string')
				str = str.toString();
			while (str.toString().length < l)
				str = str + padString;
			return str;
		},

		find : (field, value, arr) => {
			if(!arr)
				return null;

			var i,len = arr.length
			for(i=0;i<len;i++)
				if(arr[i] && arr[i][field] == value){
					return arr[i];
				}
			return null;
		},

		indexOfSortedI : (val, arr) => {
			if(!arr) return -1;
			var isdesc = false;
			if(arr.length> 1)
				isdesc = arr[0] > arr[1];
			return indexOfSortedI(val, arr, 0, arr.length-1, isdesc);
		},

		findSortedI : (field, val, arr) => {
			if(!arr) return -1;
			var isdesc = false;
			if(arr.length> 1)
				isdesc = arr[0][field] > arr[1][field];
			return getSortedI(field, val, arr, 0, arr.length-1, isdesc);
		},

		findSorted : (field, val, arr) => {
			if(!arr) return null;
			var isdesc = false;
			if(arr.length> 1)
				isdesc = arr[0][field] > arr[1][field];
			var i = getSortedI(field, val, arr, 0, arr.length-1, isdesc);
			if(i < 0) return null;
			return arr[i];
		},

		findIndex : (field, value, data) => {
			if(!data)
				return -1;

			var i,len = data.length
			for(i=0;i<len;i++)
				if(data[i] && data[i][field] == value){
					return i;
				}
			return -1;
		},

		inArray : (val, arr) => {
			for(var i in arr){
				if(arr[i]==val)
					return true;
			}
			return false;
		},

	  setCookie :  (name,value,days) => {
			var expires = ""
      if (days) {
        var date = new Date();
        date.setTime(date.getTime()+(days*86400000));
        expires = "; expires="+date.toGMTString();
      }
      document.cookie = name+"="+value+expires+"; path=/";
	  },

		getCookie :  (name) => {
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
	  },

	  deleteCookie : (name) => {
	      util.SetCookie(name,"",-1);
	  },

	  copy : (o) => {
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
	  },

		$ : {
			get : function(sel, cls) {
				var el;
				if(typeof(sel) === 'string'){
					var tmp = sel.split(' ');
					if(tmp[0][0] === '#'){
						el = document.getElementById(tmp[0].substr(1));
					}
					if(tmp.length == 2){
						cls = tmp[1];
					}
				} else {
					el = sel;
				}

				if(!cls){
					return el;
				}

		    var arr = cls.split('.'), tag = arr[0], t;
				if(tag === ''){
					t = el.children;
				}else{
		    	t = el.getElementsByTagName(tag);
				}

				if(arr.length == 1){
					return t;
				}
				return this.childrenFilter(t, arr[1]);
			},

			childrenFilter : function(list, cls) {
				var eles = [], item;
				for (var i = 0; i < list.length; i++) {
					item = list[i];
					if(typeof(item.className)!=='string') continue;
					if(this.hasClass(item, cls)){
						eles.push(item);
					}
					if(item.children.length > 0){
						var items = this.childrenFilter(item.children, cls)
						if(items) eles = eles.concat();
					}
				}

				return eles.length > 0 ? eles: null;
			},

			each : function(el, f) {
				el = this.get(el);
				if(el.length > 0){
					for (var i = 0; i < el.length; i++) {
						f(el[i]);
					}
					return;
				}
				return f(el);
			},

			addClass : function(el, clas) {
				el = this.get(el);
				if(this.hasClass(el, clas)) return;
				if(el.className.length==0){
					el.className = clas;
				}else{
					el.className += ' ' +clas;
				}
			},

			hasClass : function(el, clas) {
				el = this.get(el);
				if(!el.className){
					return false;
				}
				var arr = el.className.split(' ');
				for (var i = 0; i < arr.length; i++) {
					if(arr[i]==clas) return true;
				}
				return false;
			},

			removeClass : function(el, clas)  {
				el = this.get(el);
				if(!el.className){
					return;
				}
				var arr = el.className.split(' ');
				for (var i = 0; i < arr.length; i++) {
					if(arr[i]==clas) {
						arr.splice(i, 1);
					};
				}
				el.className = arr.join(' ');
			},

			show : function(el){
				el = this.get(el);

				if(this.hasClass(el, 'hide')){
					this.removeClass(el, 'hide');

					if(this.hasClass(el, 'fade')){
						var ths = this;
						setTimeout(function(){
							ths.addClass(el, 'in');
						}, 50)
					}
				}

				if(el.style.display == 'none'){
					el.style.display = '';
				}
			},

			hide : function(el){
				el = this.get(el);

				if(this.hasClass(el, 'fade')){
					this.removeClass(el, 'in');
					var ths = this;
					setTimeout(function(){
						ths.addClass(el, 'hide');
					}, 200)
				}else{
					this.addClass(el, 'hide');
				}
			}
		}
	}

	util.tool = {
		showLoading : function(n){
			this.loading(0, n);
		},
		showSuccess : function(n){
			this.loading(1, n);
		},
		loading : function(itype, n){
			var el = util.$.get('#toast');
			if(!el){
				el = document.createElement("DIV");
				el.id = 'toast';
				el.className = 'weui-toast fade hide';
				el.innerHTML = `<div class="weui-mask"></div>
<div class="weui-toast">
  <i class="toast-success weui-icon-success-no-circle weui-icon_toast hide"></i>
  <p class="toast-success weui-toast__content hide">已完成</p>
  <i class="toast-loading weui-loading weui-icon_toast hide"></i>
  <p class="toast-loading weui-toast__content hide">数据加载中</p>
</div>`
				document.body.appendChild(el);
			}

			util.$.show(el);
			if(itype==1){
				util.$.each(util.$.get(el, '.toast-loading'), e => {
					util.$.hide(e)
				})
				util.$.each(util.$.get(el, '.toast-success'), e => {
					util.$.show(e)
				})
			}else{
				util.$.each(util.$.get(el, '.toast-success'), e => {
					util.$.hide(e)
				})
				util.$.each(util.$.get(el, '.toast-loading'), e => {
					util.$.show(e)
				})
			}

			n = n || 1;
			setTimeout(function(){
				if(itype==1){
					util.$.each(util.$.get(el, '.toast-success'), e => {
						util.$.hide(e)
					})
				}else{
					util.$.each(util.$.get(el, '.toast-loading'), e => {
						util.$.hide(e)
					})
				}
				util.$.hide(el);
			}, n * 1000);
		},

		hideToast : ()=>{
			var el = util.$.get('#toast');
			util.$.each(util.$.get(el, '.toast-success'), e => {
				util.$.hide(e)
			})
			util.$.each(util.$.get(el, '.toast-loading'), e => {
				util.$.hide(e)
			})
			util.$.hide(el);
		},

		showBusy : function(el, n){
			var t = util.$.get(el, 'div.el-loading-mask');
			if(!t){
				t = document.createElement("DIV");
				t.className = 'el-loading-mask';
				t.innerHTML = `<div class="el-loading-spinner"><svg viewBox="25 25 50 50" class="circular"><circle cx="50" cy="50" r="20" fill="none" class="path"></circle></svg></div>`
				el.appendChild(t);
			}else{
				t = t[0];
			}

			util.$.show(t);
			n = n || 3;
			setTimeout(function(){
				util.$.hide(t);
			}, n * 1000)
		},

		hideBusy : function(el){
			var t = util.$.get(el, 'div.el-loading-mask');
			if(!t) return;
			util.$.hide(t[0]);
		},

	  viewImage : function(url) {
			var el = util.$.get('#viewImage');
	    if(!el){
				el = document.createElement("DIV");
				el.id = 'viewImage';
				el.className = 'fade hide';
				el.innerHTML = '<div style="z-index:1000;position:fixed;width:100%;height:100%;top:0;left:0;text-align:center;top:0;left:0;"><div style="width:100%;height:100%;background: #000;opacity: 0.6;" class="view-image fade in"></div><img src="" style="z-index:1001;max-width:98%;max-height:98%;transform: translate(-50%, -50%);top:50%;position:absolute;"></div>';
				document.body.appendChild(el);

	      el.addEventListener('click', function(e){
					util.$.hide(el);
	      })
	    }

			var imgs = util.$.get(el, 'img');
			imgs[0].src = url;
			util.$.show(el);
	  },

	  taggle : function(e){
	    var el = e.currentTarget;
	    if(e.currentTarget){
				if(util.$.hasClass(el, 'parent-pp')){
					el = el.parentElement.parentElement;
				}else if(util.$.hasClass(el, 'parent-p')){
					el = el.parentElement;
				}
	    }else{
	      el = e;
	    }
			var box = el.nextElementSibling;
			if(util.$.hasClass(el, 'open')){
				util.$.hide(box);
				util.$.removeClass(el, 'open');
				return 'close';
			}
			util.$.show(box);
			util.$.addClass(el, 'open');
			return 'open'
	  }

	}


  // util.getElement = ($el, name) => {
  //   var arr = name.split('.')
  //   if(arr.length!=2) return
  //   var $t = $el.getElementsByTagName(arr[0]);
  //   for (var i in $t) {
  //     if($t[i].className == arr[1]){
  //       return $t[i];
  //     }
  //   }
  // }

	// util.arrayEq = (arr1, arr2) => {
	// 	if(typeof(arr1)!='object' || !arr1.length){
	// 		return arr1 === arr2;
	// 	}
	// 	if(arr1.length!= arr2.length) return false;
	//
	// 	for (var i = 0; i < arr1.length; i++) {
	// 		if(!util.objectEq(arr1[i], arr2[i])){
	// 			return false;
	// 		}
	// 	}
	// 	return true;
	// }
	//
	// util.objectEq = (o1, o2) => {
	// 	if(o1 == null && o2 == null){
	// 		return true;
	// 	}
	// 	if(!o1 || !o2){
	// 		return false;
	// 	}
	// 	if(typeof(o1)!='object'){
	// 		return o1 === o2;
	// 	}
	// 	if (o1.hasOwnProperty('length')) {
	// 		return util.arrayEq(o1, o2);
	// 	}
	// 	// 检查key是否对应
	// 	var k ;
	// 	for (k in o2) {
	// 		if(k.substr(0,2)=='__'){
	// 			continue;
	// 		}
	// 		if (!o1.hasOwnProperty(k)) {
	// 			return false;
	// 		}
	// 	}
	//
	// 	for (k in o1) {
	// 		if(k.substr(0,2)=='__'){
	// 			continue;
	// 		}
	// 		if(typeof(o1[k])=='object') {
	// 			if(util.objectEq(o1[k], o2[k])){
	// 				continue;
	// 			}else{
	// 				return false;
	// 			}
	// 		}
	// 		if (o1[k] != o2[k]) {
	// 			return false;
	// 		}
	// 	}
	// 	return true;
	// }


  // util.weekCH = (v) => {
	// 	if(v==null || typeof v =='undefined') return '';
  //   var weeks = ['天', '一', '二', '三', '四', '五', '六' ];
  //   switch (typeof v) {
  //     case "number":
  //       return weeks[parseInt(v)];
  //     default:
  //       return weeks[util.str2date(v).getDay()];
  //   }
  // }

	// util.getUrlParameter = (sParam) => {
  //   var sPageURL = window.location.search.substring(1);
  //   var sURLVariables = sPageURL.split('&');
  //   for (var i = 0; i < sURLVariables.length; i++)
  //   {
  //     var sParameterName = sURLVariables[i].split('=');
  //     if (sParameterName[0] == sParam)
  //     {
  //         return decodeURI(sParameterName[1]);
  //     }
  //   }
	// };

	// util.getUrlRouterParam = (index) => {
  //   var arr = window.location.pathname.split('/'),
  //     l = arr.length,
  //     i = l-1-index;
  //   if (i<0) return null;
	//
  //   return arr[i];
  // }

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


	// // in2Array array inclued val
	// util.in2Array = (val, arr) => {
	// 	for(var i in arr){
	// 		if(val.toString().indexOf(arr[i].toString())>-1)
	// 			return true;
	// 	}
	// 	return false;
	// }



	// randomElement 随机element
  // util.randomElement = (arr) => {
  //   return arr[Math.floor(Math.random() * arr.length)]
  // }

	// // 洗牌
  // util.shuffle = (array)  => {
  //   var currentIndex = array.length, temporaryValue, randomIndex;
  //   // While there remain elements to shuffle...
  //   while (0 !== currentIndex) {
	//
  //     // Pick a remaining element...
  //     randomIndex = Math.floor(Math.random() * currentIndex);
  //     currentIndex -= 1;
	//
  //     // And swap it with the current element.
  //     temporaryValue = array[currentIndex];
  //     array[currentIndex] = array[randomIndex];
  //     array[randomIndex] = temporaryValue;
  //   }
	//
  //   return array;
  // }



	// util.clone = (obj) => {
  //   var newObj = {}
  //   for (let key in obj) {
  //       if (typeof obj[key] !== 'object') {
  //           newObj[key] = obj[key];
  //       } else {
  //           newObj[key] = util.clone(obj[key]);
  //       }
  //   }
  //   return newObj;
  // }

	return util;
});
