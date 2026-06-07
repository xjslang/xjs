for (let i = 0; i < 10; i++) console.log(i);
for (let i = 0; i < 10; i++) {
  console.log(i);
}

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
