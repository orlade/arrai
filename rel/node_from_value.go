package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
)

func ASTNodeFromValue(v Value, attr string) ast.Node {
	switch v := v.(type) {
	case Tuple:
		result := ast.Branch{}

	outer:
		for i := v.Enumerator(); i.MoveNext(); {
			name, value := i.Current()
			var children ast.Children
			switch name {
			case "@choice":
				values := value.(*genericSet).OrderedValues()
				ints := make(ast.Many, 0, len(values))
				for _, v := range values {
					ints = append(ints, ast.Extra{Data: int(v.(Number).Float64())})
				}
				children = ints
			case "@rule":
				children = ast.One{Node: ast.Extra{Data: wbnf.Rule(value.(String).String())}}
			case "@skip":
				children = ast.One{Node: ast.Extra{Data: int(value.(Number).Float64())}}
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
				}
			}
			result[name] = children

			switch value := value.(type) {
			case Set:
				if value.Bool() {
					for j := value.Enumerator(); j.MoveNext(); {
						if _, _, is := isStringTuple(j.Current()); !is {
							// Not a string. Must be an array.
							array := make(ast.Many, value.Count())
							for j := value.Enumerator(); j.MoveNext(); {
								index, item, _ := isArrayTuple(j.Current())
								array[index] = ASTNodeFromValue(item, name)
							}
							result[name] = array
							continue outer
						}
					}
				}

				// Not an array. Must be a string.
				result[name] = ast.One{Node: ast.Leaf(*parser.NewBareScanner(stringFromSet(value)))}
				// case ast.One:
				// 	result = result.With(name, nodeToValue(c.Node))
			}
		}
		return result
		// case ast.Leaf:
		// 	s := n.Scanner()
		// 	return NewOffsetString([]rune(s.String()), s.Offset())
		// case ast.Extra:
		// 	switch e := n.Data.(type) {
		// 	case int:
		// 		return NewNumber(float64(e))
		// 	case wbnf.Rule:
		// 		return NewString([]rune(string(e)))
	}
	panic(fmt.Errorf("unhandled node: %v %[1]T", v))
}

func stringFromSet(set Set) (int, string) {
	s := set.(String)
	return s.offset, s.String()
}
