package server

type Config struct {
	Port        int
	ContentDir  string
	OutputDir   string
	StaticDir   string
	TemplateDir string
	DevMode     bool
	BaseURL     string
	Profile     ProfileConfig
}

type ProfileConfig struct {
	Photo    string `yaml:"photo"`
	Bio      string `yaml:"bio"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	Email    string `yaml:"email"`
}
