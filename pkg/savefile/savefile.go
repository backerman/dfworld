// © 2014 Brad Ackerman. Licensed under the WTFPL.
//

// Package savefile provides methods for reading Dwarf Fortress save files.
package savefile

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
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
	Name string // offset: 138 (DF2014)
	Year int
	//	CivName string
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

type versString struct {
	version uint32
	str     string
}

// versionStrings maps the header version uint16 to its human-readable equivalent.
var versionStrings = []versString{
	{1287, "0.31.01"}, {1288, "0.31.02"}, {1289, "0.31.03"}, {1292, "0.31.04"},
	{1295, "0.31.05"}, {1297, "0.31.06"}, {1300, "0.31.08"}, {1304, "0.31.09"},
	{1305, "0.31.10"}, {1310, "0.31.11"}, {1311, "0.31.12"}, {1323, "0.31.13"},
	{1325, "0.31.14"}, {1326, "0.31.15"}, {1327, "0.31.16"}, {1340, "0.31.17"},
	{1341, "0.31.18"}, {1351, "0.31.19"}, {1353, "0.31.20"}, {1354, "0.31.21"},
	{1359, "0.31.22"}, {1360, "0.31.23"}, {1361, "0.31.24"}, {1362, "0.31.25"},
	{1372, "0.34.01"}, {1374, "0.34.02"}, {1376, "0.34.03"}, {1377, "0.34.04"},
	{1378, "0.34.05"}, {1382, "0.34.06"}, {1383, "0.34.07"}, {1400, "0.34.08"},
	{1402, "0.34.09"}, {1403, "0.34.10"}, {1404, "0.34.11"}, {1441, "0.40.01"},
	{1442, "0.40.02"}, {1443, "0.40.03"}, {1444, "0.40.04"}, {1445, "0.40.05"},
	{1446, "0.40.06"}, {1448, "0.40.07"}, {1449, "0.40.08"},
}

type versMap struct {
	version uint32
	offset  int64
}

// worldHeaderLen returns the length of this save's
// world header; its value depends on the Dwarf Fortess
// version that created the save.
//
// The world header comes after the save header.
func (f *file) worldHeaderLen() (l int64) {
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

// GetInfo returns a FileInfo struct for this world.
func (f *file) GetInfo() (i FileInfo) {
	for _, v := range versionStrings {
		if v.version == f.header.Version {
			i.Version = v.str
		}
	}
	if i.Version == "" {
		i.Version = fmt.Sprintf("Unidentified version %v", f.header.Version)
	}
	f.Seek(0, os.SEEK_CUR)
	r, err := f.DecompressedReader()
	if err != nil {
		log.Fatalf("Couldn't open savefile: %v", err)
	}
	defer r.Close()

	// Save header has already been read; skip over it.
	_ = readHeader(r)
	// Advance over the world header.
	io.CopyN(ioutil.Discard, r, f.worldHeaderLen())
	if f.activeFortress {
		fort := new(FortInfo)
		fort.Name = readString(r)
		i.Fort = fort
	}
	i.WorldName = readString(r)
	if f.activeFortress {
		var year uint32
		binary.Read(r, binary.LittleEndian, &year)
		i.Fort.Year = int(year)
	}
	return
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
