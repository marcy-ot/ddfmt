package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/marcy-ot/ddfmt/internal/config"
	"github.com/marcy-ot/ddfmt/internal/convertor"
)

type ExporterNumber int

func (en ExporterNumber) String() string {
	switch en {
	case Csv:
		return "csv"
	default:
		return ""
	}
}

func newExporterNumberFromString(v string) (ExporterNumber, error) {
	switch v {
	case "csv":
		return Csv, nil
	default:
		return -1, fmt.Errorf("undefined export file extension: %s", v)
	}
}

const (
	Csv ExporterNumber = iota
)

type Exporter interface {
	Export(fileName string) error
}

func NewExporter(config *config.Config, output convertor.OutputData, stderr io.Writer) Exporter {
	extension, err := newExporterNumberFromString(config.ExportFileExtension)
	if err != nil {
		fmt.Fprintln(stderr, err)
		fmt.Fprintln(stderr, "The specified export file extension was invalid, so the CSV format was automatically selected.")
		extension = Csv
	}
	switch extension {
	case Csv:
		return &csvExporter{
			exporterNumber: extension,
			output:         output,
			stderr:         stderr,
		}
	default:
		return &csvExporter{
			exporterNumber: extension,
			output:         output,
			stderr:         stderr,
		}
	}
}

type csvExporter struct {
	exporterNumber ExporterNumber
	output         convertor.OutputData
	stderr         io.Writer
}

func (ce *csvExporter) Export(fileName string) error {

	// CSV ファイル名
	csvFileName := fileName

	for i, rows := range ce.output.FileData {
		if i == 0 {
			csvFileName = fmt.Sprintf("%v", csvFileName)
		} else {
			csvFileName = fmt.Sprintf("%v_%d", csvFileName, i)
		}
		if err := ce.writeCsv(csvFileName, rows); err != nil {
			fmt.Fprintln(ce.stderr, err)
			return err
		}
	}
	return nil
}

func (ce *csvExporter) writeCsv(csvFileName string, rows [][]string) error {
	cf := fmt.Sprint(csvFileName, ".csv")
	f, err := os.Create(cf)
	if err != nil {
		return fmt.Errorf("error create csv file: %v\n", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(ce.stderr, "error csv file close : %v\n", err)
		}
	}()
	cw := csv.NewWriter(f)
	defer cw.Flush()

	cw.Write(ce.output.Header)
	for _, row := range rows {
		cw.Write(row)
	}

	return nil
}
