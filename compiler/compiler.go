package compiler

import (
	st "KotlinToJavaConverter/structures"
	"fmt"
	"strings"
	"unicode/utf8"
)

const eof = -1 // end of the file marker

// Lexer holds the state of the scanner.
type Lexer struct {
	Name  string       // used only for error reports
	Input string       // the string being scanned
	Start int          // start pos of this Item
	Pos   int          // current pos in input
	Width int          // width of last rune read
	Items chan st.Item // channel of scanned Item
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// Compile function takes a list of all kotlin files and analyze
// the correctness of file's content. The compiler consists of four steps:
// tokenization, parsing, transformation, converter
func Compile(files st.DataFiles) error {
	javaFiles := st.DataFiles{}
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

// Tokenization reads the file's content and translates it into a slice of Items.
func Tokenization(file st.DataFile) ([]st.Item, error) {
	l := &Lexer{
		Input: string(file.Content),
		Items: make(chan st.Item),
	}
	go l.Run() // Run the lexer in a goroutine

	var tokens []st.Item
	for {
		item := <-l.Items
		tokens = append(tokens, item)
		if item.Typ == st.ItemEOF || item.Typ == st.ItemError {
			break
		}
	}

	if tokens[len(tokens)-1].Typ != st.ItemEOF {
		return nil, fmt.Errorf("lexing error: %s at position %d", tokens[len(tokens)-1].Val, tokens[len(tokens)-1].Pos)
	}
	return tokens, nil
}

type TreeAST struct {
}

func Parser(tokens []st.Item) (TreeAST, error) {
	return TreeAST{}, nil
}

func Transformation(oldTreeAST TreeAST) (TreeAST, error) {
	return oldTreeAST, nil
}

func Converter(newTreeAST TreeAST) (st.DataFile, error) {
	return st.DataFile{}, nil
}

// Lexer support functions //

// Run lexes the input by executing state functions until the state is nil.
func (l *Lexer) Run() {
	for state := LexText; state != nil; {
		state = state(l)
	}
	close(l.Items) // No more tokens will be delivered
}

// Emit passes an item back to the client.
func (l *Lexer) Emit(t st.ItemType) {
	l.Items <- st.Item{Typ: t, Val: l.Input[l.Start:l.Pos], Pos: l.Start}
	l.Start = l.Pos
}

func (l *Lexer) Next() rune {
	if l.Pos >= len(l.Input) {
		l.Width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Width += w
	l.Pos += l.Width
	return r
}

func LexText(l *Lexer) stateFn {
	for {
		if strings.HasPrefix(l.Input[l.Pos:], "{") {
			if l.Pos > l.Start {
				l.Emit(st.ItemText)
			}
			//return lexLeftMeta
		}
	}
}
