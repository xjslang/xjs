// params destructuring
function foo({ a }, [b], c) {}

// params with default values
function foo(a = 1 + 2, b = {}, c) {}

// rest params
function foo(a, b, ...rest) {}
function foo(a, b, ...[c, d]) {}
function foo(a, b, ...{ c, d }) {}

// expressions
(function*() {});
let x = function() {};

f(a ? b : c);