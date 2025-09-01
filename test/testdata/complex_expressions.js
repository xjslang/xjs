function calculator() {
    return {
        add: function(a, b) { return a + b; },
        multiply: function(a, b) { return a * b; },
        subtract: function(a, b) { return a - b; }
    };
}

let calc = calculator();
console.log(calc.add(5, 3));
console.log(calc.multiply(4, 7));
console.log(calc.subtract(10, 4));

function outer(x) {
    function inner(y) {
        return x + y;
    }
    return inner;
}

let addTen = outer(10);
console.log(addTen(5));

function double(n) {
    return n * 2;
}

function square(n) {
    return n * n;
}

let result = double(square(3) + 1);
console.log(result);

function factorial(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}

console.log(factorial(5));
