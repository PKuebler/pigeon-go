const fs = require('node:fs');

const Pigeon=(()=>{var y=(t,e)=>()=>(e||t((e={exports:{}}).exports,e),e.exports);var v=y(($e,F)=>{var w={strict:!0,getObjectId:t=>t.id||t._id||t.uuid||t.slug,getTimestamp:Date.now};function ne(t){Object.assign(w,t)}function ie(t,e,i){if(i){let s=O(i);s&&(e=`[${s}]`)}return ce(t,e)}function se(t){return typeof t=="string"&&(t.indexOf("/")!==-1||t.indexOf("~")!==-1)?t.replace(/~/g,"~0").replace(/\//g,"~1"):t}function re(t){return typeof t=="string"&&(t.indexOf("~1")!==-1||t.indexOf("~0")!==-1)?t.replace(/~1/g,"/").replace(/~0/g,"~"):t}function oe(t){return t.split("/").map(e=>re(e))}function ce(t,e){return e=se(e),[t,e].filter(i=>i!=null).join("/").replace("//","/")}function d(t){return Array.isArray(t)?"array":t===null?"null":typeof t}function q(t){let e=d(t);return e==="number"||e==="null"||e==="boolean"||e=="string"}function ae(t){return d(t)==="array"&&t.every(e=>q(e))}function H(t){let e=d(t);if(e=="array"){let i=Array(t.length);for(let s=0;s<t.length;s++)i[s]=H(t[s]);return i}else if(e=="object"){if(t.toJSON)return t.toJSON();{let i={};for(let s in t)i[s]=H(t[s]);return i}}else if(q(t))return typeof t=="number"?isFinite(t)?t:null:t}function fe(t,e){if(q(t))return t===e;if(d(t)=="object")return O(t)===O(e);if(d(t)=="array")throw new Error("can't compare arrays of arrays")}function O(t){if(d(t)=="object"){let e=w.getObjectId(t);if(e!=null)return e;if(w.strict)throw new Error("couldn't find id for object",{cause:t});return D(_(t))}else return null}function ue(t,e,i){let s={op:t,path:e};return Object.assign(s,i),s}function _(t){return d(t)=="array"?`[${t.map(_).join(",")}]`:d(t)=="object"?`{${Object.keys(t).sort().map(e=>`${JSON.stringify(e)}:${_(t[e])}`).join(",")}}`:JSON.stringify(t)}function D(t){return Math.abs([].reduce.call(t,(e,i,s,n)=>(e<<5)-e+n.charCodeAt(s),0))}function le(t){return D(_(t))}F.exports={_path:ie,_typeof:d,_isPrimitive:q,_isArrayOfPrimitives:ae,_clone:H,_entangled:fe,_objId:O,_op:ue,_stable:_,_crc:le,_decodePath:oe,_config:w,_configure:ne}});var S=y((Ae,L)=>{var{_path:p,_typeof:m,_isPrimitive:M,_isArrayOfPrimitives:N,_clone:I,_entangled:pe,_op:C,_config:he}=v();function ge(t,e){let i=m(t);if(i!==m(e))throw new Error("can't diff different types");if(M(t))return k(t,e);if(N(t))return K(t,e);if(i=="array")return G(t,e);if(i=="object")return P(t,e);throw new Error("unsupported type")}function k(t,e,i="/"){return t!==e?[C("replace",p(i),{value:e,_prev:t})]:[]}function K(t,e,i="/"){let s=!N(e)||t.length!==e.length;if(!s)for(let n=0;n<t.length;n++)s=s||t[n]!==e[n];return s?[C("replace",p(i),{value:e,_prev:t})]:[]}function G(t,e,i="/"){let s={},n={},r=[];for(let a=0;a<t.length;a++)for(let o=0;o<e.length;o++)o in n||a in s||(!he.strict&&m(t[a])=="array"&&m(e[o])=="array"&&a==o||pe(t[a],e[o]))&&(s[a]=o,n[o]=a);let c=[];for(let a=0,o=0;o<e.length||a<t.length;){if(o in e&&a in t&&n[o]==a){m(e[o])==="object"&&c.push(...P(t[a],e[o],p(i,a,e[o]))),o++,a++;continue}if(a<t.length&&!(a in s)){c.push(C("remove",p(i,a,t[a]),{_prev:t[a]})),a++;continue}if(o<e.length&&!(o in n)){r.unshift(C("add",p(i,o,e[o+1]),{value:e[o]})),o++;continue}if(o<e.length&&o in n){let l=p(i,n[o],t[n[o]]),h=p(i,o);h!=l&&(c.push({op:"move",from:l,path:h}),m(n[o])=="object"&&c.push(...P(t[n[o]],e[o],i))),a++,o++;continue}throw new Error("couldn't create diff")}return c.concat(r)}function P(t,e,i="/",s){let n=[],r=Object.keys(t),c=r.length,a=0;for(let h=0;h<c;h++){let f=r[h];if(!e.hasOwnProperty(f)){a++,n.push({op:"remove",path:p(i,f),_prev:I(t[f])});continue}if(t[f]===e[f])continue;let g=m(t[f]);M(t[f])?n.push(...k(t[f],e[f],p(i,f),s)):N(t[f])?n.push(...K(t[f],e[f],p(i,f),s)):g!==m(e[f])?n.push({op:"replace",path:p(i,f),value:I(e[f]),_prev:I(t[f])}):g==="array"?n.push(...G(t[f],e[f],p(i,f))):g==="object"&&n.push(...P(t[f],e[f],p(i,f),s))}let o=Object.keys(e),l=o.length;if(l>c-a)for(let h=0;h<l;h++){let f=o[h];t.hasOwnProperty(f)||n.push({op:"add",path:p(i,f),value:I(e[f])})}return n}L.exports=ge});var W=y((Ee,z)=>{var{_typeof:de,_clone:b,_objId:R,_decodePath:me}=v();function ye(t,e){e=b(e);let i=[],s=null;e:for(let[n,r]of e.entries()){let c=me(r.path),a=c.shift(),o=c.pop(),l=t;for(let g of c){if(!l){i.push(r);continue e}let j=Y(g);j?l=l.find(te=>R(te)==j):l=l[g]}let h=Y(o);if(h){let g=l.findIndex(j=>R(j)==h);~g?o=g:i.push(r)}let f=de(l);if(r.op=="replace")l[o]=b(r.value);else if(r.op=="move"){s={};let g=[{op:"remove",path:r.from},{op:"add",path:r.path,value:s}];e.splice(n+1,0,...g)}else if(r.op=="remove"){if(f=="object")s&&(s.value=b(l[o])),delete l[o];else if(f=="array"){let g=l.splice(o,1);s&&([s.value]=g)}}else r.op=="add"&&(f=="object"?l[o]=b(r.value):f=="array"&&(s&&r.value===s?(l.splice(o,0,s.value),s=null):l.splice(o,0,b(r.value))))}return t}function Y(t){if(t===void 0)return;let e=t.match(/^\[(.+)\]$/);if(e)return e[1]}z.exports=ye});var $=y((Je,B)=>{var{_clone:_e,_objId:ve}=v();function be(t){let e=_e(t).reverse();for(let n of e){if(n.op=="add"){n.op="remove";let r=ve(n.value);if(r){let c=n.path.split("/"),a=c[c.length-1];!a.startsWith("[")&&!a.endsWith("]")?(c[c.length-1]=`[${r}]`,n.path=c.join("/")):a.startsWith("[")&&a.endsWith("]")&&(c[c.length-1]=`[${r}]`,n.path=c.join("/"))}}else if(n.op=="remove"){n.op="add";let r=n.path.split("/"),c=r[r.length-1];c.startsWith("[")&&c.endsWith("]")&&(r[r.length-1]="0",n.path=r.join("/"))}if("_prev"in n)var i=n._prev;if("value"in n)var s=n.value;i===void 0?delete n.value:n.value=i,s===void 0?delete n._prev:n._prev=s}return e}B.exports=be});var Z=y((Te,X)=>{var je=S(),A=W(),we=$(),{_clone:E,_crc:Oe,_configure:qe,_config:Q}=v(),J=1e3,u=new WeakMap,U=V(),T=class t{constructor(){u.set(this,{history:[],stash:[],warning:null,changeIds:{}})}static from(e,i=U){let s=new t;return u.get(s).client_id=i,s=t.change(s,n=>Object.assign(n,e)),s}static _forge(e,i=U){let s=new t;return u.get(s).client_id=i,Object.assign(s,E(e)),s}static alias(e){let i=new t;return u.set(i,u.get(e)),Object.assign(i,e),i}static init(){return t.from({})}static clone(e,i=J){let s=t._forge(e);return u.get(s).history=u.get(e).history,u.get(s).changeIds=E(u.get(e).changeIds),t.pruneHistory(u.get(s),i),s}static pruneHistory(e,i){let s=e.history.length;if(s>i){let n=e.history.slice(0,s-i);for(let r of n)delete e.changeIds[r.change_id]}e.history=e.history.slice(-i)}static getChanges(e,i){return{diff:je(e,i),client_id:u.get(e).client_id,timestamp_ms:Q.getTimestamp(),seq:Ce(),change_id:V()}}static rewindChanges(e,i,s){let{history:n}=u.get(e);for(;!(n.length<=1);){let r=n[n.length-1];if(r.timestamp_ms>i||r.timestamp_ms==i&&r.client_id>s){let c=u.get(e).history.pop();A(e,we(c.diff)),delete u.get(e).changeIds[c.change_id],u.get(e).stash.push(c);continue}break}}static fastForwardChanges(e){let{stash:i,history:s}=u.get(e),n;for(;n=i.pop();)A(e,n.diff),u.get(e).changeIds[n.change_id]=1,s.push(n)}static applyChangesInPlace(e,i){return t.applyChanges(e,i,!0)}static applyChanges(e,i,s){u.get(e).warning=null;let n=s?e:t.clone(e);if(u.get(e).changeIds[i.change_id])return n;try{t.rewindChanges(n,i.timestamp_ms,i.client_id)}catch(a){u.get(n).warning="rewind failed: "+a}try{A(n,i.diff),u.get(n).changeIds[i.change_id]=1}catch(a){u.get(n).warning="patch failed: "+a}try{t.fastForwardChanges(n)}catch(a){u.get(n).warning="forward failed: "+a}let r=u.get(n).history,c=r.length;for(;c>1&&r[c-1].timestamp_ms>i.timestamp_ms;)c--;return r.splice(c,0,i),n}static change(e,i){let s=E(e);i(s);let n=t.getChanges(e,s);return t.applyChanges(e,n)}static getHistory(e){return u.get(e).history}static merge(e,i){let s=t.from({}),n=t.getHistory(e),r=t.getHistory(i),c=[];for(;n.length||r.length;)r.length?n.length?n[0].change_id===r[0].change_id?c.push(n.shift()&&r.shift()):n[0].timestamp_ms<r[0].timestamp_ms||n[0].timestamp_ms==r[0].timestamp_ms&&n[0].seq<r[0].seq?c.push(n.shift()):c.push(r.shift()):c.push(r.shift()):c.push(n.shift());for(let a of c)s=t.applyChanges(s,a);return s}static getWarning(e){return u.get(e).warning}static getMissingDeps(e){return!1}static setHistoryLength(e){J=e}static setTimestamp(e){Q.getTimestamp=e}static crc(e){return Oe(e)}static load(e,i=J){let{meta:s,data:n}=JSON.parse(e);t.pruneHistory(s,i);let r=t.from(n);return Object.assign(u.get(r),s),r}static save(e){let{client_id:i,...s}=u.get(e);return JSON.stringify({meta:s,data:e})}static configure(e){qe(e)}};function V(){return Math.random().toString(36).substring(2)}var Ie=0;function Ce(){return Ie++}X.exports=T});var Se=y((De,ee)=>{var Pe=S(),He=W(),Ne=$(),x=Z();ee.exports=Object.assign(x,{auto:x,diff:Pe,patch:He,reverse:Ne})});return Se();})();

Pigeon.configure({
    strict: false,
    getObjectId: (x) => x.attrs?.id || x.id || x._id || x.uuid || x.slug,
    getTimestamp: Date.now,
})

const testCases = [
    {
        oldDocument: {
            id: '0',
            name: 'Test',
        },
        newDocument: {
            id: '0',
            name: 'Foo',
        }
    },
    {
        oldDocument: {
            id: '0',
            name: 'Test',
        },
        newDocument: {
            id: '1',
            name: 'Test',
        }
    },
    {
        oldDocument: {
            id: '0',
            name: 'Test',
        },
        newDocument: {
            id: '0',
        }
    },
    {
        oldDocument: {
            id: '0',
        },
        newDocument: {
            id: '0',
            name: 'Test',
        }
    },
    {
        oldDocument: {
            id: '0',
            name: null,
        },
        newDocument: {
            id: '0',
            name: null,
        }
    },
    {
        oldDocument: {
            id: '0',
            name: 1234,
        },
        newDocument: {
            id: '0',
            name: "Name",
        }
    },
    {
        oldDocument: {
            id: 135,
            name: 1234,
        },
        newDocument: {
            id: 135,
            name: "Name",
        }
    },
    {
        oldDocument: {
            channel: {
                id: '2',
                name: 'Test Channel',
                subchannel: {
                    id: '3',
                    name: 'Test Subchannel',
                }
            }
        },
        newDocument: {
            channel: {
                id: '2',
                name: 'Test Channel',
                subchannel: {
                    id: '0',
                    name: 'No Channel',
                }
            }
        },
    },
    {
        oldDocument: {
            channel: {
                id: '2',
                name: 'Test Channel',
                subchannel: {},
            }
        },
        newDocument: {
            channel: {
                id: '2',
                name: 'Test Channel',
                subchannel: {
                    id: '0',
                    name: 'No Channel',
                }
            }
        },
    },
    {
        oldDocument: {
            array: [0, 10, 20, 30, 40]
        },
        newDocument: {
            array: [0, 10, 20]
        },
    },
    {
        oldDocument: {
            array: [0, 10, 20, 30, 40]
        },
        newDocument: {
            array: [0, 30, 40, 10]
        },
    },
    {
        oldDocument: {
            array: [0, 10, 20, 30, 40]
        },
        newDocument: {
            array: [0, 10, 40, 30]
        },
    },
    {
        oldDocument: {
            array: ["one", "two", "three"],
        },
        newDocument: {
            array: ["one", "two"]
        },
    },
    {
        oldDocument: {
            array: ["one", "two", "three", "four"],
        },
        newDocument: {
            array: ["one", "two", "three"]
        },
    },
    {
        oldDocument: {
            array: ["one", "two", "three", "four"],
        },
        newDocument: {
            array: ["two", "four", "one"]
        },
    },
    {
        oldDocument: {
            array: [
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "Foo" },
                { id: "hgevcx9ds", name: "Baa" },
            ],
        },
        newDocument: {
            array: [
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "Foo" }
            ]
        },
    },
    {
        oldDocument: {
            array: [
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "Foo" },
                { id: "hgevcx9ds", name: "Baa" },
            ],
        },
        newDocument: {
            array: [
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "123" },
                { id: "hgevcx9ds", name: "World" },
            ]
        },
    },
    {
        oldDocument: {
            array: [
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "Foo" },
                { id: "hgevcx9ds", name: "Baa" },
            ],
        },
        newDocument: {
            array: [
                { id: "hgevcx9ds", name: "World" },
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "123" },
            ]
        },
    },
    {
        oldDocument: {
            array: [
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "Foo" },
                { id: "hgevcx9ds", name: "Baa" },
            ],
        },
        newDocument: {
            array: [
                { id: "hgevcx9ds", name: "World" },
                { id: "ret34sdf4", name: "123" },
            ]
        },
    },
    {
        oldDocument: {
            aaaa: "Change me",
            array: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
                { id: "hgevcx9ds", name: "Baa" },
                { id: "ret34sdf4", name: "World" },
            ],
            strings: [
                "Foo",
                "Bar",
                "Baz",
            ],
            title: "Hallo World"
        },
        newDocument: {
            aaaa: "Ok!",
            array: [
                { id: "sdfw4ldhf", name: "Me" },
                { id: "ret34sdf4", name: "World" },
                { id: "434dfsdsf", name: "Hello" },
                { id: "hgevcx9ds", name: "Baa" },
            ],
            strings: [
                "Foo",
                "Welt",
                "Baz"
            ],
            title: "World Hallo"
        },
    },
    {
        oldDocument: {
            arrayOne: [
                { id: "hgevcx9ds", name: "Baa" },
                { id: "ret34sdf4", name: "World" },
            ],
            arrayTwo: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
            ],
        },
        newDocument: {
            arrayOne: [
                { id: "ret34sdf4", name: "World" },
            ],
            arrayTwo: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "hgevcx9ds", name: "Baa" },
                { id: "434dfsdsf", name: "Hello" },
            ],
        },
    },
    {
        oldDocument: {
            arrayOne: [
                { id: "hgevcx9ds", name: "Baa" },
                { id: "ret34sdf4", name: "World" },
            ],
            arrayTwo: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
            ],
        },
        newDocument: {
            arrayOne: [
                { id: "ret34sdf4", name: "World" },
            ],
            arrayTwo: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
                { id: "hgevcx9ds", name: "Baa" },
            ],
        },
    },
    {
        oldDocument: {
            arrayOne: [
                { id: "hgevcx9ds", name: "Baa" },
                { id: "ret34sdf4", name: "World" },
            ],
            arrayTwo: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
            ],
        },
        newDocument: {
            arrayOne: [
                { id: "ret34sdf4", name: "World" },
            ],
            arrayTwo: [
                { id: "434dfsdsf", name: "Hello" },
                { id: "hgevcx9ds", name: "Baa" },
            ],
        },
    },
    {
        oldDocument: {
            arrayOne: [
                { id: "ret34sdf4", name: "World" },
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
            ],
        },
        newDocument: {
            arrayOne: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "World" },
            ],
        },
    },
    {
        oldDocument: {
            arrayOne: [
                { id: "ret34sdf4", name: "World" },
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
            ],
        },
        newDocument: {
            arrayOne: [
                { id: "hgevcx9ds", name: "Baa" },
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
                { id: "ret34sdf4", name: "World" },
            ],
        },
    },
    {
        oldDocument: {
            arrayOne: [
                { id: "ret34sdf4", name: "World" },
                { id: "sdfw4ldhf", name: "Change" },
            ],
        },
        newDocument: {
            arrayOne: [
                { id: "sdfw4ldhf", name: "Change" },
                { id: "434dfsdsf", name: "Hello" },
                { id: "hgevcx9ds", name: "Baa" },
                { id: "ret34sdf4", name: "World" },
            ],
        },
    }
];

for (let i = 0; i < testCases.length; ++i) {
    const testCase = testCases[i];
    const oldDocument = testCase.oldDocument;
    const newDocument = testCase.newDocument;

    const doc = Pigeon.from(oldDocument)
    const change = Pigeon.getChanges(doc, newDocument)
    testCases[i].change = change;
    Pigeon.applyChangesInPlace(doc, change);
    testCases[i].result = doc;
    testCases[i].name = `testcase ${i}`
}

fs.writeFileSync("testfiles/test-cases.json", JSON.stringify(testCases));