import 'library';
import * as lib from 'library';
import lib from 'library';
import { c1, c2 as c3, c4 } from 'library';
import {} from 'library';

// with comments
import /*c1*/ * //c2
as lib /*c3*/ from 'library';
import // c1
{
  c1,
  c2 as c3,
  c4
// c2
} // c3
from 'library' // c4
;
