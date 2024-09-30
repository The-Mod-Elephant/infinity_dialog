package translation

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrIdentifier = errors.New("identifier error")
	ErrValue      = errors.New("value error")
)

type Varriable struct {
	Identifier string
	Value      string
}

func ParseVarriable(s string, start *int) (*Varriable, error) {
	var i, v string
	var flag bool
	for *start < len(s) {
		r := rune(s[*start])
		if r == '@' {
			flag = true
		} else if flag && !unicode.IsDigit(r) {
			flag = false
			break
		} else if flag {
			i += string(r)
		}
		*start++
	}
	if i == "" {
		return nil, fmt.Errorf("%w empty identifier", ErrIdentifier)
	}
	for *start < len(s) {
		// "~" is 126
		if flag && s[*start] == 126 {
			flag = false
			break
		} else if s[*start] == 126 {
			flag = true
		} else if flag {
			v += string(s[*start])
		}
		*start++
	}
	if flag {
		return nil, fmt.Errorf("%w flag not closed, got %s", ErrValue, v)
	}
	return &Varriable{
		Identifier: i,
		Value:      v,
	}, nil
}

func FromFileContents(fileContents string) (*[]Varriable, error) {
	out := []Varriable{}
	n := 0
	for n < len(fileContents) {
		if n+2 < len(fileContents) && fileContents[n] == '/' && fileContents[n+1] == '/' {
			n += 2
			for n < len(fileContents) {
				s := fileContents[n]
				if s == '\n' {
					break
				}
				n++
			}
		}
		v, err := ParseVarriable(fileContents, &n)
		if err == nil {
			out = append(out, *v)
		}
	}
	return &out, nil
}
