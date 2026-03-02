package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "check log messages for special rules and safety",
	Run:  run,
}

func calledFunc(pass *analysis.Pass, call *ast.CallExpr) (*types.Func, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	// метод: logger.Info()
	if s := pass.TypesInfo.Selections[sel]; s != nil {
		if f, ok := s.Obj().(*types.Func); ok {
			return f, true
		}
		return nil, false
	}

	// функция: slog.Info()
	if f, ok := pass.TypesInfo.Uses[sel.Sel].(*types.Func); ok {
		return f, true
	}

	return nil, false
}

func isSlogOrZapLogMethod(fn *types.Func) bool {
	if fn == nil || fn.Pkg() == nil {
		return false
	}
	pkgPath := fn.Pkg().Path()
	name := fn.Name()

	isLevel := name == "Info" || name == "Error" || name == "Debug" || name == "Warn"
	if !isLevel {
		return false
	}

	return pkgPath == "log/slog" || pkgPath == "go.uber.org/zap"
}

func startsWithLower(msg string) bool {
	var first rune
	trim := strings.TrimLeftFunc(msg, unicode.IsSpace)
	if trim == "" {
		return true
	}
	for _, r := range trim {
		first = r
		break
	}
	if unicode.IsLetter(first) && unicode.IsUpper(first) {
		return false
	}
	return true
}

func isOnlyEngLetters(msg string) bool {
	for _, r := range msg {
		if unicode.In(r, unicode.Cyrillic) {
			return false
		}
	}
	return true
}

func isNotSpecSymbols(msg string) bool {
	for _, r := range msg {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r)) {
			return false
		}
	}
	return true
}

func isNotPrivateInfo(msg string) bool {
	s := strings.ToLower(msg)

	keys := []string{"password", "pass", "token", "api_key", "apikey", "apikey", "private_key", "private"}
	adds := []string{":", ": ", " =", "= ", "=", " is "}

	for _, key := range keys {
		for _, add := range adds {
			if strings.Contains(s, key+add) {
				return false
			}
		}

		if strings.Contains(s, key) && strings.Contains(s, "%") {
			return false
		}
	}

	return true
}

func extractConstantStringExpr(expr ast.Expr) (text string, full bool, ok bool) {
	switch lit := expr.(type) {
	case *ast.BasicLit:
		if lit.Kind != token.STRING {
			return "", false, false
		}
		msg, err := strconv.Unquote(lit.Value)
		if err != nil {
			return "", false, false
		}
		return msg, true, true
	case *ast.BinaryExpr:
		if lit.Op != token.ADD {
			return "", false, false
		}

		ltext, lfull, lok := extractConstantStringExpr(lit.X)
		rtext, rfull, rok := extractConstantStringExpr(lit.Y)

		if !lok && !rok {
			return "", false, false
		}

		text := ltext + rtext

		full := lok && rok && lfull && rfull

		return text, full, true
	default:
		return "", false, false
	}

}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			//fmt.Printf("in call: ", call)
			if !ok {
				return true
			}

			//var name string
			//var typ string
			//var pkgOrVar string
			//var msg string

			// switch fun := call.Fun.(type) {
			// case *ast.SelectorExpr:
			// 	name = fun.Sel.Name
			// 	id, ok := fun.X.(*ast.Ident)
			// 	if !ok {
			// 		return true
			// 	}
			// 	pkgOrVar = id.Name
			// case *ast.Ident:
			// 	name = fun.Name
			// }

			if len(call.Args) == 0 {
				return true
			}

			fn, ok := calledFunc(pass, call)
			if !ok || fn == nil {
				return true
			}

			if !isSlogOrZapLogMethod(fn) {
				return true
			}

			msg, full, ok := extractConstantStringExpr(call.Args[0])
			if !ok {
				return true
			}

			//pass.Reportf(call.Args[0].Pos(), "DEBUG msg=%q full=%v", msg, full)

			if full {
				if !startsWithLower(msg) {
					pass.Reportf(call.Args[0].Pos(), "log message must start with a lowercase letter")
				}

				if !isOnlyEngLetters(msg) {
					pass.Reportf(call.Args[0].Pos(), "log message must be in English")
				}

				if !isNotSpecSymbols(msg) {
					pass.Reportf(call.Args[0].Pos(), "log message must not contain special symbols or emoji")
				}
			}

			if !isNotPrivateInfo(msg) {
				pass.Reportf(call.Args[0].Pos(), "log message must not contain private info")
			}

			//fmt.Printf("call: %s, text: %s, typ: %s\n", pkgOrVar+"."+name, msg, typ)

			_ = call

			//pass.Reportf(call.Pos(), "call: ", call, "name: ", name, "type: ", typ)

			return true
		})
	}
	return nil, nil
}
