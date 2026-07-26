// const, let, var
for (const i in rows);
for (let i in rows);
for (var i in rows);

// destructuring
for (let [a] in obj);
for (let { a } in obj);

// iterate expressions
for (const i in rows.get()) {
  console.log(i);
}

// with indentation
for (
  const i
  in [1, 2, 3]
) {
  console.log(i);
}

// with comments
for /*c1*/ (
  //c2
  const i /*c3*/ in rows
//c4
);