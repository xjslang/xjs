let matrix = [[1, 2], [3, 4], [5, 6]];
console.log(matrix[0][0]);
console.log(matrix[1][1]);
console.log(matrix[2][0]);

let person = {
    name: "John",
    address: {
        street: "123 Main St",
        city: "Anytown"
    },
    hobbies: ["reading", "coding", "gaming"]
};

console.log(person.name);
console.log(person.address.street);
console.log(person.address.city);
console.log(person.hobbies[0]);
console.log(person.hobbies[2]);

let key = "name";
console.log(person[key]);

let index = 1;
console.log(person.hobbies[index]);
