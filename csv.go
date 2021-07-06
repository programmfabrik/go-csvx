package csvx

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
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

type CSV struct {
	comma, comment   rune
	trimLeadingSpace bool

	isTyped bool
}

func NewCSV(comma, comment rune, trimLeadingSpace bool) *CSV {
	return &CSV{
		comma:            comma,
		comment:          comment,
		trimLeadingSpace: trimLeadingSpace,
	}
}

// ToMap unmarshals the data into a slice of map[string]interface{}
func (c *CSV) ToMap(data []byte) ([]map[string]interface{}, error) {
	c.isTyped = false
	return c.parseToCSV(data)
}

// ToTypedMap unmarshals the typed data into a slice of map[string]interface{}
//
// In this case, the second column of the csv must contain the field types, otherwise it will throw an error
func (c *CSV) ToTypedMap(data []byte) ([]map[string]interface{}, error) {
	c.isTyped = true
	return c.parseToCSV(data)
}

// checkForNilOrDefault checks if the runes are set.
// If the runes are not set, the default values are used.
//
// Default values:
//  comma: ','
//  comment: '#'
func (c *CSV) checkForNilOrDefault() {
	if c.comma == *new(rune) {
		c.comma = ','
	}

	if c.comment == *new(rune) {
		c.comment = '#'
	}
}

func (c *CSV) readCSV(data []byte) ([][]string, error) {
	csvR := csv.NewReader(bytes.NewReader(data))
	csvR.Comma = c.comma
	csvR.Comment = c.comment
	csvR.TrimLeadingSpace = c.trimLeadingSpace

	records, err := csvR.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (c *CSV) parseToCSV(data []byte) ([]map[string]interface{}, error) {
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

		headerInfo = c.extractHeaderInformations([2][]string{records[0], records[1]})
		records = records[2:]
	} else {
		if len(records) < 1 {
			return nil, ErrDataIsNil
		}

		headerInfo = c.extractHeaderInformations([2][]string{records[0]})
		records = records[1:]
	}

	return c.csvToMap(headerInfo, records)
}

// extractHeaderInformations reads the header information and returns it as map of field
func (c *CSV) extractHeaderInformations(data [2][]string) map[int]field {
	headFields := map[int]field{}

	for idx, value := range data[0] {
		headFields[idx] = field{
			Name: value,
		}
	}

	for idx, value := range data[1] {
		field := headFields[idx]
		field.Type = value
		headFields[idx] = field
	}

	return headFields
}

// csvToMap build the data columns based on the typed or untyped fields
func (c *CSV) csvToMap(headerInfo map[int]field, records [][]string) ([]map[string]interface{}, error) {
	rslt := []map[string]interface{}{}

	// skip first row
	for _, value := range records {
		myColumn := make(map[string]interface{})
		for idx, v2 := range value {
			if len(headerInfo) < idx {
				// the column contains more data than we expected, break out of it
				break
			}

			// is typed
			if headerInfo[idx].Type != "" {
				// is typed
				typed, err := c.toTyped(v2, headerInfo[idx].Type)
				if err != nil {
					return nil, err
				}
				myColumn[headerInfo[idx].Name] = typed
				continue
			}

			myColumn[headerInfo[idx].Name] = v2
		}
		rslt = append(rslt, myColumn)
	}

	return rslt, nil
}

func (c *CSV) toTyped(value, format string) (interface{}, error) {
	switch format {
	case "string":
		return value, nil
	case "int64":
		if value == "" {
			return int64(0), nil
		}
		return strconv.ParseInt(value, 10, 64)
	case "int":
		if value == "" {
			return int(0), nil
		}
		return strconv.Atoi(value)
	case "float64":
		if value == "" {
			return float64(0), nil
		}
		return strconv.ParseFloat(value, 64)
	case "bool":
		if value == "" {
			return false, nil
		}
		return strconv.ParseBool(value)
	case "string,array":
		if value == "" {
			return []string{}, nil
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
		return retArray, nil
	case "int64,array":
		if value == "" {
			return []int64{}, nil
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
		return retArray, nil
	case "float64,array":
		if value == "" {
			return []float64{}, nil
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
		return retArray, nil
	case "bool,array":
		if value == "" {
			return []bool{}, nil
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
		return data, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, format)
	}
}
