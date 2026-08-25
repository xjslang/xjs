// const, let, var
for (const i in rows);
for (let i in rows);
for (var i in rows);

// destructuring
for (let [a] in obj);
for (let { a } in obj);
for (let {} in rows);

// iterate expressions
for (const i in rows.get()) {
  console.log(i);
}

// with indentation
for (const i in [1, 2, 3]) {
  console.log(i);
}

// with comments
//c2
for ( /*c1*/const i /*c3*/ in rows) //c4
;

// without declaring var
for (a[b in c] in d);
for (a(b in c)[1] in d);
for ({ a = 1 } in b);
