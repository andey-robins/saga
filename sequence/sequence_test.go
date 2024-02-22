package sequence

import "testing"

func TestLoadSequence(t *testing.T) {
	tests := []struct {
		sequenceString  string
		sequenceNumbers []int
	}{
		{
			"Operations 4\n4\n6\n5\n7\n",
			[]int{4, 6, 5, 7},
		},
	}

	for _, test := range tests {
		sequence := LoadSequenceFromString(test.sequenceString)
		if sequence == nil || sequence.Sequence == nil {
			t.Errorf("Expected sequence to not be nil")
		}

		if len(sequence.Sequence) != len(test.sequenceNumbers) {
			t.Errorf("Expected sequence to have length %d, got %d", len(test.sequenceNumbers), len(sequence.Sequence))
		}

		for i, v := range sequence.Sequence {
			if v != test.sequenceNumbers[i] {
				t.Errorf("Expected sequence to have value %d at index %d, got %d", test.sequenceNumbers[i], i, v)
			}
		}

		if sequence.ToString() != test.sequenceString {
			t.Errorf("Expected sequence to have string representation %s, got %s", test.sequenceString, sequence.ToString())
		}
	}
}
