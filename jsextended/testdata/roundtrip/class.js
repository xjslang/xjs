class MyClass {
  constructor(a, b, c, ...rest) {
    console.log("init");
  }

  method1() {}
  method2() {}
}

// with extends
class MyClass extends BaseClass {
}
class foo extends (a, b) {
}
(class extends (a, b) {
});

// with properties
class MyClass {
  prop1 = 100;
  prop2;
  method1() {}
}

// computed members
class MyClass {
  'foo' = 100;
  100 = 200;
  ['boo'] = function() {
  };
  200 = () => {};
  [{ a }] = 200;
  static [expr] = () => {};
}

// with statics
class MyClass {
  prop1;
  static prop2 = 200;
  static static = "I'm static!";
  method1() {}
  static method2() {}
  static() {}
  static [expr] = () => {};
  static "name" = () => {};
  static *gen() {}
}

// with static initializers
class MyClass {
  prop1;
  static;
  static = 1;
  static prop2 = 200;
  static prop3;
  static static = "I'm static!";
  static {
    this.static = "I am static";
  }
  method1() {}
  static method2() {}
  static get foo() {}
  static() {}
  static {
    this.prop3 = "I'm static too";
  }
}

// with getters/setters
class MyClass {
  _msg = "";
  get msg() {
    return `msg: ${this._msg}`;
  }
  set msg(value) {
    this._msg = value;
  }
  get;
  set;
  get() {}
  set() {}
}

let c = (class {
  prototype() {}
});

(class foo {
  prototype() {}
});
