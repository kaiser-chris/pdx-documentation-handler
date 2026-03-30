package main

import (
	"flag"
	"os"
	"slices"

	"bahmut.de/pdx-documentation-manager/cwt"
	"bahmut.de/pdx-documentation-manager/digest"
	"bahmut.de/pdx-documentation-manager/logging"
)

const (
	FlagVersion = "version"
)

func main() {
	logging.Info("Starting documentation handler")

	if slices.Contains(os.Args, "digest") {
		version := flag.String(FlagVersion, "", "Game Version the digest is for")
		flag.Parse()
		if version == nil || *version == "" {
			logging.Fatalf("The parameter %s%s%s is required.\n", logging.AnsiBoldOn, FlagVersion, logging.AnsiAllDefault)
			return
		}

		logging.Info("Generating digest")
		err := digest.Generate(*version)
		if err != nil {
			logging.Fatal(err.Error())
			return
		}
	}

	if slices.Contains(os.Args, "cwt") {
		logging.Info("Generating CWT files")
		err := cwt.Generate()
		if err != nil {
			logging.Fatal(err.Error())
			return
		}
	}

}
