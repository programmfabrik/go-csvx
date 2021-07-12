package csvx

import (
	"reflect"
	"testing"
)

func TestCSV_Untyped(t *testing.T) {
	type args struct {
		data          []byte
		withSkipEmpty bool
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "test_untyped_1",
			args: args{
				data: []byte(`
				foo,bar
				first,second`),
				withSkipEmpty: false,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},
		{
			name: "test_untyped_2",
			args: args{
				data: []byte(`
				foo,bar
				first,second
				third,fourth`),
				withSkipEmpty: false,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
				{
					"foo": "third",
					"bar": "fourth",
				},
			},
			wantErr: false,
		},
		{
			name: "test_untyped_with_empty_column",
			args: args{
				data: []byte(`foo,bar
				,
				first,second
				third,fourth`),
				withSkipEmpty: false,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
				{
					"foo": "third",
					"bar": "fourth",
				},
			},
			wantErr: false,
		},
		{
			name: "test_untyped_with_empty_column_row",
			args: args{
				data: []byte(`foo,placeholder,bar
				,,
				first,,second
				third,,fourth`),
				withSkipEmpty: false,
			},
			want: []map[string]interface{}{
				{
					"foo":         "first",
					"placeholder": "",
					"bar":         "second",
				},
				{
					"foo":         "third",
					"placeholder": "",
					"bar":         "fourth",
				},
			},
			wantErr: false,
		},
		{
			name: "test_with_comment",
			args: args{
				data: []byte(`foo,bar
				first,second
				#third,fourth`),
				withSkipEmpty: false,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},
		{
			name: "test_with_comment_and_space",
			args: args{
				data: []byte(`foo,bar
				first,second
				# third,fourth`),
				withSkipEmpty: false,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},

		{
			name: "test_with_clear",
			args: args{
				data: []byte(`
				foo,bar

				# third,fourth`),
				withSkipEmpty: false,
			},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := CSVParser{Comma: ',', Comment: '#', TrimLeadingSpace: true, SkipEmptyColumns: tt.args.withSkipEmpty}

			rslt, err := csv.Untyped(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestToMap() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestToMap() is not equal. \ngot = %+#v\nwant = %+#v", rslt, tt.want)
			}
		})
	}
}

func TestCSV_Typed(t *testing.T) {
	type args struct {
		data             []byte
		skipEmptyColumns bool
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "test_typed_1",
			args: args{
				data: []byte(`
				foo,bar
				string,string
				first,second`),
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_2",
			args: args{
				data: []byte(`
				foo,bar
				string,string
				first,second
				third,fourth`),
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
				{
					"foo": "third",
					"bar": "fourth",
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_3",
			args: args{
				data: []byte(`
				foo,bar
				string,int64
				first,10
				third,20`),
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": int64(10),
				},
				{
					"foo": "third",
					"bar": int64(20),
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_json",
			args: args{
				data: []byte(`
				foo,bar,subtype
				string,int,json
				first,10,{"key": 10}`),
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": int(10),
					"subtype": map[string]interface{}{
						"key": float64(10),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_json_ptr",
			args: args{
				data: []byte(`
				foo,bar,subtype
				*string,*int,*json
				first,10,{"key": 10}`),
			},
			want: []map[string]interface{}{
				{
					"foo":     func(msg string) *string { return &msg }("first"),
					"bar":     func(i int) *int { return &i }(int(10)),
					"subtype": &map[string]interface{}{"key": float64(10)},
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_empty_row",
			args: args{
				data: []byte(`
				foo,bar
				string,int64
				,
				first,10`),
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": int64(10),
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_empty_type_column",
			args: args{
				data: []byte(`
				foo,placeholder,bar
				string,,int64
				,,
				first,test,10`),
				skipEmptyColumns: true,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": int64(10),
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_clear",
			args: args{
				data: []byte(`
				foo,placeholder,bar
				string,,int64

				first,test,10`),
				skipEmptyColumns: true,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": int64(10),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := CSVParser{Comma: ',', Comment: '#', TrimLeadingSpace: true, SkipEmptyColumns: tt.args.skipEmptyColumns}

			rslt, err := csv.Typed(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestToTypedMap() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestToTypedMap() is not equal.\ngot  = %#v\nwant = %#v", rslt, tt.want)
			}
		})
	}
}

func TestCSV_checkForNilOrDefault(t *testing.T) {
	type args struct {
		csv *CSVParser
	}
	tests := []struct {
		name string
		args args
		want *CSVParser
	}{
		{
			name: "test_success_with_comma_and_comment",
			args: args{
				csv: &CSVParser{
					Comma:            ',',
					Comment:          '#',
					TrimLeadingSpace: true,
				},
			},
			want: &CSVParser{
				Comma:            ',',
				Comment:          '#',
				TrimLeadingSpace: true,
			},
		},
		{
			name: "test_success_with_comment",
			args: args{
				csv: &CSVParser{
					Comment:          '#',
					TrimLeadingSpace: true,
				},
			},
			want: &CSVParser{
				Comma:            ',',
				Comment:          '#',
				TrimLeadingSpace: true,
			},
		},
		{
			name: "test_success_with_comma",
			args: args{
				csv: &CSVParser{
					Comma:            ',',
					TrimLeadingSpace: true,
				},
			},
			want: &CSVParser{
				Comma:            ',',
				Comment:          '#',
				TrimLeadingSpace: true,
			},
		},
		{
			name: "test_success_without_comma_and_comment",
			args: args{
				csv: &CSVParser{
					TrimLeadingSpace: true,
				},
			},
			want: &CSVParser{
				Comma:            ',',
				Comment:          '#',
				TrimLeadingSpace: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.csv.checkForNilOrDefault()

			if !reflect.DeepEqual(*tt.args.csv, *tt.want) {
				t.Errorf("TestCheckForNilOrDefault() is not equal. \ngot = %+#v\nwant = %+#v", *tt.args.csv, *tt.want)
			}
		})
	}
}

func TestCSV_readCSV(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name: "test_one_line",
			args: args{
				data: []byte(`foo,bar`),
			},
			want:    [][]string{{"foo", "bar"}},
			wantErr: false,
		},
		{
			name: "test_with_one_empty",
			args: args{
				data: []byte(`
				foo,bar
				,`),
			},
			want:    [][]string{{"foo", "bar"}, {"", ""}},
			wantErr: false,
		},
		{
			name: "test_with_one_clear",
			args: args{
				data: []byte(`
				foo,bar

				first,second`),
			},
			want:    [][]string{{"foo", "bar"}, {"first", "second"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := CSVParser{Comma: ',', Comment: '#', TrimLeadingSpace: true}

			strS, err := csv.readCSV(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReadCSV() got error: %v", err)
			}

			if !reflect.DeepEqual(strS, tt.want) {
				t.Errorf("TestReadCSV() is not equal. \ngot = %+#v\nwant = %+#v", strS, tt.want)
			}
		})
	}
}

func TestCSV_parseToCSV(t *testing.T) {
	type args struct {
		data    []byte
		isTyped bool
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "test_untyped_1",
			args: args{
				data: []byte(`
				foo,bar
				first,second`),
				isTyped: false,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},
		{
			name: "test_untyped_2",
			args: args{
				data: []byte(`
				foo,bar
				first,second
				third,fourth`),
				isTyped: false,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
				{
					"foo": "third",
					"bar": "fourth",
				},
			},
			wantErr: false,
		},

		{
			name: "test_typed_1",
			args: args{
				data: []byte(`
				foo,bar
				string,string
				first,second`),
				isTyped: true,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_2",
			args: args{
				data: []byte(`
				foo,bar
				string,string
				first,second
				third,fourth`),
				isTyped: true,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
				{
					"foo": "third",
					"bar": "fourth",
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_clear_column",
			args: args{
				data: []byte(`
				foo,bar
				string,string

				first,second`),
				isTyped: true,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},
		{
			name: "test_typed_empty_column",
			args: args{
				data: []byte(`
				foo,bar
				string,string
				,
				first,second`),
				isTyped: true,
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "second",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := CSVParser{Comma: ',', Comment: '#', TrimLeadingSpace: true}

			csv.isTyped = tt.args.isTyped

			rslt, err := csv.parseToCSV(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestParseToCSV() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestParseToCSV() is not equal. \ngot = %+#v\nwant = %+#v", rslt, tt.want)
			}
		})
	}
}

func TestCSV_extractHeaderInformation(t *testing.T) {
	type args struct {
		names []string
		types []string
	}
	tests := []struct {
		name string
		args args
		want map[int]field
	}{
		{
			name: "test_success_untyped",
			args: args{
				names: []string{"foo", "bar"},
				types: nil,
			},
			want: map[int]field{
				0: {
					Name: "foo",
					Type: "",
				},
				1: {
					Name: "bar",
					Type: "",
				},
			},
		},
		{
			name: "test_success_string_string",
			args: args{
				names: []string{"foo", "bar"},
				types: []string{"string", "string"},
			},
			want: map[int]field{
				0: {
					Name: "foo",
					Type: "string",
				},
				1: {
					Name: "bar",
					Type: "string",
				},
			},
		},
		{
			name: "test_success_string_int",
			args: args{
				names: []string{"foo", "bar"},
				types: []string{"string", "int"},
			},
			want: map[int]field{
				0: {
					Name: "foo",
					Type: "string",
				},
				1: {
					Name: "bar",
					Type: "int",
				},
			},
		},
		{
			name: "test_success_string_json",
			args: args{
				names: []string{"foo", "bar"},
				types: []string{"string", "json"},
			},
			want: map[int]field{
				0: {
					Name: "foo",
					Type: "string",
				},
				1: {
					Name: "bar",
					Type: "json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := CSVParser{Comma: ',', Comment: '#', TrimLeadingSpace: true}

			rslt := csv.extractHeaderInformation(tt.args.names, tt.args.types)

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestExtractHeaderInformations() is not equal. \ngot = %+#v\nwant = %+#v", rslt, tt.want)
			}
		})
	}
}

func TestCSV_csvToMap(t *testing.T) {
	type args struct {
		headerInfo map[int]field
		records    [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "test_untyped_1",
			args: args{
				headerInfo: map[int]field{
					0: {Name: "foo"},
					1: {Name: "bar"},
				},
				records: [][]string{
					{
						"first",
						"first",
					},
				},
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "first",
				},
			},
		},
		{
			name: "test_untyped_2",
			args: args{
				headerInfo: map[int]field{
					0: {Name: "foo"},
					1: {Name: "bar"},
				},
				records: [][]string{
					{
						"first",
						"first",
					},
					{
						"second",
						"second",
					},
				},
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": "first",
				},
				{
					"foo": "second",
					"bar": "second",
				},
			},
		},

		{
			name: "test_typed_1",
			args: args{
				headerInfo: map[int]field{
					0: {Name: "foo", Type: "string"},
					1: {Name: "bar", Type: "int64"},
				},
				records: [][]string{
					{
						"first",
						"10",
					},
				},
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": int64(10),
				},
			},
		},
		{
			name: "test_typed_2",
			args: args{
				headerInfo: map[int]field{
					0: {Name: "foo", Type: "string"},
					1: {Name: "bar", Type: "int64"},
				},
				records: [][]string{
					{
						"first",
						"10",
					},
					{
						"second",
						"100",
					},
				},
			},
			want: []map[string]interface{}{
				{
					"foo": "first",
					"bar": int64(10),
				},
				{
					"foo": "second",
					"bar": int64(100),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := CSVParser{Comma: ',', Comment: '#', TrimLeadingSpace: true}

			rslt, err := csv.csvToMap(tt.args.headerInfo, tt.args.records)
			if err != nil {
				t.Errorf("TestCsvToMap() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestCsvToMap() is not equal. \ngot = %+#v\nwant = %+#v", rslt, tt.want)
			}
		})
	}
}

func TestCSV_toTyped(t *testing.T) {
	type args struct {
		value, format string
		isPointerType bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test_string",
			args: args{
				value:         "hello",
				format:        "string",
				isPointerType: false,
			},
			want:    "hello",
			wantErr: false,
		},
		{
			name: "test_string_array",
			args: args{
				value:         "hello,world",
				format:        "string,array",
				isPointerType: false,
			},
			want:    []string{"hello", "world"},
			wantErr: false,
		},

		{
			name: "test_bool",
			args: args{
				value:         "false",
				format:        "bool",
				isPointerType: false,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "test_bool_array",
			args: args{
				value:         "false,false",
				format:        "bool,array",
				isPointerType: false,
			},
			want:    []bool{false, false},
			wantErr: false,
		},

		{
			name: "test_int",
			args: args{
				value:         "10",
				format:        "int",
				isPointerType: false,
			},
			want:    10,
			wantErr: false,
		},

		{
			name: "test_int64",
			args: args{
				value:         "10",
				format:        "int64",
				isPointerType: false,
			},
			want:    int64(10),
			wantErr: false,
		},
		{
			name: "test_int64_array",
			args: args{
				value:         "10,11",
				format:        "int64,array",
				isPointerType: false,
			},
			want:    []int64{10, 11},
			wantErr: false,
		},

		{
			name: "test_float64",
			args: args{
				value:         "10.1",
				format:        "float64",
				isPointerType: false,
			},
			want:    float64(10.1),
			wantErr: false,
		},
		{
			name: "test_float64_array",
			args: args{
				value:         "10.1,11.4",
				format:        "float64,array",
				isPointerType: false,
			},
			want:    []float64{10.1, 11.4},
			wantErr: false,
		},

		{
			name: "test_json",
			args: args{
				value:         `{"name": "value"}`,
				format:        "json",
				isPointerType: false,
			},
			want:    map[string]interface{}{"name": "value"},
			wantErr: false,
		},

		{
			name: "test_string_json",
			args: args{
				value:         `{"name": "value"}`,
				format:        "string,json",
				isPointerType: false,
			},
			want:    nil,
			wantErr: true,
		},

		// pointer

		{
			name: "test_string_ptr",
			args: args{
				value:         "hello",
				format:        "string",
				isPointerType: true,
			},
			want: func(i string) *string {
				return &i
			}("hello"),
			wantErr: false,
		},
		{
			name: "test_string_array_ptr",
			args: args{
				value:         "hello,world",
				format:        "string,array",
				isPointerType: true,
			},
			want:    &[]string{"hello", "world"},
			wantErr: false,
		},

		{
			name: "test_bool_ptr",
			args: args{
				value:         "false",
				format:        "bool",
				isPointerType: true,
			},
			want: func(i bool) *bool {
				return &i
			}(false),
			wantErr: false,
		},
		{
			name: "test_bool_array_ptr",
			args: args{
				value:         "false,false",
				format:        "bool,array",
				isPointerType: true,
			},
			want:    &[]bool{false, false},
			wantErr: false,
		},

		{
			name: "test_int_ptr",
			args: args{
				value:         "10",
				format:        "int",
				isPointerType: true,
			},
			want: func(i int) *int {
				return &i
			}(10),
			wantErr: false,
		},

		{
			name: "test_int64_ptr",
			args: args{
				value:         "10",
				format:        "int64",
				isPointerType: true,
			},
			want: func(i int64) *int64 {
				return &i
			}(10),
			wantErr: false,
		},
		{
			name: "test_int64_array_ptr",
			args: args{
				value:         "10,11",
				format:        "int64,array",
				isPointerType: true,
			},
			want:    &[]int64{10, 11},
			wantErr: false,
		},

		{
			name: "test_float64_ptr",
			args: args{
				value:         "10.1",
				format:        "float64",
				isPointerType: true,
			},
			want: func(i float64) *float64 {
				return &i
			}(10.1),
			wantErr: false,
		},
		{
			name: "test_float64_array_ptr",
			args: args{
				value:         "10.1,11.4",
				format:        "float64,array",
				isPointerType: true,
			},
			want:    &[]float64{10.1, 11.4},
			wantErr: false,
		},

		{
			name: "test_json_ptr",
			args: args{
				value:         `{"name": "value"}`,
				format:        "json",
				isPointerType: true,
			},
			want:    &map[string]interface{}{"name": "value"},
			wantErr: false,
		},

		{
			name: "test_prt_string_json",
			args: args{
				value:         `value,{"name": "value"}`,
				format:        "string,json",
				isPointerType: true,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := CSVParser{Comma: ',', Comment: '#', TrimLeadingSpace: true}

			rslt, err := csv.toTyped(tt.args.value, tt.args.format, tt.args.isPointerType)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestToTyped() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestToTyped() is not equal. \ngot = %+#v\nwant = %+#v", rslt, tt.want)
			}
		})
	}
}
