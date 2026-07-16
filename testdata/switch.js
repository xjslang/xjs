switch (val) {
}

switch (val) {
  case 1:
    console.log("case 1");
    console.log("case 1");
    console.log("case 1");
    break;
  case 2:
    console.log("case 2");
    console.log("case 2");
    break;
  case 3:
    console.log("case 3");
    break;
  case 4:
    break;
  case 5:
  default:
    console.log("default");
}

// default clause in the middle!
switch (val) {
  case 1:
    console.log("case 1");
    break;
  default:
    console.log("default");
    break;
  case 2:
    console.log("case 2");
    break;
}

// with comments
switch /*c1*/ (val /*c2*/) //c3
{
  //c4
  case 1 /*c5*/:
    console.log("case 1");
    console.log("case 1");
    console.log("case 1");
    break;
  case 2 //c6
  :
    console.log("case 2");
    console.log("case 2");
    break;
  case 3:
    console.log("case 3");
    break;
  case 4:
    break;
  case 5:
  //c7
  default //c8
  :
    console.log("default");
//c9
}
