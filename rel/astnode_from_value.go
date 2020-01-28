package rel

import (
	"github.com/arr-ai/wbnf/ast"
	"github.com/arr-ai/wbnf/parser"
)

func ASTNodeFromValue(value Value) ast.Node {
	switch value := value.(type) {
	case Tuple:
		return ASTBranchFromValue(value)
	case Set:
		return ASTLeafFromValue(value)
	}
	panic("wat?")
}

func ASTLeafFromValue(l Set) ast.Leaf {
	return ast.Leaf(*parser.NewBareScanner(stringFromSet(l)))
}

func ASTBranchFromValue(b Value) ast.Branch {
	panic("unfinished")
	// 	result := ast.Branch{}
	// outer:
	// 	for i := value.Enumerator(); i.MoveNext(); {
	// 		name, value := i.Current()
	// 		var children ast.Children
	// 		switch name {
	// 		case "@choice":
	// 			values := value.(*genericSet).OrderedValues()
	// 			ints := make(ast.Many, 0, len(values))
	// 			for _, v := range values {
	// 				ints = append(ints, ast.Extra{Data: int(v.(Number).Float64())})
	// 			}
	// 			children = ints
	// 		case "@rule":
	// 			children = ast.One{Node: ast.Extra{Data: wbnf.Rule(value.(String).String())}}
	// 		case "@skip":
	// 			children = ast.One{Node: ast.Extra{Data: int(value.(Number).Float64())}}
	// 		default:
	// 			switch value.(type) {
	// 			case String:
	// 				// s := c.Scanner()
	// 				// value = NewOffsetString([]rune(s.String()), s.Offset())
	// 				// children = parser
	// 			case Tuple:
	// 				// value = ASTNodeToValue(c.Node)
	// 				// children = One() ASTNodeToValue(c.Node)
	// 			case Set:
	// 				// values := make([]Value, 0, c.Count())
	// 				// for e := c.Enumerator(); e.MoveNext(); {
	// 				// 	values = append(values, ASTNodeToValue(child))
	// 				// }
	// 				// value = NewArray(values...)
	// 			}
	// 		}
	// 		result[name] = children

	// 		switch value := value.(type) {
	// 		case Set:
	// 			if value.Bool() {
	// 				for j := value.Enumerator(); j.MoveNext(); {
	// 					if _, _, is := isStringTuple(j.Current()); !is {
	// 						// Not a string. Must be an array.
	// 						array := make(ast.Many, value.Count())
	// 						for j := value.Enumerator(); j.MoveNext(); {
	// 							index, item, _ := isArrayTuple(j.Current())
	// 							array[index] = ASTNodeFromValue(item)
	// 						}
	// 						result[name] = array
	// 						continue outer
	// 					}
	// 				}
	// 			}

	// 			// Not an array. Must be a string.
	// 			result[name] = ast.One{Node: ast.Leaf(*parser.NewBareScanner(stringFromSet(value)))}
	// 			// case ast.One:
	// 			// 	result = result.With(name, nodeToValue(c.Node))
	// 		}
	// 	}
	// 	return result
	// 	// case ast.Leaf:
	// 	// 	s := n.Scanner()
	// 	// 	return NewOffsetString([]rune(s.String()), s.Offset())
	// 	// case ast.Extra:
	// 	// 	switch e := n.Data.(type) {
	// 	// 	case int:
	// 	// 		return NewNumber(float64(e))
	// 	// 	case wbnf.Rule:
	// 	// 		return NewString([]rune(string(e)))
	// }
	// panic(fmt.Errorf("unhandled node: %v %[1]T", value))
}

func stringFromSet(set Set) (int, string) {
	s := set.(String)
	return s.offset, s.String()
}
