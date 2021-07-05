package csvx

import (
	"reflect"
	"testing"
)

func TestNewCSV(t *testing.T) {
	type args struct {
		comma, comment   rune
		trimLeadingSpace bool
	}
	tests := []struct {
		name string
		args args
		want *CSV
	}{
		{
			name: "test_success",
			args: args{
				comma:            ',',
				comment:          '#',
				trimLeadingSpace: true,
			},
			want: &CSV{
				comma:            ',',
				comment:          '#',
				trimLeadingSpace: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCSV(tt.args.comma, tt.args.comment, tt.args.trimLeadingSpace)
			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("NewCSV() does not equal. got = %v, want = %v", *got, *tt.want)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	type args struct {
		data []byte
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
				data: []byte("foo,bar\nfirst,second"),
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
				data: []byte("foo,bar\nfirst,second\nthird,fourth"),
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := NewCSV(',', '#', true)

			rslt, err := csv.ToMap(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestToMap() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestToMap() does not equal. got = %v, want = %v", rslt, tt.want)
			}
		})
	}
}

func TestToTypedMap(t *testing.T) {
	type args struct {
		data []byte
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
				data: []byte("foo,bar\nstring,string\nfirst,second"),
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
				data: []byte("foo,bar\nstring,string\nfirst,second\nthird,fourth"),
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
				data: []byte("foo,bar\nstring,int64\nfirst,10\nthird,20"),
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := NewCSV(',', '#', true)

			rslt, err := csv.ToTypedMap(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestToTypedMap() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestToTypedMap() does not equal. got = %v, want = %v", rslt, tt.want)
			}
		})
	}
}

func TestCheckForNilOrDefault(t *testing.T) {
	type args struct {
		csv *CSV
	}
	tests := []struct {
		name string
		args args
		want *CSV
	}{
		{
			name: "test_success_1",
			args: args{
				csv: &CSV{
					comma:            ',',
					comment:          '#',
					trimLeadingSpace: true,
				},
			},
			want: &CSV{
				comma:            ',',
				comment:          '#',
				trimLeadingSpace: true,
			},
		},
		{
			name: "test_success_2",
			args: args{
				csv: &CSV{
					comment:          '#',
					trimLeadingSpace: true,
				},
			},
			want: &CSV{
				comma:            ',',
				comment:          '#',
				trimLeadingSpace: true,
			},
		},
		{
			name: "test_success_3",
			args: args{
				csv: &CSV{
					comma:            ',',
					trimLeadingSpace: true,
				},
			},
			want: &CSV{
				comma:            ',',
				comment:          '#',
				trimLeadingSpace: true,
			},
		},
		{
			name: "test_success_4",
			args: args{
				csv: &CSV{
					trimLeadingSpace: true,
				},
			},
			want: &CSV{
				comma:            ',',
				comment:          '#',
				trimLeadingSpace: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.csv.checkForNilOrDefault()

			if !reflect.DeepEqual(*tt.args.csv, *tt.want) {
				t.Errorf("TestCheckForNilOrDefault() does not equal. got = %v, want = %v", *tt.args.csv, *tt.want)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
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
			name: "test_success",
			args: args{
				data: []byte("foo,bar"),
			},
			want:    [][]string{{"foo", "bar"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := NewCSV(',', '#', true)

			strS, err := csv.readCSV(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReadCSV() got error: %v", err)
			}

			if !reflect.DeepEqual(strS, tt.want) {
				t.Errorf("TestReadCSV() does not equal. got = %v, want = %v", strS, tt.want)
			}
		})
	}
}

func TestParseToCSV(t *testing.T) {
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
				data:    []byte("foo,bar\nfirst,second"),
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
				data:    []byte("foo,bar\nfirst,second\nthird,fourth"),
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
				data:    []byte("foo,bar\nstring,string\nfirst,second"),
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
				data:    []byte("foo,bar\nstring,string\nfirst,second\nthird,fourth"),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := NewCSV(',', '#', true)

			csv.isTyped = tt.args.isTyped

			rslt, err := csv.parseToCSV(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestParseToCSV() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestParseToCSV() does not equal. got = %v, want = %v", rslt, tt.want)
			}
		})
	}
}

func TestExtractHeaderInformations(t *testing.T) {
	type args struct {
		data [2][]string
	}
	tests := []struct {
		name string
		args args
		want map[int]field
	}{
		{
			name: "test_success_untyped",
			args: args{
				data: [2][]string{{"foo", "bar"}},
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
			name: "test_success_typed_1",
			args: args{
				data: [2][]string{{"foo", "bar"}, {"string", "string"}},
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
			name: "test_success_typed_2",
			args: args{
				data: [2][]string{{"foo", "bar"}, {"string", "int"}},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := NewCSV(',', '#', true)

			rslt := csv.extractHeaderInformations(tt.args.data)

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestExtractHeaderInformations() does not equal. got = %v, want = %v", rslt, tt.want)
			}
		})
	}
}

func TestCsvToMap(t *testing.T) {
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
			csv := NewCSV(',', '#', true)

			rslt, err := csv.csvToMap(tt.args.headerInfo, tt.args.records)
			if err != nil {
				t.Errorf("TestCsvToMap() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestCsvToMap() does not equal. got = %v, want = %v", rslt, tt.want)
			}
		})
	}
}

func TestToTyped(t *testing.T) {
	type args struct {
		value, format string
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
				value:  "hello",
				format: "string",
			},
			want:    "hello",
			wantErr: false,
		},
		{
			name: "test_string_array",
			args: args{
				value:  "hello,world",
				format: "string,array",
			},
			want:    []string{"hello", "world"},
			wantErr: false,
		},

		{
			name: "test_bool",
			args: args{
				value:  "false",
				format: "bool",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "test_bool_array",
			args: args{
				value:  "false,false",
				format: "bool,array",
			},
			want:    []bool{false, false},
			wantErr: false,
		},

		{
			name: "test_int",
			args: args{
				value:  "10",
				format: "int",
			},
			want:    10,
			wantErr: false,
		},

		{
			name: "test_int64",
			args: args{
				value:  "10",
				format: "int64",
			},
			want:    int64(10),
			wantErr: false,
		},
		{
			name: "test_int64_array",
			args: args{
				value:  "10,11",
				format: "int64,array",
			},
			want:    []int64{10, 11},
			wantErr: false,
		},

		{
			name: "test_float64",
			args: args{
				value:  "10.1",
				format: "float64",
			},
			want:    float64(10.1),
			wantErr: false,
		},
		{
			name: "test_float64_array",
			args: args{
				value:  "10.1,11.4",
				format: "float64,array",
			},
			want:    []float64{10.1, 11.4},
			wantErr: false,
		},

		{
			name: "test_json",
			args: args{
				value:  `{"name": "value"}`,
				format: "json",
			},
			want:    map[string]interface{}{"name": "value"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csv := NewCSV(',', '#', true)

			rslt, err := csv.toTyped(tt.args.value, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestToTyped() received error = %v", err)
			}

			if !reflect.DeepEqual(rslt, tt.want) {
				t.Errorf("TestToTyped() does not equal. got = %v, want = %v", rslt, tt.want)
			}
		})
	}
}
