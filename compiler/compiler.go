package compiler

import (
	"KotlinToJavaConverter/structures"
	"fmt"
)

// Node are pointers types to what
// would otherwise be recursive types in Go. e.g.
//
// callee Node
//
// Would cause the Go compiler to complain about a recursive type. When we want
// to use one of these types to pass through to a function, for example, we'd
// use `&` as it'd be a reference. But we'll come to that a bit later on.
type Node struct {
	kind       string
	value      string
	name       string
	callee     *Node
	expression *Node
	body       []Node
	params     []Node
	arguments  *[]Node
	context    *[]Node
}

// TreeAST is just another alias type. I find this makes part of the code
// more readable, as you'll come to see that there are a ton of references to
// `Node`.
type TreeAST Node

// Compile function takes a list of all kotlin files and analyze
// the correctness of file's content. The compiler consists of four steps:
// tokenization, parsing, transformation, converter
func Compile(files structures.DataFiles) error {
	javaFiles := structures.DataFiles{}
	for _, file := range files.Files {

		tokens, err := Tokenization(file)
		if err != nil {
			_ = fmt.Errorf("error tokenizing file %v", file.Name)
		}

		treeAST, err := Parser(tokens)
		if err != nil {
			_ = fmt.Errorf("error parsing file %v", file.Name)
		}

		newTreeAST, err := Transformation(treeAST)
		if err != nil {
			_ = fmt.Errorf("error transforming file %v", file.Name)
		}

		javaFile, err := Converter(newTreeAST)
		if err != nil {
			_ = fmt.Errorf("error converting file %v", file.Name)
		}
		javaFiles.Files = append(javaFiles.Files, javaFile)
	}
	return nil
}

//TODO: continue implement the "Compile" function with support services

func Tokenization(file structures.DataFile) ([]structures.Item, error) {
	return make([]structures.Item, 0), nil
}

func Parser(tokens []structures.Item) (TreeAST, error) {
	return TreeAST{}, nil
}

func Transformation(oldTreeAST TreeAST) (TreeAST, error) {
	return oldTreeAST, nil
}

func Converter(newTreeAST TreeAST) (structures.DataFile, error) {
	return structures.DataFile{}, nil
}
