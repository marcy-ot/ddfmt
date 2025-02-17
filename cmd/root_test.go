package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 引数のテスト
// file 引数は必須
// config 引数は任意
// config 指定がなければそのままCSV変換される
// config があるけど設定がおかしい場合のテスト
// default test
// config フルセット

// - sheet_name: sheet1 									　シート名が存在しない、Value未設定
// - # export_file_name: export 							　そのファイル名で出力されるか、Value未設定
// - overwrite_columns:										　Value未設定
//   - column: 4											　	範囲外の値(0, -1, greater)
//     value: "2000"
// unique_columns: 											　Value未設定
//  - 1 													　	範囲外の値(0, -1, greater)
// file_split: 												　Value未設定
//   row: 2													　	Value未設定、0, -1
// distinct_column: 2										　Value未設定 範囲外の値(0, -1, greater)
// completion_message: "文字列を出力{$distinct_column}します。" Value未設定

var inputPath = func(name string) string {
	return fmt.Sprintf("testdata/%s/", name)
}

var expectFile = func(name string) string {
	return fmt.Sprintf("testdata/%s/expect/", name)
}

func Test_ddfmt(t *testing.T) {
	tests := []struct {
		name            string
		testPath        string
		inputFile       string
		configFile      string
		outputsFileName []string
		wantStdout      string
	}{
		{
			name:            "full_config_test",
			inputFile:       "testdata.xlsx",
			configFile:      "ddfmt.yaml",
			outputsFileName: []string{"testdata.csv", "testdata_1.csv"},
			wantStdout:      "文字列を出力\nLaptop\nKeyboard\nmause\nします。\n",
		},
		{
			name:            "no_config_test",
			inputFile:       "testdata.xlsx",
			configFile:      "",
			outputsFileName: []string{"testdata.csv"},
			wantStdout:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// flag のセット
			cmdArg := []string{"--file", fmt.Sprint(inputPath(tt.name), tt.inputFile)}
			if tt.configFile != "" {
				cmdArg = append(cmdArg, "--config", fmt.Sprint(inputPath(tt.name), tt.configFile))
			}

			defer func() {
				for _, fileName := range tt.outputsFileName {
					file := fmt.Sprint(inputPath(tt.name), fileName)
					if err := os.Remove(file); err != nil {
						log.Fatal(err)
					}
				}
			}()

			stdOutStr := captureStdout(func() {
				Do(cmdArg, os.Stdin, os.Stdout, os.Stderr)
			})

			// 出力チェック
			assert.Equal(t, tt.wantStdout, stdOutStr)

			// ファイルが存在するかチェック
			compareContent(t, tt.outputsFileName, inputPath(tt.name), expectFile(tt.name))
		})
	}
}

func captureStdout(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w

	f()

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func compareContent(t *testing.T, outputsFileName []string, actualPath string, expectPath string) {
	for _, fileName := range outputsFileName {
		file := fmt.Sprint(actualPath, fileName)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// ファイルが存在するかチェック
			assert.NoError(t, err)
		}

		ef, err := os.Open(file)
		assert.NoError(t, err)
		actual, err := io.ReadAll(ef)
		assert.NoError(t, err)

		expectFile := fmt.Sprint(expectPath, fileName)
		af, err := os.Open(expectFile)
		assert.NoError(t, err)
		expect, err := io.ReadAll(af)
		assert.NoError(t, err)

		defer func() {
			ef.Close()
			af.Close()
		}()

		assert.Equal(t, expect, actual)
	}
}
