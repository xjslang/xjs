new Foo;
new Foo();
new Foo(1, 2);
new function () {};

// these errors are syntactically valid
// although they are not semantically valid
new "error!";
new 123;
new null;
