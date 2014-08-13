// © 2014 Brad Ackerman. Licensed under the WTFPL.
//

package savefile_test

import (
	"testing"

	"github.com/backerman/dfworld/pkg/savefile"
	. "github.com/onsi/gomega"
)

type testVector struct {
	cp437 []byte
	utf8  string
}

var testVectors = []testVector{
	{[]byte{0xae, 0x9c, 0x6f, 0x76, 0x65, 0x6c, 0x79, 0x20, 0x70, 0x69, 0xa4,
		0x61, 0x74, 0x61, 0xaf}, "«£ovely piñata»"},
	{[]byte{}, ""},
	{[]byte{0xad, 0x78, 0x21}, "¡x!"},
}

func TestCP437(t *testing.T) {
	RegisterTestingT(t)
	for _, v := range testVectors {
		actual := savefile.Convert437String(v.cp437)
		Ω(actual).Should(Equal(v.utf8))
	}
}
