package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

// Hashes maps original asset paths to hashed paths
type Hashes map[string]string

// Processor handles static asset processing
type Processor struct {
	staticDir string
	outputDir string
	devMode   bool
	minifier  *minify.M
}

// NewProcessor creates a new asset processor
func NewProcessor(staticDir, outputDir string, devMode bool) *Processor {
	// Setup minifier
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("application/javascript", js.Minify)
	m.AddFunc("text/html", html.Minify)

	return &Processor{
		staticDir: staticDir,
		outputDir: outputDir,
		devMode:   devMode,
		minifier:  m,
	}
}

// ProcessAll processes all static files and returns asset hash mappings
func (p *Processor) ProcessAll() (Hashes, error) {
	hashes := make(Hashes)

	err := filepath.Walk(p.staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(p.staticDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			destPath := filepath.Join(p.outputDir, relPath)
			return os.MkdirAll(destPath, 0755)
		}

		// Process file
		return p.processFile(path, relPath, hashes)
	})

	return hashes, err
}

// processFile processes a single file
func (p *Processor) processFile(srcPath, relPath string, hashes Hashes) error {
	// Read source file
	src, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Minify based on file extension
	ext := strings.ToLower(filepath.Ext(srcPath))
	var output []byte
	shouldHash := false

	switch ext {
	case ".css":
		output, err = p.minifier.Bytes("text/css", src)
		if err != nil {
			output = src
		}
		shouldHash = !p.devMode
	case ".js":
		output, err = p.minifier.Bytes("text/javascript", src)
		if err != nil {
			output = src
		}
		shouldHash = !p.devMode
	default:
		output = src
	}

	// Generate hash and rename if needed
	destPath := filepath.Join(p.outputDir, relPath)

	if shouldHash {
		// Generate SHA256 hash of content
		hash := sha256.Sum256(output)
		hashStr := hex.EncodeToString(hash[:])[:8] // Use first 8 chars

		// Insert hash before extension: style.css -> style.a3f5b8c2.css
		filename := filepath.Base(relPath)
		nameWithoutExt := strings.TrimSuffix(filename, ext)
		hashedFilename := fmt.Sprintf("%s.%s%s", nameWithoutExt, hashStr, ext)

		// Update destination path
		destPath = filepath.Join(filepath.Dir(destPath), hashedFilename)

		// Store mapping for templates (use forward slashes for web paths)
		webPath := filepath.ToSlash(relPath)
		hashedWebPath := filepath.ToSlash(filepath.Join(filepath.Dir(relPath), hashedFilename))
		hashes[webPath] = hashedWebPath
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Write to destination
	return os.WriteFile(destPath, output, 0644)
}
