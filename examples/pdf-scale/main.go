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
		"-dDEVICEHEIGHT=843", // Obtained from MediaBox
		"-dDEVICEWIDTH=596",
		"-sOutputFile=scaled.pdf",
		"-c",
		"<</BeginPage{0.5 0.5 scale 298 421.5 translate}>> setpagedevice",
		"-funscaled.pdf",
	}

	if err := gs.Init(args); err != nil {
		panic(err)
	}

	defer func() {
		gs.Exit()
		gs.Destroy()
	}()
}
