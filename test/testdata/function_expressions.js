let multiply = function(a, b) {
    return a * b;
};

let result1 = multiply(4, 5);
console.log(result1);

function createAdder(x) {
    return function(y) {
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
