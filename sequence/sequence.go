package sequence

import (
	"fmt"
)

type Sequence struct {
	sequence []int
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
	seq = append(seq, s.sequence...)
	return seq
}

func (s *Sequence) ToString() string {
	text := fmt.Sprintf("Operations %d", len(s.sequence))
	for _, v := range s.sequence {
		text = fmt.Sprintf("%s\n%d", text, v)
	}
	text = fmt.Sprintf("%s\n", text)
	return text
}
