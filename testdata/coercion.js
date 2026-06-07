let num = 5;
let str = "5";
let bool = true;
let zero = 0;
let emptyStr = "";
let nullVal = null;

console.log("Number vs String:");
console.log(num == str);

console.log("Boolean vs Number:");
console.log(bool == 1);
console.log(bool == zero);

console.log("Null comparisons:");
console.log(nullVal == zero);
console.log(nullVal == emptyStr);

console.log("Same type comparisons:");
console.log(num == 5);
console.log(str == "5");
console.log(bool == true);
console.log(num != str);
console.log(num != 10);
