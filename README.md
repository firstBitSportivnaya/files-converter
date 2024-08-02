# Files Converter

## Overview

**files-converter** - is a tool specifically designed to convert source configuration files or *.cf files from one format to another. This project uses Cobra for advanced parameter handling, simplified command-line interface, and argument and parameter management.

## Features
**Configuration File Conversion:** Convert configuration files or *.cf files between different formats.  
**Cobra Parameter Handling:** Simplified command-line interface for managing parameters.  
**Customizable Options:** Specify output formats, file destinations, and more.  

## Functionality

- Conversion of the *.cf form to *.cfe is not implemented.
- Conversion of source files to *.cfe is implemented.

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
Convert a source files to *.cfe file:
``` shell
files-converter srcConvert --input="C:\path\to\source" --output="C:\path\to\output"
# or short form
files-converter srcConvert -i="C:\path\to\source" -o="C:\path\to\output"
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
