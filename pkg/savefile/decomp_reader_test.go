// © 2014 Brad Ackerman. Licensed under the WTFPL.
//

package savefile_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/backerman/dfworld/pkg/savefile"
	. "github.com/onsi/gomega"
)

func TestSomething(t *testing.T) {
	RegisterTestingT(t)
	compressed, err := savefile.NewFileFromPath("testdata/test-world.sav")
	Ω(err).ShouldNot(HaveOccurred())
	actualR, err := compressed.DecompressedReader()
	Ω(err).ShouldNot(HaveOccurred())
	expectedR, err := os.Open("testdata/test-world-decomp.sav")
	Ω(err).ShouldNot(HaveOccurred())
	actual, err := ioutil.ReadAll(actualR)
	Ω(err).ShouldNot(HaveOccurred())
	expected, err := ioutil.ReadAll(expectedR)
	Ω(err).ShouldNot(HaveOccurred())
	Ω(actual).Should(Equal(expected))
	Ω(actual).ShouldNot(BeEmpty())
	Ω(expected).ShouldNot(BeEmpty())
	actualR.Close()
	expectedR.Close()
	compressed.Close()
}
