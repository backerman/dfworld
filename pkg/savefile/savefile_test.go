// © 2014 Brad Ackerman. Licensed under the WTFPL.
//

package savefile_test

import (
	"testing"

	"github.com/backerman/dfworld/pkg/savefile"
	. "github.com/onsi/gomega"
)

type world struct {
	filename  string
	worldname string
	version   string
	save      savefile.File
	fort      *fortress
}

type fortress struct {
	name    string
	civname string
}

var saves = []*world{
	{
		filename:  "testdata/savtest/world.sav",
		worldname: "Thur Minbaz",
		fort: &fortress{
			name:    "Avuzdakost",
			civname: "Vesalath",
		},
	},
}

func SetupSavefiles() {
	for _, s := range saves {
		var err error
		s.save, err = savefile.NewFileFromPath(s.filename)
		Ω(err).Should(BeNil())
	}
}

func TestSavefileParsing(t *testing.T) {
	RegisterTestingT(t)
	SetupSavefiles()
	for _, s := range saves {
		i := s.save.GetInfo()
		Ω(i.Version).Should(Equal(s.version))
		Ω(i.WorldName).Should(Equal(s.worldname))
		if s.fort != nil {
			Ω(i.Fort.Name).Should(Equal(s.fort.name))
			Ω(i.Fort.CivName).Should(Equal(s.fort.civname))
		}
	}
}
