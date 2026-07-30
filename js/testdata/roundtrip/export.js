export function log(msg) {
  console.log(msg);
}
export let pi = 3.14159;
export { _pi as pi, e };

// with comments
export let e = 2.71828;
export //c
{
  _pi /*c*/ as pi,
  b, c
// c
} /*c*/;
