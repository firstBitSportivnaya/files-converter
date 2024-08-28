# Files Converter

## Overview

**files-converter** - is a tool specifically designed to convert source configuration files or *.cf files from one format to another. This project uses Cobra for advanced parameter handling, simplified command-line interface, and argument and parameter management.

## Features
**Configuration File Conversion:** Convert configuration files or *.cf files between different formats.  
**Customizable Options:** Specify output formats, file destinations, and more.  

## Functionality
- Conversion of the *.cf form to *.cfe is implemented.
- Conversion of source files to *.cfe is implemented.

## Configuration File
The **files-converter** tool can be configured using a JSON configuration file, allowing you to predefine settings such as input paths, output paths, and conversion types.

### Configuration Parameters

- **`platform_version`**: *(string)* Specifies the required platform version.
  - **Example**: `"platform_version": "8.3.23"`

- **`extension`**: *(string)* Specifies the name of the extension.
  - **Example**: `"extension": "PSSL"`

- **`input_path`**: *(string)* Path to the directory or file to be converted.
  - **Example**: `"input_path": "C:/path/to/input"`

- **`output_path`**: *(string)* Path where the converted files will be saved.
  - **Example**: `"output_path": "C:/path/to/output"`

- **`conversion_type`**: *(string)* Specifies the type of conversion to perform. Valid values are `"srcConvert"` and `"cfConvert"`.
  - **Example**: `"conversion_type": "srcConvert"`

- **`xml_files`**: *(array)* A list of XML files and their associated operations.
  - **`file_name`**: *(string)* The name of the XML file to operate on.
    - **Example**: `"file_name": "example.xml"`
  - **`element_operations`**: *(array)* A list of operations to perform on elements within the XML file.
    - **`element_name`**: *(string)* The name of the XML element to modify.
      - **Example**: `"element_name": "SampleElement"`
    - **`operation`**: *(string)* The type of operation (`"add"`, `"delete"`, `"modify"`).
      - **Example**: `"operation": "modify"`
    - **`value`**: *(string, optional)* The new value to set for the element (used with `add` and `modify` operations).
      - **Example**: `"value": "NewValue"`

### Example Configuration File

Here's an example of a configuration file (`config.json`):

```json
{
  "platform_version": "8.3.23",
  "extension": "ПроектнаяБиблиотекаПодсистем",
  "input_path": "C:/path/to/source",
  "output_path": "C:/path/to/output",
  "conversion_type": "srcConvert",
  "xml_files": [
    {
      "file_name": "example.xml",
      "element_operations": [
        {
          "element_name": "SampleElement",
          "operation": "modify",
          "value": "NewValue"
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
Convert a source files to *.cfe file:
``` shell
files-converter srcConvert --input="C:\path\to\source" --output="C:\path\to\output"
# or short form
files-converter srcConvert -i="C:\path\to\source" -o="C:\path\to\output"
```
Convert *.cf file to *.cfe file:
``` shell
files-converter cfConvert --input="C:\path\to\source\PSSL.cf" --output="C:\path\to\output"
# or short form
files-converter cfConvert -i="C:\path\to\source\PSSL.cf" -o="C:\path\to\output"
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
