/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcy-ot/ddfmt/internal/config"
	"github.com/marcy-ot/ddfmt/internal/convertor"
	"github.com/marcy-ot/ddfmt/internal/exporter"
	"github.com/spf13/cobra"
)

// var defaultConfigFileName = "ddfmt.yaml"
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	rootCmd := &cobra.Command{
		Use:   "ddfmt",
		Short: "Convert Excel to CSV",
		Run: func(cmd *cobra.Command, args []string) {
			// 取り込みファイル名取得
			fn := cmd.Flag("file")
			inputFileName := fn.Value.String()

			// 設定ファイルの読み込み
			confFileName := ""
			cfc := cmd.Flag("config")
			if cfc != nil && cfc.Changed {
				confFileName = cfc.Value.String()
			}
			stderr := cmd.ErrOrStderr()
			configFilePath := getFilePath(stderr, confFileName)
			config, err := readConfig(stderr, configFilePath)
			if err != nil {
				os.Exit(1)
			}

			// 対象ファイルの読み込み
			convertible := convertor.NewConvertable(inputFileName)
			cfilePath := getFilePath(stderr, inputFileName)
			if err := convertible.Read(stderr, cfilePath, config); err != nil {
				os.Exit(1)
			}

			// 変換処理
			convertor := convertor.NewConvertor(convertible)
			if err := convertor.SetConfig(stderr, config); err != nil {
				os.Exit(1)
			}
			output := convertor.Convert()

			// TODO: 引数から出力する形式を変更できるようにする
			// 出力ファイル名を取得
			exporter := exporter.NewExporter(config, output, stderr)
			if err := exporter.Export(exportFileName(inputFileName)); err != nil {
				os.Exit(1)
			}

			if output.Message != "" {
				fmt.Println(output.Message)
			}
		},
	}

	rootCmd.Flags().StringP("file", "f", "", "Specify the path of the file to be processed")
	rootCmd.Flags().StringP("config", "c", "", "Specify the path of the config file")
	rootCmd.MarkFlagRequired("file")

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// var rootCmd = &cobra.Command{
// 	Use:   "ddfmt",
// 	Short: "Convert Excel to CSV",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// 取り込みファイル名取得
// 		fn := cmd.Flag("file")
// 		inputFileName := fn.Value.String()

// 		// 設定ファイルの読み込み
// 		confFileName := ""
// 		cfc := cmd.Flag("config")
// 		if cfc != nil && cfc.Changed {
// 			confFileName = cfc.Value.String()
// 		}
// 		stderr := cmd.ErrOrStderr()
// 		configFilePath := getFilePath(stderr, confFileName)
// 		config, err := readConfig(stderr, configFilePath)
// 		if err != nil {
// 			os.Exit(1)
// 		}

// 		// 対象ファイルの読み込み
// 		convertible := convertor.NewConvertable(inputFileName)
// 		cfilePath := getFilePath(stderr, inputFileName)
// 		if err := convertible.Read(stderr, cfilePath, config); err != nil {
// 			os.Exit(1)
// 		}

// 		// 変換処理
// 		convertor := convertor.NewConvertor(convertible)
// 		if err := convertor.SetConfig(stderr, config); err != nil {
// 			os.Exit(1)
// 		}
// 		output := convertor.Convert()

// 		// TODO: 引数から出力する形式を変更できるようにする
// 		// 出力ファイル名を取得
// 		exporter := exporter.NewExporter(config, output, stderr)
// 		if err := exporter.Export(exportFileName(inputFileName)); err != nil {
// 			os.Exit(1)
// 		}

// 		if output.Message != "" {
// 			fmt.Println(output.Message)
// 		}
// 	},
// }

func exportFileName(inputFileName string) string {
	ext := filepath.Ext(inputFileName)
	return strings.TrimSuffix(inputFileName, ext)
}

func readConfig(stderr io.Writer, configPath string) (*config.Config, error) {
	if configPath == "" {
		return config.DefaultConfig(), nil
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Fprintf(stderr, "error exist configuration: %v\n", err)
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error open configuration: %v\n", err)
		return nil, err
	}
	defer file.Close()

	// parse config yaml file
	conf, err := config.ParseConfig(file)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing configuration: %v\n", err)
		return conf, err
	}

	return conf, nil
}

var workDir = os.Getwd

func getFilePath(stderr io.Writer, fileName string) string {
	if fileName == "" {
		return ""
	}

	wd, err := workDir()
	if err != nil {
		fmt.Fprintln(stderr, "error pasing config")
		os.Exit(1)
	}
	return filepath.Join(wd, fileName)
}

// func Execute() {
// 	err := rootCmd.Execute()
// 	if err != nil {
// 		os.Exit(1)
// 	}
// }

// func init() {
// 	rootCmd.Flags().StringP("file", "f", "", "Specify the path of the file to be processed")
// 	rootCmd.Flags().StringP("config", "c", "", "Specify the path of the config file")
// 	rootCmd.MarkFlagRequired("file")
// }
