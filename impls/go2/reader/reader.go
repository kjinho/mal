package reader

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var tokenizer *regexp.Regexp
var strparse *regexp.Regexp

func init() {
	tokenizer = regexp.MustCompile("[\\s,]*(?P<capture>~@|[\\[\\]{}()'`~^@]|\"(?:\\\\.|[^\\\\\"])*\"?|;.*|[^\\s\\[\\]{}('\"`,;)]*)")
	strparse = regexp.MustCompile(`^"`)
}

// EOFError signals an EOF error
type EOFError struct{}

func (m *EOFError) Error() string {
	return "EOF"
}

// Reader encapsulates a stateful reader
type Reader struct {
	Tokens   []string
	Position int
	EOF      int
}

// MalType is the basic type from which all others flow
type MalType interface {
	isMalType()
	String() string
}

// MalList is a list structure
type MalList []MalType

func (MalList) isMalType() {}

// String returns a string
func (m MalList) String() string {
	length := len(m)
	res := make([]string, length)
	for idx, val := range m {
		res[idx] = val.String()
	}
	return "(" + strings.Join(res, " ") + ")"
}

// MalVec is the vector type
type MalVec []MalType

func (MalVec) isMalType() {}

// String returns a string
func (m MalVec) String() string {
	length := len(m)
	res := make([]string, length)
	for idx, val := range m {
		res[idx] = val.String()
	}
	return "[" + strings.Join(res, " ") + "]"
}

// MalHash is the hash type
type MalHash map[MalType]MalType

func (MalHash) isMalType() {}

// String returns a string
func (m MalHash) String() string {
	res := make([]string, 0, 10)
	for key, value := range m {
		res = append(res, key.String()+" "+value.String())
	}
	return "{" + strings.Join(res, " ") + "}"
}

// MalNumber is a number
type MalNumber int

func (MalNumber) isMalType() {}

// String returns a string
func (m MalNumber) String() string {
	return strconv.Itoa(int(m))
}

// MalString is a number
type MalString string

func (MalString) isMalType() {}

// String returns a string
func (m MalString) String() string {
	return string(m)
}

// MalSymbol is a number
type MalSymbol struct {
	name string
}

func (MalSymbol) isMalType() {}

// String returns a string
func (m MalSymbol) String() string {
	return m.name
}

// Next returns the next token and increments the position
func Next(r *Reader) (string, error) {
	res, err := Peek(r)
	if err == nil {
		Advance(r)
	}
	return res, err
}

// Peek returns the next token without incrementing the position
func Peek(r *Reader) (string, error) {
	if r.Position < r.EOF {
		return r.Tokens[r.Position], nil
	}
	return "", new(EOFError)
}

// Advance moves the position up one
func Advance(r *Reader) {
	r.Position = r.Position + 1
}

// ReadStr takes a string and returns a MalType
func ReadStr(s string) (MalType, error) {
	tokens := tokenize(s)
	r := Reader{Tokens: tokens, Position: 0, EOF: len(tokens)}
	res, err := readForm(&r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func tokenize(s string) []string {
	res := tokenizer.FindAllString(s, -1)
	for idx, val := range res {
		res[idx] = strings.TrimLeftFunc(val, func(r rune) bool {
			return unicode.IsSpace(r) || r == ','
		})
	}
	return res
}

func readForm(r *Reader) (MalType, error) {
	res, err := Peek(r)
	if err != nil {
		return nil, err
	}
	switch res {
	case "(":
		ml, err := readList(r, ")")
		return MalList(ml), err
	case "[":
		ml, err := readList(r, "]")
		return MalVec(ml), err
	case "{":
		return readHash(r)
	case "'":
		return readQuote(r, "quote")
	case "`":
		return readQuote(r, "quasiquote")
	case "~":
		return readQuote(r, "unquote")
	case "~@":
		return readQuote(r, "splice-unquote")
	case "@":
		return readQuote(r, "deref")
	case "^":
		return readMetaData(r)
	default:
		return readAtom(r)
	}
}

func readMetaData(r *Reader) (MalList, error) {
	Advance(r) // advance past "^"
	var ml = make([]MalType, 3)
	ml[0] = MalSymbol{name: "with-meta"}
	mm, err := readHash(r)
	if err != nil {
		return nil, err
	}
	ml[2] = mm
	ma, err := readForm(r)
	if err != nil {
		return nil, err
	}
	ml[1] = ma
	return ml, nil
}

func readQuote(r *Reader, quote string) (MalList, error) {
	Advance(r) // advance past "'"
	var ml = make([]MalType, 2)
	ml[0] = MalSymbol{name: quote}
	tmp, err := readForm(r)
	if err != nil {
		return ml, err
	}
	ml[1] = tmp
	return ml, nil
}

func readHash(r *Reader) (MalHash, error) {
	var res MalHash = make(MalHash)
	Advance(r) // make past initial "{"
	var tmp1, tmp2 MalType
	var err error
	str, err := Peek(r)
	for str != "}" {
		tmp1, err = readForm(r)
		if err != nil {
			return nil, err
		}
		str, err = Peek(r)
		if err != nil {
			return nil, err
		} else if str == "}" {
			return nil, new(EOFError)
		}
		tmp2, err = readForm(r)
		if err != nil {
			return nil, err
		}
		res[tmp1] = tmp2
		str, err = Peek(r)
		if err != nil {
			return nil, err
		}
	}
	Advance(r) // need to advance past the closing "}"
	return res, nil
}

func readList(r *Reader, closing string) ([]MalType, error) {
	var res MalList = make([]MalType, 0, 10)
	var tmp MalType
	Advance(r) // move past initial "("
	str, err := Peek(r)
	for str != closing {
		if err != nil {
			return nil, new(EOFError)
		}
		tmp, err = readForm(r) // this will advance the position
		if err != nil {
			return nil, err
		}
		res = append(res, tmp)
		str, err = Peek(r)
	}
	Advance(r) // need to advance past the closing ")"
	return res, nil
}

func balancedQuotes(s string) bool {
	prevEscape := false
	quoteLevel := 0
	openAlready := false
	for _, val := range s {
		if val == '"' {
			if prevEscape {
				prevEscape = false
			} else {
				if openAlready {
					quoteLevel--
					openAlready = false
				} else {
					quoteLevel++
					openAlready = true
				}
			}
		} else if val == '\\' {
			if prevEscape {
				prevEscape = false
			} else {
				prevEscape = true
			}
		}
	}
	return quoteLevel == 0
}

func properlyQuotedString(s string) bool {
	prevEscape := false
	init := true
	closeQuoted := false
	for _, val := range s {
		if init {
			if val != '"' {
				return false
			}
			init = false
		} else if closeQuoted {
			return false
		} else {
			if val == '\\' {
				if prevEscape {
					prevEscape = false
				} else {
					prevEscape = true
				}
			} else if prevEscape {
				prevEscape = false
			} else if val == '"' {
				closeQuoted = true
			}
			// } else if val == '\n' {
			// 	return false
			// }
		}
	}
	return closeQuoted
}

func readAtom(r *Reader) (MalType, error) {
	val, err := Next(r)
	if err != nil {
		return nil, new(EOFError)
	}
	if s, err := strconv.Atoi(val); err == nil {
		return MalNumber(s), nil
	} else if strparse.MatchString(val) {
		if !properlyQuotedString(val) {
			return nil, new(EOFError)
		}
		// val = strings.TrimPrefix(val, `"`)
		// val = strings.TrimSuffix(val, `"`)
		return MalString(val), nil
	} else {
		return MalSymbol{name: val}, nil
	}
}
