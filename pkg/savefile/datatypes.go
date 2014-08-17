// © 2014 Brad Ackerman. Licensed under the WTFPL.
//
// Data types used by the Dwarf Fortress save file and read/write methods.

package savefile

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"

	"code.google.com/p/go-charset/charset"
	// import charset data into the binary
	_ "code.google.com/p/go-charset/data"
)

// Convert437String reads a CP-437 string and returns it
// as a UTF-8 string.
func Convert437String(b []byte) string {
	strReader, err := charset.NewReader("ibm437", bytes.NewReader(b))
	if err != nil {
		log.Fatalf("Unable to open reader - dying horribly: %v", err)
	}
	utfStringBuf, err := ioutil.ReadAll(strReader)
	if err != nil {
		log.Fatalf("Unable to read string - dying horribly: %v", err)
	}
	return string(utfStringBuf)
}

// readString reads a CP437 string with prepended 16-bit, little-endian
// length.
func readString(r io.Reader) string {
	var strlen uint16
	err := binary.Read(r, binary.LittleEndian, &strlen)
	if err != nil {
		log.Fatalf("Error when reading string length: %v", err)
	}
	strBuf := make([]byte, strlen)
	_, err = io.ReadFull(r, strBuf)
	if err != nil {
		log.Fatalf("Unable to open reader - dying horribly: %v", err)
	}
	return Convert437String(strBuf)
}
