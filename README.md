# Tidy

![](https://img.shields.io/github/v/release/purposed/tidy?style=flat-square) ![](https://img.shields.io/github/go-mod/go-version/purposed/tidy?style=flat-square) ![](https://img.shields.io/github/license/purposed/tidy?style=flat-square)

Tidy is a configurable tool for automating boring and/or repetitive file managment tasks.

## Installation
### With binman
```bash
$ binman install tidy
```

### With go get 
```bash
$ go get -u github.com/purposed/tidy
```

## Usage
Tidy is simple to use
```bash
$ tidy run [-cfg CONFIG_FILE]
```

## Configuration
Tidy is configured via a JSON file (contributions are welcome for generating configuration files from a GUI).

## Sample Configuration
```json
{
  "monitors": [
    {
      "root_path": "/home/username/Downloads",
      "refresh_interval_s": 15,
      "rules": [
        {
          "name": "Delete files older than a week",
          "condition": "age > 7d",
          "action": {"type":  "delete"}
        }
      ]
    },
    {
      "root_path": "/home/username/misc_files",
      "refresh_interval_s": 3600,
      "rules": [
        {
          "name": "Find new movie files",
          "condition": "(extension = mov or extension = avi)",
          "action": {
            "type": "move",
            "parameters": {
              "name_template": "Movie-{name}.{extension}",
              "to_directory": "/home/username/Movies"
            }
          }
        }
      ]
    }
  ]
}
```