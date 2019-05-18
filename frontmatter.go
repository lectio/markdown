package markdown

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v2"
)

// A Property expresses a single front matter variable
type Property interface {
	Name() string
	Value() interface{}
	Write(context.Context, io.Writer) error
	Init(context.Context, io.Reader) error
}

// FrontMatter is what precedes the body of markdown content
type FrontMatter interface {
	Properties() []Property
	ForEach()
	Write(context.Context, io.Writer) error
	Init(context.Context, io.Reader) error
}

// ParseYAMLFrontMatter will convert an input byte array like ---<stuff>---\n<body> into v as YAML and <body> as return value
func ParseYAMLFrontMatter(b []byte, v interface{}) ([]byte, bool, error) {
	buf := bytes.NewBuffer(b)

	var insideFrontMatter bool
	var yamlStartIndex int
	var yamlEndIndex int

	for {
		line, err := buf.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, false, err
		}

		if strings.TrimSpace(line) != "---" {
			continue
		}

		if !insideFrontMatter {
			insideFrontMatter = true
			yamlStartIndex = len(b) - buf.Len()
		} else {
			yamlEndIndex = len(b) - buf.Len()
			break
		}
	}

	// if we get to here and we're not inside front matter then the entire string is body
	if !insideFrontMatter {
		return b, false, nil
	}

	if insideFrontMatter && yamlEndIndex == 0 {
		return nil, true, fmt.Errorf("Unexplained front matter parser error; insideFrontMatter: %v, yamlStartIndex: %v, yamlEndIndex: %v", insideFrontMatter, yamlStartIndex, yamlEndIndex)
	}

	err := yaml.Unmarshal(b[yamlStartIndex:yamlEndIndex], v)

	if err != nil {
		return nil, false, err
	}

	return bytes.TrimSpace(b[yamlEndIndex:]), true, nil
}
