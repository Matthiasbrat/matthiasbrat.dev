package main

// site build --dry-run
// site dev -port 3000
// site serve -port 8080
// site help

import (
	"flag"
	"fmt"
	"os"

	"site/internal/build"
	"site/internal/db"
	"site/internal/server"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type siteConfig struct {
	BaseURL    string                `yaml:"base_url"`
	DevBaseURL string                `yaml:"dev_base_url"`
	Profile    server.ProfileConfig `yaml:"profile"`
}

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		cmdBuild(os.Args[2:])
	case "dev":
		cmdDev(os.Args[2:])
	case "serve":
		cmdServe(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func cmdBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	contentDir := fs.String("content", "content", "Content directory")
	outputDir := fs.String("output", "dist", "Output directory")
	baseURL := fs.String("base-url", "", "Base URL for canonical links (defaults to site.yml)")
	fs.Parse(args)

	database, err := db.New("data/sqlite.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	cfg := build.Config{
		ContentDir:  *contentDir,
		OutputDir:   *outputDir,
		BaseURL:     *baseURL,
		StaticDir:   "static",
		TemplateDir: "templates",
		DB:          database,
	}

	if err := build.Build(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Build complete")
}

func cmdDev(args []string) {
	fs := flag.NewFlagSet("dev", flag.ExitOnError)
	port := fs.Int("port", 3000, "Port to serve on")
	contentDir := fs.String("content", "content", "Content directory")
	outputDir := fs.String("output", "dist", "Output directory")
	baseURL := fs.String("base-url", "", "Base URL (defaults to site.yml or localhost)")
	fs.Parse(args)

	siteCfg := loadSiteConfig()
	finalBaseURL := resolveBaseURL(*baseURL, siteCfg.DevBaseURL, *port)

	cfg := server.Config{
		Port:        *port,
		ContentDir:  *contentDir,
		OutputDir:   *outputDir,
		StaticDir:   "static",
		TemplateDir: "templates",
		DevMode:     true,
		BaseURL:     finalBaseURL,
		Profile:     siteCfg.Profile,
	}

	if err := server.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func cmdServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "Port to serve on")
	outputDir := fs.String("output", "dist", "Output directory")
	baseURL := fs.String("base-url", "", "Base URL (defaults to site.yml base_url)")
	fs.Parse(args)

	siteCfg := loadSiteConfig()

	// Resolve base URL: flag > site.yml > localhost fallback
	finalBaseURL := *baseURL
	if finalBaseURL == "" {
		finalBaseURL = siteCfg.BaseURL
	}
	if finalBaseURL == "" {
		finalBaseURL = fmt.Sprintf("http://localhost:%d", *port)
	}

	cfg := server.Config{
		Port:      *port,
		OutputDir: *outputDir,
		StaticDir: "static",
		DevMode:   false,
		BaseURL:   finalBaseURL,
		Profile:   siteCfg.Profile,
	}

	if err := server.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func loadSiteConfig() siteConfig {
	var cfg siteConfig
	if data, err := os.ReadFile("site.yml"); err == nil {
		_ = yaml.Unmarshal(data, &cfg)
	}
	return cfg
}

func resolveBaseURL(flagValue, configValue string, port int) string {
	if flagValue != "" {
		return flagValue
	}
	if configValue != "" {
		return configValue
	}
	return fmt.Sprintf("http://localhost:%d", port)
}

func printUsage() {
	fmt.Println(`Usage: site <command> [options]

Commands:
  build     Build static site to output directory
  dev       Development server with hot reload
  serve     Production server with reactions API
  help      Show this message

Build Options:
  -content   Content directory (default: content)
  -output    Output directory (default: dist)
  -base-url  Base URL for canonical links

Dev Options:
  -port      Port to serve on (default: 8080)
  -content   Content directory (default: content)
  -output    Output directory (default: dist)

Serve Options:
  -port      Port to serve on (default: 8080)
  -output    Output directory (default: dist)
  -base-url  Base URL for production`)
}
