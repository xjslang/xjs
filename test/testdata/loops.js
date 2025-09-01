let i = 0;
while (i < 3) {
    console.log("while: " + i);
    i = i + 1;
}

for (let j = 0; j < 3; j = j + 1) {
    console.log("for: " + j);
}

let numbers = [10, 20, 30];
for (let k = 0; k < 3; k = k + 1) {
    console.log("array[" + k + "] = " + numbers[k]);
}
