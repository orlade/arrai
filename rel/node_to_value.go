package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/wbnf"
)

func ASTNodeToValue(n ast.Node) Value {
	switch n := n.(type) {
	case ast.Branch:
		result := EmptyTuple

		for name, children := range n {
			var value Value
			switch name {
			case "@choice":
				ints := children.(ast.Many)
				values := make([]Value, 0, len(ints))
				for _, i := range ints {
					values = append(values, NewNumber(float64(i.(ast.Extra).Data.(int))))
				}
				value = NewArray(values...)
			case "@rule":
				value = NewString([]rune(string(children.(ast.One).Node.(ast.Extra).Data.(wbnf.Rule))))
			case "@skip":
				value = NewNumber(float64(children.(ast.One).Node.(ast.Extra).Data.(int)))
			default:
				switch c := children.(type) {
				case ast.One:
					value = ASTNodeToValue(c.Node)
				case ast.Many:
					values := make([]Value, 0, len(c))
					for _, child := range c {
						values = append(values, ASTNodeToValue(child))
					}
					value = NewArray(values...)
				default:
					panic("wat?")
				}
			}
			result = result.With(name, value)
		}

		return result
	case ast.Leaf:
		s := n.Scanner()
		return NewOffsetString([]rune(s.String()), s.Offset())
	case ast.Extra:
		panic(fmt.Errorf("unexpected extra: %v", n))
	}
	panic(fmt.Errorf("unhandled node: %v %[1]T", n))
}