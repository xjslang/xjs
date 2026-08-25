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
  let name = "John Smith";
  let age = 32;
  return {
    name: name,
    age: age // age
  // comments here
  };
}

let x = {
  'name': John,
  ['age']: 32,
  3.14: 'PI approx.'
};

function foo() {
  let name = "John Smith";
  let age = 32;
  return {
    name,
    age // age
  // comments here
  };
}

let x = {
  addRow: () => {},
  name: John,
  ["age"]: 32,
  3.14: "PI approx."
};

// with accessors
let y = {
  _msg: '',
  get msg() {
    return this._msg;
  },
  prop,
  set msg(val) {
    this._msg = val;
  },
  log(...msgs) {
    console.log(msgs);
  },
  get(...args) {},
  set(...args) {},
  *gen() {},
  get: 'get',
  set: 'set',
  default: 'default',
  case: 'case',
  new: 'new'
};
