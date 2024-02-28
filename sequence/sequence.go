package sequence

import (
	"fmt"
)

type Sequence struct {
	Sequence []int `json:"sequence"`
}

func NewSequence(sequence []int) *Sequence {
	return &Sequence{sequence}
}

// GetSequence returns a copy of the sequence
// so that the internal sequence cannot be modified
// and we can mutate the copy without affecting the
// original object
func (s *Sequence) GetSequence() []int {
	seq := make([]int, 0)
	seq = append(seq, s.Sequence...)
	return seq
}

func (s *Sequence) ToString() string {
	text := fmt.Sprintf("Operations %d", len(s.Sequence))
	for _, v := range s.Sequence {
		text = fmt.Sprintf("%s\n%d", text, v)
	}
	text = fmt.Sprintf("%s\n", text)
	return text
}
