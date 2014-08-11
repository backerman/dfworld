// © 2014 Brad Ackerman. Licensed under the WTFPL.
//
//

package main

import (
	"io"
	"log"
	"os"

	"github.com/backerman/dfworld/pkg/savefile"
	"github.com/spf13/cobra"
)

var inFilename string
var outFilename string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dfworld",
		Short: "dfworld is tar for Dwarf Fortress save files",
		Long:  "dfworld is tar for Dwarf Fortress save files.",
	}

	rootCmd.PersistentFlags().StringVarP(&inFilename, "in", "f", "world.sav", "Specify the save file to use")

	var decompressCmd = &cobra.Command{
		Use:   "decompress",
		Short: "Decompress a compressed save file",
		Run:   decompress,
	}

	decompressCmd.Flags().StringVarP(&outFilename, "out", "o", "world-out.sav", "Specify the new save file to write to")

	rootCmd.AddCommand(decompressCmd)
	rootCmd.Execute()

}

func decompress(cmd *cobra.Command, args []string) {
	in, err := savefile.NewFileFromPath(inFilename)
	if err != nil {
		log.Fatalf("Unable to open input file: %v", err)
	}
	if in.Header().IsCompressed != 1 {
		log.Fatalf("Input file %v is already compressed; exiting.", inFilename)
	}
	out, err := os.Create(outFilename)
	if err != nil {
		log.Fatalf("Unable to open output file: %v", err)
	}

	inReader, err := in.DecompressedReader()
	if err != nil {
		log.Fatalf("Unable to open output file: %v", err)
	}
	io.Copy(out, inReader)
	out.Close()
	inReader.Close()
	in.Close()
}
