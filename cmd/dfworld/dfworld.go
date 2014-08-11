// © 2014 Brad Ackerman. Licensed under the WTFPL.
//
//

package main

import (
	"fmt"
	"os"

	"github.com/backerman/dfworld/pkg/savefile"
)

func main() {
	f, err := savefile.NewFileFromPath("world.sav")

	if err != nil {
		fmt.Printf("Error %v!\n", err)
		os.Exit(42)
	}
	fmt.Printf("Header: %v\n", f.Header())
	f.Close()
}
