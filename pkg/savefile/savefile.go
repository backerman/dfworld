// © 2014 Brad Ackerman. Licensed under the WTFPL.
//

// Package savefile provides methods for reading Dwarf Fortress save files.
package savefile

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// File represents our file; it mostly wraps the reader
type File interface {
	Header() FileHeader
	ListChunks() []int
	GetInfo() FileInfo
	DecompressedReader() (io.ReadCloser, error)
	io.Closer
}

type file struct {
	header         FileHeader
	activeFortress bool
	*os.File
}

// FileHeader is the DF world.dat/.sav file header.
type FileHeader struct {
	Version      uint32
	IsCompressed uint32
}

// FileInfo is all the information about a file that we're
// currently extracting.
type FileInfo struct {
	Version   string // header
	WorldName string
	Fort      *FortInfo
}

// FortInfo provides information about an active fortress.
type FortInfo struct {
	Name    string // offset: 138 (DF2014)
	CivName string
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

type versMap struct {
	version uint32
	offset  int
}

// worldHeaderLen returns the length of this save's
// world header; its value depends on the Dwarf Fortess
// version that created the save.
//
// The world header comes after the save header.
func (f *file) worldHeaderLen() (l int) {
	// meh, these are all off by constant... go though rawextract again
	var offsets []versMap

	if f.activeFortress {
		// world.sav
		offsets = []versMap{
			{0, 86},
			{1372, 106},
			{1400, 110},
			{1441, 130},
		}
	} else {
		// world.dat
		offsets = []versMap{
			{0, 70},
			{1372, 90},
			{1400, 94},
			{1441, 114},
		}
	}
	for _, o := range offsets {
		if f.header.Version >= o.version {
			l = o.offset
		}
	}
	return
}

func (f *file) GetInfo() FileInfo {
	// Get the decompressed reader, seek to the start of what we
	// care about, read what we want, close the reader.
	return FileInfo{}
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
	fileExtension := strings.ToLower(filepath.Ext(f.File.Name()))
	switch fileExtension {
	case ".sav":
		f.activeFortress = true
	case ".dat":
		f.activeFortress = false
	default:
		log.Fatalf("File %v extension invalid (not .sav or .dat)", fileExtension)
	}
	f.Seek(0, 0) // go back to start
	return f, nil
}
