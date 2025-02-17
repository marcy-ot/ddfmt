package config

import (
	"io"

	"gopkg.in/yaml.v2"
)

type ColumnValue struct {
	Col int    `yaml:"column"`
	Val string `yaml:"value"`
}

type Config struct {
	SheetName           string        `yaml:"sheet_name"`
	ExportFileExtension string        `yaml:"export_file_extension"`
	OverwriteCols       []ColumnValue `yaml:"overwrite_columns"`
	UniqueCols          []int         `yaml:"unique_columns"`
	FileSplit           struct {
		Row int `yaml:"row"`
	} `yaml:"file_split"`
	DistinctCol       int    `yaml:"distinct_column"`
	CompletionMessage string `yaml:"completion_message"`
}

var defaultSheetName = "sheet1"
var defaultExportFileExtension = "csv"

func ParseConfig(file io.Reader) (*Config, error) {
	var config *Config
	dec := yaml.NewDecoder(file)
	if err := dec.Decode(&config); err != nil {
		return config, err
	}

	setDefault(config)

	return config, nil
}

func DefaultConfig() *Config {
	config := &Config{}
	setDefault(config)
	return config
}

func setDefault(conf *Config) {
	// Mapping default value
	if conf.SheetName == "" {
		conf.SheetName = defaultSheetName
	}
	if conf.ExportFileExtension == "" {
		conf.ExportFileExtension = defaultExportFileExtension
	}
}

func (c *Config) HasSplitRow() bool {
	return c.FileSplit.Row != 0
}

func (c *Config) GteSplitRow(n int) bool {
	return c.HasSplitRow() && n <= c.FileSplit.Row
}
