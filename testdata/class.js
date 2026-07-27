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

// with properties
class MyClass {
  prop1 = 100;
  prop2;
  method1() {}
}

// with statics
class MyClass {
  prop1;
  static prop2 = 200;
  static static = "I'm static!";
  method1() {}
  static method2() {}
  static() {}
}

// with static initializers
class MyClass {
  prop1;
  static prop2 = 200;
  static prop3;
  static static = "I'm static!";
  static {
    this.static = "I am static";
  }
  method1() {}
  static method2() {}
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