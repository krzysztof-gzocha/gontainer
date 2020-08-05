package definition

type Service struct {
	Getter      string   `yaml:"getter"`
	Type        string   `yaml:"type"`
	Constructor string   `yaml:"constructor"`
	WithError   bool     `yaml:"withError"`
	Disposable  bool     `yaml:"disposable"`
	Args        []string `yaml:"args"`
	Tags        []string `yaml:"tags"`
}

type Services map[string]Service
