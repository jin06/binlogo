(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-17c63f4a"],{2794:function(t,a,n){"use strict";n("a39f")},"7af2":function(t,a,n){"use strict";n("a7b3")},a39f:function(t,a,n){},a7b3:function(t,a,n){},b23f:function(t,a,n){"use strict";n.r(a);var s=function(){var t=this,a=t.$createElement,n=t._self._c||a;return n("div",{staticClass:"components-container board"},[n("Kanban",{key:1,staticClass:"kanban todo",attrs:{list:t.list1,group:t.group,"header-text":"Todo"}}),n("Kanban",{key:2,staticClass:"kanban working",attrs:{list:t.list2,group:t.group,"header-text":"Working"}}),n("Kanban",{key:3,staticClass:"kanban done",attrs:{list:t.list3,group:t.group,"header-text":"Done"}})],1)},e=[],i=function(){var t=this,a=t.$createElement,n=t._self._c||a;return n("div",{staticClass:"board-column"},[n("div",{staticClass:"board-column-header"},[t._v(" "+t._s(t.headerText)+" ")]),n("draggable",t._b({staticClass:"board-column-content",attrs:{list:t.list,"set-data":t.setData}},"draggable",t.$attrs,!1),t._l(t.list,(function(a){return n("div",{key:a.id,staticClass:"board-item"},[t._v(" "+t._s(a.name)+" "+t._s(a.id)+" ")])})),0)],1)},o=[],r=n("1980"),l=n.n(r),d={name:"DragKanbanDemo",components:{draggable:l.a},props:{headerText:{type:String,default:"Header"},options:{type:Object,default:function(){return{}}},list:{type:Array,default:function(){return[]}}},methods:{setData:function(t){t.setData("Text","")}}},c=d,u=(n("7af2"),n("2877")),b=Object(u["a"])(c,i,o,!1,null,"083991bb",null),m=b.exports,f={name:"DragKanbanDemo",components:{Kanban:m},data:function(){return{group:"mission",list1:[{name:"Mission",id:1},{name:"Mission",id:2},{name:"Mission",id:3},{name:"Mission",id:4}],list2:[{name:"Mission",id:5},{name:"Mission",id:6},{name:"Mission",id:7}],list3:[{name:"Mission",id:8},{name:"Mission",id:9},{name:"Mission",id:10}]}}},p=f,g=(n("2794"),Object(u["a"])(p,s,e,!1,null,null,null));a["default"]=g.exports}}]);