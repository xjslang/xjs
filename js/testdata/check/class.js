// extra semicolons
class A {b(){};c(){};}
class B { ; x=1; y(){} }
class C { x=1;; y(){} }
class D { ; }
class E { static {} ; }

(class {b(){};c(){};})
(class A { ; x=1; y(){} })
(class { x=1;; y(){} })
(class { ; })
(class { static {} ; })
