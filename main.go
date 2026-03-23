package main

import (
	"os"

	"bahmut.de/pdx-documentation-manager/comparer"
	"bahmut.de/pdx-documentation-manager/logging"
	"bahmut.de/pdx-documentation-manager/parser"
)

func main() {
	logging.Info("Starting documentation handler")
	oldFolder := "docs\\old"
	newFolder := "docs\\new"
	documentationOld, err := parser.ParseScriptDocumentation(oldFolder)
	if err != nil {
		logging.Fatal(err.Error())
		return
	}
	documentationNew, err := parser.ParseScriptDocumentation(newFolder)
	if err != nil {
		logging.Fatal(err.Error())
		return
	}

	compare := comparer.Compare(documentationOld, documentationNew)

	err = os.WriteFile("result.md", []byte(compare.Print()), 0644)
	if err != nil {
		logging.Fatal(err.Error())
		return
	}
}
