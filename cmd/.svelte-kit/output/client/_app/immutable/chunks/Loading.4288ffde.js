import{S as y,i as v,s as L,e as d,b as p,v as q,d as f,f as C,g as _,h as k,K as g,L as b,C as M,D as S,E as D,F as E,k as F,l as G,M as H,n as c,H as h}from"./index.880e7ca1.js";import{s as I}from"./spinner.6475b7ac.js";function K(n){let e;const l=n[3].default,s=M(l,n,n[2],null);return{c(){s&&s.c()},l(t){s&&s.l(t)},m(t,i){s&&s.m(t,i),e=!0},p(t,i){s&&s.p&&(!e||i&4)&&S(s,l,t,t[2],e?E(l,t[2],i,null):D(t[2]),null)},i(t){e||(_(s,t),e=!0)},o(t){f(s,t),e=!1},d(t){s&&s.d(t)}}}function N(n){let e,l,s;return{c(){e=F("img"),this.h()},l(t){e=G(t,"IMG",{src:!0,alt:!0,class:!0}),this.h()},h(){H(e.src,l=I)||c(e,"src",l),c(e,"alt","spinner"),c(e,"class",s=`animate-spin w-10 ${n[1].class} fixed top-0 bottom-0 left-0 right-0 my-auto mx-auto`)},m(t,i){p(t,e,i)},p(t,i){i&2&&s!==(s=`animate-spin w-10 ${t[1].class} fixed top-0 bottom-0 left-0 right-0 my-auto mx-auto`)&&c(e,"class",s)},i:h,o:h,d(t){t&&k(e)}}}function j(n){let e,l,s,t;const i=[N,K],o=[];function m(a,r){return a[0]?0:1}return e=m(n),l=o[e]=i[e](n),{c(){l.c(),s=d()},l(a){l.l(a),s=d()},m(a,r){o[e].m(a,r),p(a,s,r),t=!0},p(a,[r]){let u=e;e=m(a),e===u?o[e].p(a,r):(q(),f(o[u],1,1,()=>{o[u]=null}),C(),l=o[e],l?l.p(a,r):(l=o[e]=i[e](a),l.c()),_(l,1),l.m(s.parentNode,s))},i(a){t||(_(l),t=!0)},o(a){f(l),t=!1},d(a){o[e].d(a),a&&k(s)}}}function x(n,e,l){let{$$slots:s={},$$scope:t}=e,{loading:i=!1}=e;return n.$$set=o=>{l(1,e=g(g({},e),b(o))),"loading"in o&&l(0,i=o.loading),"$$scope"in o&&l(2,t=o.$$scope)},e=b(e),[i,e,t,s]}class B extends y{constructor(e){super(),v(this,e,x,j,L,{loading:0})}}export{B as L};
