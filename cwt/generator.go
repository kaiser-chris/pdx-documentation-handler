package cwt

import (
	"fmt"
	"os"
	"path"

	"bahmut.de/pdx-documentation-manager/parser"
)

const (
	documentationFolder = "docs"
	outputFolder        = "out"
	cwtFolder           = "cwt"
)

func Generate() error {
	err, cwtPath := setupFolders()
	if err != nil {
		return err
	}

	documentationPath := path.Join(documentationFolder, "new")
	documentation, err := parser.ParseScriptDocumentation(documentationPath)
	if err != nil {
		return err
	}

	iterators := PrintIterators(documentation)
	err = os.WriteFile(path.Join(cwtPath, "lists_generic.cwt"), []byte(iterators), 0644)
	if err != nil {
		return err
	}

	return nil
}

func setupFolders() (error, string) {
	digestPath := path.Join(outputFolder, cwtFolder)
	err := os.MkdirAll(digestPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create output directory: %v", err), digestPath
	}
	return nil, digestPath
}
