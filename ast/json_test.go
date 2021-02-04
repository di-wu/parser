package ast

import (
	"fmt"
	"github.com/di-wu/parser/op"
)

func ExampleNode_MarshalJSONString() {
	typeStrings := []string{"A", "NL"}
	p, _ := New([]byte("aaaaa\n"))
	n, _ := p.Expect(op.And{
		op.MinOne(Capture{
			Type:        0,
			TypeStrings: typeStrings,
			Value:       'a',
		}),
		Capture{
			Type:        1,
			TypeStrings: typeStrings,
			Value:       '\n',
		},
	})
	fmt.Println(n.MarshalJSONString())
	// Output:
	// [[0,"a"],[0,"a"],[0,"a"],[0,"a"],[0,"a"],[1,"\n"]] <nil>
}
