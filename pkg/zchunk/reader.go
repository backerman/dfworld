// © 2014 Brad Ackerman. Licensed under the WTFPL.
//
//

// Package zchunk handles streams of concatenated zlib-compressed
// data prepended by a 32-bit, little-endian length counter.
package zchunk

import (
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

type reader struct {
	hasFinished bool             // Whether we've finished reading the block
	zOutReader  io.Reader        // Read the decompressed block
	zInReader   io.LimitedReader // Passed to zlib
	rawReader   io.Reader        // Input from our caller
}

// NewReader creates a new chunk reader from an existing Reader.
func NewReader(r io.Reader) io.ReadCloser {
	chunkReader := new(reader)
	chunkReader.rawReader = r
	return chunkReader
}

// Read reads from the compressed input stream, up to the end of the
// current block.
func (r *reader) Read(buf []byte) (int, error) {
	if r.hasFinished {
		// Nothing more to read.
		fmt.Println("Returning w/ EOF")
		return 0, io.EOF
	}

	var err error

	if r.zInReader.N == 0 {
		// Haven't started reading yet -- get the length and create
		// a zlib reader to read the compressed string.
		var toRead uint32
		err = binary.Read(r.rawReader, binary.LittleEndian, &toRead)
		if err == nil {
			r.zInReader = io.LimitedReader{R: r.rawReader, N: int64(toRead)}
			r.zOutReader, err = zlib.NewReader(&r.zInReader)
		}
	}

	if err == io.EOF {
		// We were expecting a block and we actually got an EOF where the header
		// should have been.
		return 0, io.ErrUnexpectedEOF
	}

	if err != nil {
		fmt.Printf("Returning w/ error %v", err)
		return 0, err
	}

	// Now that we have a stream, read it.
	var n int
	n, err = r.zOutReader.Read(buf)
	// Check to see if we've read the entire block.
	if r.zInReader.N == 0 {
		r.hasFinished = true
		err = io.EOF
	}
	return n, err
}

func (r *reader) Close() error {
	var err error
	if r.zOutReader != nil {
		// Close the zlib reader. Since we created it, we know it's a Closer
		// as well as a Reader; otherwise, we'd need to check.
		err = r.zOutReader.(io.Closer).Close()
	}

	return err
}
