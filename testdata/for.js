for (let i = 0; i < 10; i++) console.log(i);
for (let i = 0; i < 10; i++) {
  console.log(i);
}

let i = 0;
for (i = 0; i < 10; i++) {
  console.log(i);
}

// with comments before semicolons
for (
  let i = 0; // init clause
  i < 10; // cond clause
  i++ // after clause
) {
  console.log(i);
}

// omit init clause
let i = 0;
for (; i < 10; i++) {
  console.log(i);
}

// omit cond clause
for (let i = 0;; i++) {
  if (i >= 10) {
    break;
  }
  console.log(i);
}

// omit after clause
for (let i = 0; i < 10;) {
  console.log(i);
  i++;
}

// omit all
for (;;);

// indented version
for (
  let i = 0;
  i < 10;
  i++
) {
  console.log(i);
}

for (let i = 0; i < 10; i++) {
  break;
}

outer: for (let i = 0; i < 5; i++) {
  for (let j = 5; j >= 0; j--) {
    if (i == 2 && j == 3) {
      break outer; // Exits the entire outer loop
    }
  }
}

label: console.log("statement is required after a label");

for (let i = 0; i < 10; i++) {
  if (i == 3) {
    continue;
  }
  text = text + i;
}

// The first for statement is labeled "loop1"
loop1: for (let i = 0; i < 3; i++) {
  // The second for statement is labeled "loop2"
  loop2: for (let j = 3; j >= 0; j--) {
    if (i == 1 && j == 1) {
      continue loop1;
    }
  }
}

// For loop
for (let j = 0; j < 3; j++) {
  console.log("for: " + j);
}

// For loop with array
let numbers = [10, 20, 30];
for (let k = 0; k < 3; k++) {
  console.log("array[" + k + "] = " + numbers[k]);
}
