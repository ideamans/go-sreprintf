# go-sreprintf

[日本語版 README](README_ja.md)

A Go library for reverse engineering sprintf-formatted strings to extract parameters and reapply them to different templates (e.g., for translation).

## Installation

```bash
go get github.com/miyanaga/go-sreprintf
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/miyanaga/go-sreprintf"
)

func main() {
    template := "Hello %s, you have %d new messages"
    message := "Hello Alice, you have 5 new messages."
    translation := "こんにちは %sさん、あなたには%d件の新しいメッセージがあります"
    
    result, err := sreprintf.Sreprintf(template, message, translation)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result)
    // Output: こんにちは Aliceさん、あなたには5件の新しいメッセージがあります
}
```

## Supported Format Specifiers

The library supports the following format specifiers:

- `%s` - String
- `%d` - Integer
- `%f` - Floating-point number
- `%v` - Generic format (Go-specific)
- `%t` - Boolean (true/false)
- `%%` - Literal % sign

### Format Options

The library correctly handles format options such as:

- Width specification: `%10s`, `%5d`
- Left alignment: `%-10s`
- Sign display: `%+d`
- Alternative forms: `%#x`, `%#o`
- Precision specification: `%.2f`, `%.5s`
- Zero padding: `%05d`
- Complex specifications: `%+10.2f`, `%-20.10s`

When parsing, format options are ignored for value extraction. When reapplying values, the format options in the translation template are used.

## How It Works

1. **Template Parsing**: The library parses the template string to identify all format specifiers
2. **Value Extraction**: Using regex patterns, it extracts values from the formatted message
3. **Type Conversion**: Extracted string values are converted to appropriate types (int, float, bool)
4. **Reapplication**: The extracted values are applied to the translation template using `fmt.Sprintf`

## Error Handling

The library follows a principle of minimal error generation:

- If the message doesn't match the template structure, an error is returned
- Type mismatches are handled gracefully - values are kept as strings if conversion fails
- Missing or extra placeholders are handled by padding with empty strings or ignoring extra values

## Testing

Run tests with:

```bash
go test -v ./...
```

## License

MIT License