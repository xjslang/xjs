function* foo(index) {
  while (index < 2) {
    yield index;
    index++;
  }
}

function* func1() {
  yield 42;
}

function* func2() {
  yield* func1();
}

function* a() {
  yield b = c, yield* d = e, f;
}