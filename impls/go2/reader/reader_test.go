package reader

import (
	"log"
	"testing"
)

func TestTokenizer(t *testing.T) {
	teststr := `this is a test`
	res := tokenize(teststr)
	if len(res) != 4 {
		t.Errorf(`Length of tokenization of %v should be 4
		Tokenized output: %v | %v`, teststr, res, len(res))
	}
	teststr = `this is "another set of" tests`
	res = tokenize(teststr)
	if len(res) != 4 {
		t.Errorf(`Length of tokenization of %v should be 4
		Tokenized output: %v | %v`, teststr, res, len(res))
	}
	teststr = `this is "another set of" (tests)`
	res = tokenize(teststr)
	if len(res) != 6 {
		t.Errorf(`Length of tokenization of %v should be 6
		Tokenized output: %v | %v`, teststr, res, len(res))
	}
	teststr = ` this is "another set of" (tests)`
	res = tokenize(teststr)
	if len(res) != 6 {
		t.Errorf(`Length of tokenization of %v should be 6
		Tokenized output: %v | %v`, teststr, res, len(res))
	}
	teststr = ` this     is  "another set of"   (    tests)`
	res = tokenize(teststr)
	if res[0] != `this` {
		t.Errorf(`First element of tokenization of %v should be "this"
		Tokenized output: '%v'`, teststr, res[0])
	}
	if res[1] != `is` {
		t.Errorf(`First element of tokenization of %v should be "is"
		Tokenized output: '%v'`, teststr, res[1])
	}
	teststr = ` this is "another set of"wae (tests)`
	res = tokenize(teststr)
	if len(res) != 8 {
		t.Errorf(`Length of tokenization of %v should be 6
		Tokenized output: %v | %v`, teststr, res, len(res))
	}
}

func TestAdvance(t *testing.T) {
	testReader := Reader{Tokens: []string{"a", "b", "c"}, Position: 0, EOF: 3}
	Advance(&testReader)
	if testReader.Position != 1 {
		t.Errorf(`Position should be 1, but was "%v" instead!`, testReader.Position)
	}
}

func TestReadList(t *testing.T) {
	teststr := "(awef)"
	res, err := ReadStr(teststr)
	if err != nil {
		t.Errorf("ReadStr error: %v", err)
	}
	log.Printf("%v", res)
}

func TestBalanceQuotes(t *testing.T) {
	teststr := `awefawef`
	if !balancedQuotes(teststr) {
		t.Errorf("The string `%v` should be balanced!", teststr)
	}
	teststr = `"awefawef"`
	if !balancedQuotes(teststr) {
		t.Errorf("The string `%v` should be balanced!", teststr)
	}
	teststr = `"awefawef" "waef"`
	if !balancedQuotes(teststr) {
		t.Errorf("The string `%v` should be balanced!", teststr)
	}
	teststr = `"awefawef \"waef"`
	if !balancedQuotes(teststr) {
		t.Errorf("The string `%v` should be balanced!", teststr)
	}
	teststr = `"awefawef" \"waef"`
	if balancedQuotes(teststr) {
		t.Errorf("The string `%v` should not be balanced!", teststr)
	}
	teststr = `"`
	if balancedQuotes(teststr) {
		t.Errorf("The string `%v` should not be balanced!", teststr)
	}
}

func TestProperlyQuotedString(t *testing.T) {
	var teststr string
	teststr = `Testing`
	if properlyQuotedString(teststr) {
		t.Errorf("The string `%v` should not be properly quoted", teststr)
	}
	teststr = `"`
	if properlyQuotedString(teststr) {
		t.Errorf("The string `%v` should not be properly quoted", teststr)
	}
	teststr = `"Testing"`
	if !properlyQuotedString(teststr) {
		t.Errorf("The string `%v` should be properly quoted", teststr)
	}
	teststr = `"Testing`
	if properlyQuotedString(teststr) {
		t.Errorf("The string `%v` should not be properly quoted", teststr)
	}
	teststr = `"Testing\""`
	if !properlyQuotedString(teststr) {
		t.Errorf("The string `%v` should be properly quoted", teststr)
	}
	teststr = `"Testing"test`
	if properlyQuotedString(teststr) {
		t.Errorf("The string `%v` should not be properly quoted", teststr)
	}
	teststr = `"\n"`
	if !properlyQuotedString(teststr) {
		t.Errorf("The string `%v` should be properly quoted", teststr)
	}
}
