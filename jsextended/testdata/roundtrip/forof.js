// const, let, var
for (const row of rows);
for (let row of rows);
for (var row of rows);

// iterate expressions
for (const row of rows.get()) {
  console.log(row);
}

// destructuring
for (let { a, b } of rows);
for (let [a, b] of rows);
for (let {} in rows);

// with indentation
for (
  const row
  of [1, 2, 3]
) {
  console.log(row);
}

// with comments
for /*c1*/ (
  //c2
  const row /*c3*/ of rows
//c4
);

// without declaring var
for (a[b in c] in d);
for (a(b in c)[1] in d);
for ({ a = 1 } in b);
