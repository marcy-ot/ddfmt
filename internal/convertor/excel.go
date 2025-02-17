package convertor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/marcy-ot/ddfmt/internal/config"
	"github.com/xuri/excelize/v2"
)

// 変換対象の interface
type Convertible interface {
	Read(w io.Writer, path string, config *config.Config) error
	Header() []string
	Rows() [][]string
}

func NewConvertable(fileName string) Convertible {
	// TODO: filename 拡張子から返却値を変更
	switch filepath.Ext(fileName) {
	case ".xlsx":
		return &Excel{}
	default:
		return &Excel{}
	}
}

type Excel struct {
	header []string
	rows   [][]string
}

func (ex *Excel) Read(w io.Writer, path string, config *config.Config) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(w, "error don't exist excel file: %v\n", err)
		return err
	}

	file, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Fprintf(w, "error open excel: %v\n", err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("err file close")
		}
	}()

	rows, err := file.GetRows(config.SheetName)
	if err != nil {
		fmt.Fprintf(w, "error get excel rows: %v\n", err)
		return err
	}

	header := rows[0]
	body := rows[1:]
	ex.header = header
	ex.rows = body

	return nil
}

func (ex *Excel) Header() []string {
	return ex.header
}

func (ex *Excel) Rows() [][]string {
	return ex.rows
}
