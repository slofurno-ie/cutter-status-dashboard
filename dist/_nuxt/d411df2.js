(window.webpackJsonp=window.webpackJsonp||[]).push([[3],{224:function(t,e,n){var content=n(228);content.__esModule&&(content=content.default),"string"==typeof content&&(content=[[t.i,content,""]]),content.locals&&(t.exports=content.locals);(0,n(59).default)("bad8e090",content,!0,{sourceMap:!1})},226:function(t,e,n){t.exports=n.p+"img/logo.6a2913c.svg"},227:function(t,e,n){"use strict";n(224)},228:function(t,e,n){var r=n(58)((function(i){return i[1]}));r.push([t.i,"#main{\n  min-height:calc(100vh - 80px)\n}",""]),t.exports=r},231:function(t,e,n){"use strict";n.r(e);var r=[function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("header",{staticClass:"fixed inset-x-0 h-24 p-5 bg-white shadow-md",attrs:{id:"header"}},[r("div",{staticClass:"grid h-full grid-cols-3",attrs:{id:"header-container"}},[r("img",{staticClass:"self-center h-6 xl:h-10",attrs:{id:"logo",src:n(226),alt:"ie logo"}}),t._v(" "),r("div",{staticClass:"self-center col-span-2 text-sm font-bold text-gray-500 xl:col-span-1 xl:text-3xl title justify-self-end lg:justify-self-center"},[t._v("\n        Cutter Service Status Dashboard\n      ")])])])}],l=n(7),o=(n(48),n(60),n(57)),c=n.n(o),d={data:function(){return{ret:null}},mounted:function(){var t=this;return Object(l.a)(regeneratorRuntime.mark((function e(){return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:setInterval((function(){console.log("interval = 60 seconds")}),6e4),t.getallstatuses();case 2:case"end":return e.stop()}}),e)})))()},methods:{getallstatuses:function(){var t=this;return Object(l.a)(regeneratorRuntime.mark((function e(){var data;return regeneratorRuntime.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,c.a.get("https://dashboard.dev.cutter.live/api/v1/get-all-statuses");case 3:data=e.sent,t.ret=data.data,console.log(data),e.next=11;break;case 8:e.prev=8,e.t0=e.catch(0),console.log(e.t0);case 11:case"end":return e.stop()}}),e,null,[[0,8]])})))()}}},f=(n(227),n(46)),component=Object(f.a)(d,(function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("main",{staticClass:"container flex flex-col mx-auto",attrs:{id:"main"}},[t._m(0),t._v(" "),n("div",{staticClass:"grid flex-1 w-full grid-cols-1 grid-rows-4 gap-5 p-5 text-gray-500 lg:grid-cols-4 lg:grid-rows-1 mx-a6to pt-36"},t._l(t.ret,(function(e){return n("div",{key:e.StatusId,staticClass:"relative flex items-center justify-between h-full px-2 py-5 overflow-hidden bg-white rounded-md shadow-md xl:p-5"},[n("div",[n("p",{staticClass:"text-2xl uppercase"},[t._v(t._s(e.service))])]),t._v(" "),n("div",{class:("200"===e.status?"bg-green-800":"bg-red-900")+" text-white absolute right-0 top-0 bottom-0 p-5 font-bold text-2xl"},[t._v("\n        "+t._s(e.status)+"\n      ")])])})),0)])}),r,!1,null,null,null);e.default=component.exports}}]);