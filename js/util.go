package js

import "github.com/xjslang/xjs/ast"

// a "self closing" statement doesn't need a semicolon at the end
func selfClosingStmt(node ast.Stmt) bool {
	switch v := node.(type) {
	case *SemiStmt, *FunctionDecl, *ForStmt:
		return true
	case *IfStmt:
		if _, ok := v.Else.(*BlockStmt); ok {
			return true
		}
		if v.Else == nil {
			if _, ok := v.Then.(*BlockStmt); ok {
				return true
			}
		}
	case *LabelStmt:
		return selfClosingStmt(v.Stmt)
	}
	return false
}
