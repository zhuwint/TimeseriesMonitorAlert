package lambda

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

func Eval(expr ast.Expr, data map[string]interface{}) interface{} {

	switch expr := expr.(type) {
	case *ast.BasicLit: // 匹配到数据
		return getlitValue(expr)
	case *ast.BinaryExpr: // 匹配到子树
		// 后序遍历
		x := Eval(expr.X, data) // 左子树结果
		y := Eval(expr.Y, data) // 右子树结果
		if x == nil || y == nil {
			return errors.New(fmt.Sprintf("%+v, %+v is nil", x, y))
		}
		op := expr.Op // 运算符

		// 按照不同类型执行运算
		switch x.(type) {
		case int64:
			return calculateForInt(x, y, op)
		case bool:
			return calculateForBool(x, y, op)
		case string:
			return calculateForString(x, y, op)
		case error:
			return errors.New(fmt.Sprintf("%+v %+v %+v eval failed", x, op, y))
		default:
			return errors.New(fmt.Sprintf("%+v op is not support", op))
		}
	case *ast.CallExpr: // 匹配到函数
		return calculateForFunc(expr.Fun.(*ast.Ident).Name, expr.Args, data)
	case *ast.ParenExpr: // 匹配到括号
		return Eval(expr.X, data)
	case *ast.Ident: // 匹配到变量
		fmt.Println("ident", expr.Name)
		return data[expr.Name]
	default:
		return errors.New(fmt.Sprintf("%x type is not support", expr))
	}
}

func getlitValue(expr *ast.BasicLit) interface{} {
	fmt.Println("getlitValue", expr.Value)
	return nil
}

func calculateForInt(x, y interface{}, op token.Token) int64 {
	fmt.Println("calculateForInt", x, y, op)
	return 0
}

func calculateForBool(x, y interface{}, op token.Token) bool {
	fmt.Println("calculateForBool", x, y, op)
	return true
}

func calculateForString(x, y interface{}, op token.Token) string {
	fmt.Println("calculateForString", x, y, op)
	return ""
}

func calculateForFunc(name string, args []ast.Expr, data map[string]interface{}) interface{} {
	return nil
}
