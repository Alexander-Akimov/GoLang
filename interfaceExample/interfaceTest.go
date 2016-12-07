package main

import "io"

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // Преобразование int в ByteCounter
	return len(p), nil
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	return w, new(int64)
}

func main() {
	var bc ByteCounter
	bc = 5
	wr, count := CountingWriter(&bc)

	s := make([]byte, 5, 5)
	wr.Write(s)
	print(count)
}
