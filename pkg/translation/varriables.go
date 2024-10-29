package translation

import (
	"errors"
	"fmt"
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

func (v *Variable) start() int {
	return v.identifierStart
}

func (v *Variable) end() int {
	return v.valueEnd
}

func FromString(s string, start int) (*Variable, error) {
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
			v.Value = s[v.valueStart:v.valueEnd]
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

func FromFileContents(fileContents string) (*[]Variable, error) {
	out := []Variable{}
	for n := 0; n < len(fileContents); n++ {
		// Deal with single line comment
		if n+1 < len(fileContents) && fileContents[n] == '/' && fileContents[n+1] == '/' {
			n++
			for n < len(fileContents) {
				s := fileContents[n]
				if s == '\n' {
					break
				}
				n++
			}
		}
		if v, err := FromString(fileContents, n); err == nil {
			out = append(out, *v)
			n = v.end()
		}
	}
	return &out, nil
}
