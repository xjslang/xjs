a, console.log(a), i + 1;
(!a, b.c(), d);

let x = (1, 2, 3);
let y = ({ a }, b);
let z = (100, [a, b]);

// with indentations
// TODO: (printer) fix parenthesis indentation
let a = (
[a, b],
{ c },
200
);

// with comments
let b = /*c*/ (1, 2, 3 //c
);