let person = {
  name: "Alice",
  age: 25,
  greet: function () {
    return "Hello, my name is " + this.name;
  }
};

console.log(person.name);
console.log(person.age);

let obj = {};
let entry = { name: "John Smith", age: 32 };
let item = {
  name: "John Smith",
  age: 32,
  status: "divorced"
};

function foo() {
  return {
    name: 'John Smith',
    age: 32 // age
  // comments here
  };
  let item = {
    name: "John Smith",
    age: 32,
    status: "divorced"
  };
}
