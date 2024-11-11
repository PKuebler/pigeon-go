const fs = require('node:fs');

const Pigeon=(()=>{var v=(t,e)=>()=>(e||t((e={exports:{}}).exports,e),e.exports);var b=v((Ae,F)=>{var j={strict:!0,getObjectId:t=>t.id||t._id||t.uuid||t.slug,getTimestamp:Date.now};function ne(t){Object.assign(j,t)}function ie(t,e,i){if(i){let r=O(i);r&&(e=`[${r}]`)}return ce(t,e)}function re(t){return typeof t=="string"&&(t.indexOf("/")!==-1||t.indexOf("~")!==-1)?t.replace(/~/g,"~0").replace(/\//g,"~1"):t}function se(t){return typeof t=="string"&&(t.indexOf("~1")!==-1||t.indexOf("~0")!==-1)?t.replace(/~1/g,"/").replace(/~0/g,"~"):t}function oe(t){return t.split("/").map(e=>se(e))}function ce(t,e){return e=re(e),[t,e].filter(i=>i!=null).join("/").replace("//","/")}function h(t){return Array.isArray(t)?"array":t===null?"null":typeof t}function q(t){let e=h(t);return e==="number"||e==="null"||e==="boolean"||e=="string"}function fe(t){return h(t)==="array"&&t.every(e=>q(e))}function N(t){let e=h(t);if(e=="array"){let i=Array(t.length);for(let r=0;r<t.length;r++)i[r]=N(t[r]);return i}else if(e=="object"){if(t.toJSON)return t.toJSON();{let i={};for(let r in t)i[r]=N(t[r]);return i}}else if(q(t))return typeof t=="number"?isFinite(t)?t:null:t}function ae(t,e){if(q(t))return t===e;if(h(t)=="object")return O(t)===O(e);if(h(t)=="array")throw new Error("can't compare arrays of arrays")}function O(t){if(h(t)=="object"){let e=j.getObjectId(t);if(e!=null)return e;if(j.strict)throw new Error("couldn't find id for object",{cause:t});return D(m(t))}else return null}function ue(t,e,i){let r={op:t,path:e};return Object.assign(r,i),r}function m(t){return h(t)=="array"?`[${t.map(m).join(",")}]`:h(t)=="object"?`{${Object.keys(t).sort().map(e=>`${JSON.stringify(e)}:${m(t[e])}`).join(",")}}`:JSON.stringify(t)}function D(t){return Math.abs([].reduce.call(t,(e,i,r,n)=>(e<<5)-e+n.charCodeAt(r),0))}function le(t){return D(m(t))}F.exports={_path:ie,_typeof:h,_isPrimitive:q,_isArrayOfPrimitives:fe,_clone:N,_entangled:ae,_objId:O,_op:ue,_stable:m,_crc:le,_decodePath:oe,_config:j,_configure:ne}});var S=v((Ee,L)=>{var{_path:p,_typeof:y,_isPrimitive:k,_isArrayOfPrimitives:I,_clone:C,_entangled:pe,_op:H,_config:ge}=b();function de(t,e){let i=y(t);if(i!==y(e))throw new Error("can't diff different types");if(k(t))return K(t,e);if(I(t))return W(t,e);if(i=="array")return G(t,e);if(i=="object")return P(t,e);throw new Error("unsupported type")}function K(t,e,i="/"){return t!==e?[H("replace",p(i),{value:e,_prev:t})]:[]}function W(t,e,i="/"){let r=!I(e)||t.length!==e.length;if(!r)for(let n=0;n<t.length;n++)r=r||t[n]!==e[n];return r?[H("replace",p(i),{value:e,_prev:t})]:[]}function G(t,e,i="/"){let r={},n={},s=[];for(let c=0;c<t.length;c++)for(let o=0;o<e.length;o++)o in n||c in r||(!ge.strict&&y(t[c])=="array"&&y(e[o])=="array"&&c==o||pe(t[c],e[o]))&&(r[c]=o,n[o]=c);let a=[];for(let c=0,o=0;o<e.length||c<t.length;){if(o in e&&c in t&&n[o]==c){y(e[o])==="object"&&a.push(...P(t[c],e[o],p(i,c,e[o]))),o++,c++;continue}if(c<t.length&&!(c in r)){a.push(H("remove",p(i,c,t[c]),{_prev:t[c]})),c++;continue}if(o<e.length&&!(o in n)){s.unshift(H("add",p(i,o,e[o+1]),{value:e[o]})),o++;continue}if(o<e.length&&o in n){let l=p(i,n[o],t[n[o]]),g=p(i,o);g!=l&&(a.push({op:"move",from:l,path:g}),y(n[o])=="object"&&a.push(...P(t[n[o]],e[o],i))),c++,o++;continue}throw new Error("couldn't create diff")}return a.concat(s)}function P(t,e,i="/",r){let n=[],s=Object.keys(t),a=s.length,c=0;for(let g=0;g<a;g++){let f=s[g];if(!e.hasOwnProperty(f)){c++,n.push({op:"remove",path:p(i,f),_prev:C(t[f])});continue}if(t[f]===e[f])continue;let d=y(t[f]);k(t[f])?n.push(...K(t[f],e[f],p(i,f),r)):I(t[f])?n.push(...W(t[f],e[f],p(i,f),r)):d!==y(e[f])?n.push({op:"replace",path:p(i,f),value:C(e[f]),_prev:C(t[f])}):d==="array"?n.push(...G(t[f],e[f],p(i,f))):d==="object"&&n.push(...P(t[f],e[f],p(i,f),r))}let o=Object.keys(e),l=o.length;if(l>a-c)for(let g=0;g<l;g++){let f=o[g];t.hasOwnProperty(f)||n.push({op:"add",path:p(i,f),value:C(e[f])})}return n}L.exports=de});var $=v((Je,z)=>{var{_typeof:he,_clone:_,_objId:R,_decodePath:ye}=b();function ve(t,e){e=_(e);let i=[],r=null;e:for(let[n,s]of e.entries()){let a=ye(s.path),c=a.shift(),o=a.pop(),l=t;for(let d of a){if(!l){i.push(s);continue e}let w=Y(d);w?l=l.find(te=>R(te)==w):l=l[d]}let g=Y(o);if(g){let d=l.findIndex(w=>R(w)==g);~d?o=d:i.push(s)}let f=he(l);if(s.op=="replace")l[o]=_(s.value);else if(s.op=="move"){r={};let d=[{op:"remove",path:s.from},{op:"add",path:s.path,value:r}];e.splice(n+1,0,...d)}else if(s.op=="remove"){if(f=="object")r&&(r.value=_(l[o])),delete l[o];else if(f=="array"){let d=l.splice(o,1);r&&([r.value]=d)}}else s.op=="add"&&(f=="object"?l[o]=_(s.value):f=="array"&&(r&&s.value===r?(l.splice(o,0,r.value),r=null):l.splice(o,0,_(s.value))))}return t}function Y(t){if(t===void 0)return;let e=t.match(/^\[(.+)\]$/);if(e)return e[1]}z.exports=ve});var A=v((Te,B)=>{var{_clone:me,_objId:be}=b();function _e(t){let e=me(t).reverse();for(let n of e){if(n.op=="add"){n.op="remove";let s=be(n.value);s&&(n._index=n.path.split("/").pop(),n.path=n.path.replace(/\d+$/,`[${s}]`))}else n.op=="remove"&&(n.op="add");if("_prev"in n)var i=n._prev;if("value"in n)var r=n.value;i===void 0?delete n.value:n.value=i,r===void 0?delete n._prev:n._prev=r}return e}B.exports=_e});var Z=v((Me,X)=>{var we=S(),E=$(),je=A(),{_clone:J,_crc:Oe,_configure:qe,_config:Q}=b(),T=1e3,u=new WeakMap,U=V(),M=class t{constructor(){u.set(this,{history:[],stash:[],warning:null,gids:{}})}static from(e,i=U){let r=new t;return u.get(r).cid=i,r=t.change(r,n=>Object.assign(n,e)),r}static _forge(e,i=U){let r=new t;return u.get(r).cid=i,Object.assign(r,J(e)),r}static alias(e){let i=new t;return u.set(i,u.get(e)),Object.assign(i,e),i}static init(){return t.from({})}static clone(e,i=T){let r=t._forge(e);return u.get(r).history=u.get(e).history,u.get(r).gids=J(u.get(e).gids),t.pruneHistory(u.get(r),i),r}static pruneHistory(e,i){let r=e.history.length;if(r>i){let n=e.history.slice(0,r-i);for(let s of n)delete e.gids[s.gid]}e.history=e.history.slice(-i)}static getChanges(e,i){return{diff:we(e,i),cid:u.get(e).cid,ts:Q.getTimestamp(),seq:He(),gid:V()}}static rewindChanges(e,i,r){let{history:n}=u.get(e);for(;!(n.length<=1);){let s=n[n.length-1];if(s.ts>i||s.ts==i&&s.cid>r){let a=u.get(e).history.pop();E(e,je(a.diff)),delete u.get(e).gids[a.gid],u.get(e).stash.push(a);continue}break}}static fastForwardChanges(e){let{stash:i,history:r}=u.get(e),n;for(;n=i.pop();)E(e,n.diff),u.get(e).gids[n.gid]=1,r.push(n)}static applyChangesInPlace(e,i){return t.applyChanges(e,i,!0)}static applyChanges(e,i,r){u.get(e).warning=null;let n=r?e:t.clone(e);if(u.get(e).gids[i.gid])return n;try{t.rewindChanges(n,i.ts,i.cid)}catch(c){u.get(n).warning="rewind failed: "+c}try{E(n,i.diff),u.get(n).gids[i.gid]=1}catch(c){u.get(n).warning="patch failed: "+c}try{t.fastForwardChanges(n)}catch(c){u.get(n).warning="forward failed: "+c}let s=u.get(n).history,a=s.length;for(;a>1&&s[a-1].ts>i.ts;)a--;return s.splice(a,0,i),n}static change(e,i){let r=J(e);i(r);let n=t.getChanges(e,r);return t.applyChanges(e,n)}static getHistory(e){return u.get(e).history}static merge(e,i){let r=t.from({}),n=t.getHistory(e),s=t.getHistory(i),a=[];for(;n.length||s.length;)s.length?n.length?n[0].gid===s[0].gid?a.push(n.shift()&&s.shift()):n[0].ts<s[0].ts||n[0].ts==s[0].ts&&n[0].seq<s[0].seq?a.push(n.shift()):a.push(s.shift()):a.push(s.shift()):a.push(n.shift());for(let c of a)r=t.applyChanges(r,c);return r}static getWarning(e){return u.get(e).warning}static getMissingDeps(e){return!1}static setHistoryLength(e){T=e}static setTimestamp(e){Q.getTimestamp=e}static crc(e){return Oe(e)}static load(e,i=T){let{meta:r,data:n}=JSON.parse(e);t.pruneHistory(r,i);let s=t.from(n);return Object.assign(u.get(s),r),s}static save(e){let{cid:i,...r}=u.get(e);return JSON.stringify({meta:r,data:e})}static configure(e){qe(e)}};function V(){return Math.random().toString(36).substring(2)}var Ce=0;function He(){return Ce++}X.exports=M});var Se=v((De,ee)=>{var Pe=S(),Ne=$(),Ie=A(),x=Z();ee.exports=Object.assign(x,{auto:x,diff:Pe,patch:Ne,reverse:Ie})});return Se();})();

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
    const changes = Pigeon.getChanges(doc, newDocument)
    testCases[i].changes = changes;
    Pigeon.applyChangesInPlace(doc, changes);
    testCases[i].result = doc;
    testCases[i].name = `testcase ${i}`
}

fs.writeFileSync("testfiles/test-cases.json", JSON.stringify(testCases));