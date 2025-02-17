package convertor

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/marcy-ot/ddfmt/internal/config"
)

type OutputData struct {
	Header    []string
	FileData  [][][]string
	Aggregate []string
	Message   string
}

type Convertor struct {
	*config.Config
	Output OutputData
}

func NewConvertor(convertible Convertible) *Convertor {
	return &Convertor{
		Output: OutputData{
			Header:   convertible.Header(),
			FileData: [][][]string{convertible.Rows()},
		},
	}
}

func (con *Convertor) SetConfig(stderr io.Writer, config *config.Config) error {
	if err := con.validateConfig(config); err != nil {
		fmt.Fprintf(stderr, "error config validate: %v\n", err)
		return err
	}
	con.Config = config
	return nil
}

func (con *Convertor) validateConfig(config *config.Config) error {
	hlen := len(con.Output.Header)

	for _, ow := range config.OverwriteCols {
		if !(0 <= ow.Col-1 && ow.Col-1 < hlen) {
			return fmt.Errorf("overwrite_columns is out of range.\nvalue: %v", ow.Col)
		}
	}
	if 0 < config.DistinctCol {
		if hlen < config.DistinctCol {
			return fmt.Errorf("distinct_column is out of range.\nvalue: %v", config.DistinctCol)
		}
	}

	for _, c := range config.UniqueCols {
		if !(0 <= c-1 && c-1 < hlen) {
			return fmt.Errorf("unique_columns is out of range.\nvalue: %v", c)
		}
	}

	return nil
}

func (con *Convertor) Convert() OutputData {
	// unique
	con.uniqueColumns()
	// overwrite
	con.overWrite()
	// aggregate
	con.setAggregate()

	// devide
	con.dataDivide()

	// messsage set
	con.setMessage()

	return con.Output
}

func (con *Convertor) dataDivide() {
	// 分割無しの場合
	if con.FileSplit.Row == 0 || con.GteSplitRow(len(con.Output.FileData[0])) {
		return
	}

	splitTarget := con.Output.FileData[0]
	splitRowCount := con.FileSplit.Row
	divNum := len(splitTarget) / splitRowCount
	if 0 < divNum && len(splitTarget)%splitRowCount != 0 {
		// 割り切れる場合は分割数がオーバーしてしまうので -1
		// [1,2,3,4] / 3 の場合、分割数は 1 となるが実際には 2 である
		divNum++
	}
	newFileData := make([][][]string, divNum)
	var start, end int
	for i := 0; i < divNum; i++ {
		start = i * splitRowCount
		end = start + splitRowCount

		if len(splitTarget) < end {
			newFileData[i] = splitTarget[start:]
		} else {
			newFileData[i] = splitTarget[start:end]
		}
	}

	con.Output.FileData = newFileData
}

func (con *Convertor) setAggregate() {
	if con.DistinctCol == 0 {
		return
	}

	contain := func(source []string, target string) bool {
		for _, s := range source {
			if s == target {
				return true
			}
		}
		return false
	}

	var agg []string
	for _, rows := range con.Output.FileData {
		for _, row := range rows {
			c := row[con.DistinctCol-1]
			if !contain(agg, c) {
				agg = append(agg, c)
			}
		}
	}

	con.Output.Aggregate = agg
}

func (con *Convertor) setMessage() {
	if con.CompletionMessage == "" {
		return
	}
	val := "\n" + strings.Join(con.Output.Aggregate, "\n") + "\n"
	re := regexp.MustCompile(`\{\$distinct_column\}`)
	message := re.ReplaceAllString(con.CompletionMessage, val)

	con.Output.Message = message
}

func (con *Convertor) uniqueColumns() {
	if len(con.UniqueCols) == 0 {
		return
	}

	var newExcel [][]string
	for _, rows := range con.Output.FileData {
		for _, row := range rows {
			if con.isUnique(newExcel, row) {
				newExcel = append(newExcel, row)
			}
		}
	}

	con.Output.FileData = [][][]string{newExcel}
}

func (con *Convertor) isUnique(source [][]string, target []string) bool {
	// 今回のユースケースでは、重複があるのは最新のデータの可能性が高いことから逆から探索していく
	for i := len(source) - 1; i >= 0; i-- {
		if len(source[i]) > 0 {
			isSame := true
			for _, col := range con.UniqueCols {
				if source[i][col-1] != target[col-1] {
					isSame = false
				}
			}

			if isSame {
				return false
			}
		}
	}

	return true
}

func (con *Convertor) overWrite() {
	if len(con.OverwriteCols) <= 0 {
		return
	}

	for _, rows := range con.Output.FileData {
		for _, row := range rows {
			for _, pair := range con.OverwriteCols {
				row[pair.Col-1] = pair.Val
			}
		}
	}
}
