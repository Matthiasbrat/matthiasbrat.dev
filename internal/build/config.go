package build

import (
	"errors"
	"os"

	"site/internal/build/search"
)

const (
	DefaultContentDir  = "content"
	DefaultOutputDir   = "dist"
	DefaultStaticDir   = "static"
	DefaultTemplateDir = "templates"
)

type Config struct {
	ContentDir         string
	OutputDir          string
	StaticDir          string
	TemplateDir        string
	BaseURL            string
	SiteName           string
	SiteDesc           string
	DevMode            bool
	Profile            ProfileConfig
	Referrals          []ReferralConfig
	DefaultSocialImage string
	DB                 search.DB
}

func (c *Config) Validate() error {
	if c.ContentDir == "" {
		return errors.New("content directory is required")
	}
	if c.OutputDir == "" {
		return errors.New("output directory is required")
	}
	if _, err := os.Stat(c.ContentDir); os.IsNotExist(err) {
		return errors.New("content directory does not exist")
	}
	return nil
}

type ProfileConfig struct {
	Photo    string `yaml:"photo"`
	Bio      string `yaml:"bio"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	Email    string `yaml:"email"`
}

type ReferralConfig struct {
	Name     string `yaml:"name"`
	Photo    string `yaml:"photo"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	Website  string `yaml:"website"`
	Twitter  string `yaml:"twitter"`
	Email    string `yaml:"email"`
}
