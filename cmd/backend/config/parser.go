package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Parser[T interface{}] interface {
	ParseFile(path string) (T, error)
}

func NewParser[T interface{}]() Parser[T] {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &parser[T]{
		validate: validate,
	}
}

type parser[T interface{}] struct {
	validate *validator.Validate
}

func (p *parser[T]) ParseFile(path string) (T, error) {
	fBytes, err := os.ReadFile(path)
	if err != nil {
		var t T
		return t, err
	}

	return p.parse(fBytes)
}

func (p *parser[T]) parse(b []byte) (T, error) {
	var t T

	if err := yaml.Unmarshal(b, &t); err != nil {
		return t, err
	}
	if err := p.validate.Struct(t); err != nil {
		return t, err
	}

	return t, nil
}
