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


// with default
export default 3.1419;
export default { a: 1 };
export default function a() {};
export { default } from "foo";
export { a as default, b } from "foo";

// re-export
export { a as b, c } from "foo";
export {} from "a";
