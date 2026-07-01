let matrix = [
  [1, 2],
  [3, 4],
  [5, 6]
];
console.log(matrix[0][0]);
console.log(matrix[1][1]);
console.log(matrix[2][0]);

let person = {
  if: 'reserved key', // reserved keys are allowed
  name: "John",
  address: {
    street: "123 Main St",
    city: "Anytown"
  },
  hobbies: ["reading", "coding", "gaming"],
  // numbers
  100: 'number 1',
  1.5: 'number 2',
  1e2: 'number 3',
  // strings
  'hi there!': 'string',
  // bracket expressions
  ['hi there!']: 'aaa',
  [100 + 1]: 'expression',
  [
  200 // comments are preserved
  ]: 'expression'
};

console.log(person.if); // reserved keys are allowed
console.log(person.name);
console.log(person.address.street);
console.log(person.address.city);
console.log(person.hobbies[0]);
console.log(person.hobbies[2]);

let key = "name";
console.log(person[key]);

let index = 1;
console.log(person.hobbies[index]);

foo.bar;
a.b.c;
a.
b;
