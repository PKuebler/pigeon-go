const fs = require('node:fs');

const Pigeon=(()=>{var v=(t,e)=>()=>(e||t((e={exports:{}}).exports,e),e.exports);var b=v(($e,D)=>{var w={strict:!0,getObjectId:t=>t.id||t._id||t.uuid||t.slug,getTimestamp:Date.now};function ne(t){Object.assign(w,t)}function ie(t,e,i){if(i){let r=O(i);r&&(e=`[${r}]`)}return ce(t,e)}function re(t){return typeof t=="string"&&(t.indexOf("/")!==-1||t.indexOf("~")!==-1)?t.replace(/~/g,"~0").replace(/\//g,"~1"):t}function se(t){return typeof t=="string"&&(t.indexOf("~1")!==-1||t.indexOf("~0")!==-1)?t.replace(/~1/g,"/").replace(/~0/g,"~"):t}function oe(t){return t.split("/").map(e=>se(e))}function ce(t,e){return e=re(e),[t,e].filter(i=>i!=null).join("/").replace("//","/")}function d(t){return Array.isArray(t)?"array":t===null?"null":typeof t}function q(t){let e=d(t);return e==="number"||e==="null"||e==="boolean"||e=="string"}function fe(t){return d(t)==="array"&&t.every(e=>q(e))}function N(t){let e=d(t);if(e=="array"){let i=Array(t.length);for(let r=0;r<t.length;r++)i[r]=N(t[r]);return i}else if(e=="object"){if(t.toJSON)return t.toJSON();{let i={};for(let r in t)i[r]=N(t[r]);return i}}else if(q(t))return typeof t=="number"?isFinite(t)?t:null:t}function ae(t,e){if(q(t))return t===e;if(d(t)=="object")return O(t)===O(e);if(d(t)=="array")throw new Error("can't compare arrays of arrays")}function O(t){if(d(t)=="object"){let e=w.getObjectId(t);if(e!=null)return e;if(w.strict)throw new Error("couldn't find id for object",{cause:t});return M(m(t))}else return null}function ue(t,e,i){let r={op:t,path:e};return Object.assign(r,i),r}function m(t){return d(t)=="array"?`[${t.map(m).join(",")}]`:d(t)=="object"?`{${Object.keys(t).sort().map(e=>`${JSON.stringify(e)}:${m(t[e])}`).join(",")}}`:JSON.stringify(t)}function M(t){return Math.abs([].reduce.call(t,(e,i,r,n)=>(e<<5)-e+n.charCodeAt(r),0))}function le(t){return M(m(t))}D.exports={_path:ie,_typeof:d,_isPrimitive:q,_isArrayOfPrimitives:fe,_clone:N,_entangled:ae,_objId:O,_op:ue,_stable:m,_crc:le,_decodePath:oe,_config:w,_configure:ne}});var S=v((Ae,L)=>{var{_path:p,_typeof:y,_isPrimitive:F,_isArrayOfPrimitives:I,_clone:C,_entangled:pe,_op:P,_config:ge}=b();function he(t,e){let i=y(t);if(i!==y(e))throw new Error("can't diff different types");if(F(t))return k(t,e);if(I(t))return K(t,e);if(i=="array")return G(t,e);if(i=="object")return H(t,e);throw new Error("unsupported type")}function k(t,e,i="/"){return t!==e?[P("replace",p(i),{value:e,_prev:t})]:[]}function K(t,e,i="/"){let r=!I(e)||t.length!==e.length;if(!r)for(let n=0;n<t.length;n++)r=r||t[n]!==e[n];return r?[P("replace",p(i),{value:e,_prev:t})]:[]}function G(t,e,i="/"){let r={},n={},s=[];for(let f=0;f<t.length;f++)for(let o=0;o<e.length;o++)o in n||f in r||(!ge.strict&&y(t[f])=="array"&&y(e[o])=="array"&&f==o||pe(t[f],e[o]))&&(r[f]=o,n[o]=f);let c=[];for(let f=0,o=0;o<e.length||f<t.length;){if(o in e&&f in t&&n[o]==f){y(e[o])==="object"&&c.push(...H(t[f],e[o],p(i,f,e[o]))),o++,f++;continue}if(f<t.length&&!(f in r)){c.push(P("remove",p(i,f,t[f]),{_prev:t[f]})),f++;continue}if(o<e.length&&!(o in n)){s.unshift(P("add",p(i,o,e[o+1]),{value:e[o]})),o++;continue}if(o<e.length&&o in n){let l=p(i,n[o],t[n[o]]),g=p(i,o);g!=l&&(c.push({op:"move",from:l,path:g}),y(n[o])=="object"&&c.push(...H(t[n[o]],e[o],i))),f++,o++;continue}throw new Error("couldn't create diff")}return c.concat(s)}function H(t,e,i="/",r){let n=[],s=Object.keys(t),c=s.length,f=0;for(let g=0;g<c;g++){let a=s[g];if(!e.hasOwnProperty(a)){f++,n.push({op:"remove",path:p(i,a),_prev:C(t[a])});continue}if(t[a]===e[a])continue;let h=y(t[a]);F(t[a])?n.push(...k(t[a],e[a],p(i,a),r)):I(t[a])?n.push(...K(t[a],e[a],p(i,a),r)):h!==y(e[a])?n.push({op:"replace",path:p(i,a),value:C(e[a]),_prev:C(t[a])}):h==="array"?n.push(...G(t[a],e[a],p(i,a))):h==="object"&&n.push(...H(t[a],e[a],p(i,a),r))}let o=Object.keys(e),l=o.length;if(l>c-f)for(let g=0;g<l;g++){let a=o[g];t.hasOwnProperty(a)||n.push({op:"add",path:p(i,a),value:C(e[a])})}return n}L.exports=he});var W=v((Ee,z)=>{var{_typeof:de,_clone:j,_objId:R,_decodePath:ye}=b();function ve(t,e){e=j(e);let i=[],r=null;e:for(let[n,s]of e.entries()){let c=ye(s.path),f=c.shift(),o=c.pop(),l=t;for(let h of c){if(!l){i.push(s);continue e}let _=Y(h);_?l=l.find(te=>R(te)==_):l=l[h]}let g=Y(o);if(g){let h=l.findIndex(_=>R(_)==g);~h?o=h:i.push(s)}let a=de(l);if(s.op=="replace")l[o]=j(s.value);else if(s.op=="move"){r={};let h=[{op:"remove",path:s.from},{op:"add",path:s.path,value:r}];e.splice(n+1,0,...h)}else if(s.op=="remove"){if(a=="object")r&&(r.value=j(l[o])),delete l[o];else if(a=="array"){let h=l.splice(o,1);r&&([r.value]=h)}}else s.op=="add"&&(a=="object"?l[o]=j(s.value):a=="array"&&(r&&s.value===r?(l.splice(o,0,r.value),r=null):l.splice(o,0,j(s.value))))}return t}function Y(t){if(t===void 0)return;let e=t.match(/^\[(.+)\]$/);if(e)return e[1]}z.exports=ve});var $=v((Je,B)=>{var{_clone:me,_objId:be}=b();function je(t){let e=me(t).reverse();for(let n of e){if(n.op=="add"){n.op="remove";let s=be(n.value);if(s){let c=n.path.split("/"),f=c[c.length-1];!f.startsWith("[")&&!f.endsWith("]")?(c[c.length-1]=`[${s}]`,n.path=c.join("/")):f.startsWith("[")&&f.endsWith("]")&&(c[c.length-1]=`[${s}]`,n.path=c.join("/"))}}else if(n.op=="remove"){n.op="add";let s=n.path.split("/"),c=s[s.length-1];c.startsWith("[")&&c.endsWith("]")&&(s[s.length-1]="0",n.path=s.join("/"))}if("_prev"in n)var i=n._prev;if("value"in n)var r=n.value;i===void 0?delete n.value:n.value=i,r===void 0?delete n._prev:n._prev=r}return e}B.exports=je});var Z=v((Te,X)=>{var _e=S(),A=W(),we=$(),{_clone:E,_crc:Oe,_configure:qe,_config:Q}=b(),J=1e3,u=new WeakMap,U=V(),T=class t{constructor(){u.set(this,{history:[],stash:[],warning:null,gids:{}})}static from(e,i=U){let r=new t;return u.get(r).cid=i,r=t.change(r,n=>Object.assign(n,e)),r}static _forge(e,i=U){let r=new t;return u.get(r).cid=i,Object.assign(r,E(e)),r}static alias(e){let i=new t;return u.set(i,u.get(e)),Object.assign(i,e),i}static init(){return t.from({})}static clone(e,i=J){let r=t._forge(e);return u.get(r).history=u.get(e).history,u.get(r).gids=E(u.get(e).gids),t.pruneHistory(u.get(r),i),r}static pruneHistory(e,i){let r=e.history.length;if(r>i){let n=e.history.slice(0,r-i);for(let s of n)delete e.gids[s.gid]}e.history=e.history.slice(-i)}static getChanges(e,i){return{diff:_e(e,i),cid:u.get(e).cid,ts:Q.getTimestamp(),seq:Pe(),gid:V()}}static rewindChanges(e,i,r){let{history:n}=u.get(e);for(;!(n.length<=1);){let s=n[n.length-1];if(s.ts>i||s.ts==i&&s.cid>r){let c=u.get(e).history.pop();A(e,we(c.diff)),delete u.get(e).gids[c.gid],u.get(e).stash.push(c);continue}break}}static fastForwardChanges(e){let{stash:i,history:r}=u.get(e),n;for(;n=i.pop();)A(e,n.diff),u.get(e).gids[n.gid]=1,r.push(n)}static applyChangesInPlace(e,i){return t.applyChanges(e,i,!0)}static applyChanges(e,i,r){u.get(e).warning=null;let n=r?e:t.clone(e);if(u.get(e).gids[i.gid])return n;try{t.rewindChanges(n,i.ts,i.cid)}catch(f){u.get(n).warning="rewind failed: "+f}try{A(n,i.diff),u.get(n).gids[i.gid]=1}catch(f){u.get(n).warning="patch failed: "+f}try{t.fastForwardChanges(n)}catch(f){u.get(n).warning="forward failed: "+f}let s=u.get(n).history,c=s.length;for(;c>1&&s[c-1].ts>i.ts;)c--;return s.splice(c,0,i),n}static change(e,i){let r=E(e);i(r);let n=t.getChanges(e,r);return t.applyChanges(e,n)}static getHistory(e){return u.get(e).history}static merge(e,i){let r=t.from({}),n=t.getHistory(e),s=t.getHistory(i),c=[];for(;n.length||s.length;)s.length?n.length?n[0].gid===s[0].gid?c.push(n.shift()&&s.shift()):n[0].ts<s[0].ts||n[0].ts==s[0].ts&&n[0].seq<s[0].seq?c.push(n.shift()):c.push(s.shift()):c.push(s.shift()):c.push(n.shift());for(let f of c)r=t.applyChanges(r,f);return r}static getWarning(e){return u.get(e).warning}static getMissingDeps(e){return!1}static setHistoryLength(e){J=e}static setTimestamp(e){Q.getTimestamp=e}static crc(e){return Oe(e)}static load(e,i=J){let{meta:r,data:n}=JSON.parse(e);t.pruneHistory(r,i);let s=t.from(n);return Object.assign(u.get(s),r),s}static save(e){let{cid:i,...r}=u.get(e);return JSON.stringify({meta:r,data:e})}static configure(e){qe(e)}};function V(){return Math.random().toString(36).substring(2)}var Ce=0;function Pe(){return Ce++}X.exports=T});var Se=v((Me,ee)=>{var He=S(),Ne=W(),Ie=$(),x=Z();ee.exports=Object.assign(x,{auto:x,diff:He,patch:Ne,reverse:Ie})});return Se();})();

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