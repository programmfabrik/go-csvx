package csvx

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrDataIsNil                          = errors.New("data is nil")
	ErrOnlyOneRowIsAllowedForStringArray  = errors.New("only one row is allowed for type 'string,array'")
	ErrOnlyOneRowIsAllowedForInt64Array   = errors.New("only one row is allowed for type 'int64,array'")
	ErrOnlyOneRowIsAllowedForFloat64Array = errors.New("only one row is allowed for type 'float64,array'")
	ErrOnlyOneRowIsAllowedForBoolArray    = errors.New("only one row is allowed for type 'bool,array'")
	ErrInEmbeddedJSON                     = errors.New("unable to parse json in csv")
	ErrUnsupportedType                    = errors.New("unsupported type format type")
)

type field struct {
	Name string
	Type string
}

type CSVParser struct {
	// Comma defines the rune with which the entries in the csv file are separated from each other.
	Comma rune
	// Comment defines the rune used to mark comment strings within the CSV.
	// If the line starts with this rune, the whole line is ignored.
	Comment rune
	// TrimLeadingSpace specifies whether leading spaces should be trimmed or not.
	TrimLeadingSpace bool
	// SkipEmptyColumns defines whether empty rows should be ignored or not.
	SkipEmptyColumns bool
	// isTyped defines whether the user expected to receive a typed or untyped response.
	isTyped bool
}

// Untyped unmarshals the data into a slice of map[string]interface{}
func (c *CSVParser) Untyped(data []byte) ([]map[string]interface{}, error) {
	c.isTyped = false
	return c.parseToCSV(data)
}

// Typed unmarshals the typed data into a slice of map[string]interface{}
//
// In this case, the second column of the csv must contain the field types, otherwise it will throw an error
func (c *CSVParser) Typed(data []byte) ([]map[string]interface{}, error) {
	c.isTyped = true
	return c.parseToCSV(data)
}

// checkForNilOrDefault checks if the runes are set.
// If the runes are not set, the default values are used.
//
// Default values:
//  comma: ','
//  comment: '#'
func (c *CSVParser) checkForNilOrDefault() {
	if c.Comma == *new(rune) {
		c.Comma = ','
	}

	if c.Comment == *new(rune) {
		c.Comment = '#'
	}
}

// readCSV delegates the read command to csv.NewReader (stdlib) and writes it to a two-dimensional string slice that is returned.
func (c *CSVParser) readCSV(data []byte) ([][]string, error) {
	csvR := csv.NewReader(bytes.NewReader(data))
	csvR.Comma = c.Comma
	csvR.Comment = c.Comment
	csvR.TrimLeadingSpace = c.TrimLeadingSpace
	csvR.FieldsPerRecord = -1
	csvR.LazyQuotes = true

	records, err := csvR.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// parseToCSV extracts the header information from the byte slice and generates a map based on the format (typed or untyped).
func (c *CSVParser) parseToCSV(data []byte) ([]map[string]interface{}, error) {
	c.checkForNilOrDefault()

	records, err := c.readCSV(data)
	if err != nil {
		return nil, err
	}

	var headerInfo map[int]field
	if c.isTyped {
		if len(records) < 2 {
			return nil, ErrDataIsNil
		}

		headerInfo = c.extractHeaderInformation(records[0], records[1])
		records = records[2:]
	} else {
		if len(records) < 1 {
			return nil, ErrDataIsNil
		}

		headerInfo = c.extractHeaderInformation(records[0], nil)
		records = records[1:]
	}

	return c.csvToMap(headerInfo, records)
}

// extractHeaderInformation reads the header information and returns it as map of field
func (c *CSVParser) extractHeaderInformation(names, types []string) map[int]field {
	headFields := map[int]field{}

	// extract field names
	for idx, value := range names {
		headFields[idx] = field{
			Name: value,
		}
	}

	// extract field types
	for idx, value := range types {
		field := headFields[idx]
		field.Type = value
		headFields[idx] = field
	}

	return headFields
}

// csvToMap builds the data columns based on the typed or untyped fields
func (c *CSVParser) csvToMap(headerInfo map[int]field, records [][]string) ([]map[string]interface{}, error) {
	rslt := []map[string]interface{}{}

	// skip first row
	for _, value := range records {
		skipColumn := true

		myColumn := make(map[string]interface{})
		for idx, v2 := range value {
			if len(headerInfo) < idx {
				// the column contains more data than we expected, break out of it
				break
			}

			// checks if the first entry of the row and the first character of the string matches the comment character.
			// If it matches, this row is skipped.
			// This is necessary because csvR.ReadAll() ignores some cases that contain such a comment rune
			if idx == 0 && len(v2) > 0 {
				if rune(v2[0]) == c.Comment {
					// the column contains the comment rune, skip it
					break
				}
			}

			// check whether v2 contains a value or not
			// set skip column to false, if a value was set
			if len(v2) > 0 {
				skipColumn = false
			}

			// check whether isTyped is true, the header info is not set and skip columns is set
			// then this row should be skipped
			if c.isTyped && headerInfo[idx].Type == "" && c.SkipEmptyColumns {
				continue
			}

			// check whether the type was set for the row
			if headerInfo[idx].Type != "" {
				// toTyped returns the
				typed, err := c.toTyped(v2, strings.TrimPrefix(headerInfo[idx].Type, "*"), strings.HasPrefix(headerInfo[idx].Type, "*"))
				if err != nil {
					return nil, err
				}

				// type is not a pointer
				myColumn[headerInfo[idx].Name] = typed
				continue
			}

			myColumn[headerInfo[idx].Name] = v2
		}
		if !skipColumn {
			rslt = append(rslt, myColumn)
		}
	}

	return rslt, nil
}

// toTyped takes the value and the format and converts the value into the desired format.
func (c *CSVParser) toTyped(value, format string, isPointer bool) (interface{}, error) {
	switch format {
	case "string":
		if value == "" && !isPointer {
			return "", nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		if isPointer {
			return &value, nil
		}

		return value, nil
	case "int64":
		if value == "" && !isPointer {
			return int64(0), nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		val, err := strconv.ParseInt(value, 10, 64)
		if isPointer {
			return &val, err
		}

		return val, err
	case "int":
		if value == "" && !isPointer {
			return int(0), nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		val, err := strconv.Atoi(value)
		if isPointer {
			return &val, err
		}

		return val, err
	case "float64":
		if value == "" && !isPointer {
			return float64(0), nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		val, err := strconv.ParseFloat(value, 64)
		if isPointer {
			return &val, err
		}
		return val, err
	case "bool":
		if value == "" && !isPointer {
			return false, nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		val, err := strconv.ParseBool(value)
		if isPointer {
			return &val, err
		}
		return val, err
	case "string,array":
		if value == "" && !isPointer {
			return []string{}, nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		records, err := c.readCSV([]byte(value))
		if err != nil {

			return nil, err
		}

		//Check if we only have one row. If not return error
		if len(records) > 1 {
			return nil, ErrOnlyOneRowIsAllowedForStringArray
		}

		retArray := make([]string, 0)
		retArray = append(retArray, records[0]...)

		if isPointer {
			return &retArray, nil
		}

		return retArray, nil
	case "int64,array":
		if value == "" && !isPointer {
			return []int64{}, nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		records, err := c.readCSV([]byte(value))
		if err != nil {
			return nil, err
		}

		//Check if we only have one row. If not return error
		if len(records) > 1 {
			return nil, ErrOnlyOneRowIsAllowedForInt64Array
		}

		retArray := make([]int64, 0)
		for _, v := range records[0] {
			vi := int64(0)
			if v != "" {
				vi, err = strconv.ParseInt(strings.TrimSpace(v), 10, 64)
				if err != nil {
					return nil, err
				}
			}
			retArray = append(retArray, vi)
		}

		if isPointer {
			return &retArray, nil
		}

		return retArray, nil
	case "float64,array":
		if value == "" && !isPointer {
			return []float64{}, nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		records, err := c.readCSV([]byte(value))
		if err != nil {
			return nil, err
		}

		//Check if we only have one row. If not return error
		if len(records) > 1 {

			return nil, ErrOnlyOneRowIsAllowedForFloat64Array
		}

		retArray := make([]float64, 0)
		for _, v := range records[0] {
			vi := float64(0)
			if v != "" {
				vi, err = strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					return nil, err
				}
			}
			retArray = append(retArray, vi)
		}

		if isPointer {
			return &retArray, nil
		}

		return retArray, nil
	case "bool,array":
		if value == "" && !isPointer {
			return []bool{}, nil
		} else if value == "" && isPointer {
			return nil, nil
		}

		records, err := c.readCSV([]byte(value))
		if err != nil {
			return nil, err
		}

		//Check if we only have one row. If not return error
		if len(records) > 1 {
			return nil, ErrOnlyOneRowIsAllowedForBoolArray
		}

		retArray := make([]bool, 0)
		for _, v := range records[0] {
			retArray = append(retArray, strings.TrimSpace(v) == "true")
		}

		if isPointer {
			return &retArray, nil
		}

		return retArray, nil
	case "json":
		if value == "" {
			return nil, nil
		}

		var data interface{}
		err := json.Unmarshal([]byte(value), &data)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInEmbeddedJSON, err)
		}

		if isPointer {
			p := reflect.New(reflect.TypeOf(data))
			p.Elem().Set(reflect.ValueOf(data))
			return p.Interface(), nil
		}

		return data, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, format)
	}
}
