package main

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

func main() {
	astFile, err := hcl.Parse(`
service {
	Hello = "green"
}

handle {
}
`)
	if err != nil {
		panic(err)
	}

	s, err := toJson(astFile)
	fmt.Println(s)
}

func toJson(astFile *ast.File) (string, error) {
	return "", nil
}

func handleObjectList(list *ast.ObjectList) {
	for _, item := range list.Items {
		printKeys(item)
	}
}

func printKeys(item *ast.ObjectItem) {
	for _, key := range item.Keys {
		fmt.Println(key.Token.Text)
	}
}
