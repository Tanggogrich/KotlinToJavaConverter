package structures

// This implementation referred to the official lexer in Go https://go.dev/src/text/template/parse/lex.go

import "fmt"

type Pos int

// Item represents a token or text string returned from the scanner.
type Item struct {
	Typ  ItemType // The type of this Item.
	Pos  Pos      // The starting position, in bytes, of this Item in the input string.
	Val  string   // The value of this Item.
	Line int      // The line number at the start of this Item.
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
	ItemError        ItemType = iota // error occurred; value is text of error
	ItemBool                         // boolean constant
	ItemChar                         // printable ASCII character; grab bag for comma etc.
	ItemCharConstant                 // character constant
	ItemComment                      // comment text
	ItemComplex                      // complex constant (1+2i); imaginary is just a number
	ItemAssign                       // equals ('=') introducing an assignment
	ItemDeclare                      // colon-equals (':=') introducing a declaration
	ItemEOF
	ItemField      // alphanumeric identifier starting with '.'
	ItemIdentifier // alphanumeric identifier not starting with '.'
	ItemLeftDelim  // left action delimiter
	ItemLeftParen  // '(' inside action
	ItemNumber     // simple number, including imaginary
	ItemPipe       // pipe symbol
	ItemRawString  // raw quoted string (includes quotes)
	ItemRightDelim // right action delimiter
	ItemRightParen // ')' inside action
	ItemSpace      // run of spaces separating arguments
	ItemString     // quoted string (includes quotes)
	ItemText       // plain text
	ItemVariable   // variable starting with '$', such as '$' or  '$1' or '$hello'
	// keywords appear after all the rest.
	ItemKeyword  // used only to delimit the keywords
	ItemBlock    // block keyword
	ItemBreak    // break keyword
	ItemContinue // continue keyword
	ItemDot      // the cursor, spelled '.'
	ItemDefine   // define keyword
	ItemElse     // else keyword
	ItemEnd      // end keyword
	ItemIf       // if keyword
	ItemNil      // the untyped nil constant, easiest to treat as a keyword
	ItemRange    // range keyword
	ItemTemplate // template keyword
	ItemWith     // with keyword
)
