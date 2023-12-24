package sequence

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// LoadSequenceFromFile panics if there is an error
// with the file since it's beyond the scope of our program
// to handle this.
func LoadSequenceFromFile(path string) *Sequence {
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	return loadSequence(f)
}

// WriteToFile writes a sequence to a file in the same format
// as the files when they are parsed. Cascades any write errors
// to the caller for handling as close to the user as possible
func (s *Sequence) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(s.ToString())
	return err
}

func LoadSequenceFromString(sequenceString string) *Sequence {
	return loadSequence(strings.NewReader(sequenceString))
}

// loadSequence reads a sequence from a reader. It's used by
// the exported functions in this package to read from files
// and strings. It contains the parsing logic for sequence files
func loadSequence(encoding io.Reader) *Sequence {
	sequence := make([]int, 0)

	scanner := bufio.NewScanner(encoding)
	scanner.Split(bufio.ScanLines)

	// skip first line as all it does is declare the number of elements
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		for _, sequenceNumber := range strings.Split(line, " ") {
			id, err := strconv.Atoi(sequenceNumber)
			check(err)
			sequence = append(sequence, id)
		}
	}

	return NewSequence(sequence)
}
