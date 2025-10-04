let person = {
    name: 'Alice',
    age: 25,
    greet: function() {
        return 'Hello, my name is ' + this.name;
    }
};

console.log(person.name);
console.log(person.age);
