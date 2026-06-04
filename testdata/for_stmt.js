for (let i = 0; i < 10; i = i + 1) {
  print(i);
}

// indented version
for (
  let i = 0;
  i < 10;
  i = i + 1
) {
  print(i);
}

for (let i = 0; i < 10; i = i + 1) {
  break;
}

outer: for (let i = 0; i < 5; i = i + 1) {
  for (let j = 0; j < 5; j = j + 1) {
    if (i == 2 && j == 3) {
      break outer; // Exits the entire outer loop
    }
  }
}

label: print("statement is required after a label");

for (let i = 0; i < 10; i = i + 1) {
  if (i == 3) {
    continue;
  }
  text = text + i;
}

// The first for statement is labeled "loop1"
loop1: for (let i = 0; i < 3; i = i + 1) {
  // The second for statement is labeled "loop2"
  loop2: for (let j = 0; j < 3; j = j + 1) {
    if (i == 1 && j == 1) {
      continue loop1;
    }
  }
}
