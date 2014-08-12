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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dfworld",
		Short: "dfworld is tar for Dwarf Fortress save files",
		Long:  "dfworld is tar for Dwarf Fortress save files.",
	}

	var decompressCmd = &cobra.Command{
		Use:   "decompress [infile] [outfile]",
		Short: "Decompress a compressed save file",
		Long:  `Decompress infile, saving the decompressed version as outfile.`,
		Run:   decompress,
	}

	rootCmd.AddCommand(decompressCmd)
	rootCmd.Execute()
}

func decompress(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("decompress requires exactly two arguments.")
	}
	inFilename, outFilename := args[0], args[1]
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
