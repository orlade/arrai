package syntax

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
)

var stdScopeOnce sync.Once
var stdScopeVar rel.Scope

func stdScope() rel.Scope {
	stdScopeOnce.Do(func() {
		stdScopeVar = rel.EmptyScope.
			With(".", rel.NewTuple(
				rel.NewAttr("math", rel.NewTuple(
					rel.NewAttr("pi", rel.NewNumber(math.Pi)),
					rel.NewAttr("e", rel.NewNumber(math.E)),
					newFloatFuncAttr("sin", math.Sin),
					newFloatFuncAttr("cos", math.Cos),
				)),
				rel.NewAttr("grammar", rel.NewTuple(
					rel.NewNativeFunctionAttr("parse", parseGrammar),
					rel.NewAttr("lang", rel.NewTuple(
						rel.NewAttr("arrai", rel.ASTNodeToValue(arraiParsers.Node().(ast.Node))),
						rel.NewAttr("wbnf", rel.ASTNodeToValue(wbnf.Core().Node().(ast.Node))),
					)),
				)),
				rel.NewAttr("fn", rel.NewTuple(
					rel.NewAttr("fix", parseLit(`(\f f(f))(\f \g \n g(f(f)(g))(n))`)),
					rel.NewAttr("fixt", parseLit(`(\f f(f))(\f \t t :> \g \n g(f(f)(t))(n))`)),
				)),
				rel.NewAttr("log", rel.NewTuple(
					rel.NewNativeFunctionAttr("print", func(value rel.Value) rel.Value {
						log.Print(value)
						return value
					}),
					createNestedFuncAttr("printf", 2, func(args ...rel.Value) rel.Value {
						format := args[0].(rel.String).String()
						strs := make([]interface{}, 0, args[1].(rel.Set).Count())
						for i, ok := args[1].(rel.Set).ArrayEnumerator(); ok && i.MoveNext(); {
							strs = append(strs, i.Current())
						}
						log.Printf(format, strs...)
						return args[1]
					}),
				)),
				stdArchive(),
				stdReflect(),
				stdRel(),
				stdStr(),
				stdOs(),
			))
	})
	return stdScopeVar
}

func createNestedFunc(name string, nArgs int, f func(...rel.Value) rel.Value, args ...rel.Value) rel.Value {
	if nArgs == 0 {
		return f(args...)
	}

	return rel.NewNativeFunction(name+strconv.Itoa(nArgs), func(parent rel.Value) rel.Value {
		return createNestedFunc(name, nArgs-1, f, append(args, parent)...)
	})
}

func createNestedFuncAttr(name string, nArgs int, f func(...rel.Value) rel.Value) rel.Attr {
	return rel.NewAttr(name, createNestedFunc(name, nArgs, f))
}

func parseLit(s string) rel.Value {
	v, err := MustCompile(NoPath, s).Eval(rel.EmptyScope)
	if err != nil {
		panic(err)
	}
	return v
}

func newFloatFuncAttr(name string, f func(float64) float64) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(value rel.Value) rel.Value {
		return rel.NewNumber(f(value.(rel.Number).Float64()))
	})
}

func parseGrammar(v rel.Value) rel.Value {
	astNode := rel.ASTNodeFromValue(v).(ast.Branch)
	g := wbnf.NewFromAst(astNode)
	parsers := g.Compile(astNode)
	return rel.NewNativeFunction("parse(<grammar>)", func(v rel.Value) rel.Value {
		rule := v.String()
		return rel.NewNativeFunction(fmt.Sprintf("parse(%s)", rule), func(v rel.Value) rel.Value {
			node, err := parsers.Parse(parser.Rule(rule), parser.NewScanner(v.String()))
			if err != nil {
				panic(err)
			}
			return rel.ASTNodeToValue(ast.FromParserNode(parsers.Grammar(), node))
		})
	})
}
