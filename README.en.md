# Files Converter

## Overview

**files-converter** - This is a tool that converts source configuration files, or a *.cf file, to a *.cfg file, and vice versa.

## Functionality

- Conversion from *.cf to *.cfe is implemented.
- Conversion from source files (configuration) to *.cfe is implemented.
- Conversion from *.cfe to *.cf is not implemented.
- Conversion from source files (extension) to *.cf is not implemented.

## Configuration file usage

### Required

- **`platform_version`**: *(string)* - specifies the platform version.
  - **Example**: `"platform_version": "8.3.23"`

- **`extension`**: *(string)* - extension name.
  - **Example**: `"extension": "PSSL"`

- **`prefix`**: *(string)* - prefix used for metadata object names.
  - **Example**: `"prefix": "pssl_"`

- **`input_path`**: *(string)* - path to the directory or file to be converted.
  - **Example**: `"input_path": "C:/path/to/input"`

- **`output_path`**: *(string)* - path where the converted files will be saved.
  - **Example**: `"output_path": "C:/path/to/output"`

- **`conversion_type`**: *(string)* - specifies the type of conversion to perform. Supported values:
  - `"srcConvert"` - convert source files to *.cfe.
  - `"cfConvert"` - convert *.cf configuration file to *.cfe.
  - **Example**: `"conversion_type": "srcConvert"`

### Optional

- **`xml_files`**: *(array)* - a list of XML files and their associated operations. For example, to change the flag in a cpmmon module or add a configuration prefix.
  - **`file_name`**: *(string)* - the name of the XML file to operate on.
    - **Example**: `"file_name": "example.xml"`
  - **`element_operations`**: *(array)* - a list of operations to perform on elements within the XML file.
    - **`element_name`**: *(string)* - the name of the XML element to modify.
      - **Example**: `"element_name": "Global"`
    - **`operation`**: *(string)* - the type of operation. Supported values: 
      - `"add"` - add element.
      - `"delete"` - delete element.
      - `"modify"` - change element.
    - **`value`**: *(string, optional)* - the new value to set for the element (used with `add` and `modify` operations).
      - **Example**: `"value": "false"`

### Example Configuration File

Here's an example of a configuration file (`config.json`):

```json
{
  "platform_version": "8.3.23",
  "extension": "ProjectSubsystemsLibrary",
  "prefix": "pssl_",
  "input_path": "C:/path/to/source",
  "output_path": "C:/path/to/output",
  "conversion_type": "srcConvert",
  "xml_files": [
    {
      "file_name": "example.xml",
      "element_operations": [
        {
          "element_name": "Global",
          "operation": "modify",
          "value": "false"
        }
      ]
    }
  ]
}
```

## Usage

It can be installed by running:

``` shell
go install github.com/firstBitSportivnaya/files-converter@latest
```

**Note:**  The use of this program requires the appropriate platform (8.3.23).

## Help Command

For more information on available commands and options:

``` shell
files-converter --help
```

## Examples

Using the configuration file:

``` shell
files-converter --config="configs/config.json"
```

If no --config flag is provided, the program will use the default configuration file located at `$HOME/.files-converter/configs/config.json`:

``` shell
files-converter
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
