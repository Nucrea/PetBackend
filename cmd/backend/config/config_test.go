package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testParserConf = `
token: "test_token"
url: test.url
port: 5010
`

type TestConf struct {
	Token string `yaml:"token"`
	Url   string `yaml:"url"`
	Port  int    `yaml:"port"`
}

func TestParser(t *testing.T) {
	file, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(testParserConf)
	if err != nil {
		t.Fatal(err)
	}

	p := NewParser[TestConf]()
	conf, err := p.ParseFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test_token", conf.Token)
	assert.Equal(t, "test.url", conf.Url)
	assert.Equal(t, 5010, conf.Port)
}
