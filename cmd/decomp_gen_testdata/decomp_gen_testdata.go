// © 2014 Brad Ackerman. Licensed under the WTFPL.
//
// decomp_gen_testdata generates a test artifact file for
// compressed-world reading.
//
// Yes, this isn't even a valid save file; will change the artifact appropriately
// as we actually implement more functionality.

package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"os"

	"github.com/backerman/dfworld/pkg/savefile"
)

var testBlocks = []string{
	"Test one two three\n",
	"Four five six seven eight ",
	"",
	"Nineteneleventwelve",
}

func main() {
	testfile, _ := os.Create("test-world.sav")
	testfileDecomp, _ := os.Create("test-world-decomp.sav")
	// Write a header.
	header := savefile.FileHeader{Version: 1446, IsCompressed: 1}
	dHeader := savefile.FileHeader{Version: 1446, IsCompressed: 0}
	binary.Write(testfile, binary.LittleEndian, &header)
	binary.Write(testfileDecomp, binary.LittleEndian, &dHeader)
	for _, s := range testBlocks {
		// Compress the string.
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		w.Write([]byte(s))
		w.Close()
		// Write it to our test file.
		l := uint32(b.Len())
		binary.Write(testfile, binary.LittleEndian, &l)
		testfile.Write(b.Bytes())
		// Now the decompressed version.
		testfileDecomp.Write([]byte(s))
	}
	testfile.Close()
	testfileDecomp.Close()
}
