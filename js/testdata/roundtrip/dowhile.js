do foo(); while (cond);

let i = 0;
do {
  console.log(i);
  i++;
} while (i > 10);

// with comments
let j = 0;
// c1
do {
  console.log(i);
  j++;
} /*c2*/ while // c3
(j > 10 /*c4*/);
