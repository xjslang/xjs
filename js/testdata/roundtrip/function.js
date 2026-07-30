let x = function () {};
let y = function () {
  console.log("hello!");
};
let z = 2 * (function () {
  return Math.PI;
})();
let int = setInterval(function () {
  console.log('tick!');
}, 1000);

// IIFE
(function () {
  console.log('init...');
})();

let multiply = function (a, b) {
  return a * b;
};

let result1 = multiply(4, 5);
console.log(result1);

function createAdder(x) {
  return function (y) {
    return x + y;
  };
}

let addFive = createAdder(5);
let result2 = addFive(3);
console.log(result2);

function sayHello() {
  return "Hello from function!";
}

console.log(sayHello());

function foo() {}

function boo() {
  let x = 100;
  let y = 200;
}

function printText(txt) {
  console.log(txt);
}

function printPosition(
  x, // x-coord
  y // y-coord
) {
  console.log(x, y);
}

function add(
  a,
  b
) {
  return a + b;
}

let result = add(10, 20);
console.log(result);

function foo() {
  console.log("calling foo");
  return;
}

foo();
