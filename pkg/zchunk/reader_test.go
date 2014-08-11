// © 2014 Brad Ackerman. Licensed under the WTFPL.
//
//
// The zchunk package handles streams of concatenated zlib-compressed
// data prepended by a 32-bit, little-endian length counter.

package zchunk_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/backerman/dfworld/pkg/zchunk"
	. "github.com/onsi/gomega"
)

type testVector struct {
	plain      string
	compressed []byte
}

var testVectors = []testVector{
	{"test",
		[]byte{0x0c, 0x00, 0x00, 0x00,
			0x78, 0x9c, 0x2b, 0x49, 0x2d, 0x2e, 0x01, 0x00, 0x04, 0x5d, 0x01, 0xc1}},
}

func TestRead(t *testing.T) {
	RegisterTestingT(t)
	for _, test := range testVectors {
		inReader := bytes.NewReader(test.compressed)
		outReader := zchunk.NewReader(inReader)
		outWriter := new(bytes.Buffer)
		io.Copy(outWriter, outReader)
		Ω(outWriter.String()).Should(Equal(test.plain))
		outReader.(io.Closer).Close()
	}
}
