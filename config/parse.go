package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Parser interface {
	Parse(configFile string) (ServiceConfig, error)
}

type ParserFunc func(string) (ServiceConfig, error)

func (f ParserFunc) Parse(configFile string) (ServiceConfig, error) { return f(configFile) }

func NewParser() Parser {
	return NewParserWithFileReader(os.ReadFile)
}

func NewParserWithFileReader(f FileReaderFunc) Parser {
	return parser{fileReader: f}
}

type parser struct {
	fileReader FileReaderFunc
}

func (p parser) Parse(configFile string) (ServiceConfig, error) {
	var result ServiceConfig
	var cfg parseableServiceConfig
	data, err := p.fileReader(configFile)
	if err != nil {
		return result, CheckErr(err, configFile)
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return result, CheckErr(err, configFile)
	}

	result = cfg.normalize()
	return result, nil
}

type FileReaderFunc func(string) ([]byte, error)

type parseableServiceConfig struct {
	Name string `json:"name"`
}

// CheckErr returns a proper documented error
func CheckErr(err error, configFile string) error {
	switch e := err.(type) {
	case *json.SyntaxError:
		return NewParseError(err, configFile, int(e.Offset))
	case *json.UnmarshalTypeError:
		return NewParseError(err, configFile, int(e.Offset))
	case *os.PathError:
		return fmt.Errorf(
			"'%s' (%s): %s",
			configFile,
			e.Op,
			e.Err.Error(),
		)
	default:
		return fmt.Errorf("'%s': %v", configFile, err)
	}
}

func (p *parseableServiceConfig) normalize() ServiceConfig {
	return ServiceConfig(*p)
}

func NewParseError(err error, configFile string, offset int) *ParseError {
	b, _ := os.ReadFile(configFile)
	row, col := getErrorRowCol(b, offset)
	return &ParseError{
		ConfigFile: configFile,
		Err:        err,
		Offset:     offset,
		Row:        row,
		Col:        col,
	}
}

type ParseError struct {
	ConfigFile string
	Offset     int
	Row        int
	Col        int
	Err        error
}

func (p *ParseError) Error() string {
	return fmt.Sprintf(
		"'%s': %v, offset: %v, row: %v, col: %v",
		p.ConfigFile,
		p.Err.Error(),
		p.Offset,
		p.Row,
		p.Col,
	)
}
func getErrorRowCol(source []byte, offset int) (row, col int) {
	if len(source) < offset {
		offset = len(source) - 1
	}
	for i := 0; i < offset; i++ {
		v := source[i]
		if v == '\r' {
			continue
		}
		if v == '\n' {
			col = 0
			row++
			continue
		}
		col++
	}
	return
}
