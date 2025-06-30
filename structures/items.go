package structures

// This implementation referred to the official lexer in Go https://go.dev/src/text/template/parse/lex.go

import "fmt"

type Pos int

// Item represents a token or text string returned from the scanner.
type Item struct {
	Typ ItemType // The type of this Item.
	Val string   // The value of this Item.
}

type IdentifierType int

const (
	idError IdentifierType = iota
	idUndefined
	idFunction
	idVariable
)

type Identifier struct {
	typ IdentifierType
	val string
}

func (i Item) ToString() string {
	switch {
	case i.Typ == ItemEOF:
		return "EOF"
	case i.Typ == ItemError:
		return i.Val
	case i.Typ > ItemKeyword:
		return fmt.Sprintf("<%s>", i.Val)
	case len(i.Val) > 10:
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

// ItemType identifies the type of lex items.
type ItemType int

const (
	ItemError     ItemType = iota // error occurred; value is text of error
	ItemNumeric                   // arithmetic symbols
	ItemColon                     // colon keyword
	ItemSemicolon                 // semicolon (optional)
	ItemComment                   // comment text
	ItemCompare                   // compare symbols
	ItemAssign                    // equals ('=') introducing an assignment
	ItemDot                       // dot
	ItemEOF
	ItemLogical    // and '&&', or '||', not '!'
	ItemField      // alphanumeric identifier starting with '.'
	ItemIdentifier // alphanumeric identifier not starting with '.'
	ItemLeftDelim  // left action delimiter
	ItemLeftParen  // '(' inside action
	ItemNumber     // simple number, including imaginary
	ItemRawString  // raw quoted string (includes quotes)
	ItemRightDelim // right action delimiter
	ItemRightParen // ')' inside action
	ItemSpace      // run of spaces separating arguments

	// literals appear after
	ItemLiteral
	ItemBoolean      // boolean constant
	ItemInteger      // integer numbers
	ItemFloat        // floating point numbers
	ItemChar         // printable ASCII character; grab bag for comma etc.
	ItemCharConstant // escape sequences e.g. '\n'
	ItemString       // quoted string (includes quotes)
	ItemText         // plain text
	ItemNull         // the untyped null constant, easiest to treat as a keyword

	// keywords appear after all the rest.
	ItemKeyword
	ItemBreak     // break keyword
	ItemClass     // class keyword
	ItemContinue  // continue keyword
	ItemDo        // do keyword
	ItemDefine    // define keyword
	ItemElse      // else keyword
	ItemFor       // for loop keyword
	ItemIf        // if keyword
	ItemInterface // interface keyword
	ItemPackage   // package keyword
	ItemRange     // range keyword 'in'
	ItemThrow     // throw keyword
	ItemFunction  // function keyword 'fun'
	ItemWhile     // while loop keyword
	ItemReturn    // return call keyword
	ItemSuper     // super keyword
	ItemVariable  // variable keyword 'var'
	ItemConstant  // constant variable keyword 'val'
	ItemTry       // try keyword
	ItemTypeof    // typeof keyword
)
