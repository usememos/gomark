package main

import (
	"fmt"

	"github.com/usememos/gomark"
	"github.com/usememos/gomark/renderer"
)

func main() {
	nodes, err := gomark.Parse("Here is a #tag and some **emboldened** text.")
	if err != nil {
		panic(err)
	}

	result := renderer.NewHTMLRenderer().Render(nodes)

	fmt.Println(result)
}
