// © 2014 Brad Ackerman. Licensed under the WTFPL.
//

// Package savefile provides methods for reading Dwarf Fortress save files.
package savefile

import (
	"encoding/binary"
	"io"
	"os"
)

// File represents our file; it mostly wraps the reader
type File interface {
	Header() FileHeader
	ListChunks() []int
	DecompressedReader() (io.ReadCloser, error)
	io.Closer
}

type file struct {
	header FileHeader
	*os.File
}

// FileHeader is the DF world.dat/.sav file header.
type FileHeader struct {
	Version      uint32
	IsCompressed uint32
}

func readHeader(r io.Reader) FileHeader {
	var header FileHeader
	binary.Read(r, binary.LittleEndian, &header)
	return header
}

func (f *file) Header() FileHeader {
	return f.header
}

func (f *file) ListChunks() []int {
	return []int{0}
}

func (f *file) Close() error {
	return f.File.Close()
}

// NewFileFromPath is a convenience method for NewFile.
func NewFileFromPath(path string) (File, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewFile(r)
}

// NewFile gets a new File object.
func NewFile(r *os.File) (File, error) {
	f := new(file)
	f.File = r
	f.header = readHeader(f)
	f.Seek(0, 0) // go back to start
	return f, nil
}
