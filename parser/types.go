package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"bahmut.de/pdx-documentation-manager/util"
)

const (
	dataTypesCommon     = "data_types_common.txt"
	dataTypesGui        = "data_types_gui.txt"
	dataTypesClausewitz = "data_types_internalclausewitzgui.txt"
	dataTypesScript     = "data_types_script.txt"
	dataTypesOther      = "data_types_uncategorized.txt"
)

const (
	dataTypeSeparator   = "-----------------------"
	dataTypeDefinition  = "Definition type: "
	dataTypeReturn      = "Return type: "
	dataTypeDescription = "Description: "
	dataTypeMacro       = "Macro replacement: "
)

type DataTypeDocumentation struct {
	DataTypes      []*DataType         `json:"data-types"`
	GlobalPromotes []*DataTypeFunction `json:"global-promotes"`
}

type DataType struct {
	Name        string
	Description string
}

type DataTypeFunction struct {
	Name        string
	Description string
	ReturnType  string
}

func ParseDataTypeDocumentation(folder string) (*DataTypeDocumentation, error) {
	documentation := &DataTypeDocumentation{
		DataTypes:      make([]*DataType, 0),
		GlobalPromotes: make([]*DataTypeFunction, 0),
	}

	err := parseDataTypeFile(path.Join(folder, dataTypesCommon), documentation)
	if err != nil {
		return nil, err
	}

	err = parseDataTypeFile(path.Join(folder, dataTypesGui), documentation)
	if err != nil {
		return nil, err
	}

	err = parseDataTypeFile(path.Join(folder, dataTypesClausewitz), documentation)
	if err != nil {
		return nil, err
	}

	err = parseDataTypeFile(path.Join(folder, dataTypesScript), documentation)
	if err != nil {
		return nil, err
	}

	err = parseDataTypeFile(path.Join(folder, dataTypesOther), documentation)
	if err != nil {
		return nil, err
	}

	return documentation, nil
}

func parseDataTypeFile(file string, documentation *DataTypeDocumentation) error {
	if !util.Exists(file) {
		return fmt.Errorf("common data type documentation does not exist: %s", file)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	var dataType *DataType = nil
	var function *DataTypeFunction = nil
	for scanner.Scan() {
		cleanLine := strings.TrimSpace(scanner.Text())
		if isDataTypeName(cleanLine) && (dataType != nil && function != nil) {
			return fmt.Errorf("unterminated data type: %s, %s", dataType.Name, cleanLine)
		}
		if isDataTypeName(cleanLine) {
			dataType = &DataType{
				Name: strings.TrimSuffix(cleanLine, ":"),
			}
			function = &DataTypeFunction{
				Name: strings.TrimSuffix(cleanLine, ":"),
			}
			continue
		}
		if cleanLine == terminator && (dataType != nil || function != nil) {
			if dataType != nil {
				documentation.DataTypes = append(documentation.DataTypes, dataType)
				function = nil
				dataType = nil
			}
			if function != nil {
				documentation.GlobalPromotes = append(documentation.GlobalPromotes, function)
				function = nil
				dataType = nil
			}
			continue
		}
		if dataType == nil && function == nil {
			continue
		}
		if strings.HasPrefix(cleanLine, dataTypeDescription) {
			if dataType != nil {
				dataType.Description = strings.TrimPrefix(cleanLine, dataTypeDescription)
			}
			if function != nil {
				function.Description = strings.TrimPrefix(cleanLine, dataTypeDescription)
			}
			continue
		}
		if strings.HasPrefix(cleanLine, dataTypeDefinition) {
			definition := strings.TrimPrefix(cleanLine, dataTypeDefinition)
			if definition == "Type" {
				function = nil
				continue
			}
			if isGlobalPromote(definition) {
				dataType = nil
				continue
			}
			dataType = nil
			function = nil
			continue
		}
		if strings.HasPrefix(cleanLine, dataTypeReturn) {
			if function != nil {
				function.ReturnType = strings.TrimPrefix(cleanLine, dataTypeReturn)
			}
			continue
		}
	}

	return nil
}

func isDataTypeName(line string) bool {
	return !strings.HasPrefix(line, dataTypeDefinition) &&
		!strings.HasPrefix(line, dataTypeReturn) &&
		!strings.HasPrefix(line, dataTypeDescription) &&
		!strings.HasPrefix(line, dataTypeSeparator) &&
		!strings.HasPrefix(line, dataTypeMacro) &&
		!(strings.TrimSpace(line) == terminator)
}

func isGlobalPromote(definition string) bool {
	return definition == "Global promote" ||
		definition == "Global function"
}
