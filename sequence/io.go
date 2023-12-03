package sequence

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// standard check helper function
func check(e error) {
	if e != nil {
		panic(e)
	}
}

/// FILE IO

// LoadSequenceFromFile reads a sequence from a file. Returns
// a pointer to a sequence object. Panics if there is an error
// with the file since it's beyond the scope of our program
// to handle this.
func LoadSequenceFromFile(path string) *Sequence {
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	return LoadSequence(f)
}

// WriteToFile writes a sequence to a file in the same format
// as is described by the assignment. Returns an error if there
// is a problem writing to the file.
func (s *Sequence) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(s.ToString())
	return err
}

/// STRING IO

// LoadSequenceFromString reads a sequence from a string
// and returns a pointer to a sequence object
func LoadSequenceFromString(sequenceString string) *Sequence {
	return LoadSequence(strings.NewReader(sequenceString))
}

/// PRIMARY PARSER

// LoadSequence reads a sequence from a reader, see wrappers
// for file and string reading above. Returns a pointer to
// a sequence object
func LoadSequence(encoding io.Reader) *Sequence {
	sequence := make([]int, 0)

	scanner := bufio.NewScanner(encoding)
	scanner.Split(bufio.ScanLines)
	log.Println("Setup scanner for reading sequence from file")

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
