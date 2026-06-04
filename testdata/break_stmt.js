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
