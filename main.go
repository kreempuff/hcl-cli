package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	_ "github.com/hashicorp/hcl/hcl/token"
	"strings"
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
	buf := bytes.Buffer{}
	src := handleNode(astFile.Node)
	src = fmt.Sprintf("{%s}", src)
	err := json.Compact(&buf, []byte(src))
	return string(buf.Bytes()), err
}

func handleNode(node ast.Node) string {
	switch node.(type) {
	case *ast.ObjectItem:
		return handleObjectItem(node.(*ast.ObjectItem))
	case *ast.ObjectList:
		return handleObjectList(node.(*ast.ObjectList))
	case *ast.LiteralType:
		return handleLiteral(node.(*ast.LiteralType))
	case *ast.ObjectType:
		return handleObjectType(node.(*ast.ObjectType))
	case *ast.ListType:
		return handleListType(node.(*ast.ListType))
	}
	return ""
}

func handleListType(listType *ast.ListType) string {
	list := []string{}
	for _, item := range listType.List {
		list = append(list, handleNode(item))
	}

	return fmt.Sprintf("[%s]", strings.Join(list, ","))
}

func handleObjectType(objectType *ast.ObjectType) string {
	return fmt.Sprintf("{%s}", handleObjectList(objectType.List))
}
func handleObjectList(list *ast.ObjectList) string {
	s := []string{}
	for _, item := range list.Items {
		s = append(s, handleObjectItem(item))
	}
	return strings.Join(s, ",")
}

func handleObjectItem(item *ast.ObjectItem) string {
	//TODO handle duplicate keys and nested objects
	// nested objects are:
	// "free" "res" {
	// }
	// duplicate keys are
	// "free" "res" {
	// }
	// "free" "hey" {
	// }
	// These should both go under "free"

	for _, key := range item.Keys {
		return fmt.Sprintf("\"%s\":%s", key.Token.Text, handleNode(item.Val))
	}
	return "{}"
}

func handleLiteral(item *ast.LiteralType) string {
	switch item.Token.Type {
	case token.HEREDOC:
		var tokenIdentifer string
		firstNewline := strings.Index(item.Token.Text, "\n")
		tokenIdentifer = item.Token.Text[2:firstNewline]
		end := strings.LastIndex(item.Token.Text, tokenIdentifer)
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(item.Token.Text[firstNewline+1:end], "\n", "\\n"))
	}

	return item.Token.Text
}
