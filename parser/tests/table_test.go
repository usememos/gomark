package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/internal"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestTableParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "| header |\n| --- |\n| cell |\n",
			node: &ast.Table{
				Header: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{Content: "header"},
						},
					},
				},
				Delimiter: []string{"---"},
				Rows: [][]ast.Node{
					{
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "cell"},
							},
						},
					},
				},
			},
		},
		{
			text: "| **header1** | header2 |\n| --- | ---- |\n| cell1 | cell2 |\n| cell3 | cell4 |",
			node: &ast.Table{
				Header: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Bold{
								Symbol: "*",
								Children: []ast.Node{
									&ast.Text{Content: "header1"},
								},
							},
						},
					},
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{Content: "header2"},
						},
					},
				},
				Delimiter: []string{"---", "----"},
				Rows: [][]ast.Node{
					{
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "cell1"},
							},
						},
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "cell2"},
							},
						},
					},
					{
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "cell3"},
							},
						},
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "cell4"},
							},
						},
					},
				},
			},
		},
		{
			text: "|A|B|\n|:---|---:|\n|C|D|\n",
			node: &ast.Table{
				Header: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{Content: "A"},
						},
					},
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{Content: "B"},
						},
					},
				},
				Delimiter: []string{":---", "---:"},
				Rows: [][]ast.Node{
					{
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "C"},
							},
						},
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "D"},
							},
						},
					},
				},
			},
		},
		{
			text: "Header|Column\n---|---\nValue|Data\n",
			node: &ast.Table{
				Header: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{Content: "Header"},
						},
					},
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{Content: "Column"},
						},
					},
				},
				Delimiter: []string{"---", "---"},
				Rows: [][]ast.Node{
					{
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "Value"},
							},
						},
						&ast.Paragraph{
							Children: []ast.Node{
								&ast.Text{Content: "Data"},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := internal.NewTableParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
