!function(e,t){"object"==typeof exports&&"undefined"!=typeof module?t(exports):"function"==typeof define&&define.amd?define(["exports"],t):t((e="undefined"!=typeof globalThis?globalThis:e||self).bridge={})}(this,(function(e){"use strict";function t(){}const n=e=>e;function o(e){return e()}function r(){return Object.create(null)}function s(e){e.forEach(o)}function i(e){return"function"==typeof e}function l(e,t){return e!=e?t==t:e!==t||e&&"object"==typeof e||"function"==typeof e}function a(e,n,o){e.$$.on_destroy.push(function(e,...n){if(null==e)return t;const o=e.subscribe(...n);return o.unsubscribe?()=>o.unsubscribe():o}(n,o))}const c="undefined"!=typeof window;let u=c?()=>window.performance.now():()=>Date.now(),d=c?e=>requestAnimationFrame(e):t;const f=new Set;function A(e){f.forEach((t=>{t.c(e)||(f.delete(t),t.f())})),0!==f.size&&d(A)}function p(e,t){e.appendChild(t)}function g(e,t,n){e.insertBefore(t,n||null)}function m(e){e.parentNode.removeChild(e)}function h(e){return document.createElement(e)}function y(e){return document.createTextNode(e)}function b(){return y("")}function w(e,t,n){null==n?e.removeAttribute(t):e.getAttribute(t)!==n&&e.setAttribute(t,n)}const v=new Set;let x,$=0;function k(e,t,n,o,r,s,i,l=0){const a=16.666/o;let c="{\n";for(let e=0;e<=1;e+=a){const o=t+(n-t)*s(e);c+=100*e+`%{${i(o,1-o)}}\n`}const u=c+`100% {${i(n,1-n)}}\n}`,d=`__svelte_${function(e){let t=5381,n=e.length;for(;n--;)t=(t<<5)-t^e.charCodeAt(n);return t>>>0}(u)}_${l}`,f=e.ownerDocument;v.add(f);const A=f.__svelte_stylesheet||(f.__svelte_stylesheet=f.head.appendChild(h("style")).sheet),p=f.__svelte_rules||(f.__svelte_rules={});p[d]||(p[d]=!0,A.insertRule(`@keyframes ${d} ${u}`,A.cssRules.length));const g=e.style.animation||"";return e.style.animation=`${g?`${g}, `:""}${d} ${o}ms linear ${r}ms 1 both`,$+=1,d}function _(e,t){const n=(e.style.animation||"").split(", "),o=n.filter(t?e=>e.indexOf(t)<0:e=>-1===e.indexOf("__svelte")),r=n.length-o.length;r&&(e.style.animation=o.join(", "),$-=r,$||d((()=>{$||(v.forEach((e=>{const t=e.__svelte_stylesheet;let n=t.cssRules.length;for(;n--;)t.deleteRule(n);e.__svelte_rules={}})),v.clear())})))}function O(e){x=e}const C=[],E=[],I=[],q=[],B=Promise.resolve();let S=!1;function D(e){I.push(e)}let T=!1;const j=new Set;function R(){if(!T){T=!0;do{for(let e=0;e<C.length;e+=1){const t=C[e];O(t),N(t.$$)}for(O(null),C.length=0;E.length;)E.pop()();for(let e=0;e<I.length;e+=1){const t=I[e];j.has(t)||(j.add(t),t())}I.length=0}while(C.length);for(;q.length;)q.pop()();S=!1,T=!1,j.clear()}}function N(e){if(null!==e.fragment){e.update(),s(e.before_update);const t=e.dirty;e.dirty=[-1],e.fragment&&e.fragment.p(e.ctx,t),e.after_update.forEach(D)}}let P;function L(e,t,n){e.dispatchEvent(function(e,t){const n=document.createEvent("CustomEvent");return n.initCustomEvent(e,!1,!1,t),n}(`${t?"intro":"outro"}${n}`))}const M=new Set;let U;function X(){U={r:0,c:[],p:U}}function z(){U.r||s(U.c),U=U.p}function Q(e,t){e&&e.i&&(M.delete(e),e.i(t))}function Z(e,t,n,o){if(e&&e.o){if(M.has(e))return;M.add(e),U.c.push((()=>{M.delete(e),o&&(n&&e.d(1),o())})),e.o(t)}}const V={duration:0};function Y(e,o,r,l){let a=o(e,r),c=l?0:1,p=null,g=null,m=null;function h(){m&&_(e,m)}function y(e,t){const n=e.b-c;return t*=Math.abs(n),{a:c,b:e.b,d:n,duration:t,start:e.start,end:e.start+t,group:e.group}}function b(o){const{delay:r=0,duration:i=300,easing:l=n,tick:b=t,css:w}=a||V,v={start:u()+r,b:o};o||(v.group=U,U.r+=1),p||g?g=v:(w&&(h(),m=k(e,c,o,i,r,l,w)),o&&b(0,1),p=y(v,i),D((()=>L(e,o,"start"))),function(e){let t;0===f.size&&d(A),new Promise((n=>{f.add(t={c:e,f:n})}))}((t=>{if(g&&t>g.start&&(p=y(g,i),g=null,L(e,p.b,"start"),w&&(h(),m=k(e,c,p.b,p.duration,0,l,a.css))),p)if(t>=p.end)b(c=p.b,1-c),L(e,p.b,"end"),g||(p.b?h():--p.group.r||s(p.group.c)),p=null;else if(t>=p.start){const e=t-p.start;c=p.a+p.d*l(e/p.duration),b(c,1-c)}return!(!p&&!g)})))}return{run(e){i(a)?(P||(P=Promise.resolve(),P.then((()=>{P=null}))),P).then((()=>{a=a(),b(e)})):b(e)},end(){h(),p=g=null}}}function F(e,t,n){const{fragment:r,on_mount:l,on_destroy:a,after_update:c}=e.$$;r&&r.m(t,n),D((()=>{const t=l.map(o).filter(i);a?a.push(...t):s(t),e.$$.on_mount=[]})),c.forEach(D)}function H(e,t){const n=e.$$;null!==n.fragment&&(s(n.on_destroy),n.fragment&&n.fragment.d(t),n.on_destroy=n.fragment=null,n.ctx=[])}function J(e,t){-1===e.$$.dirty[0]&&(C.push(e),S||(S=!0,B.then(R)),e.$$.dirty.fill(0)),e.$$.dirty[t/31|0]|=1<<t%31}function W(e,n,o,i,l,a,c=[-1]){const u=x;O(e);const d=e.$$={fragment:null,ctx:null,props:a,update:t,not_equal:l,bound:r(),on_mount:[],on_destroy:[],before_update:[],after_update:[],context:new Map(u?u.$$.context:[]),callbacks:r(),dirty:c,skip_bound:!1};let f=!1;if(d.ctx=o?o(e,n.props||{},((t,n,...o)=>{const r=o.length?o[0]:n;return d.ctx&&l(d.ctx[t],d.ctx[t]=r)&&(!d.skip_bound&&d.bound[t]&&d.bound[t](r),f&&J(e,t)),n})):[],d.update(),f=!0,s(d.before_update),d.fragment=!!i&&i(d.ctx),n.target){if(n.hydrate){const e=function(e){return Array.from(e.childNodes)}(n.target);d.fragment&&d.fragment.l(e),e.forEach(m)}else d.fragment&&d.fragment.c();n.intro&&Q(e.$$.fragment),F(e,n.target,n.anchor),R()}O(u)}class K{$destroy(){H(this,1),this.$destroy=t}$on(e,t){const n=this.$$.callbacks[e]||(this.$$.callbacks[e]=[]);return n.push(t),()=>{const e=n.indexOf(t);-1!==e&&n.splice(e,1)}}$set(e){var t;this.$$set&&(t=e,0!==Object.keys(t).length)&&(this.$$.skip_bound=!0,this.$$set(e),this.$$.skip_bound=!1)}}const G=[];function ee(e,n=t){let o;const r=[];function s(t){if(l(e,t)&&(e=t,o)){const t=!G.length;for(let t=0;t<r.length;t+=1){const n=r[t];n[1](),G.push(n,e)}if(t){for(let e=0;e<G.length;e+=2)G[e][0](G[e+1]);G.length=0}}}return{set:s,update:function(t){s(t(e))},subscribe:function(i,l=t){const a=[i,l];return r.push(a),1===r.length&&(o=n(s)||t),i(e),()=>{const e=r.indexOf(a);-1!==e&&r.splice(e,1),0===r.length&&(o(),o=null)}}}}function te(e){console.log("%c wails bridge %c "+e+" ","background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem","background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem")}const ne=ee(!1);function oe(){ne.set(!0)}const re=ee(!0),se=ee([]);function ie(e,{delay:t=0,duration:o=400,easing:r=n}={}){const s=+getComputedStyle(e).opacity;return{delay:t,duration:o,easing:r,css:e=>"opacity: "+e*s}}function le(e){let t,n,o;return{c(){t=h("div"),t.innerHTML='<div class="wails-reconnect-overlay-content svelte-9nqyfr"><div class="wails-reconnect-overlay-loadingspinner svelte-9nqyfr"></div></div>',w(t,"class","wails-reconnect-overlay svelte-9nqyfr")},m(e,n){g(e,t,n),o=!0},i(e){o||(D((()=>{n||(n=Y(t,ie,{duration:200},!0)),n.run(1)})),o=!0)},o(e){n||(n=Y(t,ie,{duration:200},!1)),n.run(0),o=!1},d(e){e&&m(t),e&&n&&n.end()}}}function ae(e){let t,n,o=e[0]&&le();return{c(){o&&o.c(),t=b()},m(e,r){o&&o.m(e,r),g(e,t,r),n=!0},p(e,[n]){e[0]?o?1&n&&Q(o,1):(o=le(),o.c(),Q(o,1),o.m(t.parentNode,t)):o&&(X(),Z(o,1,1,(()=>{o=null})),z())},i(e){n||(Q(o),n=!0)},o(e){Z(o),n=!1},d(e){o&&o.d(e),e&&m(t)}}}function ce(e,t,n){let o;return a(e,ne,(e=>n(0,o=e))),[o]}class ue extends K{constructor(e){var t;super(),document.getElementById("svelte-9nqyfr-style")||((t=h("style")).id="svelte-9nqyfr-style",t.textContent=".wails-reconnect-overlay.svelte-9nqyfr{position:fixed;top:0;left:0;width:100%;height:100%;backdrop-filter:blur(20px) saturate(160%) contrast(45%) brightness(140%);z-index:999999\n    }.wails-reconnect-overlay-content.svelte-9nqyfr{position:relative;top:50%;transform:translateY(-50%);margin:0;background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAC8AAAAuCAMAAACPpbA7AAAAflBMVEUAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQAAAAAAAAAAAABAQEEBAQAAAAAAAAEBAQAAAADAwMAAAABAQEAAAAAAAAAAAAAAAAAAAACAgICAgIBAQEAAAAAAAAAAAAAAAAAAAACAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAFBQWCC3waAAAAKXRSTlMALgUMIBk0+xEqJs70Xhb3lu3EjX2EZTlv5eHXvbarQj3cdmpXSqOeUDwaqNAAAAKCSURBVEjHjZTntqsgEIUPVVCwtxg1vfD+L3hHRe8K6snZf+KKn8OewvzsSSeXLruLnz+KHs0gr6DkT3xsRkU6VVn4Ha/UxLe1Z4y64i847sykPBh/AvQ7ry3eFN70oKrfcBJYvm/tQ1qxP4T3emXPeXAkvodPUvtdjbhk+Ft4c0hslTiXVOzxOJ15NWUblQhRsdu3E1AfCjj3Gdm18zSOsiH8Lk4TB480ksy62fiqNo4OpyU8O21l6+hyRtS6z8r1pHlmle5sR1/WXS6Mq2Nl+YeKt3vr+vdH/q4O68tzXuwkiZmngYb4R8Co1jh0+Ww2UTyWxBvtyxLO7QVjO3YOD/lWZpbXDGellFG2Mws58mMnjVZSn7p+XvZ6IF4nn02OJZV0aTO22arp/DgLPtrgpVoi6TPbZm4XQBjY159w02uO0BDdYsfrOEi0M2ulRXlCIPAOuN1NOVhi+riBR3dgwQplYsZRZJLXq23Mlo5njkbY0rZFu3oiNIYG2kqsbVz67OlNuZZIOlfxHDl0UpyRX86z/OYC/3qf1A1xTrMp/PWWM4ePzf8DDp1nesQRpcFk7BlwdzN08ZIALJpCaciQXO0f6k4dnuT/Ewg4l7qSTNzm2SykdHn6GJ12mWc6aCNj/g1cTXpB8YFfr0uVc96aFkkqiIiX4nO+salKwGtIkvfB+Ja8DxMeD3hIXP5mTOYPB4eVT0+32I5ykvPZjesnkGgIREgYnmLrPb0PdV3hoLup2TjcGBPM4mgsfF5BrawZR4/GpzYQzQfrUZCf0TCWYo2DqhdhTJBQ6j4xqmmLN5LjdRIY8LWExiFUsSrza/nmFBqw3I9tEZB9h0lIQSO9if8DkISDAj8CDawAAAAASUVORK5CYII=);background-repeat:no-repeat;background-position:center\n    }.wails-reconnect-overlay-loadingspinner.svelte-9nqyfr{pointer-events:none;width:2.5em;height:2.5em;border:.4em solid transparent;border-color:#f00 #eee0 #f00 #eee0;border-radius:50%;animation:svelte-9nqyfr-loadingspin 1s linear infinite;margin:auto;padding:2.5em\n    }@keyframes svelte-9nqyfr-loadingspin{100%{transform:rotate(360deg)}}",p(document.head,t)),W(this,e,ce,ae,l,{})}}function de(e){let n,o,r,s=e[0].Label+"";return{c(){n=h("span"),o=h("span"),r=y(s),w(o,"class","label"),w(n,"class","tray-menu svelte-1vq8f40")},m(e,t){g(e,n,t),p(n,o),p(o,r)},p(e,[t]){1&t&&s!==(s=e[0].Label+"")&&function(e,t){t=""+t,e.wholeText!==t&&(e.data=t)}(r,s)},i:t,o:t,d(e){e&&m(n)}}}function fe(e,t,n){let{tray:o}=t;return console.log({tray:o}),e.$$set=e=>{"tray"in e&&n(0,o=e.tray)},[o]}class Ae extends K{constructor(e){var t;super(),document.getElementById("svelte-1vq8f40-style")||((t=h("style")).id="svelte-1vq8f40-style",t.textContent=".tray-menu.svelte-1vq8f40{padding-left:0.5rem;padding-right:0.5rem;overflow:visible;font-size:14px}",p(document.head,t)),W(this,e,fe,de,l,{tray:0})}}function pe(e,t,n){const o=e.slice();return o[2]=t[n],o}function ge(e){let t,n,o,r,s=e[1],i=[];for(let t=0;t<s.length;t+=1)i[t]=me(pe(e,s,t));const l=e=>Z(i[e],1,1,(()=>{i[e]=null}));return{c(){t=h("div"),n=h("span");for(let e=0;e<i.length;e+=1)i[e].c();w(n,"class","tray-menus svelte-iy5wor"),w(t,"class","wails-menubar svelte-iy5wor")},m(e,o){g(e,t,o),p(t,n);for(let e=0;e<i.length;e+=1)i[e].m(n,null);r=!0},p(e,t){if(2&t){let o;for(s=e[1],o=0;o<s.length;o+=1){const r=pe(e,s,o);i[o]?(i[o].p(r,t),Q(i[o],1)):(i[o]=me(r),i[o].c(),Q(i[o],1),i[o].m(n,null))}for(X(),o=s.length;o<i.length;o+=1)l(o);z()}},i(e){if(!r){for(let e=0;e<s.length;e+=1)Q(i[e]);D((()=>{o||(o=Y(t,ie,{},!0)),o.run(1)})),r=!0}},o(e){i=i.filter(Boolean);for(let e=0;e<i.length;e+=1)Z(i[e]);o||(o=Y(t,ie,{},!1)),o.run(0),r=!1},d(e){e&&m(t),function(e,t){for(let n=0;n<e.length;n+=1)e[n]&&e[n].d(t)}(i,e),e&&o&&o.end()}}}function me(e){let t,n;return t=new Ae({props:{tray:e[2]}}),{c(){var e;(e=t.$$.fragment)&&e.c()},m(e,o){F(t,e,o),n=!0},p(e,n){const o={};2&n&&(o.tray=e[2]),t.$set(o)},i(e){n||(Q(t.$$.fragment,e),n=!0)},o(e){Z(t.$$.fragment,e),n=!1},d(e){H(t,e)}}}function he(e){let t,n,o=e[0]&&ge(e);return{c(){o&&o.c(),t=b()},m(e,r){o&&o.m(e,r),g(e,t,r),n=!0},p(e,[n]){e[0]?o?(o.p(e,n),1&n&&Q(o,1)):(o=ge(e),o.c(),Q(o,1),o.m(t.parentNode,t)):o&&(X(),Z(o,1,1,(()=>{o=null})),z())},i(e){n||(Q(o),n=!0)},o(e){Z(o),n=!1},d(e){o&&o.d(e),e&&m(t)}}}function ye(e,t,n){let o,r;return a(e,re,(e=>n(0,o=e))),a(e,se,(e=>n(1,r=e))),[o,r]}class be extends K{constructor(e){var t;super(),document.getElementById("svelte-iy5wor-style")||((t=h("style")).id="svelte-iy5wor-style",t.textContent=".tray-menus.svelte-iy5wor{display:flex;flex-direction:row;justify-content:flex-end}.wails-menubar.svelte-iy5wor{position:relative;display:block;top:0;height:2rem;width:100%;border-bottom:1px solid #b3b3b3;box-shadow:antiquewhite;box-shadow:0px 0px 10px 0px #33333360}",p(document.head,t)),W(this,e,ye,he,l,{})}}let we,ve=null,xe=null;function $e(e){xe=e,window.onbeforeunload=function(){ve&&(ve.onclose=function(){},ve.close(),ve=null)},Oe()}function ke(){te("Connected to backend"),window.webkit={messageHandlers:{external:{postMessage:e=>{ve.send(e)}},windowDrag:{postMessage:()=>{}}}},ne.set(!1),clearInterval(we),ve.onclose=_e,ve.onmessage=Ce}function _e(){te("Disconnected from backend"),ve=null,oe(),Oe()}function Oe(){we=setInterval((function(){null==ve&&(ve=new WebSocket("ws://"+window.location.hostname+":34115/bridge"),ve.onopen=ke,ve.onerror=function(e){return e.stopImmediatePropagation(),e.stopPropagation(),e.preventDefault(),ve=null,!1})}),1e3)}function Ce(e){switch(e.data[0]){case"b":(function(e,t){const n=document.createElement("script");n.setAttribute("type","text/javascript"),n.textContent=e,document.head.appendChild(n),t&&n.parentNode.removeChild(n)})(e=e.data.slice(1)),te("Loaded Wails Runtime"),window.webkit.messageHandlers.external.postMessage("SS"),xe&&(te("Notifying application"),xe(window.wails));break;case"c":const t=e.data.slice(1);window.wails._.Callback(t);break;case"T":const n=e.data.slice(1);switch(n[0]){case"S":const t=n.slice(1);!function(e){te("Set Tray:"+JSON.stringify(e)),se.update((t=>{const n=t.findIndex((t=>t.ID===e.ID));return-1===n?t.push(e):t[n]=e,t}))}(JSON.parse(t));break;case"U":const o=n.slice(1);!function(e){te("Update Tray Label:"+JSON.stringify(e)),se.update((t=>{const n=t.findIndex((t=>t.ID===e.ID));return-1===n?te("ERROR: Attempted to update tray index ",e.ID):(t[n].Label=e.Label,t)}))}(JSON.parse(o));break;default:te("Unknown tray message: "+e.data)}break;default:te("Unknown message: "+e.data)}}e.InitBridge=function(e){new be({target:document.body}),new ue({target:document.body,anchor:document.querySelector("#wails-bridge")}),oe(),$e(e)},Object.defineProperty(e,"__esModule",{value:!0})}));