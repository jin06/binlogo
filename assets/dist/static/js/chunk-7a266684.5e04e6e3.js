(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-7a266684"],{"3b1d":function(t,e,n){"use strict";n("bb70")},"565d":function(t,e,n){},"568c":function(t,e,n){},"9ed6":function(t,e,n){"use strict";n.r(e);var s=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"login-container"},[n("el-form",{ref:"loginForm",staticClass:"login-form",attrs:{model:t.loginForm,autocomplete:"on","label-position":"left"}},[n("div",{staticClass:"title-container"},[n("h3",{staticClass:"title"},[t._v(" "+t._s(t.$t("login.title"))+" ")]),n("lang-select",{staticClass:"set-language"})],1),n("el-form-item",{attrs:{prop:"username"}},[n("span",{staticClass:"svg-container"},[n("svg-icon",{attrs:{"icon-class":"user"}})],1),n("el-input",{ref:"username",attrs:{placeholder:t.$t("login.username"),name:"username",type:"text",tabindex:"1",autocomplete:"on"},model:{value:t.loginForm.username,callback:function(e){t.$set(t.loginForm,"username",e)},expression:"loginForm.username"}})],1),n("el-tooltip",{attrs:{content:"Caps lock is On",placement:"right",manual:""},model:{value:t.capsTooltip,callback:function(e){t.capsTooltip=e},expression:"capsTooltip"}},[n("el-form-item",{attrs:{prop:"password"}},[n("span",{staticClass:"svg-container"},[n("svg-icon",{attrs:{"icon-class":"password"}})],1),n("el-input",{key:t.passwordType,ref:"password",attrs:{type:t.passwordType,placeholder:t.$t("login.password"),name:"password",tabindex:"2",autocomplete:"on"},on:{blur:function(e){t.capsTooltip=!1}},nativeOn:{keyup:[function(e){return t.checkCapslock(e)},function(e){return!e.type.indexOf("key")&&t._k(e.keyCode,"enter",13,e.key,"Enter")?null:t.handleLogin(e)}]},model:{value:t.loginForm.password,callback:function(e){t.$set(t.loginForm,"password",e)},expression:"loginForm.password"}}),n("span",{staticClass:"show-pwd",on:{click:t.showPwd}},[n("svg-icon",{attrs:{"icon-class":"password"===t.passwordType?"eye":"eye-open"}})],1)],1)],1),n("el-button",{staticStyle:{width:"100%","margin-bottom":"30px"},attrs:{loading:t.loading,type:"primary"},nativeOn:{click:function(e){return e.preventDefault(),t.handleLogin(e)}}},[t._v(" "+t._s(t.$t("login.logIn"))+" ")]),n("div",{staticStyle:{position:"relative"}},[n("div",{staticClass:"tips"},[n("span",[t._v(t._s(t.$t("login.authType"))+" : "+t._s(t.auth.type))])])])],1),n("el-dialog",{attrs:{title:t.$t("login.thirdparty"),visible:t.showDialog},on:{"update:visible":function(e){t.showDialog=e}}},[t._v(" "+t._s(t.$t("login.thirdpartyTips"))+" "),n("br"),n("br"),n("br"),n("social-sign")],1)],1)},o=[],i=(n("13d5"),n("d3b7"),n("b64b"),n("1131")),a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"social-signup-container"},[n("div",{staticClass:"sign-btn",on:{click:function(e){return t.wechatHandleClick("wechat")}}},[n("span",{staticClass:"wx-svg-container"},[n("svg-icon",{staticClass:"icon",attrs:{"icon-class":"wechat"}})],1),t._v(" WeChat ")]),n("div",{staticClass:"sign-btn",on:{click:function(e){return t.tencentHandleClick("tencent")}}},[n("span",{staticClass:"qq-svg-container"},[n("svg-icon",{staticClass:"icon",attrs:{"icon-class":"qq"}})],1),t._v(" QQ ")])])},r=[],c={name:"SocialSignin",methods:{wechatHandleClick:function(t){alert("ok")},tencentHandleClick:function(t){alert("ok")}}},l=c,u=(n("aa05"),n("2877")),p=Object(u["a"])(l,a,r,!1,null,"7309fbbb",null),d=p.exports,h=n("c24f"),g={name:"Login",components:{LangSelect:i["a"],SocialSign:d},data:function(){return{loginForm:{username:"admin",password:""},passwordType:"password",capsTooltip:!1,loading:!1,showDialog:!1,redirect:void 0,otherQuery:{},auth:{}}},watch:{$route:{handler:function(t){var e=t.query;e&&(this.redirect=e.redirect,this.otherQuery=this.getOtherQuery(e))},immediate:!0}},created:function(){this.getAuthType()},mounted:function(){""===this.loginForm.username?this.$refs.username.focus():""===this.loginForm.password&&this.$refs.password.focus()},destroyed:function(){},methods:{checkCapslock:function(t){var e=t.key;this.capsTooltip=e&&1===e.length&&e>="A"&&e<="Z"},showPwd:function(){var t=this;"password"===this.passwordType?this.passwordType="":this.passwordType="password",this.$nextTick((function(){t.$refs.password.focus()}))},handleLogin:function(){var t=this;this.$refs.loginForm.validate((function(e){if(!e)return console.log("error submit!!"),!1;t.loading=!0,t.$store.dispatch("user/login",t.loginForm).then((function(){t.$router.push({path:t.redirect||"/",query:t.otherQuery}),t.loading=!1})).catch((function(){t.loading=!1}))}))},getAuthType:function(){var t=this;Object(h["a"])().then((function(e){t.auth=e.data}))},getOtherQuery:function(t){return Object.keys(t).reduce((function(e,n){return"redirect"!==n&&(e[n]=t[n]),e}),{})}}},f=g,m=(n("3b1d"),n("e499"),Object(u["a"])(f,s,o,!1,null,"912c8410",null));e["default"]=m.exports},aa05:function(t,e,n){"use strict";n("565d")},bb70:function(t,e,n){},e499:function(t,e,n){"use strict";n("568c")}}]);