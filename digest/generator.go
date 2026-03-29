package digest

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"bahmut.de/pdx-documentation-manager/comparer"
	"bahmut.de/pdx-documentation-manager/parser"
)

const (
	documentationFolder = "docs"
	outputFolder        = "out"
	digestFolder        = "digest"
)

func Generate() error {
	err, digestPath := setupFolders()
	if err != nil {
		return err
	}

	err = generateScriptDocumentationChanges(digestPath)
	if err != nil {
		return err
	}

	documentation, err := parser.ParseDataTypeDocumentation(path.Join(documentationFolder, "old"))
	print(documentation)

	err = generateDocumentationFolder(digestPath)
	if err != nil {
		return err
	}

	return nil
}

func generateScriptDocumentationChanges(digestPath string) error {
	oldFolder := path.Join(documentationFolder, "old")
	newFolder := path.Join(documentationFolder, "new")

	documentationOld, err := parser.ParseScriptDocumentation(oldFolder)
	if err != nil {
		return err
	}
	documentationNew, err := parser.ParseScriptDocumentation(newFolder)
	if err != nil {
		return err
	}

	compare := comparer.Compare(documentationOld, documentationNew)

	err = os.WriteFile(path.Join(digestPath, "changes_script_docs.md"), []byte(compare.Print()), 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateDocumentationFolder(digestPath string) error {
	docsPath := path.Join(digestPath, documentationFolder)
	newFolder := path.Join(documentationFolder, "new")

	err := os.MkdirAll(docsPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create docs directory: %v", err)
	}

	err = filepath.Walk(newFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".log") && !strings.HasSuffix(info.Name(), ".txt") {
			return nil
		}
		target := filepath.Join(docsPath, info.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		err = os.WriteFile(target, data, 0644)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not copy documentation into digest: %v", err)
	}

	return nil
}

func setupFolders() (error, string) {
	digestPath := path.Join(outputFolder, digestFolder)
	err := os.MkdirAll(digestPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create output directory: %v", err), digestPath
	}
	return nil, digestPath
}
