package translation

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	ErrIdentifier = errors.New("identifier error")
	ErrValue      = errors.New("value error")
)

type Variable struct {
	Identifier      string
	Value           string
	recording       bool
	identifierStart int
	identifierEnd   int
	valueStart      int
	valueEnd        int
}

func (v *Variable) toggleRecording() {
	v.recording = !v.recording
}

func (v *Variable) isRecording() bool {
	return v.recording
}

func ToASCII(str string) string {
	result, _, err := transform.String(transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn))), str)
	if err != nil {
		return ""
	}
	return result
}

func FromString(s string) (*Variable, error) {
	start := 0
	v := &Variable{}
	for i := start; i < len(s); i++ {
		// 48->57 is 0-9 in ascii
		if v.isRecording() && (s[i] < 48 || s[i] > 57) {
			v.identifierEnd = i
			v.Identifier = s[v.identifierStart:v.identifierEnd]
			break
		}
		// "@" is 64 and "#" is 35
		if s[i] == 64 || s[i] == 35 {
			v.identifierStart = i + 1
			v.toggleRecording()
		}
	}
	if v.Identifier == "" {
		return nil, fmt.Errorf("%w empty identifier", ErrIdentifier)
	}
	v.toggleRecording()
	for i := v.identifierEnd; i < len(s); i++ {
		// "~" is 126
		if v.isRecording() && s[i] == 126 {
			v.valueEnd = i
			v.Value = ToASCII(s[v.valueStart:v.valueEnd])
			v.toggleRecording()
			break
		}
		// "~" is 126
		if s[i] == 126 {
			v.valueStart = i + 1
			v.toggleRecording()
		}
	}
	if v.isRecording() {
		return nil, fmt.Errorf("%w still recording should be closed, got %+v", ErrValue, v)
	}
	return v, nil
}

func FromFileContents(fileContents *[]string) (*[]Variable, error) {
	out := []Variable{}
	multi := false
	buffer := ""
	for _, line := range *fileContents {
		// Deal with single line comment
		if len(line) < 2 && line[0:1] == "//" {
			continue
		}
		switch {
		case !multi && strings.Count(line, "~") != 2:
			multi = true
			buffer += line
		case strings.Count(line, "~") == 1:
			multi = false
			buffer += line
			if v, err := FromString(buffer); err == nil {
				out = append(out, *v)
			}
		default:
			if v, err := FromString(line); err == nil {
				out = append(out, *v)
			}
		}
	}
	return &out, nil
}
