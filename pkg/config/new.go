package config

func NewFromFile[T interface{}](filePath string) (T, error) {
	p := NewParser[T]()

	config, err := p.ParseFile(filePath)
	if err != nil {
		var t T
		return t, err
	}

	return config, nil
}
