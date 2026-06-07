let x = function () {};
let y = function () {
  console.log("hello!");
};
let z = 2 * (function () {
  return Math.PI;
})();

// IIFE
(function () {
  console.log('init...');
})();