package lambda

import (
	"go/ast"
	"go/parser"
	"testing"
)

func TestEval(t *testing.T) {
	type args struct {
		expr ast.Expr
		data map[string]interface{}
	}

	expr := `a > threshold_lower &&  a < threshold_upper`
	exprAst, err := parser.ParseExpr(expr)
	if err != nil {
		panic(err.Error())
	}

	data := map[string]interface{}{
		"a":               1,
		"b":               2,
		"c":               3,
		"threshold_lower": 0.0,
		"threshold_upper": 2.0,
	}

	Eval(exprAst, data)
}
