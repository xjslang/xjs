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
  get: 'get',
  set: 'set',
  default: 'default',
  case: 'case',
  new: 'new'
};