new;
aaa bbb // ; expected
x =; // unknown expression
{ aaa bbb } // ; expected
for; // ( expected
for (?; // unknown expression
for (let i = 0; i < 10; i++; // ) expected
function; // identifier expected
function foo; // ( expected
function foo(; // identifier expected
function foo(a,; // identifier expected
function foo(p0; // ) expected
if; // ( expected
if (; // unknown expression
let; // identifier expected
let x =; // unknown expression
[; // unknown expression
[1,; // unknown expression
[1; // ] expected
1+; // unknown expression
!; // unknown expression
print(; // unknown expression
print(1,; // unknown expression
print(1; // ) expected
(function; // ( expected
(function(; // identifier expected
(function(a,; // identifier expected
(function(a; // ) expected
(100; // ) expected
({; // key expected
({name; // : expected
({name:; // unknown expression
({name: 100,; // key expected
({name: 100; // } expected
1x123; // invalid number
2O123; // invalid number
0X; // invalid number
0o; // invalid number
a.100; // ; expected
let if = 100; // unknown expression
.1.2;
