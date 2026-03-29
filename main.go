package main

import (
	"os"
	"slices"

	"bahmut.de/pdx-documentation-manager/cwt"
	"bahmut.de/pdx-documentation-manager/digest"
	"bahmut.de/pdx-documentation-manager/logging"
)

func main() {
	logging.Info("Starting documentation handler")

	if slices.Contains(os.Args, "digest") {
		logging.Info("Generating digest")
		err := digest.Generate()
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
