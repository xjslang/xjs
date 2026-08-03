aaa bbb // ; expected
x =; // expression expected
{ aaa bbb } // ; expected
for; // ( expected
for (?; // expression expected
for (let i = 0; i < 10; i++; // ) expected
function; // identifier expected
function foo; // ( expected
function foo(; // identifier expected
function foo(a,; // identifier expected
function foo(p0; // ) expected
if; // ( expected
if (; // expression expected
let; // identifier expected
let x =; // expression expected
[; // expression expected
[1,; // expression expected
[1; // ] expected
1+; // expression expected
!; // expression expected
print(; // expression expected
print(1,; // expression expected
print(1; // ) expected
(function; // ( expected
(function(; // identifier expected
(function(a,; // identifier expected
(function(a; // ) expected
(100; // ) expected
({; // key expected
({name; // : expected
({name:; // expression expected
({name: 100,; // key expected
({name: 100; // } expected
1x123; // ; expected
2O123; // ; expected
0X; // expression expected
0o; // expression expected
a.100; // ; expected
let if = 100; // identifier expected
.1.2;