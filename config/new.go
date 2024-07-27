package config

func NewFromFile(path string) (IConfig, error) {
	p := NewParser[Config]()
	config, err := p.ParseFile(path)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
