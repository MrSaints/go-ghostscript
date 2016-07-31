package main

import (
	"github.com/mrsaints/go-ghostscript/ghostscript"
	"log"
)

func main() {
	rev, err := ghostscript.GetRevision()
	if err != nil {
		panic(err)
	}
	log.Printf("Revision: %+v\n", rev)

	gs, err := ghostscript.NewInstance()
	if err != nil {
		panic(err)
	}

	args := []string{
		"gs", // This will be ignored
		"-q",
		"-dBATCH",
		"-dColorConversionStrategy=/LeaveColorUnchanged",
		"-dCompatibilityLevel=1.5",
		"-dEmbedAllFonts=true",
		"-dNOPAUSE",
		"-dPDFSETTINGS=/printer",
		"-dSAFER",
		"-dSubsetFonts=true",
		"-sDEVICE=pdfwrite",
		"-sOutputFile=concat.pdf",
		"page-1.pdf",
		"page-2.pdf",
	}

	if err := gs.Init(args); err != nil {
		panic(err)
	}

	defer func() {
		gs.Exit()
		gs.Destroy()
	}()
}
