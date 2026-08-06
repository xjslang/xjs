; // File: 005dc7dff71d4b97.js
[ 1 ]
; // File: 006949a4f1471866.js
if (a) {
}

; // File: 00b851b06af02cc0.js
a.b('c').
    d('e',
        /*@ngInject*/
        function(f) {
            return f;
        }).
    g('h',
        /*@ngInject*/
        function(i) {
            return i;
        })
; // File: 00bd68a9d0203f10.js
if (a) {
    b();
    c();
    d();
} else {
    e();
    f();
    g();
}
; // File: 00c79d09c52df3ec.js
for([a,b[a],{c,d=e,[f]:[g,h().a,(1).i,...j[2]]}] in 3);

; // File: 0140c25a4177e5f7.module.js
export default (1 + 2);

; // File: 01533b37d1d9ede8.js
1 - 2

; // File: 017a45a1919f4006.js
a: while (true) { continue a }
; // File: 01f27ee3c1bb68e5.js
a >>= 1
; // File: 01fd8e8a0a42307b.js
(function* () { yield *a })
; // File: 02028e3b961bfee0.js
({ get: 1 })
; // File: 0228be549a7706e7.js
(class {prototype() {}})
; // File: 023e4178e1ad1a82.module.js
import * as a from "a"
; // File: 024f7b95336f7fad.js
a = (b, c)

; // File: 0262c247b28885e2.js
({ if: 1 })
; // File: 0266b93cf3014995.js
!a
; // File: 027abe815032df72.js
/p/;
; // File: 028846a58c67687f.js
{;}
a();
{};
{
    {};
};
b();
{}
; // File: 02b924339f85fe00.js
for (var {a, b} in c);

; // File: 02cf1a37af2403fe.js
a: for (;;) break a
; // File: 02dad3c9ec38d3c7.js
({a: b = c = 1} = 2)
; // File: 0339fa95c78c11bd.js
(a, ...[]) => 1

; // File: 034ded949b5c2fa3.js
function *a(){yield ++a;}
; // File: 03608b6e222ae700.js
a && (() => {});
; // File: 0371eb8b8c28569d.js
`$$$`
; // File: 037ecd1db38c230c.module.js
export const a = 1;

; // File: 03d1cf071a76d061.js
[...a[1]] = 2;

; // File: 040001f3b0eb3bde.js
function eval() { function a() { "use strict" } }
; // File: 0426f15dac46e92d.js
(function () {
    var a = {
        '1e2000': 1
    };
}());

; // File: 0453974dd98e662d.js
a = { get: 1 }
; // File: 0458e0c30e8e6fb0.module.js
import a, * as b from "foo";

; // File: 0466764f0fb9af62.js
function* a(){yield}
; // File: 046a0bb70d03d0cc.js
T‌
; // File: 046b1012ef9b0e26.js
/[a-c]/i
; // File: 04b26d042948d474.js
(function() {
    a(), 1, 2;
}());

; // File: 04df09188055748a.js
if (a) {
    b();
} else if (c) {
    d();
} else if (e) {
    f();
}

if (a) {
    b();
} else if (c) {
    d();
} else if (e) {
    f();
} else {
    g();
}
; // File: 0507b18e39d58a9f.js
var a = class b extends 1{}
; // File: 05089e6cc717523e.js
(function () {
}(1,2,3))

; // File: 051696d4c46ad99b.js
a: while (true) { break a }
; // File: 053480e541f54faf.js
/*a
b*/ 1
; // File: 053c0475e49bd36b.js
for (const a in b) c(a);

; // File: 05448bc107f9b759.js
class a {b(){};c(){};}
; // File: 054620d2d7fbe8fb.js
function a() {
    if (false) {
        // because test is not referenced
        var a = 1;
    }
}

; // File: 058c33e92f0d37a5.js
var a;
if (b()) {
    new a(1);
} else {
    a(2);
}
; // File: 059b850298ae3352.js
``
; // File: 05b9c5f007cbaa56.js
a >= b
; // File: 05d268921a1f6899.js
do continue; while (true)  // should be empty statement

; // File: 05d93894463f57ca.js
// ContinueStatement should be removed.
// And label is not used, then label also should be removed.
a: for(;;) continue a;

; // File: 05fcc31bfd8d3e60.js
function a({ b, c }){}
; // File: 066b76285ce79182.js
class a { set b(c) {} get b() {} }
; // File: 066e2ec2de8a7c6e.js
c: {
    a();
    switch (1) {
      case 2:
        b();
        if (a) break c;
        d();
      case 3+4:
        e();
        break;
      default:
        f();
    }
}
; // File: 0671ec3d0b8ded79.js
(1, a)();
(2, b.a)();
; // File: 068fd501eb381dba.js
(a) => b;
; // File: 06981f39d0844079.js
function *a(){yield/=3/}
; // File: 06c7efc128ce74a0.js
(function(){ a() })();
; // File: 06d84c003dc8a3af.js
"use strict";a={b:1,b:2}
; // File: 06f7278423cef571.js
({2e308:1})
; // File: 070d82d1b3b3a975.js
(function () {
    var a;
    eval('a');
    function b() {
        a = a += 1;  // eval makes dynamic
    }
}());

; // File: 071f05b40ea0163f.js
new a("aa, [bb]", 'return aa;');
new a("aa, {bb}", 'return aa;');
new a("[[aa]], [{bb}]", 'return aa;');
; // File: 075c7204d0b0af60.js
({ get a() {} })
; // File: 079b7b699d0cacab.js
(function () {
    arguments[1] = 2;
    var a =3;  // should not hoist to parameter
}());

; // File: 07a74deab99e85eb.js
var a = class extends b {}
; // File: 07bce073a241288b.js
``;
`xx\`x`;
`${ a + 1 }`;
` foo ${ b + `baz ${ c }` }`;
; // File: 07cfd31162dc117a.js
let {a,b=1,c:d,e:f=2,[g]:[h]}=3

; // File: 07d4bedb35fb60b6.js
(1, a.a)();

; // File: 0813adc754c82a98.js
class a {'constructor'() {}}
; // File: 0821d3a84023aca2.js
var [{a},b] = c;
; // File: 0827a8316cca777a.js
(class {get a() {}})
; // File: 0860caf88460e363.js
((a,a),(a,a))
; // File: 0889a34434e586e9.js
1;
; // File: 08a39e4289b0c3f3.js
T‍
; // File: 08ba81b9af0132ea.js
(function () {
    for (var a; a < 1; ++a);
}());

; // File: 091d00847cbf8a9d.js
var {a, b} = {a:1, b:2};
; // File: 09245ed873c9e7ea.js
a(...b)
; // File: 0986e63317738f46.js
var a, b;
if (a && !(a + "1") && b) { // 1
    var c;
    d();
} else {
    e();
}

if (a || !!(a + "1") || b) { // 2
    d();
} else {
    var f;
    e();
}
; // File: 098e1fe1335e222b.js
[a, ...{0: b}] = 1
; // File: 09be3a3198b40536.js
function a([]) {}

; // File: 09c1c4b95bf0df77.js
({ __proto__, __proto__: 1 })
; // File: 0a068bc70fe14c94.js
("a")
; // File: 0a2fc93b6a63bbd3.js
[...a] = b
; // File: 0a38bb9fff27bc21.js
(function () {
    var a;
    function b() {
        a = a += 1;
    }
}());

; // File: 0a616ee6dd067bc6.js
// global, do not optimize
(function () {
  a("b");
}());

; // File: 0a9e4cbb36d95f7c.js
function a() {
    var b = 1;
    c();
    {
        c();
        c();
    }
}

; // File: 0aa6aab640155051.js
class a extends b { constructor() { () => { super(); } } }
; // File: 0aa9242278e1393b.js
class a extends b {
    c() {
        return super.d
    }
}

; // File: 0aeb95f62766e684.js
"\u0061"
; // File: 0b1fc7208759253b.js
'\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b\x2c\x2d\x2e\x2f\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3a\x3b\x3c\x3d\x3e\x3f\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4a\x4b\x4c\x4d\x4e\x4f\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5a\x5b\x5c\x5d\x5e\x5f\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6a\x6b\x6c\x6d\x6e\x6f\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7a\x7b\x7c\x7d\x7e\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe' + a

; // File: 0b2804600405dbf6.js
var a, b, c, d;
a = (b(), c(), d()) ? 1 : 2;
; // File: 0b281915a3227177.js
"Hello\012World"
; // File: 0b4d61559ccce0f9.js
({a:(b.c)} = 1)

; // File: 0b4e932ec15cdae4.js
for(var a = 1;;) { let a; }
; // File: 0b50309b4112013e.js
var a; (a);

; // File: 0b5f023129f23abf.js
a = [ 1 ]
; // File: 0b6dfcd5427a43a6.js
0008
; // File: 0b881b80b7220fad.js
a = { b(...c) { } }
; // File: 0ba326a76aa2a0ae.js
let {a:{}} = 1

; // File: 0bbda5d7d8ae8990.js
function a() { /* infinite */ while (true) { } /* bar */ var b; }
; // File: 0bef54f61acccbe5.js
(function() {
    return 1;
    var a = 2;
}());

; // File: 0c44152a1a3e2f90.js
- (1 - 2 - 3)

; // File: 0c7719169ed21a87.js
class a { b() { new super.c(); } }
; // File: 0c8a07486c1ff18e.js
(a) => b
; // File: 0cd7b76d7d1431d6.js
function a(b, c) { return b-- >= c; }
; // File: 0ce4fa8fdf700065.js
new (function () {
    var a = 1;
});

; // File: 0cf1df0ef867a7f4.js
({yield})
; // File: 0d137e8a97ffe083.js
function *a(){({get b(){yield}})}

; // File: 0d7e3e1647af9ba6.js
do {} while (true)
; // File: 0d9f26fe4d91ad07.js
[{a=1},{a=2}] = 3
; // File: 0da4b57d03d33129.js
typeof (1, a, 2)

; // File: 0da6496ed75822b1.module.js
export function a () {} false

; // File: 0de4ef1344cbb907.js
function a(){return {} / 1}
; // File: 0de707242475664c.js
for (;a();) {
    b();
    c();
    if (d()) break;
    e();
    f();
}
; // File: 0de805d0c921e235.js
a <<= 1
; // File: 0def12c63f682470.js
[] = 1
; // File: 0e22e969622bf137.js
({ *a() {} })

; // File: 0e3ca454ddfb4729.js
for (a[b in c] in d);
; // File: 0eb53d0e06cd5417.js
function *a() { yield b=c, yield* d=e, f }

; // File: 0f18951fd55b8c07.js
switch (a) {
  case b: {
    c;
  }
}

; // File: 0f59aedfe2c7682c.js
({ null: 1 })
; // File: 0f630e67e4542867.js
var [a]=[1];
; // File: 0f809258920b3469.js
c: {
    a();
    switch (1) {
      case 2:
        b();
        if (a) break c;
        for (var b = 3; b < 4; b++) {
            if (b > 5) break; // this break refers to the for, not to the switch; thus it
                              // shouldn't ruin our optimization
            d.e(b);
        }
        f();
      case 6+7:
        g();
        break;
      default:
        h();
    }
}
; // File: 0f88c334715d2489.js
function *a(){yield void 1}

; // File: 0f9f10c894a7d811.js
(0o0)
; // File: 0fa2102f53acd283.js
(a)++
; // File: 0fc7d5705a324efb.js
(function () {
    a = 1;
    for (b = 2;;);
}());

; // File: 0fe1f55610641156.js
a["b"] = "c";
; // File: 0fe2654034a20f6b.module.js
import {a, b} from "foo";

; // File: 0ffdc03e2ffcb5dc.js
void /test/

; // File: 10786cdac00d0c02.module.js
import { a, b } from "c"
; // File: 10857a84ed2962f1.js
class a { get b() {} set b(c) {}}
; // File: 1093d98f5fc0758d.js
(function a({ b: { c, a }, d: [e, f] }, ...[b, d, g]){})
; // File: 10d6486502949e74.js
({ __proto__: null, get __proto__(){} })

; // File: 10f0ef998c05c611.js
do {} while (false) a();
; // File: 10fda5cd119b39a5.js
[...eval] = a
; // File: 110fa1efdd0868b8.js
"Hello\
world"
; // File: 111668493e3e0823.js
/[a-z]/i
; // File: 1145e94ad27e7ba6.js
function a() {
    b();
}
if (a() || true) {
    c();
}
; // File: 116cacc3c80a5a3e.js
do {
    a();
} while (false)b()
; // File: 119e9dce4feae643.js
var a = {
    b: null,
    set c(d) {
    },
    get c() {
        return this.b;
    }
}
; // File: 11a021c9efe0e432.js
--eval
; // File: 121491a690a13543.js
a(() => {})
; // File: 1223609b0f7a2129.js
!function(){a()}(),!function(){b()}(),c()+1

; // File: 123bfcc3f6cf379f.js
function a(b, ...c) { }
; // File: 123f89c06747ced2.js
`\u{000042}\u0042\x42u0\A`
; // File: 124490e0f2dbbac7.js
a % b
; // File: 126a6455f0f721fe.js
({a(b){}})
; // File: 12752899d5c5eb00.js
`\n\r\b\v\t\f\
\
`
; // File: 129c95a57d234b7b.js
a * b % c
; // File: 12d4b327a5e20850.js
if (a) { b() // Some comment
 }
; // File: 12d5bedf1812952a.js
class a { static get b() {} static get c() {} }
; // File: 12e59b6d403833ae.js
'use strict';
a();

; // File: 12ea3bf0653f8409.js
a = {"__proto__": 1 }
; // File: 12edb6ae55d95b59.js
{ let a = 1, b = 2, c = 3 }
; // File: 12ef713cb7737bdd.js
class a {b(){};c(){}}
; // File: 13045bfdda0434e0.js
({set a(b) {}})
; // File: 1325417193f50cc3.js
if (a) {
} else {
    b();
}

; // File: 1426cb41eb6d515f.js
while (true) { break /* Multiline
Comment */a; }
; // File: 14360fa75e6ae228.js
function a({b} = {b: 1}) {}
; // File: 1450a897a4ba83a7.js
'use strict';
var a = {
    '0': 'b'
};

; // File: 14551b80fa8a0ce1.module.js
export function a () {}

; // File: 14a62ce75845f5dd.js
[a, b] = [b, a]
; // File: 14bb381a17b683e3.js
for (var [a, b] in c);

; // File: 14df05a5ad02af18.js
class a { ; }
; // File: 14f95b3c9a9e7480.js
/(()(?:\2)((\4)))/;
; // File: 14fb22cf10e20236.js
typeof (1, a)  // Don't transform to 0,typeof ident

; // File: 1530c2c5484d867f.js
(a) => 1
; // File: 153688477d7e69ba.js
// TODO(Constellation):
// This transformation sometimes make script bigger size.
// So we should handle it in post processing pass.
(function () {
    while (!a || !b()) {
        c();
    }
}());

; // File: 153bd6819f5fa69b.js
({ yield() {} })

; // File: 159c17331c90a465.module.js
export default class a {}
; // File: 15a12468ff312d51.js
6.02214179e+23
; // File: 15d072c60817cdca.js
new a(...b, c, d);

; // File: 15d9592709b947a0.js
function a(...[ b, c ]){}
; // File: 15dfd62aa10c8b18.js
function a() {'use strict'; 0O0; }
; // File: 1623cc76ec1fb540.js
(function () {
    if (a != true) {
        b();
    }
    if (a != false) {
        b();
    }
}());

; // File: 163e6a68a09abaed.js
function *a(){yield null}
; // File: 166431dca77feba6.module.js
import a, { b, d as c } from "a"
; // File: 16b9227a4a41bc7e.js
function a() {
    // If foo is null or undefined, this should be an exception
    var {a,b} = c;
}
; // File: 16c7073c546fdd58.js
// prevent optimization because of this.constructor.arguments access
new function () {
    var a = 1;
    this.arguments;
};

; // File: 16d0c12aad83f9b3.js
function a() {
}
function b() {
    return c;
}
function d() {
    return void 1;
}
function e() {
    return void 2;
}
function f() {
    return;
}
function g(h, i) {
    j.k(h, i);
    l(h);
    return;
}
function m(h, i) {
    j.k(h, i);
    if (h) {
        n(i);
        l(h);
        return h + i;
    }
    return c;
}
function o(h, i) {
    j.k(h, i);
    if (h) {
        n(i);
        l(h);
        return void 3;
    }
    return h + i;
}
function p(h, i) {
    n(h);
    q(i);
    return void 4;
}
function r(h, i) {
    n(h);
    q(i);
    return c;
}
function s() {
    return false;
}
function t() {
    return null;
}
function u() {
    return 5;
}
; // File: 1714b06e6a415766.js
a = { false: 1 }
; // File: 1717229250780255.js
function a() {
    (class b { });
    class c {};
}
; // File: 17302b9b0cab0c69.module.js
export {};1
; // File: 17326734a7bf9629.js
// adapted from http://asmjs.org/spec/latest/
function a(b, c, d) {
  "use asm";
  var e = b.f.e;
  var g = b.f.g;
  var h = new b.i(d);
  function j(k, l) {
    k = k|1;
    l = l|2;
    var m = 0.0, n = 3, o = 4;
    // asm.js forces byte addressing of the heap by requiring shifting by 3
    for (n = k << 5, o = l << 6; (n|7) < (o|8); n = (n + 9)|10) {
      m = m + +g(h[n>>11]);
    }
    return +m;
  }
  function p(k, l) {
    k = k|12;
    l = l|13;
    return +e(+j(k, l) / +((l - k)|14));
  }
  return { p: p };
}
function q(b, c, d) {
  var e = b.f.e;
  var g = b.f.g;
  var h = new b.i(d);
  function j(k, l) {
    k = k|15;
    l = l|16;
    var m = 0.0, n = 17, o = 18;
    // asm.js forces byte addressing of the heap by requiring shifting by 3
    for (n = k << 19, o = l << 20; (n|21) < (o|22); n = (n + 23)|24) {
      m = m + +g(h[n>>25]);
    }
    return +m;
  }
  function p(k, l) {
    k = k|26;
    l = l|27;
    return +e(+j(k, l) / +((l - k)|28));
  }
  return { p: p };
}
; // File: 174d05abbd69a960.js
for (var {a, b} of c);

; // File: 175a032b2252eb0d.js
(function () {
    var a;
    (1, a)();
}());

; // File: 177fef3d002eb873.js
function* a(){(class extends (yield) {});}

; // File: 17a2de2c9e102bba.js
switch (a) { case 1: /* perfect */ b() }
; // File: 17bd5dc47ec4a3ba.js
new a(b, ...c, d)
; // File: 17bd95dfa6a302f2.js
for(a.b in c);
; // File: 17cc7c10e02028be.js
if(1)/  foo/
; // File: 17d63bb0b9482189.js
var a = 'very cute';

; // File: 17d881105a9a6c85.js
for (var a in b)
  // do not optimize it
  (function () {
    c('d');
  }());

; // File: 1819ffb142e9c5ea.js
0O0
; // File: 185dc3ee443bb737.js
a[b] = b
; // File: 18cc9a6b7038070f.js
(function ( a ) { });
(function ( [ a ] ) { });
(function ( [ a, b ] ) { });
(function ( [ [ a ] ] ) { });
(function ( [ [ a, b ] ] ) { });
(function ( [ a, [ b ] ] ) { });
(function ( [ [ b ], a ] ) { });

(function ( { a } ) { });
(function ( { a, b } ) { });

(function ( [ { a } ] ) { });
(function ( [ { a, b } ] ) { });
(function ( [ a, { b } ] ) { });
(function ( [ { b }, a ] ) { });

( [ a ] ) => { };
( [ a, b ] ) => { };

( { a } ) => { };
( { a, b, c, d, e } ) => { };

( [ a ] ) => b;
( [ a, b ] ) => c;

( { a } ) => b;
( { a, b } ) => c;
; // File: 18e32b70e6a5574c.js
Infinity.a();
NaN.a();
; // File: 18f05b95a72dffa1.js
function *a(){yield/a/}
; // File: 18f731daf0845475.js
a[b, c]
; // File: 1908280b73954ef7.js
/foo\/bar/
; // File: 1938db3bb862ded1.js
function *a({yield: b}){}

; // File: 1972b64c4704a1eb.js
((a))()
; // File: 19d1d07fe88ec849.js
({a:yield} = 1);
; // File: 19ffea7e9e887e08.js
({ set null(a) { a } })
; // File: 1a0dac12dbd33ef6.js
a = {}
; // File: 1a1c717109ab67e1.js
var a = /[\uD834\uDF06-\uD834\uDF08a-z]/u

; // File: 1a7800a74a866638.js
var eval = 1, arguments = 2
; // File: 1b0c0fc32b9e5e35.js
// reported from issue #60
void function () {
  var a;  // this foo should be dropped
  a = function () {  // this should be transformed to non-assignment expression
    return 1;
  };
}.b(this);

; // File: 1b542dd79e4444c7.module.js
import {} from "foo";

; // File: 1b6e33ab982844af.js
class a {;b(){};c(){};}
; // File: 1b87f88ae8ea1cb1.js
switch(a){case 1:default:}
; // File: 1b884461ff1acfc6.js
(function({a}){})
; // File: 1ba78d63a36ea567.js
for(let a in a);
; // File: 1bbe65871120530b.js
({ a: b } = c)
; // File: 1c055d256ec34f17.js
(a, b, c, d) + e;

; // File: 1c1e2a43fe5515b6.js
if (1) function a(){} else function b(){}
; // File: 1c2f680b78692645.js
for(var a = 1, b = 2;;);
; // File: 1c6424d9a7209f81.js
a > b
; // File: 1c6c67fcd71f2d08.js
({set a(b=1){}})
; // File: 1c7e1e347f726166.js
{ a }
; // File: 1ca991b39b6e7754.js
void 'test string'

; // File: 1cb2c267c552028f.js
b: while (1) { continue 
 a; }
; // File: 1cdce2d337e64b4f.js
b: while (1) { continue /* */ a; }
; // File: 1ce4afd9b35e3312.js
[a.b=b] = c
; // File: 1d1037fcfa0c7958.js
[ 1, 2, 3, ]
; // File: 1d1ac5ee0d1a9bd4.js
let [[]]=1

; // File: 1d3dd296a717e478.js
(eval = 1) => 2
; // File: 1db0d98ff1726af8.js
class a {static ["prototype"](){}}
; // File: 1db76a05c7b9a090.js
a >>> b
; // File: 1de765c987733026.js
/[\s-\w]/;
; // File: 1e3f57c4ec83f5bc.js
a(
	b(c + 'd'),
	b('d' + c)
);
; // File: 1e61843633dcb483.js
(function () { 'use\nstrict'; with (a); }())
; // File: 1ea254c74f1071de.js
({*a(){}})
; // File: 1ef3bdd7e919cca8.js
function* a() { b.c(yield); }
; // File: 1efde9ddd9d6e6ce.module.js
(function* () { yield
a })
; // File: 1f039e0eeb1bc271.js
a = { get b () {} } 
; // File: 1f3808cbdfab97e4.js
a: c: b: while (true) { continue a; }
; // File: 1f5de1d7092dcd82.js
var a
;
; // File: 1f89cd96db326f7a.js
new a(...b, c, ...d)
; // File: 1f988cc22167927b.js
(function() {
    1, 2, 3;
}());

; // File: 1fbf374c8a04fb23.js
0o10
; // File: 1fc4349ef394b505.js
while (true) { continue; }
; // File: 1fd743f03945fd05.js
a(b, c)
; // File: 2010526ea64db82e.js
class a extends b { constructor() { super() } }
; // File: 20644d335e3cd008.js
"Hello\312World"
; // File: 206ebb4e67a6daa9.js
var a = /[\u{0000000000000061}-\u{7A}]/u
; // File: 2072cb8131a4ae2b.js
void ((a) ? b : 1);

; // File: 209fc98bea7b9d67.js
if (true) a()
; else;
; // File: 20aca21e32bf7772.js
('\u{10FFFF}')
; // File: 20b873ad024b210f.js
(class {})
; // File: 20f9bec9f3215688.js
/**/ function a() {/**/function b() {}}

; // File: 2100bec1b92b51ae.js
(function () {
    if (false) {
        var a = 1;
    }
    b();
}());

; // File: 212d2ca66d97a90f.js
const {a:b} = {}
; // File: 213e3455c8f8ceb2.js
(class a {})
; // File: 2160fc99c3589501.js
/*Venus*/ debugger; // Mars
; // File: 2179895ec5cc6276.js
(...a) => 1

; // File: 218ca74570bf06b5.module.js
export {a as default, b} from "foo";

; // File: 218e751b8b453b9b.js
if (a) {
  // optimize it
  (function () {
    b('c');
  }());
  try {
    b("d");
  } catch (e) {
  }
}

; // File: 21ebb8746371268b.module.js
import {b as a, c} from "foo";

; // File: 21f1173fff072ee5.js
let a
; // File: 21f5ce68788d4ffa.js
while (a < 1) { a++; b--; }
; // File: 2207b24e625f30db.js
(function () {
    var a = {};
    a.b = (c(), 1);  // ok
}());

; // File: 22119a1d30256255.js
a | b
; // File: 224d4dca3d98b618.js
 /****/
; // File: 227118dffd2c9935.js
class a extends b { static get c() {}}
; // File: 22b24d1deb35baf3.js
(a => a)
; // File: 22dc0bb1d4e8d89f.js
function a({yield: b}){}

; // File: 22eba6e3841edeec.js
switch (a) {}
; // File: 230da70c908c1859.js
'use strict';
var a = {
    delete: 1
};

; // File: 23869c020fc2cb0f.js
({ "__proto__": null, __proto__(){}, })

; // File: 23d6a92eed7f18fa.js
for (a.in in a);
; // File: 2418fddf06e515f8.js
a in b
; // File: 242ede66951e11b1.js
// optimize this
(function () {
  a('b');
}());
try {
} catch (c) {
}

; // File: 24557730b5076325.js
(function() {
    a(), 1, b();
}());

; // File: 247a3a57e8176ebd.js
({a, a: 1, b})
; // File: 24e299720285b6c1.js
var a = /foo\/bar/
; // File: 24fa28a37061a18f.js
"use strict"; ({ yield() {} })

; // File: 250ced8c8e83b389.js
var [,a] = 1;
; // File: 25296359c69440e8.js
for(let a in [1,2]) 3
; // File: 25542e65ad9d2bf1.js
(function () {
    var a;
    b(typeof a === 'c');
}());

; // File: 2565ae4b2f2956b0.module.js
import * as a from "foo";

; // File: 257f15ea5c44a423.js
for(let a = 1, b = 2;;);
; // File: 25824f6a683e7467.module.js
export class a {}
;
; // File: 25fd48ccc3bef96a.js
(function () {
    a(typeof b !== 'c');
}());

; // File: 2619be6c7f521c49.js
function a() {
    b = {
        c() { return 1; }
    }
}
; // File: 264266c68369c672.js
function a() { var b; if (b = 'b') { return b; } else { return b; } }; a();
; // File: 26974bc54e93b191.js
a[b]
; // File: 26998ded3750f7d8.js
// To avoid JSC bug, we don't distinguish FunctionExpression name scope and it's function scope
(function a() {
    var b = 1;  // Don't rename this variable to a name that is the same to function's name
    c(b);
}());

; // File: 26a4b2dddf53ab39.js
/\uD834\uDF06\u{1d306}/u
; // File: 26aa785e12e00fb1.js
(function () {
    if (a) return;
    else return;
}());

; // File: 26aa8b685715d445.module.js
export const a = { }
; // File: 26b946d7cc01c226.js
var yield;
; // File: 26edf4bcd3ed9e74.js
(function () {
    void ((a) ? 1 : b);
}());

; // File: 26f27d747e98d3eb.js
function a() {
    if (b) {
        let c;
        let d;
        var e;
        var f;
    }
}
; // File: 26f632a0a4d60150.js
switch (1) {
  case 2: a();
  case 3+4: b(); break;
  case 5+6+7: c();
}
; // File: 27409f5b7b692b24.js
(function({a = 1}){})
; // File: 2754a9872f3512ed.js
function a() { "use strict" + 1; }
; // File: 27ac24465c731ff9.js
var a = /[x-z]/i
; // File: 27b5d00cc75de02f.js
function a(b) {
    if (b) {
        return;
    } else {
        return 1;
    }
}
; // File: 27ca96102da82628.js
"Hello\712World"
; // File: 27ed2f0fdb7f53f6.module.js
export {a, b} from "foo";

; // File: 284ca169f09605be.js
[a, ...b] = c
; // File: 285648c16156804f.js
(function(){ return/* Multiline
Comment */a; })
; // File: 286876d6fdab22d7.js
__proto__: while (true) { continue __proto__; }
; // File: 28a54e6410ad3f19.js
// Surpress reducing because of alternate
for (;;) {
    if (a) {
        if (b) {
            continue;
        } else {
            ;
        }
    } else {
        ;
    }
}

; // File: 290fdc5a2f826ead.js
(function a(b, c) { })
; // File: 2935f62bfe48ca1b.js
function a() {
    return (a, void 1);
}

; // File: 2976a1598d3a75e1.js
let {a} = {}
; // File: 29e41f46ede71f11.js
(yield) => 1;

; // File: 29ef8a7a1cbfda7f.js
a.in
{}
/foo/
; // File: 2a11bb318142547e.js
if (a()) {
    if (b()) {
        c();
    } else {
        d();
    }
} else {
    d();
}

if (a()) {
    if (b()) {
        c();
    }
}
; // File: 2a327fdbcc6cb870.js
a({ ["b" + "b"]: 1 });
; // File: 2a7d131074016ba6.js
'use strict';
var a = function(b) {
    b();
    a()
}

; // File: 2aa1db78027ba395.js
() => (a) = 1
; // File: 2ab9ca1f6a30c203.js
// 1.
if (a) {
    {{{}}}
    if (b) { c(); }
    {{}}
} else {
    d();
}

// 2.
if (a) {
    for (var e = 1; e < 2; ++e)
        if (b) c();
} else {
    d();
}
; // File: 2acc83b037420689.js
class a extends b {
    constructor() {
        () => super()
    }
}

; // File: 2afc5d4b75dbf12d.js
(function () {
    void 1;
}());

; // File: 2b0727c871857af5.js
let a, b;
; // File: 2b1f4f042cff07a3.js
(function () {
    if (true) {
        var a = 1;
    }
}());

; // File: 2b393e093a0e2fb3.module.js
export {a as default} from "foo";

; // File: 2b478bb5ceb2e18b.js
a++
; // File: 2b83dea123ed2e2e.js
({ a: [a, b] }, ...c) => {}
; // File: 2b9d4a632590814a.js
function a() {'use strict'; 0o0; }
; // File: 2ba11d8ca169ab6c.js
b: while (a) break b;
c: while (a) break;

; // File: 2bdb271c1ff34f35.js
for(let a = 1;;);
; // File: 2c16af589c5c8535.js
class a { static *b(c) { yield c; }}
; // File: 2c4b264884006a8e.js
(function() {
    throw 'a';
    with (b);  // This should be removed.
}());

; // File: 2c7a69627f1d8062.js
while (a) {
    {
        b();
        b();
    }
}

; // File: 2c7e2fecbc1cb477.js
(function () {
  for (!!!a&&a();!!!b&&a();!!!b&&a()) {
  }
}());

; // File: 2ccf4707fe3749ff.js
((1, 2), 3);

; // File: 2cda5eb51a2d97e7.js
class a extends b { constructor() { super.c } }
; // File: 2cdf798a24c241e3.js
a * b * c
; // File: 2d10fed2af94fbd1.js
(a)=>{'use strict';}
; // File: 2d1ecf6fb0d1afe2.js
({0: a, 1: a} = 1)
; // File: 2d3273e0386e9cb8.js
function* a() {}
; // File: 2d614e07c62fc32d.js
while (true) { break }
; // File: 2db5219f0ac5dd71.js
a(
    `<span>${b}</span>`,
    `<a href="${c}">${d}</a>`
);
; // File: 2dc1c08d0bff6eba.js
([a,...b])=>1;

; // File: 2dd810da4984502b.js
1 /* the * answer */
; // File: 2df813ffa8a0a9e1.js
function a() { return "<!--HTML-->comment in<!--string literal-->"; }
; // File: 2e371094f1b1ac51.js
('\1111')
; // File: 2e5fbf7b1685fa1b.js
(function () {
    var a = 1;  // should not hoist this
    arguments[2] = 3;
    (function () {
        eval('');
    }());
}());

; // File: 2e7336dc8eba87ef.js
a.b
; // File: 2e75e3bd39e6df05.js
a = "b".c;
a = ("b" + "d")["e" + "f"];
a = g.c;
a = ("b" + g).c;
; // File: 2e7f443b2f555bc5.js
/**/ function a() {function b() {}}

; // File: 2e8a88da875f40c7.js
delete /test/

; // File: 2f6d8a2215407ae3.js
function a() {
    for (var b = 1, c = 2; b < 3; ++b) {
    }
}

; // File: 2f84859abd5a242c.js
(function () {
  do {
    a();
  } while (false);
}());

; // File: 300a638d978d0f2c.js
T‍ = []
; // File: 30142c5b79e4eea9.js
// line comment
1
; // File: 3085290028dd33e1.js
function a(b, c) {
    function d() {
        e();
    }
    return b + c;
}
; // File: 3097f73926c93640.js
(function(){ return })
; // File: 3098b57020860587.js
(function () {
    var a = 1;  // should not hoist to parameter
    eval('');
}());

; // File: 30aee1020fc69090.js
(class {set a(b) {'use strict';}})
; // File: 30c2911c05100e92.js
// line comment
; // File: 31232b72db0fd24f.js
a(
	b(c, c),
	d(c, c),
	e(c, c)
);
; // File: 312f85fecc352681.js
(function a() { b; c() });
; // File: 315692af7fe2aad3.js
var a = {};
a.b = 1;
a.c = 2;
d.e(a.c);
; // File: 3156a92ca5319b8b.js
throw 'a';
b();

; // File: 317532451c2ce8ff.js
'use strict';var a = function(){}(b())

; // File: 318c169a25ee42c5.js
new a(....5);

; // File: 31ad88cae27258b7.js
var a = /[\]/]/
; // File: 31cca30ad2bf696d.js
(function () {
    (function () {
    }());
}());

; // File: 323783be9a53a31e.js
08
; // File: 329bc0e532da6227.js
````
; // File: 32b4780ad9c4292a.js
a: 1; a: 2;
; // File: 32b635a9667a9fb1.js
/* header */ (function(){ var a = 1; }).b(this)
; // File: 32b6854d07aefbda.js
0B0
; // File: 32efa0efd255748a.js
class a { static *[b]() {} }

; // File: 32f782a4b16306aa.js
a << b << c
; // File: 3315c524a740fe55.js
('\\\'')
; // File: 332f0bc46d28db25.js
({a}) => 1
; // File: 3348741c8bdd4f3c.js
for (const a of b) c(a);
; // File: 33720cee3dabde0d.js
a || b && c
; // File: 338762eadb13a2f0.js
a ^= 1
; // File: 341bc3f1b434f6d1.js
``````
; // File: 345713fe7f52524a.js
({a,b,} = 1)
; // File: 347e9a5443e4cd3c.js
{ let a }
; // File: 34d5455824302935.js
function *a() { yield* 1; }

; // File: 350a7ec7041c079f.js
a = { set: 1 }
; // File: 3514acf61732f662.js
09
; // File: 35bf182594dc08ac.js
function a(b, c, d, e, f) {
    return b + c;
}
; // File: 35e730121a5e6326.js
for (;;) if (a()) b(); else break;
; // File: 35eb2e229858a6c7.js
a => "b"
; // File: 3610e596404818d6.js
with(1);
; // File: 36224cf8215ad8e4.js
(function (a, ...b){});
(function (...c){});
; // File: 366585381e4610b4.js
a+(b(), c(), d())  // do not transform because of global getter

; // File: 369b56fe359d52fc.js
a()``
; // File: 369fd0a1e40030d8.js
class a extends b {
    c() {
        return super[1]
    }
}

; // File: 36a9e7f1c95b82ff.js
 
; // File: 36fb3e9c8cedf764.js
class a {set(b) {};}
; // File: 36ff120198eea816.js
var { yield: a } = b;

; // File: 370a2bd1387fd440.js
for ({a, b} of c);

; // File: 372097a44c33daf2.js
eval => 1
; // File: 373e35460ecaccc6.js
1.492417830e-10
; // File: 3793ec99f844de1c.js
'use strict';
var a = {
    'b': 1
};

; // File: 37ac3bcee6fa89f9.js
0b10
; // File: 37c0f5275362d1c9.js
"\u{00000000034}"
; // File: 37d26e3bec6d9a0f.js
switch (a) {
default:
  // do not optimize it
  (function () {
    b('c');
  }());
}

; // File: 37e4a6eca1ece7e5.js
(function a([ b, c ]){})
; // File: 37e845e0d8283fb3.js
(function () {
    null!=(a?void 1:b)
}());

; // File: 380e999de8f31c7d.js
function a([ b, c ]){}
; // File: 3812dc38bcdc97db.js
b: for (var a = 1; a < 2; ++a) {
    if (a < 3) continue b;
    c.d(a);
}
; // File: 38284ea2d9914d86.js
({ get a() { }, get a() { } })
; // File: 3852fb3ffb8fd8d5.js
if (true) a(); else;
; // File: 38594572e7bb32f4.js
(function () {
  function arguments() {
    a(arguments);
  }
  a(arguments);
}());

; // File: 38befc89fcf92e25.js
+{} / 1
; // File: 38c0e030050edb57.js
(function () {
    -1;
}());

; // File: 38e0b9de817f645c.js
null

; // File: 38fefd37caf6f8bb.js
a | b | c
; // File: 3990bb94b19b1071.js
('\1')
; // File: 39bd53b0c3dcd639.js
let {a} = b
; // File: 3a1ccd915e97ed68.js
//
1
; // File: 3a1f039e533d1543.js
'use strict';
a.static();
; // File: 3a50539d66e7fb07.js
{
    var a = 1;
    b();
    {
        b();
        b();
    }
}

; // File: 3a707c56867f396c.js
d: while (a) {
    b();
    c();
    break d;
    e();
    f();
}
; // File: 3aa600e48cbd8a5c.js
class a {static [b](){};}
; // File: 3ae4f46daa688c58.js
'use strict'; ('\0')
; // File: 3b1ab093f7ebeb51.module.js
export var a;

; // File: 3b1fca65828182ab.js
/(?=.)*/;
; // File: 3b36d546985cd9cb.js
"Hello\
world"
; // File: 3b57183c81070eec.js
a = { set null(b) { c = b } }
; // File: 3b5d1fb0e093dab8.js
(a) => ((b, c) => (a, b, c))
; // File: 3b9779d2e19376a1.js
0o2
; // File: 3b9e8797aacce77f.js
1
;
; // File: 3bac973df7480fe9.js
(class {3() {}})
; // File: 3bc21b350f65c8f2.js
if (a) { b() /* Some comment */ }
; // File: 3c0c251ad455218d.js
do a()
;while (true)
; // File: 3c1e2ada0ac2b8e3.js
(function*(a, b, c) { yield* a; })

; // File: 3c6b557b1aa9cc05.js
for (var a = 1; a < 2; ++a) {
    if (a < 3) continue;
    b.c(a);
}
; // File: 3c895971bd50ea01.js
for({a=1} in b);
; // File: 3caf07d66e4f7b5a.js
({ __proto__() { return 1 }, __proto__: 2 })
; // File: 3cbf0138d2dc0686.js
(function() {
  // https://github.com/Constellation/esmangle/issues/65
  var a = 1;
  var b = 2;
  var c = 3;
  var d = [].e.f(arguments);
  return [a, b, c, g];
}());

; // File: 3cf11f8790169c3f.js
a !== b
; // File: 3cf53efb53099596.js
for (;a();) {
    if (b()) c();
    else break;
    d();
    e();
}
; // File: 3d0c4eda96e0412b.js
a ^ b
; // File: 3d137e7b0cb6c8bc.js
(function() {
    if (a) {
        b();
    } else {
        return 1;
        c();
    }
    return 2;
}());

; // File: 3d2ab39608730a47.js
'use strict'; ('\0x')
; // File: 3d3ddc63a85b13a0.js
let {a,} = 1

; // File: 3d9c76216b0a9d4b.js
for(a; a < 1; a++) b(a);
; // File: 3dabeca76119d501.js
try {} catch (a) { if(1) function a(){} }
; // File: 3df03e7e138b7760.js
function a() { new["b"]; }
; // File: 3e1a6f702041b599.js
b: while (1) { continue /* */ a; }
; // File: 3e3a99768a4a1502.js
a => ({ b: 1 })
; // File: 3e48826018d23c85.js
('\5a')
; // File: 3e665d875e0049a3.js
function a() { b(); }
; // File: 3e69c5cc1a7ac103.js
try {} catch ([a, ...b]) {}
; // File: 3ea15e86885d3c1a.js
// ContinueStatement should not be removed.
a: while (true) while (true) continue a;

; // File: 3eac36e29398cdc5.js
try {} catch ([a]) {}
; // File: 3eb2c2bf585c0916.js
void ('a' + 'a')

; // File: 3ec1e9982b5f4a45.js
a && b ? 1 : 2
; // File: 3ee117e37bd3bcea.js
1e100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000

; // File: 3f34ca3add7bcb9b.js
let [{a}] = 1

; // File: 3f39e406df3080dc.js
a(...b, c, d);

; // File: 3f46ee1db509d55d.js
(1, a['a'])()

; // File: 3f6fd744861ee7c3.js
arguments--
; // File: 3f8b15109761ea65.js
var a = /[a-z]/u
; // File: 3f9b0dd207c09990.js
switch (1+2) {
  case 3: a(); break;
  case 4+5: b(); break;
  case 6+7+8: c(); break;
}
; // File: 3fb07536eb5aea8d.js
"Hello\412World"
; // File: 3ff52d86c77678bd.js
(function(...a){})
; // File: 4014ec6c7931de54.js
a = { "b": 1 }
; // File: 401544b8abe9d656.js
var [{ a, b }, ...c] = d
; // File: 40215319424a8227.js
a<!--b
; // File: 402c32920b1b9991.js
1/**/
; // File: 402e8d30db64e5af.js
1 // line comment
; // File: 403a7d28c611b71b.js
a.b(b, c)
; // File: 4040d47a3534b244.js
a = { true: 1 }
; // File: 40766161d96ac708.js
({"a": b} = 1)
; // File: 4086605956ddfcbb.js
if (!a || b());

; // File: 408971d922c72ea2.js
(function () {
  for (;false;) {
    a();
  }
  b();
}());

; // File: 409f30dc7efe75d5.js
({get __proto__() {}, set __proto__(a) {}})
; // File: 40adcdf7cfe3fa0d.js
function a() {
'use strict';
/* Comment */
}
; // File: 40b9ff090910c512.js
(function () {
    if (a == true) {
        b();
    }
    if (a == false) {
        b();
    }
}());

; // File: 414b9b02f2789648.js
var a;
// compress these
a = b ? true : false;
a = !b ? true : false;
a = b() ? true : false;

a = b ? !1 : !2;
a = !b ? !null : !3;
a = b() ? !4 : !-3.5;

if (b) {
    a = true;
} else {
    a = false;
}

if (b) {
    a = !5;
} else {
    a = !6;
}

a = b ? false : true;
a = !b ? false : true;
a = b() ? false : true;

a = b ? !7 : !8;
a = !b ? !9 : !10;
a = b() ? !11 : !12;

if (b) {
    a = false;
} else {
    a = true;
}

if (b) {
    a = !13;
} else {
    a = !14;
}

a = b ? 15 : false;
a = !b ? true : 16;
a = b ? 17 : 18;
; // File: 4180f57196d0388d.js
function *a() { b.yield(); }

; // File: 41ad2d6d8414c573.js
var [a] = []
; // File: 41b805ea7ac014e2.js
;
; // File: 41e79ea43f242aed.js
(function (a) {
    switch (a) {
    case 1:
    default:
        b("c");
    }
}());

; // File: 41fc5bd8d644937c.js
var a = function b() { c() };
; // File: 420855197cbff7ce.js
(function() {
    if (a) {
        return 1;
        b();
    } else {
        c();
    }
    return 2;
}());

; // File: 424fb5db0f6734b6.js
var [{__proto__:a, __proto__:b}] = 1;
; // File: 4263e76758123044.js
a === b
; // File: 42907dc7a3d7b79b.js
while (true) { break // Comment
a; }
; // File: 43023cd549deee77.js
try { } catch (a) { }
; // File: 43163c094787d534.js
for (const {a, b} of c);

; // File: 432639592c565344.js
new a(...b);

; // File: 433859474119631f.js
[a, {b: {c = 1}}] = d
; // File: 4369559377b6394e.js
a: while (true) break a
; // File: 438521c40cf1b08b.js
a => yield* 1
; // File: 4389b59f7805c7c7.js
(a + /* assignment */b ) * c
; // File: 43bbb253d4035175.module.js
export {a as b} from "foo";

; // File: 4412172b5dc13cd6.js
/\0/
; // File: 44136fa355b3678a.js
{}
; // File: 441a92357939904a.js
[(a = 1)]

; // File: 446ffc8afda7e47f.js
(function () {
    var a,b,c=1,d,e,f=2;
    (a,b,c)+(d,e,f);
}());

; // File: 44af28febe2288cc.js
a = { set b(c) { d = c } }
; // File: 44b0c8a5a1ecb389.js
a = b => false;
a = () => false;
; // File: 44f31660bd715f05.js
T‌ = []
; // File: 45ab34717c038020.js
for(;a;b);
; // File: 45d1662a41c9a1e9.js
for (;;) continue;  // should be empty statement

; // File: 45dd9586f26a3cf4.js
018
; // File: 45ed987996568823.js
`
`
; // File: 45ff445d87e37214.js
a = a + 1, b = b in c

; // File: 46173461e93df4c2.js
for(a,b,c;;);

; // File: 46279e885d2aa853.js
(new a).b()
; // File: 46657ec13f5857d5.js
switch (a) {
case 1:
  // optimize it
  (function () {
    b("c");
  }());
  b("d");
}

; // File: 4672c2ef688237c9.js
"\n\r\t\v\b\f\\\'\"\0"
; // File: 4694af065eecd95a.js
class a {static b(){};}
; // File: 47094fe8a994b7de.js
var a = 1<!--foo
; // File: 4724023c6bb03bac.js
while (true) { continue // Comment
a; }
; // File: 472765ae4585cf8b.js
a(...b, c, ...d)
; // File: 4743508488414d6a.js
for(a = 1;;);
; // File: 4789c3375f112cd4.js
(a) = 1
; // File: 478ede4cfe7906d5.js
(function* () { yield yield 1 })
; // File: 47ddfd79dcd20fd5.js
a = function(b = 1) {}
; // File: 47ea193a5fc3f2c7.js
(function () {
    a['b'];
}());

; // File: 47f974d6fc52e3e4.js
a = { b, c }
; // File: 47fce5046a1b2098.js
// not reduce this because of ToNumber
+ /test/

; // File: 48567b651f81277e.js
({["__proto__"]:1, ["__proto__"]:2})
; // File: 4869454dd215468e.js
try{}catch(a){}finally{}
; // File: 488ae37630cb4d83.js
() => () => 1
; // File: 488cd27c94308caa.js
function a() {
    var b = {};
    b[void 1] = (c(),d);
}

; // File: 489e6113a41ef33f.js
a & b | c
; // File: 48b43f80306f5dff.js
let {} = 1

; // File: 48b6f8ce65d3b3ee.js
!function(){a()}();b()

; // File: 48bb091783df3da9.js
function a() {
}
var b = "is a valid variable name";
b = { b: "is ok" };
c.b;
b: d()
; // File: 48bb138a6b033a34.js
(class {;;;
;a(){}})
; // File: 48f39ccbea69907a.js
arguments++
; // File: 490efeb71bdb7c3b.js
true;false

; // File: 492d3fde7a53e85a.js
switch (a) { case 1: let b; }
; // File: 4933a329b80ed6ec.js
(function () {
    function a() {
        var b = 1;
        return b;
    }
}());

; // File: 495c05812d179d67.js
/[a-z]/gimuy
; // File: 49e54e5acd18a8e1.js
1 ;
; // File: 4a0d9236bc523b77.js
('\u{0000000000F8}')
; // File: 4a479db6af79906e.js
a, b
; // File: 4a56cf2dea99fcd6.js
({ get "a"() {} })
; // File: 4a5fe6bf2362352b.js
(class a extends a {})
; // File: 4a79205bd8cd49d0.js
(function a() {
    b(typeof a() == 'c');
}());

; // File: 4a807fda565547a2.js
for (let a in b) c(a);
; // File: 4ac1a1bc6b3cbe66.js
a % b * c
; // File: 4ad6e3a59e27e9b1.js
/[\uD834\uDF06-\uD834\uDF08a-z]/u
; // File: 4ada45968b9f45ec.js
+a++ / 1
; // File: 4b346e8c85a29408.js
(function () {
    var a = {
        'b': 1
    };
}());

; // File: 4b40241551a495c2.js
({[a]: 1, b: 2})
; // File: 4b6559716b2f7b21.js
a++
{}
/foo/
; // File: 4ba667b7404cc45d.js
{ const a = 1 }
; // File: 4bd3199f5a4d8e52.js
try {} catch ({a}) {}
; // File: 4bd7e14411b6a889.js
[a, a] = 1
; // File: 4beb0b6ae8b9801a.js
function a({}) {}

; // File: 4bffa044ecd9d841.js
var [a, ...[b, c]] = d
; // File: 4c2a2b32f0470048.js
0O12
; // File: 4c3a394af4d281d1.js
function a(b, c) { return b < !--c; }
; // File: 4c44d7c28ec0f6ca.js
debugger
; // File: 4c56fb063bea0ec2.js
for (var a of b);

; // File: 4c71e11fbbc56349.js
({ set a(b) { }, set a(b) { } })
; // File: 4d310ef039a7435c.js
({a(b,c){}})
; // File: 4d833cbc56caaaf9.js
({ set 10(a) { a } })
; // File: 4d88f169e3827587.js
try {} catch ({a = 1}) {}
; // File: 4dc600d5ae71e8eb.js
// false
!/test/

; // File: 4deb8938d7b36024.js
a = (b(), c(), d())

; // File: 4df14f701f9881fd.js
var [a, b] = c;
; // File: 4dfe7c0219422eff.module.js
export var a = 1;

; // File: 4e07f8992cca7db0.js
'use strict'; 0o0
; // File: 4e1a0da46ca45afe.js
function a(...{ b }){}
; // File: 4e3de59ad16a7d0f.js
throw this
; // File: 4e625840177567fc.js
while (a) {
    if (b) break;
    c.d("a");
}
e: while (a) {
    if (b) break e;
    c.d("a");
}
; // File: 4e742059e0fc3d3c.js
for(a, b;;);
; // File: 4e7c58761e24d77c.js
a = { get null() {} }
; // File: 4e8e7d6fe1e67ce5.js
({1: {} / 1})
; // File: 4eafc760484cd72b.js
function a() {
    function b() {
    }
    function c() {
    }
    function d() {
    }
}

; // File: 4ec2942a940cd0b8.js
// Surpress reducing because of alternate
for (;;) {
    if (a) {
        if (b) {
            continue;
        }
    } else {
        ;
    }
}

; // File: 4ed17e0e2686e5e5.js
(function () {
    var a = 1;
    a + (b(), c(), d());
}());

; // File: 4eee835d0ac8382a.js
a && b && c
; // File: 4efcd175a5db8b47.js
if(a)b
;else c;
; // File: 4f21a4e88694c0d8.js
a = { set "null"(b) { c = b } }
; // File: 4f24ffe2c3ebe706.js
a.b.c.d
; // File: 4f53cda18c2baa0c.js
[]
; // File: 4f60d8fbb4be1120.js
09.0
; // File: 4f805a43cc2e8854.js
/[x-z]/i
; // File: 4fa08a62c2d8c495.js
(a = b('100')) == a 
; // File: 4fa4f9e47503bc5f.js
var a = 1, b = 2, c = 3
; // File: 4fdc22a42fa0d040.js
a - b + c
; // File: 4fee4ac53bdfd7f7.js
"Hello\0World"
; // File: 500804fd29695dac.js
0X1A
; // File: 5021396f85a70480.js
switch(a){case 1:}
; // File: 503cf49b200abf64.js
(1)
;
; // File: 50ac15a08f7c812f.js
({get __proto__() {}})
; // File: 50bc1f24c865c57a.js
a(b,c)
; // File: 50c6ab935ccb020a.module.js
export default (class{});
; // File: 50cea0e25b2b707d.js
function a(b, c) {
    var d = 1, e = f, f = d + e, g = h();
    return b + c;
}
; // File: 50e04108598730ff.js
({a = 1,} = 2)
; // File: 511a2a5fd8cac64d.js
d: {
    // this 'test2' to 'a'
    b: {
        if (a) break b;
        if (a) break b;
        if (a) break b;
    }
    if (c) break d;
}

; // File: 513275ce0e3c7ef3.js
'use strict';
var a = {
    'arguments': 1,
    'eval': 2
};

; // File: 5147bda197f961c1.js
({["a"+1]:"b"})
; // File: 515825915b8d1cd8.js
while (a) {
  try { } catch (b) { }
  // do not optimize it
  (function () {
    c('d');
  }());
}

; // File: 5171e99c2d9d3e5a.js
// Do not remove first if consequent block
if (a) {
    if (b) { true; }
} else {
    false;
}

; // File: 5183eafe6b4cd6e0.js
for(var a;b;c);
; // File: 51a3505b43223a9f.js
if (a) {
} else {
}

; // File: 51b243bb5076b692.js
(...a) => 1
; // File: 51b58dc84e1fab89.js
0X04
; // File: 51ea4e18429c02e4.js
if (a) var b = 1;
; // File: 51fd2b53ad7e1581.js
class a {b(){}}
; // File: 5203633f36fbe544.js
"use strict"; var { yield: a } = b;

; // File: 5212ddf4e4b70261.js
a = b; 
; // File: 521479b987ae2d7f.js
while (a) {
    if (b) {
        switch (true) {
          case c():
            d();
        }
        continue;
    }
    e();
}
; // File: 521b6dfff0a28aa1.js
function a([a=1]) {}

; // File: 523950fa023d7305.js
{ let a = 1 }
; // File: 524172bf792ef97e.js
for(a of b);
; // File: 52aeec7b8da212a2.js
if (1) function a(){}
; // File: 52ce5853ea953f0f.js
+{}
; // File: 52f2f30356750b9b.js
var a = /[P QR]/i
; // File: 52f9245e7cd97f6a.js
a: break a;
; // File: 5317b960ad78bbfe.js
function a(b, c) {
    // circular reference
    function d() {
        return e();
    }
    function e() {
        return d();
    }
    return b + c;
}
; // File: 53645d3765e5f67f.js
(function() {
  if (a) return b;
  return c;
}());

; // File: 54032532b8655caf.js
({a: b,} = 1)
; // File: 5406bea2982a6e13.module.js
import "foo";

; // File: 54190cc5a11a0233.js
(function () {
    ((a) ? 1 : 2) != null;
}());

; // File: 541ee533b54ae664.js
function a(__proto__) { }
; // File: 54257d53a8fffe8c.js
class a extends b { c() { [super.d] = e } }
; // File: 547fa50af16beca7.js
let [a,,b]=1

; // File: 5495e25325fdd364.js
function a(b = c) {}
; // File: 54e70df597a4f9a3.js
try { } catch (eval) { }
; // File: 54fb77cb2384a86b.js
'use strict';
{
    var a = 1;
    b();
    {
        b();
        b();
    }
}

; // File: 551af1dc1686e912.module.js
export {a} from "foo";

; // File: 5526c98fdf9150c1.js
(function() {
    try {
        throw 'a';
    } catch (b) {
    } finally {
        return 1;
    }
    c();  // This should be removed.
}());

; // File: 55b74de671f60184.js
class a {"constructor"(){} ["constructor"](){}}
; // File: 55c27b3727ba1165.js
(a, ...b) => {}
; // File: 55d1482dc2d95e91.js
1 + 2 << (3)
; // File: 55d721b105cc1780.js
a();
b();
c();
; // File: 561ccbf2e5091865.js
!(a=b)
; // File: 5641ad33abcd1752.js
a = { get true() {} }
; // File: 5665da18579dd006.js
1 + (a(), b(), c())

; // File: 569a2c1bad3beeb2.js
({a,...b}) => 0;

; // File: 56dcd0733a23aa26.js
if (a) {
  try { } catch (b) { }
  // do not optimize it
  (function () {
    c('d');
  }());
} else {
  try { } catch (b) { }
  // do not optimize it
  (function () {
    c('d');
  }());
}

; // File: 56debc26cbc2e077.js
a = [ ]
; // File: 56ec311ffc030121.module.js
export let a
; // File: 56fd564979894636.js
/\.\/\\/u
; // File: 571bb9d1fdd6fcc0.js
a: !function(){ a:; };
; // File: 575306c08cc44b10.js
(function () {
    a['NaN'] = 1;
}());

; // File: 578ebe526f02ab34.js
{ var a = 1, b = 2
c; }
; // File: 57971b49e239c0ff.js
function a(a) { 'use strict'; }
; // File: 57ad28ff7d96f031.js
for (var a in b) {
  c;
}

; // File: 5829d742ab805866.js
({a:(b)} = 1)

; // File: 585130f356b0729f.js
let [a,] = 1;

; // File: 5856de37689f8db9.js
(a) => a * yield;

; // File: 585b857c11763bad.js
([[[[[[[[[[[[[[[[[[[[{a=b}]]]]]]]]]]]]]]]]]]]])=>1;

; // File: 587400d1c019785a.js
try { } finally { a(b) }
; // File: 589dc8ad3b9aa28f.js
(function() {
    a((1, 2, 3));
}());

; // File: 58a52091eaa8746c.js
(function(){a()}(),b()+1)

; // File: 58cf2c5c0cecdf0e.js
class a {}
; // File: 58d72762ccb4d31f.js
typeof a
; // File: 58ed6ffb30191684.js
({ set false(a) { a } })
; // File: 597108fd45a6e79b.js
class a extends b { constructor() { ({c: super()}); } }
; // File: 597b9759467727fc.js
({a} = 1)
; // File: 5984eac0c5c6d947.js
[ ]
; // File: 598a5cedba92154d.js
[(a) = 1] = 2

; // File: 599dff255c5ec792.js
eval => 'use strict'
; // File: 59ae0289778b80cd.js
if (1) function a(){} else;
; // File: 5a06dab3e9fd0f65.js
(function() {
    for (;;) {
        break;
        a();  // This should be removed.
    }
    b();
}());

; // File: 5a079debdfff12da.js
a - b
; // File: 5a0dcc9e43fed2c2.js
(/* comment */{
    a: null,
    b: null
})

; // File: 5a51417e1ceb294f.js
for (;;) {
    while (true) {
        continue;
    }
}

; // File: 5a54ee2c0b326b18.js
function a(b) {
    for (var c = 1, d = b.e(); ; c++) {}
}
; // File: 5a7812b78a03b937.js
new a()
; // File: 5b146261dda66d63.js
('‪')
; // File: 5b39aca97d9006f4.js
a`token ${`nested ${`deeply` + {}} blah`}`
; // File: 5b4cef6792d9462f.js
void a
; // File: 5b683275df4548d1.js
(function() {
    var a = 1;
    a;  // 'i' should remain (dynamic)
    eval('');
}());

; // File: 5b8d2b991d2c1f5b.js
({})
; // File: 5b8fad162f489b3b.js
(a)--
; // File: 5b9f113c3bdd0c49.js
        function a() {
            b();
            c = 1;
            throw "d";
            // completely discarding the `if` would introduce some
            // bugs.  UglifyJS v1 doesn't deal with this issue; in v2
            // we copy any declarations to the upper scope.
            if (c) {
                e();
                var c;
                function b(){};
                // but nested declarations should not be kept.
                (function(){
                    var f;
                    function e(){};
                })();
            }
        }
; // File: 5bae374be95382c6.js
function a(){return
{}
/foo/}
; // File: 5bb4c1e68b0925d1.js
// global getter to o
a.b = (c(), 1)

; // File: 5beffd72ddb47f13.js
('\a')
; // File: 5c57eec29a019ebb.js
class a {static[b](){}; static[c](){}}
; // File: 5c587adcfe50a8c6.js
switch(a){case 1:default:case 2:}
; // File: 5c5ef7a4bdc3e081.js
var a;
// access to global should be assumed to have side effects
if (b) {
    a = 1+2;
} else {
    a = 3;
}
; // File: 5cc7ceeebdccb6d4.js
for(const a = 1;;);
; // File: 5cf0dc4259e98c15.js
var a, b;
a.c = (a = {}, 1);
b = (b = {}, 2);
; // File: 5d0cbb3fb27c21b7.js
for (;;) {
    if (a) {
        if (b) {
            continue;
        }
        c()  // This should not removed and translation should not occur.
    }
}

; // File: 5d1a7c61bf135457.js
((a))
; // File: 5d3e89c83953788e.js
function a() {
    // Do not remove this i
    for (var b in c);
}

; // File: 5d687a45c607ea42.js
{
  {
    a;
  }
}
{
  b;
}

; // File: 5d8ab2c35c7eb883.js
var a;
if (b()) {
    a();
} else {
    a();
}
; // File: 5dd65055dace49bc.js
({a,b=b,a:c,[a]:[d]})=>1;

; // File: 5e0cab2e2e36274c.js
a(....0)
; // File: 5e1cbe1737b1bbc6.js
(eval, a = 1) => 2
; // File: 5e6d5c3edf519b99.js
a = function({b} = {b: 1}) {}
; // File: 5ec03710bd21b933.js
b: while (1) { continue /*
*/ a; }
; // File: 5ecf2f4d83e6260d.module.js
"use strict";
; // File: 5f1e0eff7ac775ee.js
delete (1, a, 2)

; // File: 5f2834246274eccc.js
(void a)
; // File: 5f5e1d12ad68e832.js
if (a && b) {
    c(a)[1].b.d = e();
} else
    c(a)[2].b.d = f();
; // File: 5f730961df66e8e8.js
a = { get false() {} }
; // File: 5f85b0b6828b081b.js
(a, ...[b]) => {}
; // File: 5f9eeac7b076f34b.js
(a, {b = 1})=>2
; // File: 5fa8c711247d70f5.js
() => {}
; // File: 5fcc16142185c87c.js
var [] = 1;

; // File: 600327b79f60606c.js
(function () {
    ((a) ? 1 : b) != null;
}());

; // File: 600fd3c4d9f2ca42.js
a -= 1
; // File: 60a1991953372b97.js
(function() {
    for (;;) {
        continue;
        a();  // This should be removed.
    }
    b();
}());

; // File: 60bb345d725fe68b.js
(function() {
    if (a) {
        b();
        return 1;
    } else {
        b();
        return 2;
    }
    c();
}());

; // File: 60c092cb83b525f2.js
a &= 1
; // File: 60dcd48a3f6af44f.js
try {} catch(a) { var a = 1; }
; // File: 610b397691988417.module.js
import {b as a} from "foo";

; // File: 612fed84b89e42a8.js
1 instanceof 2

; // File: 61ceb5809404ee85.js
(a) + (b)
; // File: 61d8a7e497b6db72.js
(function () {
    function a() {
        (function () {
            b('c');
        }());
    }
    a();
}());

; // File: 61f55d9f22cc8426.js
class a {static static(){};}
; // File: 623cec03370f088a.js
a(...b, ...c, ...d);

; // File: 624bc7f99260037f.js
function a(b = 1) {}
; // File: 62541961bcef8d79.js
0b0
; // File: 627fede559e0bcac.js
(a = b('100')) != a 
; // File: 62ab44289ebbba49.js
(function() {
    return 1;
    function a() {
    }
}());

; // File: 62c217b2844680ab.js
a + b / c
; // File: 62d0da6771d5317d.js
function* a(){yield a}
; // File: 62d7c1ee4e1626c4.js
(function(){ return
a; })
; // File: 63208a19ffb4baeb.js
({["a" + "b"]: 1})
; // File: 633fac25082a90af.js
var [{a = 1}] = 2;
; // File: 63586de6fec2e3cf.js
({[a]: 1})
; // File: 639b9076cc56e57c.js
while (true) {
    if (a) break
    ;
    else b;
}

; // File: 63ee9cd383dc68a3.js
switch (1) {
  case 2:
    a();
    if (b) break;
    c();
    break;
  case 3+4:
    d();
  default:
    e();
}
; // File: 64117d5c682ec505.js
while (a) {
    b();
    c();
    continue;
}
; // File: 641ac9060a206183.js
([a, , b]) => 1
; // File: 645e8cce491528cd.js
typeof /test/

; // File: 646c2391c11102b5.js
(a,a)
; // File: 647e21f8f157c338.js
(' ')
; // File: 64cc57f82a54b7fb.js
({ get 10() {} })
; // File: 64ff3b3ee7f636c5.js
(function a({ b, c }){})
; // File: 65047600233c760c.js
a = false; 
; // File: 65228d6a31a06406.js
var a;
// compress these

a = true     && b;
a = 1        && c.d("a");
a = 2 * 3    && 4 * b;
a = 5 == 6   && b + 7;
a = "e" && 8 - b;
a = 9 + ""   && b / 10;
a = -4.5     && 11 << b;
a = 12        && 13;

a = false     && b;
a = NaN       && c.d("f");
a = 14         && c.d("g");
a = h && 15 * b;
a = null      && b + 16;
a = 17 * 18 - 19 && 20 - b;
a = 21 == 22   && b / 23;
a = !"e" && 24 % b;
a = 25         && 26;

// don't compress these

a = b        && true;
a = c.d("a") && 27;
a = 28 - b    && "e";
a = 29 << b   && -4.5;

a = b        && false;
a = c.d("f") && NaN;
a = c.d("g") && 30;
a = 31 * b    && h;
a = b + 32    && null;

; // File: 655eab0815e0567e.js
/**/ function a() {}

; // File: 65fcb1f93f1684ef.js
// Hallo, world!

; // File: 664b0da1dd015106.js
a["b"] = "c";
a["if"] = "if";
a["*"] = "d";
a["\u0EB3"] = "e";
a[""] = "f";
; // File: 665f4940c7cf30c9.js
(function(){ return; })
; // File: 668ab87597363d53.js
// ContinueStatement should be removed.
// And label is not used, then label also should be removed.
a: while(true) continue a;

; // File: 66aabbb5b00fb1ae.js
(function () {
  var a = 1;
  a;
  eval('');
}());

; // File: 66bd9903ea05f8cc.js
function a() {} function a() {}
; // File: 66d2dbcb692491ec.module.js
export * from "a"
;
; // File: 66ea15f7de78add7.module.js
export var a
; // File: 671c914df5da04df.js
(/* comment */{
    a: null
})

; // File: 6733f491913ccff2.js
class a {constructor(){}}
; // File: 673e6f2765ef3cb3.js
a = { }
; // File: 67711cbb84083749.js
a(
	b(c),
	d(c),
	e(c),

	b(),
	d(),
	e()
);
; // File: 680880af107834e8.js
/test/ && /test/

; // File: 68125aef6f5cc46f.js
'use strict'; 0b0
; // File: 6815ab22de966de8.js
for(let();;);
; // File: 681f352b7356594c.js
a => { return 1; }
; // File: 6823058797ddd563.js
/[a-z]/gim
; // File: 684237281767d41d.js
{do ; while(false); false}
; // File: 6861bb23b186f65a.js
/=([^=\s])+/g
; // File: 687b7b904904fcfd.js
class a {prototype(){}}
; // File: 687f678cde900411.js
([1].a) = 2
; // File: 69063bc9496ea6e5.js
a`hello ${b}`
; // File: 691e1d9954f3e6e2.js
a = { __proto__: 1 }
; // File: 697b3d30c1d06918.js
function a() {'use strict'; ({ b: 1, b: 2 }) }
; // File: 698a8cfb0705c277.js
({ a: 1, get a() { } })
; // File: 69bbdc7c34ed23cc.js
({ get a() { }, a: 1 })
; // File: 69bdc785b6e244ff.js
a||(b||(c||(d||(e||f))))

; // File: 69cbe8ec2f64382d.js
var a, b;
; // File: 6a218750a221c68b.module.js
import * as a from "a"
;
; // File: 6a220df693ce521c.js
for (a(b in c)[1] in d);
; // File: 6a240463b40550d2.js
// This test is ported from uglify-js1e20,1e21
; // File: 6a323491fe75918a.js
(function*() { yield 1; })

; // File: 6a735105a5e79722.js
(function () {
    for (var a = 1; a < 2; ++a);
}());

; // File: 6a7ed6cb99ea0b81.js
// because `with` can observe i lookup
a = a += 1

; // File: 6b5e7e125097d439.js
((((((((((((((((((((((((((((((((((((((((a.a)))))))))))))))))))))))))))))))))))))))) = 1
; // File: 6b63d36394b0ffb3.js
() => "a"
; // File: 6b68aefbfbf0beb9.js
a = { set() { } }
; // File: 6b86b273ff34fce1.js
1
; // File: 6bb2a138b9eb0088.js
1 * 2

; // File: 6c27d048b07ca7e0.js
class a {
    constructor() {
    };
    b() {};
};
class c {
    constructor(...d) {
    }
    b() {}
};
class e extends a {};
var f = class g {};
var h = class {};
; // File: 6c42024bfadac21f.js
a = { if: 1 }
; // File: 6c5f0dd83c417a5a.js
function a() {
    try {
        a();
    } catch(b) {
        var c = 1;
    }
    return c;
}
; // File: 6c688efe01b3631e.js
function*a(){yield*a}
; // File: 6cfcfc99afcb6e1a.js
for(let a of [1,2]) 3
; // File: 6d1bf4c3db76b489.js
0xdef
; // File: 6d707802519c7158.js
a(`${b} + ${c} = ${b + c}`)
; // File: 6d79220c64963dad.js
do ; while (true)
; // File: 6d8728cbc7bfe6b5.js
({})=>1;

; // File: 6d8c97119162ad95.js
`
	
`
; // File: 6d981ff8b6a3faec.js
var a = { *b () { yield *c } };
; // File: 6db6e4c3ba0299b7.js
a instanceof b
; // File: 6db7dbc9b1365dfa.module.js
export default /foo/
; // File: 6dcd76e9be7c3d00.js
b: {
    if (a) break b;
    c.d("e");
}
; // File: 6e5fe0c2bb20b016.js
({ "a": 1 })
; // File: 6ec818aa7f27cdbf.js
const a = 1
; // File: 6edc155d463535cb.js
(function () {
    'use strict';

    a = 1;
    function b() {
    }
});

; // File: 6f256be2ef45a7d6.js
__proto__: while (true) { break __proto__; }
; // File: 6f6e870785069487.js
/((((((((((((.))))))))))))\12/;
; // File: 6f824ec22e22a198.js
({a(b,...c){}})
; // File: 6ffb11115fcefb96.js
// mangle to the same name 'a'
c: {
          a("b");
          break c;
}
c: {
          a("b");
          break c;
}

; // File: 6ffb1fb47c2dd12f.js
({
    a,
    a:a,
    a:a=a,
    [a]:{a},
    a:b()[a],
    a:this.a
} = 1);

; // File: 6ffc173d4e1e5158.js
(class {;;;
;
})
; // File: 6ffd0afb19f0a92c.js
var a = class extends (b,c) {};
; // File: 702e4ee53d26635a.module.js
import a from "a"
;
; // File: 7055b45fe7f74d94.js
class a extends b { "constructor"() { super() } }
; // File: 70ad5a19a1b2a4b6.js
(function() {
    a: for(;;) {
        for (;;) {
            break a;
            b();  // This should be removed.
        }
    }
}());

; // File: 70b701c0eb7d36fd.js
if (!a) debugger;

; // File: 70bf2c409480ae10.js
({ get true() {} })
; // File: 70c2ced6bad143f1.js
new a.b()
; // File: 70da848e355cdfd2.js
new a();

; // File: 7148f242d6770f89.js
while (a) {
  b;
}

; // File: 714be6d28082eaa7.js
((((((((((((((((((((((((((((((((((((((((a)))))))))))))))))))))))))))))))))))))))) = 1
; // File: 717b2f65b69e809e.js
'use strict';(a)=>1
; // File: 717def9f9459b4e1.module.js
export {} from "a"
;
; // File: 71a2d3e7d606a959.js
1 >> 2

; // File: 71bcb4b846c22c58.js
a['0'];
a['1'];
a['00'];
a['0x20'];

; // File: 71e066a0fa01825b.js
('\11')
; // File: 72286da2cadacba6.js
let [a] = []
; // File: 726ee28a1b50ff13.module.js
export default { a: 1 };

; // File: 729212ece9234c48.js
a[1].b
; // File: 72d79750e81ef03d.js
a ** b

; // File: 72e4f3f9f66a40b8.js
function *a(){yield 1}
; // File: 7305be27a0713dfa.js
do {
  // do not optimize it
  (function () {
    a('b');
  }());
} while (c);

; // File: 73298cb8636154f2.js
var a = function eval() { };
; // File: 739bef73b11c87de.js
/[--]/
; // File: 74234e98afe7498f.js
null
; // File: 748a60621d2abe2b.js
__proto__: a
; // File: 74c5ebda713c8bd7.js
a = { set false(b) { c = b } }
; // File: 74cfbae1c9639338.js
var [a, ...a] = 1;
; // File: 753a8b016a700975.js
(function(a = b){})
; // File: 756579211447db0b.js
0O2
; // File: 756e3fe0ef87b136.js
function a() {
    b();
    c();
    return d();
}
function e() {
    b();
    c();
    throw new f();
}
; // File: 757fc3fbe38b4ecb.js
a(...b);

; // File: 75969544af546abc.js
new a(...b, ...c, ...d);

; // File: 75b0eeaf3aa61e74.js
function a() {
    if (b) return;
    c();
    d();
}
function e() {
    if (b) return;
    if (c) return;
    if (d) return;
    if (f) return;
    g();
    h();
}
; // File: 75bb6594d6ad253f.module.js
export let a = { }
; // File: 75e16348fe9e6213.js
var a = "b" + "c" + d() + "e" + "b" + f() + "d" + "f" + "g" + h();
var i = "b" + 1 + d() + 2 + "j";
var k = 3 + d() + 4 + "j";

// this CAN'T safely be shortened to 1 + x() + "5boo"
var l = 5 + d() + 6 + 7 + "j";

var m = 8 + d() + 9 + "n" + 10 + "j";
; // File: 764e72657e7321b5.js
/**
 * @type {number}
 */
var a = 1;

; // File: 765a7a823aa1b070.js
(function () {
    1/2;
}());

; // File: 76703c4b987330fb.js
a => { b: 1 }
; // File: 76a46be6c2f09fa3.js
function a() {'use strict';return 1;};
; // File: 76d4858e4a60be95.js
// one\n
; // File: 771467ccdae93157.js
1 /*the*/ /*answer*/
; // File: 7716587c3d80e9ab.js
(class { constructor() { super.a } });
; // File: 7733ab7955652851.js
({a=1}, {})=>2
; // File: 776076cb09759e40.js
1 /*The*/ /*Answer*/
; // File: 7779cfcd717e97d3.js
while(true) { continue
; }
; // File: 7788d3c1e1247da9.js
[(a.b)] = 1

; // File: 779e65d6349f1616.js
a = typeof 1;
b = typeof 'c';
d = typeof [];
e = typeof {};
f = typeof /./;
g = typeof false;
h = typeof function(){};
i = typeof j;
; // File: 77a541b0502d0bde.js
('\
')
; // File: 77c661b2fbe3dd3a.js
(a, b) => { 1; }
; // File: 77db52b103913973.js
[a, ...(b=c)]
; // File: 78435241f6c87ece.js
123.+1
; // File: 784a059faa166072.js
(function a(b, b) { })
; // File: 784cbc06d5ade346.js
/[-a-b-]/
; // File: 7855fbf5ea10e622.js
if (a) (function(){})
; // File: 787170711cb8abd6.js
((a))((a))
; // File: 789af9b27c832306.js
;;;;

; // File: 78cf02220fb0937c.js
[a,a,,...a]=1;

; // File: 78e1b8a4f3318967.js
"\u{714E}\u{8336}"
; // File: 78ea6e4e98c18f91.js
function a() {} / 1 /
; // File: 78eb22badc114b6f.js
function a() {
    try {
        a();
    } catch(b) {
        var c = 1;
    }
}
; // File: 78ecd285b8b44e99.js
var {a} = {}
; // File: 78fa04077cf1950c.js
({ set a([{b = 1}]){}, })
; // File: 790a34467d7d9d58.js
function eval() { }
; // File: 7912cf1671c75406.js
/*a
c*/ 1
; // File: 791ee64772f0ea45.js
function a() {
    var b = function c() { }
}

; // File: 7993945fc0f58feb.js
(a) => { yield + a };

; // File: 799fad61dcd88f30.js
({a = 1} = 2)
; // File: 79a4d1fdd55febec.js
var a = !b &&       // should not touch this one
    (!c || d) &&
    (!e || f) &&
    g();
; // File: 79b7f48e8a6d401c.js
a - b - c
; // File: 79ea421b940c3474.js
a << b
; // File: 79f7d5d83decb768.js
(function() {
    a: for(;;) {
        for (;;) {
            continue a;
            b();  // This should be removed.
        }
    }
}());

; // File: 7a405ea1fdb6a26e.js
a: b: while (true) { continue a; }
; // File: 7a71c5f849677cd1.js
do a();while (true)
; // File: 7a815cb480c3cac2.js
({} = 1);

; // File: 7a964712d5220b79.js
eval = 1
; // File: 7ab6a1dd47c6bc1f.js
+a
; // File: 7ac0063e99bc8720.js
for (var a of b) c(a);
; // File: 7ae57d0c2d30db3a.js
a >> b
; // File: 7afd38d79e6795a8.js
(function () {
    void 1;
    "not a directive";
    a();
})();

; // File: 7b0a9215ec756496.js
{ throw a/* Multiline
Comment */a; }
; // File: 7b514406528ff126.js
"Hello\02World"
; // File: 7b71bc250036251c.js
do continue; while(1);
; // File: 7b72d7b43bedc895.js
'use strict';
var a = {
    '10': 1,
    '0x20': 2
};

; // File: 7b8a8232be18df90.js
(function () {
    if (true) {
        a();
    }
}());

; // File: 7bc8dc445fc0f1c3.js
for(;;) { continue; }

; // File: 7bdcce70c382a9a4.js
(function () {
    var a = {};
    a.b + (c(), d(), e());  // do not transform
}());

; // File: 7be9be4918d25634.js
--arguments
; // File: 7c027cdbc7f493b2.js
var a = /[a-z]/y
; // File: 7c03e5eb6a9f6f1a.js
function a() { 'use strict'; "\0"; }
; // File: 7c3fc6d2a783ecd9.js
/*a
b*/ 1
; // File: 7c46ecc8f111b567.js
do a(); while (true)
; // File: 7c508ad20a5ecbce.js
class a extends b { c() { ({d: super[e]} = f) } }
; // File: 7c6d13458e08e1f4.js
01.a
; // File: 7c9c0cce695bc705.js
(function () {
    a();
    function a() {
        b.c('d');
    }
    function a() {
        b.c('e');
    }
}());

; // File: 7cbf77c14b9c89bc.js
function a(b, c) { return b-- > c; }
; // File: 7cd7c68a6131f816.js
a = { set true(b) { c = b } }
; // File: 7d029e0be60dc821.module.js
import a from "b"
; // File: 7d7dd05015778d56.js
(a,b) => 1 + 2
; // File: 7d8b61ba2a3a275c.js
with (a)
  // do not optimize it
  (function () {
    b('c');
  }());

; // File: 7da12349ac9f51f2.js
(a, b, c && d) && e;

; // File: 7dab6e55461806c9.js
function *a(){yield ~1}

; // File: 7dde401422530d6b.js
/*42*/
; // File: 7dea677261fc5dd8.module.js
export var a = function () {};

; // File: 7df2a606ecc6cd84.js
(function(){ return // Comment
a; })
; // File: 7dfb625b91c5c879.js
(function a() {'use strict';return 1;});
; // File: 7e094109208fc749.js
/[a-z]/g
; // File: 7e28d9664deeef8a.js
[{a=b}=1]

; // File: 7e50a0527f791c52.js
2e308
; // File: 7e6e3b4c766a4d33.js
(a = b('100')) <= a 
; // File: 7e6eac5fdc429608.js
for(a; a < 1;);
; // File: 7e88047a36603238.js
a != b
; // File: 7e8f17e7be305a2a.js
var [a = b] = c
; // File: 7e99cc8b7ce365fb.js
var a = {};
a.b = 1;
a.c = 2;
d.e(a.b);
; // File: 7ebaa39b4a9b5b5b.js
var a, b, c, d;
a = !(b(), c(), d());
; // File: 7f4c40906c3ebe2b.js
var {a} = 1;
; // File: 7f88f149f16fe97a.js
var a = {
    'arguments': 1,
    'eval': 2
};

; // File: 7fac17daa2bd5186.js
for(a = 1; a < 2;);
; // File: 7fbe94acda67721e.js
({*yield(){}})
; // File: 7fdf990c6f42edcd.module.js
export * from "a"
; // File: 7fe89d8edf6e778a.js
'a\u0026b'
; // File: 801ac33e4c34efb8.js
({ false: 1 })
; // File: 802658d6ef9a83ec.js
({ ['__proto__']: 1, __proto__: 2 })
; // File: 804e022cd08b4ae1.js
[a] = 1
; // File: 807dfc91f4ed4394.js
function *a(){yield 2e308}
; // File: 80950061e291542b.js
(function(){ return a * b })
; // File: 80c6bda5cad0fbc5.js
var a, [a] = 1;
; // File: 80d2351a5ae68524.js
a = { set b (c) {} } 
; // File: 80f60039028189e4.js
a(.0)
; // File: 811b309b010a36ce.js
({[1*2]:3})
; // File: 8152f05423c90f61.js
a >>>= 1
; // File: 8179659ef4fd0965.js
({a: b} = 1)
; // File: 81a0322e554af8da.js
(function () {
    a(typeof b === 'c');
}());

; // File: 81a6472df96f185f.js
a = []
; // File: 81b986948b58ffda.js
(function () {
    while (!a || b()) {
        c();
    }
}());

; // File: 81be1572d1eebdf2.js
a = { b(c=1) {} }
; // File: 81be47a15713178e.js
a + b + c
; // File: 820521ef532dce18.js
function a([b] = [1]) {}
; // File: 8286caaa8e0196bb.js
{ function a(){} }
; // File: 8290412f79ac2bb6.js
/a/
; // File: 82c827ccaecbe22b.js
[a = (b = c)] = 1

; // File: 832ad002639ce202.js
const a = 1, b = 2, c = 3
; // File: 8340cdb8653046bb.js
for (var [a, b] in c) {}
; // File: 836158118a07b45d.js
var {a, b: {c: a}} = 1;
; // File: 8386fbff927a9e0e.js
1000000000000000000000000000000

; // File: 838d87085df03a6d.js
"Hello\nworld"
; // File: 83cea3f2e14d1e23.js
(class a extends 1{})
; // File: 83f083525ae5a0e0.js
(function () {
    if (!!a) {
        b();
    }
}());

; // File: 83fc5b5bbdb601ef.js
(function () {
    function a() {
    }
    function b() {
    }
    function c() {
    }
}());

; // File: 8411f3c15e3e8529.js
{{}
/foo/}
; // File: 84250e15785d8a9e.js
({ set true(a) { a } })
; // File: 842fe071562c1a9e.js
function* a() { yield
; }
; // File: 845368e466d341f5.js
[{ a, b }, ...c] = d
; // File: 845631d1a33b3409.js
class a { b() {} get b() {} }
; // File: 845e30448809e2bc.js
a.b.c
; // File: 8462f068b299bca2.js
var {let, yield} = 1;
; // File: 849e112b480fda30.js
void (a)
; // File: 84b2a5d834daee2f.js
0012
; // File: 84eaae502ca93891.js
class a { b() { new super.c; } }
; // File: 84f901eb37273117.js
('\0')
; // File: 850a60daa178d3b6.js
() => 1 + 2
; // File: 85263ecacc7a4dc5.js
(function () {
    'use strict';
    function a() {
    }
});

; // File: 8543b43f3c48c975.module.js
export { default } from "a"
; // File: 855b8dea36c841ed.js
(function () { 'use\x20strict'; with (a); }())
; // File: 858d6a756ff641f3.js
function a(b) { c(); }
; // File: 8597768c0fe519eb.js
`${a}$`
; // File: 85d2d93264f2672d.module.js
import a, {b} from "foo";

; // File: 85d6723f13f33101.js
({'a': b} = 1)
; // File: 85e4314fa8f0661f.js
a(b, ...c = d);

; // File: 8628cd459b39ffe8.js
try { } catch ([a = 1]) { }

; // File: 8664d1a4e7a73078.js
function a({[b]: c}) {}
; // File: 86a25a2a0e393ed6.js
switch (a) { case 1: let b = 2; break; }
; // File: 86b0ffc811e713ec.js
({ __proto__: 1 })
; // File: 86d8f1465e745b44.js
[a[b]=b] = c
; // File: 86f68610fcefaeae.js
class a {b() {} static c() {}}
; // File: 870a0b8d891753e9.js
(function(a, ...b){})
; // File: 8751eb24f903c279.js
var a = 1
; // File: 87844be2334fba9e.js
function a() { "use strict"
 + 1; }
; // File: 87a9b0d1d80812cc.js
(function() {
    a((b(), 1, 2));
}());

; // File: 87cb789c4ed2b97a.js
a: while (true) { continue a; }
; // File: 87e1d3eab9d05339.js
while (a()) b();
; // File: 88127d108648d05b.js
while (true) { continue /* Multiline
Comment */a; }
; // File: 881a7a3d4e17e621.js
function a(b, c) {
    var d = function() { return e() };
    var e = function() { return d() };
    return b + c;
}
; // File: 882910de7dd1aef9.js
((((((((((((((((((((((((((((((((((((((((a))))))))))))))))))))))))))))))))))))))))
; // File: 884e5c2703ce95f3.js
({ *a() { super.b = 1; } });
; // File: 8854cac4acddd510.js
null && (a += null)
; // File: 88561e211e862344.js
null;
; // File: 88827d8021b5b3ab.js
new a("aa, bb", 'return aa;');
; // File: 88af07b3dc006159.js
a | b ^ c
; // File: 88c21621e3e8bba0.js
('a')
; // File: 88d42455ac933ef5.js
class a {static constructor(){} static constructor(){}}
; // File: 88e99d6cd8e8d87f.js
function *a() {}

; // File: 891fc3470b618587.js
a((b, c) => {})
; // File: 892f6e09c02c35b5.js
0.
; // File: 8996d3eb07c6f7cd.js
for(let a of b);
; // File: 89a31837e6736b2a.js
{a: 1}

; // File: 89c872e56d527908.js
1+2;
; // File: 8a0fc8ea31727188.module.js
export * from "foo";

; // File: 8a40542f1f53c4f0.js
for(const a of b);
; // File: 8aa3cd2609b4f278.js
a = { b: 1 }
; // File: 8ad4edbe9317df28.js
a(...b, ...c)
; // File: 8ae0c86bd7897b7b.js
a(...b = c)
; // File: 8af69d8f15295ed2.js
(' ')
; // File: 8b4ff58f416e17b5.js
(a++)
; // File: 8b6bded4f89f89f6.js
if ((function(){ return true })()) {
    a(true);
} else {
    b(false);
}
(function(){
    c.d("e");
})();
; // File: 8b9cd46352285386.js
new a[b]
; // File: 8bd57faa6bcca5e2.js
 //

; // File: 8be0df708b9e56ca.js
a = [ 1, ]
; // File: 8bf3ec35c55ed3c0.js
(class extends a {})
; // File: 8c27fb7ef1e3ca3d.js
(function() {
if (a) return b;
else return c;
}());

; // File: 8c56513a6ac3cdff.js
(/* comment */{
    /* comment 2 */
    a: null
})

; // File: 8c56cf12f007a392.js
a <!--bar
+b
; // File: 8c5c46a300d5addb.module.js
export function a(){};1
; // File: 8c80f7ee04352eba.js
(...[]) => 1

; // File: 8c8a7a2941fb6d64.js
a = [ 1, 2,, 3, ]
; // File: 8ce30bb40ffff192.js
arguments => 1
; // File: 8d14286a8cc6ee9d.js
(function*(...a) {})

; // File: 8d2463220e3cd0d7.js
(function () {
    // 'a'
    b: {
        if (a) break b;
        (function () {
            // 'a'
            b: {
                if (a) break b;
                c("d");
            }
        }());
    }
}());

; // File: 8d2ffebc7214c34f.js
yield* 1
; // File: 8d40617aec6fabe5.js
({a: b = 1} = 2)
; // File: 8d470354c5d2e216.js
;
function a() {
}

; // File: 8d5679ec94e658e1.js
~a
; // File: 8d5aeaf1120f0897.js
throw /* ‪ */ a
; // File: 8d6435fa243d5b1a.js
var a
; // File: 8d67ad04bfc356c9.js
({set __proto__(a) {}})
; // File: 8d6ab6352a3f7fa0.js
(a = b('100')) >= a 
; // File: 8d7d59e5d573ca84.js
({ "__proto__": null, set __proto__(a){} })

; // File: 8d8913ebd8403c6a.js
// 1
a();
b();
for (; false;);
// 2
a();
b();
for (c = 1; false;);
// 3
c = (a in b);
for (; false;);
// 4
c = (a in b);
for (d = 2; false;);
; // File: 8dbf5f9322d8b5ac.js
a = { b(c) { } }
; // File: 8ddbdc8954dd7aa9.js
function a() {
    if (false) {
        var b = 1;
    } else {
        var c = 2;
    }
    d(c, b);
}

; // File: 8e37579cd5ffb2df.js
var [...a] = b
; // File: 8e3f0660b32fbfd2.js
('\7a')
; // File: 8e609bb71c20b858.module.js
export {};

; // File: 8e6c915d1746636d.js
(({a} = {a: 1}) => {})
; // File: 8e8e921a75194950.js
do {
  a;
} while(b);

; // File: 8ecaef2617d8c6a7.js
try{}catch(a){}
; // File: 8ed2a171ab34c301.js
eval++
; // File: 8ed2fce2b9b43fb7.js
a(b, ...c)
; // File: 8ef08a335a7f5966.js
        a.b('c');
        a.b.d(a, arguments);
; // File: 8f0084b1073e1877.js
(a, 1)
; // File: 8f4f97274dea4723.js
(function() {
    a: {
        break a;
        b();  // This should be removed.
    }
    c();
}());

; // File: 8f659ed872554f99.js
class a { static get() {}}
; // File: 8f8a9f6ca890939e.js
({__proto__:1})
; // File: 8f8bfb27569ac008.js
'use strict'; eval[1] = 2
; // File: 8f9e8be5a6c50e77.js
function a() {
    return function b(c) {};
}
; // File: 8fcaa7f3f8926a5e.js
0B1
; // File: 9013f39c33dc8416.js
throw a;
; // File: 9027dae72a91a9ed.js
((a)) = 1
; // File: 9036914be00c0dc7.js
(a, b, [c]) => 1
; // File: 903dd05bf49c8fac.js
a(b, ...c);
a(...d);
; // File: 906e545ceef0fcfd.js
(function() {
if (a) throw b;
else throw c;
}());

; // File: 90919cd6fd06c4d8.js
function *a(){yield a}
; // File: 90ad0135b905a622.js
function a() {
    var b = 1;
    c && (2, 3, d, b);
}

; // File: 914ae3168da48965.js
throw {}
; // File: 9159ea4175a5a021.js
for(a = 1; a < 2; a++);
; // File: 918e105a2ff6c64a.js
class a extends b {get c() {}}
; // File: 91cbb6971c86509e.js
({a = 1}, {b = 2}, {c = 3})=>4
; // File: 91f2fa0b11550b30.js
function *a(){yield "a"}
; // File: 9203cb34e9b091dc.js
a: if (true) break a;
; // File: 9208254b5f8a8481.js
a == b
; // File: 923c99b441ab5a26.js
var a = [ "b", "c", "d" ].e("");
var f = [ "b", "c", "d" ].e();
var g = [ "b", 1, 2, 3, "c" ].e("");
var h = [ i(), "b", 4, 5, 6, "c", c() ].e("");
var j = [ i(), c(), "b", 7, 8, 9, "c", c() ].e("");
var k = [ 10, 11, "b", "c", d() ].e("");
var l = [ "b", 12 + 13 + "c", "d" ].e("m");
var n = [].e(b + c);
var o = [].e("");
var p = [].e("b");
; // File: 925443c6cf79aa88.js
// ContinueStatement should not be removed.
a: do do continue a; while(true); while(true);

; // File: 927efb51d4882ccb.js
(function() {
    1/-2;
}());

; // File: 92a997b1ba17876e.js
({a = 1}) => a
; // File: 92dd079c741d2a95.js
a <= b
; // File: 92fd8e24864fde0a.js
(function eval() { });
; // File: 93108a695e5ff29d.js
(function () {
  var a
  if (b) return /* I would
                              insert something
                              there, but I'm sort
                              of lazy so whatever.
  */ a = new c()
  return a
})()

(function () {
  var a
  if (b) return /* I would insert something there, */ /*
                             but I'm sort of lazy so
  */ /*                      whatever. */ a = new c()
  return a
})()

(function () {
  var a
  if (b) return // I would insert something there, but I'm sort of lazy so whatever.
  a = new c()
  return a
})()


; // File: 9312a1adbbf0a4c0.js
debugger;
; // File: 9331c78bb0fc6a55.js
// do not transform to 0
-1

; // File: 9349f48a456341b8.js
new a().b()
; // File: 937059a5e53177e5.js
[...[...a[b]]] = c
; // File: 939479d60d564ccd.js
({ *a(b, c, d) {} })

; // File: 93bd9d668ac308a0.js
var a = function(b) {
    b();
    a()
}

; // File: 93c32bb0a4bad388.module.js
import {default as a} from "foo";

; // File: 93c75264893587a5.js
({ *yield() {} })

; // File: 93cac77bbf2242ab.js
(a(),b(),c()) + (d(),e(),f());

; // File: 94846b0ae1cac1a2.js
void void void /test/

; // File: 9495a0dcecf5713c.js
({[a]:()=>{}})
; // File: 94b8a654a87039b9.js
(function*(a, b, c) {})

; // File: 94be09b126b946b8.js
`

`
; // File: 94c72b68d8726b07.js
(1, a)
; // File: 94cb828d5dcfd136.js
try { } catch (arguments) { }
; // File: 951e1b875db534f9.js
a = (1) ? 2 : 3
; // File: 95408309cc3a1d30.js
function a() {
    if (void 1) {
        var b = 2;
    } else {
        c(a);
    }
    c(b);
}

; // File: 954a896471379dc8.js
[a, ...a] = 1
; // File: 95520bedf0fdd4c9.js
function a() {
    b,c,d;
    if (e) {
        throw 'f'
    }
}

; // File: 955c5fedb3931500.js
'a' + 'a'

; // File: 95ab0d795c04ff38.js
('\111')
; // File: 96059002704b3ac3.js
for (;a();) {
    b();
    c();
    if (d()) e();
    else break;
    f();
    g();
}
; // File: 9677a7160d769b1a.js
var a = {
    delete: 1
};

; // File: 9681f5d844d7acd0.js
(a) => b;  // 1 args
(a, b) => c;  // n args
() => b;  // 0 args
(a) => (b) => c;  // func returns func returns func
(a) => ((b) => c);  // So these parens are dropped
() => (b,c) => d;  // func returns func returns func
a=>{return b;}
a => 'e';  // Dropping the parens
; // File: 96909e1dce85ca53.module.js
export default 1;

; // File: 96941f16c2d7cec4.js
(a)=(1);
; // File: 96e1b294d19a101d.js
(a) => "b"
; // File: 96f5d93be9a54573.js
({a, b});
[{a}];
a({c});
; // File: 970fb35ce6ce89bb.js
function a() {
    var a = 1;
    function b() {
        var c = a;
    }
}

; // File: 973cbc9ece13acbc.js
function* a() {} function a() {}
; // File: 974e7275fdedce49.js
1 .a
; // File: 97593deb177d09ae.js
function a(b, c) { return b-- >>> c; }
; // File: 976afd9ae5f5d71a.module.js
export function a(){}
;
; // File: 979b36a2c530f286.js
var { a: yield } = b;

; // File: 97e246302dfe8616.js
function a(b, c) { return b << !--c; }
; // File: 982595e2af9d9703.js
for(a = 1; a < 2; a++) b(a);
; // File: 982835d8c977075c.js
/a/;
; // File: 988e362ed9ddcac5.js
(a**b).c=1

; // File: 98c7fb7947f7eae4.js
eval--
; // File: 98df58b0c40fac90.js
a = function(){}(b())

; // File: 993584ec37388320.js
0b1
; // File: 9949a2e1a6844836.module.js
(function* () { yield a })
; // File: 996001e00a0c575b.js
({a})
; // File: 9974571a855d4447.js
function *a(){yield false}
; // File: 9975820eb10bc0ff.js
for (;;) { a(); continue; b(); }

; // File: 999c1001e3761320.js
`abc`
; // File: 99cdfc40e20af6f5.js
var a = function(){}(b())

; // File: 99fceed987b8ec3d.js
function a(...[]) { }

; // File: 9a5b92dfd9d19f60.js
    1


; // File: 9a666205cafd530f.js
function *a(){({b(){yield}})}

; // File: 9a6711e879a99536.js
{ let a; }
; // File: 9a7f06880ce32bbc.js
a=1;
; // File: 9a9cb616daadf90a.js
(function(){ return a; })
; // File: 9b9d0e250e01155d.js
var a, b;
if (c()) {
    a = new b(1);
} else {
    a = new b(2);
}
; // File: 9bcae7c7f00b4e3c.js
(function a(...[ b, c ]){})
; // File: 9c30b0817f412a30.js
(function(){ return {} })().a = 1; // should not transform this one
; // File: 9cf32425f04fd865.js
`${a}${b}`
; // File: 9d0fd95dd43f59ec.js
1 /* The * answer */
; // File: 9d1320f0185b1586.js
a[b](b,c)
; // File: 9d3d960e32528788.js
(function () {} / 1)
; // File: 9d935d1b787ed251.js
//
;a;
; // File: 9db4dccf1122bfc4.js
new new a
; // File: 9db573299f02bf36.js
var a = [ "b", "c", d(), "e", "f", "g" ].h("");
var i = [ "b", "c", d(), "e", "f", "g" ].h("j");
var k = [ "b", "c", d(), "e", "f", "g" ].h("l");
var m = [ "b", "c", d(),
          [ "b", 1, 2, 3, "c" ].h("+"),
          "e", "f", "g" ].h("j");
var n = [ "b", "c", d(),
          [ "b", 4, 5, 6, "c" ].h("+"),
          "e", "f", "g" ].h("l");
var o = [ "p", "p" + q, "b", "c", "r" + b ].h("");
; // File: 9dc20e081005fba4.js
a = {
    b(c, d) {
        return a;
    }
}
e = {
    b([{c}]) {
        return c;
    },
    f(){}
}
; // File: 9dfa08b5b7ad82a9.js
var [a, ...b] = c
; // File: 9e3e46891aaf13de.js
let {a:b} = {}
; // File: 9e98dbfde77e3dfe.js
[ ,, 1 ]
; // File: 9ec644dbf797e95c.js
('\`')
; // File: 9ed0369295348e76.js
(function() {
    return 1;
    {
        var a = 2;
    }
}());

; // File: 9f23d57c37e238cb.js
/**
 * @type {number}
 */
var a = 1,
    /**
     * @type {number}
     */
    b = 2;

; // File: 9f272d23fc62842a.js
a || (b || c)

; // File: 9fd584806e085e35.js
class a { b() {} c() {}}
; // File: a0079146ab045c26.js
/[-a-]/
; // File: a022debc42a58f0c.js
([a]) => [1];

; // File: a03475b02913e16a.js
for (;a();)
    if (b()) break;
; // File: a08ed291f78352a0.js
a(function() {
  function b() {}
  if (b()) {
    b();
    return void 1;
  }
});
; // File: a0af29e4dd6d3845.js
(function () {
  while (false) {
    a();
  }
  b();
}());

; // File: a0b7bf790311b763.js
/[a-z]/u
; // File: a0fba75a10c21ac9.js
(!a)?1:2

; // File: a11e875c4dd100af.js
for (var a = 1; a < 2; ++a)
  // do not optimize it
  (function () {
    b('c');
  }());

; // File: a150e917230aa57d.js
(function () {
    if (true == a) {
        b();
    }
    if (false == a) {
        b();
    }
}());

; // File: a157424306915066.js
function *a() { (b) => { yield + b }; }

; // File: a15c90cc56980c41.js
a: while (true) { break a; }
; // File: a18a8865e65d4bdd.js
({[a]() {}})
; // File: a194909bf50b1467.js
function *a() { yield yield }

; // File: a1ab463999957845.js
(    a  )()
; // File: a1cd0d76806cef90.js
(function(){
  return typeof a() != "b";
}())

; // File: a1dde88d8a87b573.module.js
import {a} from "foo";

; // File: a2042d86c592dd55.js
// NaN => 0/0
'a' - 1

; // File: a248ef84a404262c.js
`$`
; // File: a2781f8227f7f1e6.js
[...{ a }] = b
; // File: a2798917405b080b.js
try { } catch ({}) {}

; // File: a29b007e8fb9d020.js
[a, b, ...c] = 1
; // File: a2c2339691fc48fb.js
//
; // File: a2cb5a14559c6a50.js
try {} catch (yield) {}

; // File: a2f26b79b01628f9.js
(function () {
    var a = 1;  // should not hoist to parameter
    with (b)
        arguments = 2;
}());

; // File: a33250f92a2f000e.js
function a(b, c) { return b-- >> c; }
; // File: a353c62d8ed56d6f.js
{ throw a// Comment
a; }
; // File: a378fc25898cf05b.js
( new a).b()
; // File: a3b497c58f78b243.js
(function () {
    function a() { }
}());

; // File: a43df1aea659fab8.js
a & b & c
; // File: a445a478b4ce0c58.js
({ get a() { return b } })
; // File: a454d2e2ab3484e6.js
(function*() {})

; // File: a487d6498ec0efbf.module.js
import a, {} from "a"
;
; // File: a4931f8127e03c4e.js
a.null
; // File: a4b3402765acaa0e.js
(a) ? (b) : (c)
; // File: a4d62a651f69d815.js
if (a) function b(){}
; // File: a515d7aaea7b816f.js
(function () {
    if (a) return 1;
    else return;
}());

; // File: a54cca69085ad35a.js
class a { static set b(c) {}}
; // File: a54ce6036e646e24.js
a = b % c / c * d * 1;
a = b % c * 2
; // File: a59e0d0b4d3e1b7d.js
/* multiline
comment
should
be
ignored */ 1
; // File: a5a01023fef4d506.js
a['b'];
a['in'];
a['eval'];
a['arguments'];

; // File: a5a7eb83bf27418b.js
for(a in b) c(a);
; // File: a5aaa3992025795a.js
({yield = 1} = 2);
; // File: a5b30a03e9c774af.js
try {
  // do not optimize it
  (function () {
    a('b');
  }());
} catch (c) {
  // do not optimize it
  (function () {
    a('b');
  }());
} finally {
  // do not optimize it
  (function () {
    a('b');
  }());
}

; // File: a62c6323a3696fa8.js
('\
')
; // File: a6806d6fedbf6759.js
[ 1, 2,, 3, ]
; // File: a68bf27df978135d.module.js
export default 1;2
; // File: a6b7dab7088e5269.js
var a = /[a-z]/i
; // File: a6cb605b66ef0eb5.js
(function () {
    (!!a, !b, 1)
}());

; // File: a7964b1dcfd2dc21.js
() => 1, 2
; // File: a7b8ce1d4c0f0bc2.js
b: while (1) { continue  a; }
; // File: a7c8ad2a73ed00d0.js
for (let [a, b] of c);

; // File: a830df7cf2e74c9f.js
({ get a() { super[1] = 2; } });
; // File: a8535eac4c7c9c3a.js
class a { static() {} }
; // File: a86a29773d1168d3.js
try { } catch ([]) {}

; // File: a871e54c3956acd9.js
a: while(true) { continue a
; }
; // File: a885d88fb046dea5.js
(function () {
    for (var a in []);
}());

; // File: a8a03a88237c4e8f.js
({a:(b = 1)})

; // File: a8b6c3139974f6e1.js
({ set "null"(a) { a } })
; // File: a8b832d61af9cdc4.js
/[a-z]/y
; // File: a8fea31fe6aa588e.js
0o1
; // File: a91ad31c88855e59.js
var implements, interface, package
; // File: a93f6b22796d4868.js
function* a(){ () => yield; }
; // File: a953f09a1b6b6725.js
02
; // File: a955c7a892679016.js
for (var a in b) c(a);
; // File: a9b0aedbd9f25ec9.js
var a = function() { b() };
; // File: a9e054dbd43d4b86.js
a * b / c
; // File: a9e4ff39f145a1fd.js
({a({b} = {b: 1}) {}})
; // File: a9f99e657441a735.js
(class {static(){}})
; // File: aa06ccdc7ff9e10d.js
({a: [b] = 1} = 2)
; // File: aa3b466be5c7f7e3.js
for (;;) {
    if (a) continue;
    continue;
}

; // File: aa3d1fa7a22e6460.js
while(true) var a
; // File: aa7e721756949024.js
({ set a(b) { super.c[1] = 2; } });
; // File: aaa0bc7fa72df5e4.js
// Hello, world!
1
; // File: aab08ba9fd01cbb8.js
a.b.c.d(1)
; // File: aab51bc524d9c623.js
(function () {
    a(null!=(b?c:void 1));
}());

; // File: aac70baa56299267.js
function a() {
    return (void 1, void 2, a, void 3);
}

; // File: aaf1be6cd60a9ac9.js
({"a"() {}})
; // File: ab23ca0a6e356883.js
do {
    a();
} while (false);b()
; // File: ab452fc45813857a.js
if (a) {
  try { b('try'); } catch (c) {}
  // do not optimize it
  (function () {
    b('d');
  }());
}

; // File: ab4734949243c00a.js
new a(b, ...c = d);

; // File: ab7ea8d738da7043.js
function a() { return
; }
; // File: abcfae2381708c43.module.js
export {default} from "foo";

; // File: ac09566949f0db57.js
{ do { } while (false) false }
; // File: ac112e0c69fe603e.js
(a, b) => "c"
; // File: ac1bc6b81949c063.module.js
let a

; // File: ac6bbe8465f70ebd.module.js
export {}
;
; // File: ad06370e34811a6a.js
({a:(b) = 1} = 2)

; // File: ad0fd65944942eee.js
0e+100 
; // File: ad410fec8c09f67b.js
for ([a, b] of c);

; // File: ad4414fcaaa6abb7.js
if(a)b;
; // File: ad6bf12aa7eda975.js
(function([a]){})
; // File: ad6bfcbfca5afde1.js
.14
; // File: ad7d61a732903cd8.js
({ a([ b, c ]){} })
; // File: adad54e09949b0e9.js
a && b
; // File: ade301f0d871c610.js
"Hello\122World"
; // File: ae204e41bacd8237.js
class a extends b { c() { super.yield } }

; // File: ae4bbee73a0f80a5.js
[ 1, ]
; // File: ae700e3f8ff82c6c.js
a()+(b(),c(),d(),e())  // do not transform

; // File: ae89d08bdd65b56b.js
class a {b(){}c(){}}
; // File: ae9667ad0d837abc.js
(1, 2, a)();
(3, 4, b.a)();
; // File: ae97d36bd01b43b2.js
let [...a] = 1;

; // File: ae9a8ca09473df05.js
(function () {
    void ((a) ? 1 : 2);
}());

; // File: aec65a9745669870.js
(function () {
    var a,b,c,d,e,f;
    (a,b,c()) + (d,e,f());
}());

; // File: aeca992c7be882ba.js
(function() {
    a, 1, 2;  // 'i' should remain (global variable)
}());

; // File: af17707f71e402a7.module.js
import {a, b,} from "foo";

; // File: af1d905ed056724f.js
var a = function arguments() { };
; // File: af4bbcea9802b120.js
(class {a() {}})
; // File: af5766d06630bbc5.js
[[a]] = 1
; // File: af97a3752e579223.js
(function () { yield* 1 })
; // File: afa63b136c835723.js
({a:1, get 'b'(){}, set 3(c){}})
; // File: affd557fd820e1f2.js
({get [a]() {}, set [a](b) {}})
; // File: afffb6d317e53b92.js
('\ ')
; // File: b030378ad6e36751.js
class a { static [b]() {} }
; // File: b0423be1317c7b69.js
(1) + (2  ) + 3
; // File: b05d4355cc5e2802.js
a => b => c => 1
; // File: b0659cf9cb6793a2.js
class a extends b { constructor(c = super()){} }
; // File: b06e2c3814e46579.js
function *a(){yield class{}}
; // File: b07c5fdc1003316b.js
switch (a) {
  case 'b': c(); break;
  default:
    break;
}
; // File: b0a834e1180ccd73.js
1 / 2

; // File: b0c6752e1db068ed.js
(a = yield) => {}

; // File: b0fdc038ee292aba.js
09.5
; // File: b1072e92becf06a9.js
(function () {
    a(typeof b === 'function');
}());

; // File: b13e700d2613a5a7.js
while (a) {
    b();
    c();
    break;
}
; // File: b175bdef718c4012.js
if (a) { /* Some comment */ b() }
; // File: b195e06e2ba5e787.js
a.in / b
; // File: b1b969fae2973dae.js
while (a) {
    {
        b();
        while (a);
    }
    b();
}

; // File: b1c37dedeec0b867.js
function a(...b) { }
; // File: b2048a6a14348122.js
for(a;b;c);
; // File: b205355de22689d1.js
function *a(){({set b(c){yield}})}

; // File: b2495b0864c7835a.js
a:{break a;}
; // File: b24fa2a1936d38d9.js
({3() {}})
; // File: b25057b11104844d.js
var a = ({ a : 1 });
a(({ a: 2 }));
; // File: b29070859dbeaf08.js
[a,b=1,[c,...a[2]]={}]=3;

; // File: b2a13c2c2c795427.js
for(;;);
; // File: b2a567473d09b770.js
function a() {
  debugger;
  return;
}

; // File: b2e6c124e2822117.js
({ "__proto__": null, get __proto__(){}, set __proto__(a){} })

; // File: b3093945d65d33d5.js
for(var a = 1;;);
; // File: b32aa0e4195927c1.js
(1) / 2
; // File: b363a70923be42c8.js
for(;;){}
; // File: b3717dd9314332d2.js
/(?!.){0,}?/;
; // File: b376d3924d77aa8a.js
+{"\\":1}

; // File: b3e783194b210cc3.js
var a;
// compress these

a = true     || b;
a = 1        || c.d("a");
a = 2 * 3    || 4 * b;
a = 5 == 6   || b + 7;
a = "e" || 8 - b;
a = 9 + ""   || b / 10;
a = -4.5     || 11 << b;
a = 12        || 13;

a = false     || b;
a = 14         || c.d("f");
a = NaN       || c.d("g");
a = h || 15 * b;
a = null      || b + 16;
a = 17 * 18 - 19 || 20 - b;
a = 21 == 22   || b / 23;
a = !"e" || 24 % b;
a = null      || 25;

a = c.d(h && b || null);
a = c.d(h || b && null);

// don't compress these

a = b        || true;
a = c.d("a") || 26;
a = 27 - b    || "e";
a = 28 << b   || -4.5;

a = b        || false;
a = c.d("f") || NaN;
a = c.d("g") || 29;
a = 30 * b    || h;
a = b + 31    || null;

; // File: b41dad3363eaab20.js
while(1);
; // File: b42377ca7015e7d4.js
const {a} = {}
; // File: b433b0cfef2a3cd1.js
({[a]: function() {}})
; // File: b43b759f208e1afb.js
while (true) { break
a; }
; // File: b4546b664cc70c58.js
b: while (1) { continue 
 a; }
; // File: b45d869ab09c7b00.js
for(let [a] of b);
; // File: b46a2c1b5d9a97cd.js
`42`
; // File: b471292fedf2f813.js
function a() {
    b();
    function c() {}
}
; // File: b49ca42695798691.js
(function () {
    var a;
    if (a) return;
    else return;
}());

; // File: b4a54642589bc396.js
a.$.b.c
; // File: b4bc5f201c297bca.module.js
import a, * as b from 'baz';
; // File: b506e9cc13c4ad2e.js
({0x0:1})
; // File: b549d045fc8e93bf.js
"use strict";
try {} catch (a) {}

; // File: b563460c1031daf2.js
({ a() { super.b(); } });
; // File: b5755ec32b9418af.js
({a,b} = 1)
; // File: b5a1e1a3679f81ba.js
({ a }, ...b) => {}
; // File: b5bc1ffd90912fb1.js
a = 1
; // File: b5cf21a87ec272d1.js
('\01')
; // File: b5d302467c6f2f16.js
(class {a(b) {'use strict';}})
; // File: b6145fa4a8cb8c35.js
a => b => 1
; // File: b62c6dd890bef675.js
var a = {
    '10': 1,
    '0x20': 2
};

; // File: b6b5e49c97cedebb.js
{
    'use strict';
}

; // File: b6e396c8cdf28f95.js
// Hello, world!

//   Another hello
1
; // File: b6efc88be898bda8.js
(function a() {
    var b = 1;  // should not hoist variable.
}());

; // File: b756b64f0eef72db.js
function a() {
  if (!b || c());
}

; // File: b75a0b610a41c000.js
(function () {
    if (a) return 1;
    else return 2;
}());

; // File: b77549e54bfef0f9.js
var {a:b} = {}
; // File: b79d2c4df1141981.js
function a([b]) {
    c();
    var d;
    var b;  // Because anArg is already declared, this goes away!
}
; // File: b7a5cd294201221b.js
/* block comment */ 1
; // File: b7a6a807ae6db312.js
({ set a(b) { }, a: 1 })
; // File: b7c2a3690011dd5e.js
for(var a in b);
; // File: b7cd8df8dc875529.js
        while (!((a && b) || (c + "0"))) {
            d.e("f");
            var a;
            function b() {}
        }
        for (var c = 1, g; c && (g || c) && (!typeof c); ++c) {
            h();
            a();
            var i;
        }
; // File: b7d99c0034be0ce1.js
function a(b, b) { }
; // File: b8403938b1ddd626.js
a.b,a.Infinity,a.true,a.false

; // File: b885e6a35c04d915.js
for(let a;;); let a;
; // File: b88624492a2c81d3.js
(a,b,...c) => 1;

; // File: b89aef8a4690aa20.js
while (true) { continue
a; }
; // File: b8ad1bd2ff50021f.js
0x0
; // File: b8bf39a3e60568ab.js
(eval, a) => 1
; // File: b8f8dfc41df97add.js
(function(){ return true })() ? a.b(true) : a.b(false);
(function(){
    a.b("c");
})();
; // File: b926f0fefd69158a.js
switch (a) { case 1: b(); break; default: break }
; // File: b92bdcf6c2591e5c.js
(class extends (a,b) {})
; // File: b92dd0bc25eaebe3.js
// ContinueStatement should not be removed.
d: for(var a in b) for (var c in b) continue d;

; // File: b93d116fd0409637.js
({ get false() {} })
; // File: b96ba7cdf0b42ca9.js
(function(a){}(b()))

; // File: b9a0cb6df76a73d2.js
a = { set 10(b) { c = b } }
; // File: b9a4f9232146d4d9.js
0e+100
; // File: b9a5f5c8c12525c7.js
/0/g.a
; // File: b9b8fb218e1990af.js
a&&(b=c)&&(d=e)
; // File: b9e1124424a35ad1.js
function a() {
    if (false) {
        var b = 1;
    }
    c(b);
}

; // File: ba21e63736d8fd46.js
a: function a(){}
; // File: ba4cc699857f41f2.js
function a() { function* a() {} function a() {} }
; // File: ba6624f5f448dfe4.js
for(a; a < 1; a++);
; // File: ba6dd18da17dbc10.js
for (;;) {
    if (a) {
        if (b) {
            continue;
        }
    }
}

; // File: ba9a047839eb4682.js
a = { b() { } }
; // File: bad55fbc19618df8.js
a = typeof b.c != "d" 
; // File: bb402b98f5398890.js
a || b ^ c
; // File: bb41f0778f00f131.js
let eval = 1, arguments = 2
; // File: bb447d4ed988a1cb.js
+{ }
; // File: bb7db120ad2fe995.js
({a: function({b} = {b: 1}) {}})
; // File: bb87b410a1170cf0.js
for (let a of b);

; // File: bb8b546cf9db5996.js
({ get a() { return a }, set a(a) { return a; } })
; // File: bbff5671643cc2ea.js
(class {set a(b) {}})
; // File: bbffb851469a3f0e.js
try { } catch (a) { let b; }
; // File: bc302492d441d561.js
(function(){ var a = 1; /* sync */ }).b(this)
; // File: bc89b2b2f1e19f9e.js
{}
/foo/
; // File: bcd690cfb709ffe8.js
for ({a = 1} in b);
; // File: bce83ece0ba80598.js
// Hello, world!

; // File: bd160eed5626ae7d.js
d: while (a) {
    b();
    c();
    continue d;
    e();
    f();
}
; // File: bd28a7d19ac0d50b.js
"use strict"; var a = { b: 1, get b() {} }
; // File: bd697f0fda948394.js
(function() {
    try {
        throw 'a';
    } catch (b) {
    }
    c();  // This must not be removed.
}());

; // File: bd7b54d5e0ce444b.js
(function () {
    // 'b'
    b: {
        if (a) break b;
        // 'a'
        c: {
            if (a) break c;
            if (a) break c;
            if (a) break c;
        }
        // 'a'
        c: {
            if (a) break c;
            if (a) break c;
            if (a) break c;
        }

    }
}());

; // File: bd9b563f02b80dae.js
"use strict" + 1
; // File: bdc4accd07049034.js
(class {;;;
;a(){}b(){}})
; // File: bde1a5ea9aebf9d2.js
function arguments() { }
; // File: bdfc6c05edd19925.js
/{}/;
; // File: be2c3fff6426873e.js
/*ab*/ 1
; // File: be2fd5888f434cbd.js
(function () { 'use\nstrict'; with (a); })
; // File: be6eb70d9330c165.js
0x100
; // File: be879445c87d7e72.js
var a;
; // File: be9d538d5041fd5f.js
(...[a, b]) => {}
; // File: beb5335e463d92c1.js
a;
; // File: bedf5be599c82fe8.js
({ *a() { yield 1; } })

; // File: bf1db420b006027f.js
{do ; while(false) false}
; // File: bf210a4f0cf9e352.js
class a { static get b() {} get b() {}}
; // File: bf65e886047db371.js
var a = function (...b) { c() };
; // File: bf6aaaab7c143ca1.js
012
; // File: bf8ffad512a5f568.js
class a { static get b() {} static set b(c) {} get b() {} set b(c) {}}
; // File: bf9c4d8ecd728018.js
0x10
; // File: bf9e8bd90d8537c3.js
throw { a: "b" }
; // File: bfb61863d3b10adf.js
// Infinity => 1/0
1 / 2

; // File: c06df922631aeabc.js
if (1); else function a(){}
; // File: c0740dd25c9de39b.js
3.14159
; // File: c0f5f3f7db69c5a0.js
a = { get() { } }
; // File: c1319833fc139cf8.js
var a = {
    '0': 'b'
};

; // File: c162248ee699b68f.js
var a; var b; var c;

; // File: c16d7f2993152b6b.module.js
import { b as a } from "c"
; // File: c17fd07fc9b5bf7e.js
new a(1);
new a(2)(3);
new a(4)(5)(6);
new new a(7);
new new a(8)(9);
new (new a(10))(11);
(new new a(12))(13);
; // File: c18d547cafb43e30.js
this

; // File: c1914072e996ddbe.js
/\uDF06/u
; // File: c1967c44c4179fb4.js
 /**/
; // File: c1c5c5d42a32aac1.js
(class {a({b} = {b: 1}) {}})
; // File: c1c8f5c6abfc1d72.js
({var: a} = 1)
; // File: c1dd2285cd8a959d.js
for(a of b) c(a);
; // File: c1e396cb7871b175.js
({a: [1]}+[]) / 2
; // File: c2116aecaac68db9.js
function *a() { yield *yield }

; // File: c2203cb9e7bfe40f.js
a--
; // File: c247dcc00119f19c.js
for(var a in [1,2]) 3
; // File: c24da2ce6761a80a.js
a + b - c
; // File: c2558dc20a45b0a8.js
({[1+2]() {}})
; // File: c25bf945aaff8fe1.js
"Hello"
; // File: c269a2a601c495f1.js
(function () {
    a = 1;
    for (;;);
}());

; // File: c274891790345c56.js
-a
; // File: c27ded6ec20ea305.js
while(true) { break
; }
; // File: c2d90d623b0f4c2e.js
a += 1
; // File: c2f12d66ce17d5ab.js
({ get null() {} })
; // File: c2fe8f120b796831.js
for(a in b);
; // File: c30eafd82f40470b.js
(1 + 2 ) * 3
; // File: c3172ad30aed99c8.js
const[a] = b;
; // File: c35304fa99a2c331.js
(function () {
    if (a) return;
    else return 1;
}());

; // File: c35dcf99291ec6be.js
({ __proto__: null, __proto__(){}, })

; // File: c3699b982b33926b.js
a | b & c
; // File: c3799cf68cbac258.js
({0: a} = 1)
; // File: c38644033565f7b9.js
(()=>1)
; // File: c3ce623096553057.js
a + (b < (c * d)) + e
; // File: c3dc60d438666700.js
new a()``
; // File: c3fc8ace42f3fb44.js
a.if
; // File: c412905e229d6f2b.js
delete a
; // File: c4336e0b6801c42c.js
({ *a() { 
  yield
  1
} })

; // File: c4627eaa56f73949.js
1 % 2

; // File: c4a57a72e25e042c.js
a = { get "b"() {} }
; // File: c4c51e5c6d4012ef.js
a: while(true) { break a
; }
; // File: c4f2243c81525bbd.js
var a = (1 + 2)
; // File: c52db35cba7fdbc0.js
(function() {
    try {
        throw 'a';
    } finally {
        b();
    }
    c();  // This should be removed.
}());

; // File: c5328483d3ccadd0.js
0o0
; // File: c546a199e87abaad.js
for(let [a=b in c] in null);

; // File: c5823f1dccaf9787.js
for(var a of b);
; // File: c58e9029f1fd3d1b.js
let a = 1;
; // File: c5957fd3a6d258df.js
b: {
  a;
}

; // File: c5b2ea7da55d24c1.js
function a() { function a() {} function a() {} }
; // File: c5bd72f618d7cade.js
function* a(){}
; // File: c6827eb9dd7b3dc6.js
({a = 1} = b)
; // File: c6ea3404ea5c6c91.js
a(...b, c)
; // File: c6ff61d189c5cbee.js
(a[b]||(c[d]=e))
; // File: c756f39dca1f7423.js
a = function(b, ...c) {}
; // File: c771490bbb3dd6e9.js
class a { *static() {} }
; // File: c78c8fbfbd3e779e.js
while (1) /foo/
; // File: c7dd4bc60ffb40e9.js
var {} = 1

; // File: c7e5fba8bf3854cd.js
(function () {
 'use strict';
 ' ';
}())
; // File: c83a2dcf75fa419a.js
/[\]/]/
; // File: c844c5ec9f6dbf86.js
a = "b" != typeof c.d 
; // File: c85bc4de504befc7.js
a = { get b() { return c } }
; // File: c8689b6da6fd227a.js
a = { get b() {} }
; // File: c87859666bd18c8c.js
while (1) {} /foo/
; // File: c88c5d1e7e9574b6.js
var a = /[a-c]/i
; // File: c8b9a4d186ec2eb8.js
a: do continue a; while(1);
; // File: c8dbdecbde2c1869.js
(function() {
    throw 'a';
    a();
}());

; // File: c963ac653b30699b.js
a = { 'b'() { } }
; // File: c964ed7bc2373c54.js
function a({a=(1), b}) {}
function b([b, c=(2)]) {}
var { d = (3), e } = d;
var [ d, e = (4) ] = d;
; // File: c98889d7d94a0a63.js
[a, ...[b, c]] = d
; // File: c98eba310f5568b1.js
for (;;)
    if (a()) break;
; // File: c9b780fb91a9db4e.js
/*@ngInject*/
var a = function(b) {
    return b;
}
; // File: c9d32e4fc1687f5d.js
c: switch (1) {
  case 2:
    a();
    for (;;) if (b) break c;
    d();
  case 3+4:
    e();
  default:
    f();
}
; // File: ca20c15b39c87033.js
class a {static(){};}
; // File: ca34a796e624adaf.js
switch (1) {
  case 2: a();
  case 3+4: b(); break;
  case 5+6+7: c();
  default:
    d();
}
; // File: ca39d991b4f07bf1.js
var a, b, c;
d(
	a(e),
	b(e, e),
	c(e)
);
; // File: ca450ebe11a7e7c9.module.js
import "a"
; // File: ca452a778322112a.js
// 
; // File: ca4f13a64e35195f.js
a = { null: 1 }
; // File: ca7a0ca0d22f30f8.js
(class{[1+2](){}})
; // File: ca978112ca1bbdca.js
a
; // File: caa0719b52a1409d.js
a``
; // File: caaa9f06dd52e5a5.module.js
import { null as a } from "b"
; // File: caf6539007d41b5e.js
[/q/]
; // File: cb05f3c30f5f88c0.js
new a(...b, ...c)
; // File: cb095c303f88cd0b.js
('\2111')
; // File: cb211fadccb029c7.js
function *a(){yield delete 1}

; // File: cb23f6635a581786.js
a + b
; // File: cb3316f2b008bec3.js
({ *a() { yield* 1; } })

; // File: cb4b35cf4cd815d8.js
({[a]: a} = 1)
; // File: cb625ce2970fe52a.js
(function () { 'use\x20strict'; with (a); })
; // File: cb898749d76e51fd.js
1 + 2

; // File: cbc644a20893a549.js
a & b
; // File: cbc7fdab53161051.js
debugger;
if (a) debugger;
; // File: cbccdb75b22a522c.js
function *a(){yield-1}
; // File: cbf9e832efe61a2e.js
for (let {} in 1);

; // File: cc561e319220c789.js
do /x/; while (false);
; // File: cc6ea8664124953a.js
a.true
; // File: cc793d44a11617e7.js
try { a(); } catch (b) { c(b) }
; // File: cc7b1f054147aa5c.js
switch (a) { case 1: b() /* perfect */ }
; // File: ccd1f89a0344e04e.js
(a=1) => a * a
; // File: cd136009983641b5.js
(function(){
    if(a) return b;
    if(b) return c;
}())

; // File: cd2f5476a739c80a.js
new(a in b)

; // File: cda499c521ff60c7.js
a = class {
    static b() {}
    static get c() {}
    static set d(a) {}
    static() { /* "static" can be a method name! */ }
    get() { /* "get" can be a method name! */ }
    set() { /* "set" can be a method name! */ }
}
; // File: cdb9bd6096e2732c.js
++arguments
; // File: cdbd7fe30e1c7321.js
(function () {
    var a =1;  // should hoist to parameter
}());

; // File: cdca52810bbe4532.js
function a() {
    b.c('d');
}
a();
{
}
function a() {
    b.c('e');
}

; // File: cdee1bf4a6391af8.js
++a
; // File: cdf411040ab4b29b.js
((() => {}))()
; // File: cdf43a987840ece8.js
/[]/
; // File: ce0aaec02d5d4465.js
class a { static get [b]() {} }
; // File: ce349e20cf388e87.js
var a = {
    'b': 1
};

; // File: ce3d1f8d346bb92d.js
class a {;;}
; // File: ce52f1c3d90b194a.js
// ContinueStatement should not be removed.
a: for(;;) for (;;) continue a;

; // File: ce569e89a005c02a.js
'use strict'; arguments[1] = 2
; // File: ce5f3bc27d5ccaac.js
function *a() { (b) => b * yield; }

; // File: ce6a4854c1f79924.js
// Surpress reducing because of alternate
for (;;) {
    if (a) {
        if (b) {
            continue;
        } else {
            ;
        }
    }
}

; // File: ce8c443eb361e1a2.js
(function([]){})
; // File: ce968fcdf3a1987c.js
function *a(){yield typeof 1}

; // File: cea8816bffe4238c.js
function a() {
    if (void b()) {
        c();
    }
}

; // File: cec2d94dc09a6a71.js
a = {
    catch() {},
    throw() {}
}
; // File: cefd0dd07bfa670f.js
const a = 1;
; // File: cf0eb6e6c4317c33.js
('\u{00F8}')
; // File: cf0fb26afd0eaaf1.js
new a(...b = c)
; // File: cfca620b63dd98b8.js
a.b`foo`;
c `bar`;
; // File: cfebdd6b58e65e90.js
function*a(){yield
a}
; // File: d010d377bcfd5565.js
({ __proto__: null, get __proto__(){}, set __proto__(a){} })

; // File: d038789ad15922ff.js
(function () {
    var a = {
        'Infinity': 1
    };
}());

; // File: d043d114b966415b.js
function *a(){yield --a;}
; // File: d0724a029fb7e4b1.js
var a;/* block comment 1 */ /* block comment 2 */
; // File: d082f8d1c2eec454.js
[a,] = 1
; // File: d09dbe1357abd967.js
(function () {
  (function () {
    a("b");
  }());
}());

; // File: d0dba4e03608ad64.js
`${a}`
; // File: d126aa10835287e6.js
class a {;}
; // File: d198e0d3a33b7b61.js
var a = {*[b]() { yield *c; }}
; // File: d1eafbc6bda219a7.js
class a{}
; // File: d1fea0e461717b09.js
`\``
; // File: d22f8660531e1c1a.js
var static;
; // File: d2332f9187c6a20a.module.js
export default class {}
; // File: d24d5f53dc15bcc7.module.js
export var a
;
; // File: d2ae1c7b6e55143f.js
({ set a(b) { b } })
; // File: d2af344779cc1f26.js
(a) => 00
; // File: d2c9ab6dc14dc774.js
for({a=1} of b);
; // File: d2d8885e0c00ad51.js
("\\\"")
; // File: d2fe67b1990df65c.js
(function () {
    a ? !b : !c;
}());

; // File: d33efc20e46c3961.js
(...a) => {}
; // File: d368a7bc70ca3120.js
({a(b,c){let d;}})
; // File: d3762fcf2ad7d285.js
(class extends 1{})
; // File: d37653c5aedf3d46.js
({"a"(b, c, d) {}})
; // File: d38771967621cb8e.js
('\5111')
; // File: d3bc8cc2c239b25f.js
a * b
; // File: d3c1ea553fea8944.js
class a { static set(b) {};}
; // File: d3d6ca7932414eed.js
function a(b = 1) { }
function c(b = (2 + 3)) { }
function d({ e } = {}, [ f ] = [ 4 ]) { }
; // File: d3f70f4410bb8346.js
/}/;
; // File: d4104d0ed6a07c28.js
// Do not mangle to the same name
e: {
    d: {
        a("b");
        if (c) {
            break d;
        }
        break e;
    }
}

; // File: d42cf386ef394628.js
a(1).b(2, 3, 4).c
; // File: d45f1126ef89120b.js
let [] = [];

; // File: d45fa56c26ed4a4e.js
class a extends 1 {}
; // File: d483926898410cae.js
"Hello\0122World"
; // File: d487300b8deff2ff.module.js
import "a"
;
; // File: d4b898b45172a637.js
"use strict"; var a = { get b() {}, get b() {} }
; // File: d4c979f1a92a8cac.js
(a(), b(), c()) + d();

; // File: d4e81043d808dc31.js
function *a() { yield b.yield; }

; // File: d515f6ce0c47a609.js
(a,b,...c) => 1 + 2
; // File: d51711f888aeeac9.js
while (true) continue;  // should be empty statement

; // File: d53aef16fe683218.js
switch(a) { case 1: {}
/foo/ }
; // File: d54bfed43597e9ac.js
void void a

; // File: d55a93310a309c43.js
/a/i;
; // File: d57a361bc638f38c.js
var a, b, c, d;
a = (b, c, d);
; // File: d57d9e2865e43807.js
switch (a) {
  case 'b': c();
  default:
    d();
    break;
}
; // File: d58831cddf9cbd99.js
        function a() {
            b();
            c();
            d = 1;
            return;
            if (d) {
                e();
            }
        }
; // File: d59a168fe5b7c787.js
var a = /42/g.b
; // File: d59a6667e160c0b3.js
/\uD834/u
; // File: d60426fd0160fb91.js
(function () {
  // not void context
  // do not optimize
  1 + (function () {
    return 2;
  }());
}());

; // File: d61d161a9c36fa45.js
(function() {
    (1, eval)('');  // indirect call to eval
}());

; // File: d6aed84ca98bee95.js
({a, b})
; // File: d6b2fd3884a06d56.js
var {get} = a;
; // File: d6bb7d557971a15f.js
new a`42`
; // File: d7076912d1c9786c.module.js
export default [];

; // File: d7284aa68a87bb97.js
let [a,] = [1]
; // File: d759838042f0bf78.js
a = [ 1, 2, 3, ]
; // File: d767138e133ad239.js
({ "__proto__": null, get __proto__(){} })

; // File: d79a08ea5cc1e2f6.js
const eval = 1, arguments = 2
; // File: d7c7ff252e84e81d.js
try { a(); } catch (b) { c(b) } finally { d(e) }
; // File: d7da7ccd42af2c4b.js
a.a *= 1
; // File: d7e461a3aa2cd9bc.module.js
export let a = 1;

; // File: d80edd7fb074b51d.js
6.
; // File: d818deffd07a5c3a.js
for (var a in b) continue;  // should be empty statement

; // File: d81d71f4121e3193.js
('\u{0}')
; // File: d82ae3dbc61808f8.js
function a(...yield) {}

; // File: d843ddb6cde8c408.js
(function(){
  return typeof a() !== "b";
}())

; // File: d8882ceedce6eae0.js
// ContinueStatement should be removed.
// And label is not used, then label also should be removed.
c: for(var a in b) continue c;

; // File: d8aff43ba7b44ef3.js
(function () {
    if (true != a) {
        b();
    }
    if (false != a) {
        b();
    }
}());

; // File: d8b6a56583bdefab.js
++eval
; // File: d8db2079f10d30ff.js
while (true) a()
; // File: d917a549d3f308d8.js
a / b
; // File: d93ec22aea12336a.js
function a() {
    class b {
    }
    var c = b
    var d = class e {}
}
; // File: d94d38d65e8b715f.js
function a() { return null
; }
; // File: d95b0608f939e81a.js
for(var a = 1;b;c);
; // File: d95e9ad32d562722.js
(function(){
  var a;
  return function(){
    return typeof a === 'b';
  };
}());

; // File: d96153b59454dddd.js
[1].a = 2
; // File: d97144839fbdca91.js
('\')
; // File: d99414900a405295.js
a=(function(){ return 1;})()
; // File: d99714b3c4e81b56.js
for (;a();) {
    if (b()) c();
    else break;
}
; // File: d9a0d4f0a35dc04e.js
0.0.a
; // File: d9d0b115106f376c.js
(function () {
    switch (a) {
        case 1:
            b("c");
        default:
            // drop this default clause
            // https://github.com/mishoo/UglifyJS2/issues/141
    }
}());

; // File: da3756d1f8acb3c5.js
({ a }) => {}
; // File: da4c5dd50fbdda83.js
new a``
; // File: da671a25e498bcac.module.js
export default function a() {}
; // File: da9fdc3a2d7f9452.js
function a(b, c = 1) {}
; // File: dad51383642e0d27.js
/{/;
; // File: dadae97bf343020d.js
function *a(){yield true}
; // File: dadccefeaae19dbf.js
/[-\-]/u
; // File: db1fd3f76ebc6554.js
a(b, ...c, d)
; // File: db3c01738aaf0b92.js
a **= 1

; // File: db456532eea62941.js
if (a) b(); else c()
; // File: db66e1e8f3f1faef.js
b: while (1) { continue /**/ a; }
; // File: db8fe6c7579e6ead.js
{ throw a
a; }
; // File: dc1acc240053a397.js
({a(){}})
; // File: dc2b756c7828d827.module.js
export function a() { }
; // File: dc3afa2f13259ae0.js
('\ ')
; // File: dc3e1097a489e009.js
throw 1
;
; // File: dc43022b3729abd1.js
0o12
; // File: dc9d42142b4ada05.js
if(a) { b = 1; for(var c;;); }

; // File: dcc5609dcc043200.js
(class {static constructor(){}})
; // File: dcc634c173bc704f.js
var a = !a || !b || !c || !d || !e || !f;
; // File: dcc9c2ff46392f30.js
a = { b: function(c=1) {} }
; // File: dcdf666e16667f4c.js
/**/1
; // File: dcfaa5f359400cf2.js
for (;;) {
    if (a) continue;
}

; // File: dcfb11abc780d6d9.js
{ do { } while (false);false }
; // File: dd0e8f971ab4d6ab.js
(function a() {
    b(typeof a() === 'c');
}());

; // File: dd3c63403db5c06e.js
((((((((((((((((((((((((((((((((((((((((((((((((((1))))))))))))))))))))))))))))))))))))))))))))))))))
; // File: dd500055335127b3.js
var a = 1;
this.b = 2;

; // File: dd67e8365153c4fb.js
(class extends a { constructor() { super() } });
; // File: dd80c278722f97e9.js
/*@ngInject*/
function a(b) {
    return b;
}
; // File: ddcd0bf839779a45.js
switch (a) {
  case 1:
    // do not optimize it
    (function () {
      b('c');
    }());
}

; // File: ddd3c540fa087867.js
(a(), b.b)()

; // File: ddef0827f7a75499.js
var a /* comment */;
; // File: de24062f6e293cf0.js
[[[[[[[[[[[[[[[[[[[[{a=b[1]}]]]]]]]]]]]]]]]]]]]]=2;

; // File: de25059a9dd7b618.js
(function () {
    null!=(a?void 1:void 2)
}());

; // File: de6b6c9002d2d43e.js
new function () {
    var a = 1;
    b(this.constructor.arguments.c);
};

; // File: de6dd6b2ec971861.js
class a {get() {}}
; // File: dec1ae80150e1664.js
{}/=/
; // File: dec6aac10ea17f7f.js
function a() {
    while (b) {
        {
            c();
            c();
            var d = 1;
        }
        {
            c();
        }
    }
}

; // File: decdfa7f961d283c.js
if(a)b;else c;
; // File: df20c9b7a7d534cb.js
a = { get b() { return c }, set b(b) { c = b; } }
; // File: df4eb225b4ba9ae2.js
a?b:c
; // File: df5fee9e52377ab9.module.js
export class a {}
; // File: df7e8c48ed8d9e6f.js
(function () {
    return 1;
    a();
}());

; // File: df9c60e4ff82b9d9.js
new new a()
; // File: dfa22e3eac3cd26e.js
class a {static b() {}}
; // File: dfbd1b07bd57a08d.js
function a(...b) { }

; // File: e0204155218e1d42.js
123..a(1)
; // File: e05209211a87a606.js
debugger
;
; // File: e08112a34cfea369.js
for (var [a,b] in c);
for (var [d] = 1;;);
for (var {e} of f);
; // File: e08e181172bad2b1.js
({get [1+2]() {}, set [3/4](a) {}})
; // File: e0b98eaceaaeaf9b.js
--a
; // File: e0c3f07a142a589d.js
function a() {
  if (b) {
  } else {
      c();
  }
}

; // File: e0f831f2b08fd35c.js
0b1001;
0B1001;
0o11;
0O11;
; // File: e0fc2148b455a6be.js
(function({a: b, a: c}){})
; // File: e1237566c1f89d8e.js
if (a) b()
; // File: e12aa6994333466f.js
var a; // if undeclared it's assumed to have side-effects
if (b()) {
    a(c);
} else {
    a(d);
}
; // File: e1387fe892984e2b.js
function a({ b: { c, a }, d: [e, f] }, ...[b, d, g]){}
; // File: e1820bdb79ebe44b.js
[...[a]] = 1
; // File: e18c297bf29c4b6b.js
var a, {b: {c: a}} = 1;
; // File: e1939e7cb50f65b4.js
(function*() { [...{ a = yield }] = 1; })
; // File: e1d373aa5d926fde.module.js
export var a = { }
; // File: e1dd1979a86a5f1d.js
try {} catch ([a,b, {c, d:e=1, [f]:g=2, h=i}]) {}

; // File: e23748bdbb0713dc.js
(function () {
    function a() {
        b.c('d');
    }
    {
        function a() {
            b.c('e');
        }
    }
}());

; // File: e23f481ffc072aee.js
if (a) b();
if (!a); else b();
if (a); else b();
if (a); else;
; // File: e290a32637ffdcb7.js
(function () {
    a['Infinity'] = 1;
}());

; // File: e2ac0bea41202dc9.js
({ get __proto__() { return 1 }, __proto__: 2 })
; // File: e2c7f7c0da23bc45.js
(function () {
    switch (a) {
        default:
        case 1:
            b("c");
    }
}());

; // File: e349023df8e12f2d.js
1 in 2

; // File: e374d329af31c20a.js
'\''

; // File: e3b0c44298fc1c14.js

; // File: e42f306327c0f578.js
"use strict"; var a = { set b(a) {}, b: 1 }
; // File: e463265266cee73e.js
function* a () { yield *b }
; // File: e46381af137ed2e2.js
a(1).b
; // File: e46f7944dd0d4eb4.js
1 /* block comment 1 */ /* block comment 2 */
; // File: e4bd395227b4ee8e.js
([a]) => 1
; // File: e4c6c19e4b214180.js
function *a(){yield+1}
; // File: e4cef19dab44335a.js
[1,,2]

; // File: e5204d6e30f296a8.js
switch(a){default:case 1:}
; // File: e5393f15b0e8585d.js
({let} = 1);
; // File: e54c1a2fc15cd4b8.js
class a {static b(){} static get b(){} static set b(c){} }
; // File: e5570b178254bfb9.js
// ContinueStatement should be removed.
// And label is not used, then label also should be removed.
a: do continue a; while (true);

; // File: e577d5b725159d71.js
0..a
; // File: e5951efaf0b0c5b3.js
function *a(){yield(1)}
; // File: e5a7d56b798ec7e6.js
a("\v");
; // File: e5fbf9e911ec36cd.js
a => 1
; // File: e65f3cca9a4637c3.js
'use strict';

; // File: e6643a557fe93de0.js
({yield} = 1);
; // File: e686d016100a7a08.js
class a extends b {
    c() {
        new super.d()
    }
}

; // File: e6ac25f6aa73a2be.js
/test/||1

; // File: e6b424d430520bf2.js
function a(b) {
    if (c) for (var d = 1, e = b.f(); ; d++) {}
}
; // File: e6e24cfdc6d308a2.js
{ const a = 1, b = 2, c = 3 }
; // File: e71a91c61343cdb1.js
a = { get 10() {} }
; // File: e71c1d5f0b6b833c.js
[(a)] = 1

; // File: e720d4faf2b41f42.js
function a(){ /*Jupiter*/ return; /*Saturn*/}
; // File: e748a1e428ccdf69.js
new a(b,c)
; // File: e74a8d269a6abdb7.js
var private, protected, public
; // File: e75df8aea1749780.js
a *= 1
; // File: e78c7b54fc87d08c.js
function *a(b, c, d) {}

; // File: e7c1f6f0913c4a95.js
function a(){yield*a}
; // File: e7c444fc9aed1257.js
for(const a in b);
; // File: e7fa87b10d5136a0.js
class a {get b(){} set c(d){};}
; // File: e815494eb50fa42f.js
function a() {
    var b;
}

; // File: e84ef669246313d2.js
('(');
a(')');

; // File: e877f5e6753dc7e4.js
(a,b,c,d) ? e : f

; // File: e899a2594bd5311c.js
() => 1
; // File: e8de5af87dc0004c.js
a /= 1
; // File: e8ea384458526db0.js
var {[a]: b} = {c}
; // File: e8ef6188865f9def.js
d: {
    if (a) b("c");
    else break d;
    e.f("g");
}
; // File: e8ef944fd2c2e7fa.js
a`42`
; // File: e95b9364e90a4b5c.js
`$$$${a}`
; // File: e9682c37a1a959e1.js
"\u{20BB7}\u{91CE}\u{5BB6}"
; // File: e99d260ec2ea47be.js
let a = 1, b = 2, c = 3
; // File: e9a24a964ace5330.js
if (a) {
  b;
} else {
  b;
}

; // File: e9a74729daea9b84.js
a + b * c
; // File: e9d44e4cbaf92011.js
({[a](){}})
; // File: ea2e883b50b24651.js
b: switch (1) {
  case 2:
    a();
    for (;;) break b;
    c();
    break;
  case 3+4:
    d();
  default:
    e();
}
; // File: ea3fcad439ac905f.js
({ a: 1, a: 2 })
; // File: ea54fe11ef8702f7.js
function *a() { yield 1; }

; // File: eabc983d82222f2a.js
for (let a of b) c(a);
; // File: eb4b9e8905923468.js
0xabc
; // File: eb7bb0c4a0ced2a8.js
a || b || c
; // File: ebbc09d90157cb5b.js
function a() {
    var {b, c} = a;
    var d = a;
}
; // File: ebd6534f7bb01a7a.js
({ true: 1 })
; // File: ec05d8a5722be86c.js
(function () {
    var a = 1;
    a = a += 2;
}());

; // File: ec782937135d4f32.js
/a/i
; // File: ec79f9c27c045b00.js
b: while (1) { continue /*
*/ a; }
; // File: ec97990c2cc5e0e8.js
a || b && c | d ^ e & f == g < h >>> i + j * k
; // File: ecba8fb326c2c985.js
var a, b, c, d, e;
// compress these
if (b) {
    a = 1+2;
} else {
    a = 3;
}

if (b) {
    a = 4+5;
} else if (c) {
    a = 6;
} else {
    a = 7-8;
}

a = b ? 'f' : 'g'+'h';

a = b ? 'f' : b ? 'f' : 'g'+'h';

// Compress conditions that have side effects
if (i()) {
    a = 9+10;
} else {
    a = 11;
}

if (c) {
    a = 'j';
} else if (i()) {
    a = 'k'+'l';
} else {
    a = 'j';
}

a = i() ? 'm' : 'f'+'n';

// don't compress these
a = b ? d : e;

a = b ? 'f' : 'g';
; // File: ed0783c35e43032b.js
a["b"] = "c";
a["if"] = "if";
a["*"] = "d";
a["\u0EB3"] = "e";
a[""] = "f";
a["1_1"] = "b";
; // File: ed085cb2fd0dc355.js
var a;
if (!b && !c && !d && !e) {
    a = 1;
} else {
    a = 2;
}
; // File: ed49ee70d6eabf4a.js
arguments = 1
; // File: ed65dd575be2b4ab.js
// Don't apply transformation to global code
function a() {
    b.c('d');
}
function a() {
    b.c('e');
}
a();

; // File: ed6981438ac1918b.js
a |= 1
; // File: ed894bd570d47113.js
throw a * b
; // File: eda5026c194f7279.js
0x0;
; // File: edbdeeb1761675a7.js
({"[": 1})
; // File: edd1f39f90576180.js
try { } catch (a) { b(a) }
; // File: edfe04e832b81a82.js
new a(b)
; // File: ee2342b2715c3bf0.js
do { a++; b--; } while (a < 1)
; // File: eea2875eacf36279.js
if (a) {
    b = c();
    d = e();
    for (; b < d; ++b) f.g(b);
}

; // File: eebefa78eec0af44.js
(function () {
    -((a) ? b : 1)
}());

; // File: eed97872dd924560.js
`outer${{a: {b: 1}}}bar${`nested${function(){return 2;}}endnest`}end`
; // File: eef60d36274e4ed8.js
new (a, b)
new (a || b)
new (c ? a : b)
; // File: ef086346e9707e91.js
class a {[b]() {}}
; // File: ef15294c7bc4675e.js
function a() {
    // do not remove it
    var b = 1;
    ++b;
}

; // File: ef61944dbb440b60.js
class a extends class b extends c {} {}
; // File: ef7843986fabc25d.module.js
export default a;

; // File: ef812b85ce5fbc44.js
var Infinity, NaN;
Infinity.a();
NaN.a();
; // File: efb88a0b6e2e170e.js
`

`
; // File: efe1e5c7656bf0ba.js
a < b
; // File: efef19e06f58fdd9.js
({ __proto__: null, set __proto__(a){} })

; // File: f01d9f3c7b2b2717.js
a = { get if() {} }
; // File: f062a3f543a622f8.module.js
import a from "foo";

; // File: f0a5cf41bdef6532.js
do a(); while (true);
; // File: f0bf9ec665d85fa1.js
class a { b() { () => super.c; } }
; // File: f0d9a7a2f5d42210.js
let + 1
; // File: f0f2ab32e7f42314.js
for(var a of [1,2]) 3
; // File: f0f9e218a70eba5c.js
/[\w-\s]/;
; // File: f108a85d36ec9afc.js
({ *a() { yield super.a(); } })
; // File: f1218947a6a17e65.js
(function () {
    var a = 1;  // hoist this, but it is very difficult.
    (function () {
        eval('');
    }());
}());

; // File: f139fd88bd0ad9d0.js
`{${a}}`, `}`
; // File: f13a130829aa77c5.js
class a { static *b() {} }

; // File: f1534392279bddbf.js
00
; // File: f15772354efa5ecf.js
(function() {'use strict';return 1;});
; // File: f1643d0e6c7fde9a.js
var a = /=([^=\s])+/g
; // File: f17ec9517a3339d9.js
({ set if(a) { a } })
; // File: f1bf02f18fa71ba7.js
([])=>1;

; // File: f1d7e3cc86ffc02b.js
function a() {
    b();
    var c;
    var d;
}
function e(f) {
    b();
    var c;
    var f;
}
; // File: f2113065d9111e6d.js
(() => null)();
; // File: f2142c1dabd961c1.js
/* not comment*/; a-->1
; // File: f2aa3da994da03a7.js
var {a = b} = c
; // File: f2d394b74219a023.js
typeof /test/ + ' RegExp'

; // File: f2e0a415d88b3451.js
if (a) { // Some comment
b(); }
; // File: f2ed650f15f224fa.module.js
export default 1
; // File: f30d88a123e11b55.js
for (var [a, b] of c);

; // File: f3219596b50bb381.js
{[1]}
/foo/
; // File: f3260491590325af.js
(function(){ return true })() ? a.b(true) : a.b(false);
; // File: f355802cb6d444e1.js
let a;
; // File: f3d3a0f30115de54.js
a < b < c
; // File: f404f7ff29ba5d1a.module.js
export default function () {}
; // File: f407a3693faf595b.js
([a=1], [])=>2
; // File: f43f922cccf5b9af.js
(function () {
    a = a += 1;  // This is global varible, so observable by getter.
}());

; // File: f471327b3e9b8933.js
a = { a: 1, a: 2 }
; // File: f4864ec70dd99c21.js
a.false
; // File: f4a61fcdefebb9d4.js
var private, protected, public, static
; // File: f4b2d8937ec13ab0.js
for (var a = ("b" in c), d = 1; d < 2; ++d);

; // File: f50f858c3ef003f4.js
a && (b && c)

; // File: f552daf299e1c6e5.js
({a(){let a;}})
; // File: f597b0312e2b678c.js
if (a) {
    b();
} else {
}

; // File: f5ba9f1b21487d3b.js
function a() {
    var b = 1;
}

; // File: f601e7dd0235d423.js
var {a: b = c} = d
; // File: f69b27444afab042.js
(function () {
    // 'a'
    b: {
        if (a) break b;
        if (a) break b;
        if (a) break b;
        if (a) break b;
    }
    eval(c);
}());

; // File: f6d42525cd87339b.js
(function () {
    var a, b, c;
}());

; // File: f7291c5ec70a4152.js
({ a: b, c }, [d, e], ...f) => {}
; // File: f78abc3cba581cdd.js
class a {*b(c) { yield c; }}
; // File: f7af1a6b02dbd440.js
for (;;) {
  a;
}

; // File: f7e2edf1ccb61303.js
var [[a]]=1;
; // File: f80f30fbdd7e7b19.js
/.{.}/;
; // File: f8323b3c45bd107a.js
a ^ b ^ c
; // File: f89bf797c3b1dda4.js
(a)
; // File: f8a07bd5ab703d4b.js
for (;;) {
    if (a) {
        continue;
    }
    b();  // This should not be removed.
}

; // File: f8cf06a0d5699319.js
function a(b, ...c) {}
; // File: f8d843a30c73377a.js
(class a{})
; // File: f8dc2e8bbddcdfbe.js
a || b
; // File: f94e47b7b5cfda74.js
// should be 4
1 << 2

; // File: f96c694c5a2f2be9.js
function a(b, c, d, e) { return b < !--c && d-- > e; }
; // File: f974f2619b25b027.js
({ a: 1 })
; // File: f9888fa1a1e366e7.js
a = [ ,, 1 ]
; // File: f990e76e7fcb0dd9.js
const [a] = []
; // File: f9b92700d0e68f49.js
({a,} = 1)
; // File: f9c201250f225ab9.js
delete (1, a)

; // File: f9d67ab9db16c4d5.js
var a = 1;
; // File: fa58aa963031f8df.js
 /**

**/
; // File: fa59ac4c41d26c14.js
({let})
; // File: fa5b398eeef697a6.js
({set a(eval){}})
; // File: fa6c17d9a188d0bb.js
(function () {
    if (!a || b());
}());

; // File: fa736f4b0cf19c0c.js
"Hello\1World"
; // File: fa9eaf58f51d6926.js
(function(){})
; // File: faa4a026e1e86145.js
(() => {})()
; // File: fada2c7bbfabe14a.js
(function arguments() { });
; // File: fae42f5a2ab85c1d.js
a %= 1
; // File: fb50400b4c9cf740.js
({ *a() { yield; } })

; // File: fb69459d7628ace1.js
({ set: 1 })
; // File: fb7c5656640f6ec7.js
`${/\d/.a('1')[1]}`
; // File: fb8d437ce90b1178.js
[a] = 1;

; // File: fb8db7a71f3755fc.js
class a {set b(c) {}}
; // File: fba24e17d16fd0c4.js
"\x61"
; // File: fbacebe72fb15fed.module.js
export class a{};1
; // File: fbb6b30b41732026.js
new a
; // File: fbde237f11796df9.js
({ a: 1, set a(b) { } })
; // File: fc020c065098cbd5.js
var a = /[\u{61}-b][\u0061-b][a-\u{62}][a-\u0062]\u{1ffff}/u;

; // File: fc035551a2a4c15c.js
while (a-->1) {}
; // File: fc1ba7d289fb1af1.js
a = { b: function(c, ...d) {} }
; // File: fc286bf26373db8d.js
switch (a) {
  case 'b': c();
  default:
}
; // File: fc5c8d6f6bf16121.js
{ a(); b(); }
; // File: fc9f000aa3e4bd79.js
for(let a;;);
; // File: fcb318e400b44257.js
a ? 1 : 2
; // File: fcd33c00916dd6ad.js
a.b.c(1)
; // File: fcf3738a49a5f358.js
while (true) {
  /**
   * comments in empty block
   */
}

; // File: fd0ad9026eee596b.js
(1)
; // File: fd0e7b0f778f8a3b.js
if (a) {
  a;
}

; // File: fd29828f68a7634e.js
for (;a();) {
    if (b()) break;
    c();
    d();
}
; // File: fd34477284c96cbf.js
var [a, a] = 1;
; // File: fd5ea844fcc07d3d.js
a => { 1; }
; // File: fd889a4ef6e361f1.js
function a(b, c) { d(); }
; // File: fdb684acf63f6274.js
0B10
; // File: fe03ba1b818c762e.js
(function () {
    null != (a, 1);
}());

; // File: fe24fc72de1ef7cc.js
([a,b])=>1;

; // File: fe5ae04c8d239b26.js
function a() {
  // do not concat i=20,i2=30
  var b = 1;
  c();
  var d = 2;
}

; // File: fe5f0dcb8e902857.js
while (a)
  // optimize it
  (function () {
    b('c');
  }());
try { } catch (d) { b('e'); }

; // File: fe7c2a6e1efe2cf4.js
let[a] = b;
; // File: fec4c4ff229d3fc2.js
a(b() + 1 + "c" + "d");
a(b() + (2 + "c") + "d");
a((b() + 3) + "c" + "d");
a(b() + 4 + "c" + "d" + ("e" + "f"));
a("e" + "f" + b() + 5 + "c" + "d" + ("e" + "f"));
a("c" + b() + 6 + "d");
a(b() + 'e' + (7 + g('10')));
; // File: fee1cb654a489f02.js
a - b % c
; // File: fee3f54aa720263f.js
(a ? b : c) ? d : e
; // File: fef4facb0b8479bf.js
function a() {
  return 1 ? 2 : 3;
}

; // File: ff03d6d14c3f4007.js
(function () {
    var a = {
        'NaN': 1
    };
}());

; // File: ff215f966bed2b85.js
({"__proto__": 1 })
; // File: ff488aae349cc02d.js
for (const a of b);

; // File: ff4b8762733080cb.js
while (true) { continue }
; // File: ff902593b25092d1.js
/[P QR]/i
; // File: ffaf5b9d3140465b.js
let()
; // File: ffbba9592c03baa6.js
switch (a) { case 1: b(); break; }
; // File: ffc32056a146cc9b.js
for (a of b);

; // File: fffe7e78a7ce9f9a.js
/foobar/
