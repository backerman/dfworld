// © 2014 Brad Ackerman. Licensed under the WTFPL.
//

package savefile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/backerman/dfworld/pkg/zchunk"
)

type zReader struct {
	f       *file
	hReader io.LimitedReader
	bReader io.ReadCloser
}

// ReadDecompressed returns a reader for the raw save file
// that will transparently decompress. Its behavior is
// undefined if the file is directly accessed before this reader
// is closed. When called, this function will reset the position
// to the beginning of the file.
func (f *file) DecompressedReader() (io.ReadCloser, error) {
	f.Seek(0, 0)

	if f.header.IsCompressed == 0 {
		return nil, errors.New("File is already compressed!")
	}
	zr := &zReader{f: f}
	decompHeader := f.header
	decompHeader.IsCompressed = 0
	headerLen := binary.Size(decompHeader)
	var newHeaderBuf bytes.Buffer
	binary.Write(&newHeaderBuf, binary.LittleEndian, &decompHeader)
	// I'm keeping a LimitedReader here just because it keeps track of the
	// number of bytes read. I think that makes sense.

	zr.hReader = io.LimitedReader{N: int64(headerLen), R: &newHeaderBuf}
	return zr, nil
}

func (r *zReader) Read(buf []byte) (int, error) {
	if r.hReader.N > 0 {
		// Still reading the file header.
		n, err := r.hReader.Read(buf)
		// Read same # of bytes from input and discard.
		dummyBuf := make([]byte, n)
		io.ReadFull(r.f, dummyBuf)
		return n, err
	}
	// Reading from compressed blocks.
	var freshReader bool
	if r.bReader == nil {
		r.bReader = zchunk.NewReader(r.f)
		freshReader = true
	}
	n, err := r.bReader.Read(buf)
	if err == io.EOF {
		// End of block
		r.bReader.Close()
		r.bReader = nil
		err = nil
	} else if err == io.ErrUnexpectedEOF && freshReader {
		// Tried to read a new block but we ran out of file—
		// this means we're at the end.
		err = io.EOF
	}
	return n, err
}

func (r *zReader) Close() error {
	if r.bReader != nil {
		return r.bReader.Close()
	}
	return nil
}
