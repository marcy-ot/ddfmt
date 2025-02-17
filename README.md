# ddfmt

ddfmt is a tool to convert excel files to csv files.

## Usage

```
ddfmt -f <input_file> -c <config_file>
```

## Config

CSV files will be generated based on the settings in config.yaml.
Please place config.yaml in the execution directory.

### config.yaml

```yaml
# Specify the sheet name
sheet_name: sheet1

# Override specific column values
overwrite_columns:
  - column: 4    # Column number
    value: "2000" # Value to override with

# Specify column numbers to check for uniqueness
unique_columns: 
  - 1
  - 2
  - 4

# File splitting configuration
file_split:
  row: 2    # Split file every N rows

# Column number to check for distinct values
distinct_column: 2

# Completion message
# {$distinct_column} will be replaced with the distinct column number
completion_message: "Output string {$distinct_column}"
```

Configuration options:
- `sheet_name`: Target Excel sheet name
- `overwrite_columns`: Override values in specified columns
- `unique_columns`: List of column numbers to check for unique constraints
- `file_split`: Output file splitting settings
- `distinct_column`: Column number to check for duplicate values
- `completion_message`: Completion message (supports variable expansion)

