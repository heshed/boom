package boomer
import (
	"bufio"
	"os"
)

type Reader struct {
	path string
	input *os.File
	scanner *bufio.Scanner
}

func NewReader(path string) *Reader {
	input, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	return &Reader{
		path: path,
		input: input,
		scanner: bufio.NewScanner(input),
	}
}

func (r *Reader) Next() (string, bool) {
	if r.scanner.Scan() {
		return r.scanner.Text(), true
	} else {
		r.input.Close()
		return "", false
	}
}

