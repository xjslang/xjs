// program stmt
aaa bbb // ; expected

// assign stmt
x =; // expression expected

// block stmt
{ aaa bbb } // ; expected

// for stmt
for; // ( expected
for (?; // expression expected
for (let i = 0; i < 10; i++; // ) expected

// function stmt
function; // identifier expected
function foo; // ( expected
function foo(; // identifier expected
function foo(a,; // identifier expected
function foo(p0; // ) expected

// if stmt
if; // ( expected
if (; // expression expected

// let stmt
let; // identifier expected
let x =; // expression expected

// arr expr
[; // ] expected
[1,; // expression expected
[1; // ] expected

// binary expr
1+; // expression expected

// unary expr
!; // expression expected

// call expr
print(; // expression expected
print(1,; // expression expected
print(1; // ) expected

// function expr
(function; // ( expected
(function(; // identifier expected
(function(a,; // identifier expected
(function(a; // ) expected

// group expr
(100; // ) expected

// obj expr
({; // key expected
({name; // : expected
({name:; // expression expected
({name: 100,; // key expected
({name: 100; // } expected

// numbers
.123; // expression expected (numbers cannot start with '.')
1x123; // ; expected (invalid hex)
2O123; // ; expected (invalid octal)
0X; // expression expected (incomplete hex)
0o; // expression expected (incomplete octal)

// member expr
a.100; // key expected
a.(b); // key expected
a.(b + c); // key expected

// reserved keys cannot be used as identifiers
let if = 100; // identifier expected