import{g as r}from"./navigation.2acacab1.js";async function n(o,s){try{return await fetch(o,s).then(e=>(e.status==401&&r("/"),{code:200,msg:"lets go",data:e.json()}))}catch(t){return console.log(t),{code:500,msg:"OOps",data:void 0}}}export{n};