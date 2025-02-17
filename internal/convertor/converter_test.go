package convertor

import (
	"io"
	"os"
	"testing"

	"github.com/marcy-ot/ddfmt/internal/config"
	"github.com/stretchr/testify/assert"
)

func seedConvertable(header []string, rows [][]string) Convertible {
	return seedConvertableStruct{
		header: header,
		rows:   rows,
	}
}

type seedConvertableStruct struct {
	header []string
	rows   [][]string
}

func (c seedConvertableStruct) Read(w io.Writer, path string, config *config.Config) error {
	return nil
}
func (c seedConvertableStruct) Header() []string { return c.header }
func (c seedConvertableStruct) Rows() [][]string { return c.rows }

func TestConvert(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		convertable Convertible
		want        OutputData
	}{
		{
			name:   "正常系_Configなし",
			config: config.DefaultConfig(),
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product2", "40"},
					{"3", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "product1", "20"},
						{"2", "product2", "40"},
						{"3", "product3", "80"},
					},
				},
			},
		},
		{
			name: "正常系_Configあり_Unique_single",
			config: &config.Config{
				UniqueCols: []int{2},
			},
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product3", "40"},
					{"3", "product3", "80"},
					{"3", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "product1", "20"},
						{"2", "product3", "40"},
					},
				},
			},
		},
		{
			name: "正常系_Configあり_Unique_multi",
			config: &config.Config{
				UniqueCols: []int{2, 3},
			},
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product3", "40"},
					{"3", "product3", "80"},
					{"4", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "product1", "20"},
						{"2", "product3", "40"},
						{"3", "product3", "80"},
					},
				},
			},
		},
		{
			name: "正常系_Configあり_OverWrite",
			config: &config.Config{
				UniqueCols:    []int{2, 3},
				OverwriteCols: []config.ColumnValue{{Col: 2, Val: "computer"}, {Col: 3, Val: "9999"}},
			},
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product3", "40"},
					{"3", "product3", "80"},
					{"4", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "computer", "9999"},
						{"2", "computer", "9999"},
						{"3", "computer", "9999"},
					},
				},
			},
		},
		{
			name: "正常系_Configあり_Aggregate",
			config: &config.Config{
				UniqueCols:  []int{2, 3},
				DistinctCol: 2,
			},
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product3", "40"},
					{"3", "product3", "80"},
					{"4", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "product1", "20"},
						{"2", "product3", "40"},
						{"3", "product3", "80"},
					},
				},
				Aggregate: []string{"product1", "product3"},
			},
		},
		{
			name: "正常系_Configあり_Divide",
			config: &config.Config{
				UniqueCols:  []int{2, 3},
				DistinctCol: 2,
				FileSplit: struct {
					Row int `yaml:"row"`
				}{Row: 2},
			},
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product3", "40"},
					{"3", "product3", "80"},
					{"4", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "product1", "20"},
						{"2", "product3", "40"},
					}, {
						{"3", "product3", "80"},
					},
				},
				Aggregate: []string{"product1", "product3"},
			},
		},
		{
			name: "正常系_Configあり_Message",
			config: &config.Config{
				UniqueCols:  []int{2, 3},
				DistinctCol: 2,
				FileSplit: struct {
					Row int `yaml:"row"`
				}{Row: 2},
				CompletionMessage: "output message",
			},
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product3", "40"},
					{"3", "product3", "80"},
					{"4", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "product1", "20"},
						{"2", "product3", "40"},
					}, {
						{"3", "product3", "80"},
					},
				},
				Aggregate: []string{"product1", "product3"},
				Message:   "output message",
			},
		},
		{
			name: "正常系_Configあり_Message_embed_aggregate",
			config: &config.Config{
				UniqueCols:  []int{2, 3},
				DistinctCol: 2,
				FileSplit: struct {
					Row int `yaml:"row"`
				}{Row: 2},
				CompletionMessage: "output message {$distinct_column} outputs.",
			},
			convertable: seedConvertable(
				[]string{"Product ID", "Product Name", "Stock Quantity"},
				[][]string{
					{"1", "product1", "20"},
					{"2", "product3", "40"},
					{"3", "product3", "80"},
					{"4", "product3", "80"},
				},
			),
			want: OutputData{
				Header: []string{"Product ID", "Product Name", "Stock Quantity"},
				FileData: [][][]string{
					{
						{"1", "product1", "20"},
						{"2", "product3", "40"},
					}, {
						{"3", "product3", "80"},
					},
				},
				Aggregate: []string{"product1", "product3"},
				Message:   "output message \nproduct1\nproduct3\n outputs.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			convertor := NewConvertor(tt.convertable)
			convertor.SetConfig(os.Stderr, tt.config)

			actual := convertor.Convert()
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestDataDivide(t *testing.T) {
	type args struct {
		fileSplitRow int
		data         [][][]string
	}
	tests := []struct {
		name string
		arg  args
		want [][][]string
	}{
		{
			name: "正常系_分割なし",
			arg: args{
				fileSplitRow: 0,
				data: [][][]string{
					{{"col1", "col2", "col3"}},
				},
			},
			want: [][][]string{
				{{"col1", "col2", "col3"}},
			},
		},
		{
			name: "正常系_分割なし_レコード数指定あり",
			arg: args{
				fileSplitRow: 5,
				data: [][][]string{
					{{"col1", "col2", "col3"}},
				},
			},
			want: [][][]string{
				{{"col1", "col2", "col3"}},
			},
		},
		{
			name: "正常系_分割あり_余りなし",
			arg: args{
				fileSplitRow: 2,
				data: [][][]string{
					{
						{"col1", "col2", "col3"},
						{"col4", "col5", "col6"},
						{"col7", "col8", "col9"},
						{"col10", "col11", "col12"},
					},
				},
			},
			want: [][][]string{
				{
					{"col1", "col2", "col3"},
					{"col4", "col5", "col6"},
				},
				{
					{"col7", "col8", "col9"},
					{"col10", "col11", "col12"},
				},
			},
		},
		{
			name: "正常系_分割あり_余りあり",
			arg: args{
				fileSplitRow: 2,
				data: [][][]string{
					{
						{"col1", "col2", "col3"},
						{"col4", "col5", "col6"},
						{"col7", "col8", "col9"},
					},
				},
			},
			want: [][][]string{
				{
					{"col1", "col2", "col3"},
					{"col4", "col5", "col6"},
				},
				{
					{"col7", "col8", "col9"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.Config{
				FileSplit: struct {
					Row int `yaml:"row"`
				}{
					Row: tt.arg.fileSplitRow,
				},
			}

			var output OutputData
			output.FileData = tt.arg.data
			convertor := &Convertor{
				Config: config,
				Output: output,
			}

			convertor.dataDivide()

			assert.Equal(t, tt.want, convertor.Output.FileData)
		})
	}
}
