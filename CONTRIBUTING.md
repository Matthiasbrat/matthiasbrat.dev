# Contributing to matthiasbrat.com

Thank you for your interest in contributing! This guide will help you get started with development and contributions.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Project Architecture](#project-architecture)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- A text editor or IDE (VS Code, GoLand, etc.)
- Basic understanding of Go and web development

### Setting Up Your Development Environment

1. Fork the repository on GitHub

2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/matthiasbrat.com.git
cd matthiasbrat.com
```

3. Add the upstream repository:
```bash
git remote add upstream https://github.com/Matthiasbrat/matthiasbrat.com.git
```

4. Install dependencies:
```bash
go mod download
```

5. Build the project:
```bash
go build -o site ./cmd/site
```

6. Run the development server:
```bash
./site dev
```

7. Open http://localhost:3000 in your browser

## Development Workflow

### Creating a New Feature

1. Update your local main branch:
```bash
git checkout main
git pull upstream main
```

2. Create a feature branch:
```bash
git checkout -b feature/your-feature-name
```

3. Make your changes and commit regularly:
```bash
git add .
git commit -m "feat: add your feature"
```

4. Push to your fork:
```bash
git push origin feature/your-feature-name
```

5. Create a Pull Request on GitHub

### Commit Message Convention

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

**Examples:**
```
feat: add support for custom markdown extensions
fix: resolve syntax highlighting for YAML files
docs: update installation instructions
refactor: simplify content loading logic
```

## Code Style

### Go Code Style

Follow the standard Go conventions:

- Use `gofmt` to format your code
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused

**Running gofmt:**
```bash
# Format all Go files
gofmt -w .

# Check formatting without modifying files
gofmt -l .
```

### File Organization

- Keep files focused on a single responsibility
- Group related functionality in packages
- Use internal packages for implementation details
- Keep exported APIs minimal and well-documented

### Error Handling

Always handle errors explicitly:

```go
// Good
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad
doSomething() // ignoring error
```

Wrap errors with context using `fmt.Errorf` and `%w` for error chains.

## Project Architecture

### Directory Structure

```
cmd/site/              # Application entry point
internal/
  ├── build/           # Build system
  │   ├── assets/      # Asset processing (minification, hashing)
  │   ├── content/     # Content loading and parsing
  │   ├── markdown/    # Markdown rendering and extensions
  │   ├── og/          # Open Graph image generation
  │   └── search/      # Search indexing
  ├── db/              # Database operations (SQLite)
  ├── models/          # Data models (Post, Collection, User)
  └── server/          # HTTP server and API endpoints
```

### Key Packages

#### `internal/build`

Responsible for static site generation:
- Loading content from markdown files
- Processing templates
- Generating HTML pages
- Creating sitemaps
- Processing static assets

#### `internal/server`

HTTP server for both development and production:
- Serving static files
- API endpoints for reactions and comments
- OAuth authentication
- Hot-reload in development mode
- File watching

#### `internal/models`

Data models for:
- Posts (blog posts and documentation pages)
- Collections (topics and series)
- Users (for authentication)
- Reactions and comments

#### `internal/db`

SQLite database operations:
- Full-text search indexing
- Storing reactions and comments
- User management

### Adding New Features

#### Adding a Markdown Extension

1. Create a new extension in `internal/build/markdown/extensions/`:
```go
package extensions

import "github.com/yuin/goldmark/ast"

type YourExtension struct{}

func (e *YourExtension) Extend(m goldmark.Markdown) {
    // Register parsers, renderers, etc.
}
```

2. Register in `internal/build/markdown/renderer.go`:
```go
import "site/internal/build/markdown/extensions"

func NewRenderer() *Renderer {
    md := goldmark.New(
        // ...
        goldmark.WithExtensions(
            &extensions.YourExtension{},
        ),
    )
}
```

#### Adding a New Template

1. Create template in `templates/`:
```html
{{ define "your-template.html" }}
<!-- Your template content -->
{{ end }}
```

2. Load template in `internal/build/build.go`:
```go
ts.YourTemplate, err = parseWithBase("your-template.html")
```

3. Generate pages using the template:
```go
func (s *Site) generateYourPage(tmpl *template.Template) error {
    // Prepare data
    data := struct {
        PageData
        // Your custom fields
    }{
        // ...
    }

    return s.renderPage(tmpl, "output/path.html", data)
}
```

#### Adding a New API Endpoint

1. Add handler in `internal/server/`:
```go
func (s *Server) handleYourEndpoint(w http.ResponseWriter, r *http.Request) {
    // Your logic here
}
```

2. Register route in `internal/server/server.go`:
```go
func (s *Server) setupRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    // ...
    mux.HandleFunc("/api/your-endpoint", s.handleYourEndpoint)
    return mux
}
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/build/...
```

### Writing Tests

Create test files alongside your code with `_test.go` suffix:

```go
package yourpackage

import "testing"

func TestYourFunction(t *testing.T) {
    result := YourFunction("input")
    expected := "expected output"

    if result != expected {
        t.Errorf("YourFunction() = %v, want %v", result, expected)
    }
}
```

### Test Coverage

Aim for at least 70% test coverage for critical packages:
- `internal/build/content`
- `internal/build/markdown`
- `internal/db`
- `internal/models`

## Pull Request Process

### Before Submitting

1. Ensure all tests pass:
```bash
go test ./...
```

2. Format your code:
```bash
gofmt -w .
```

3. Update documentation if needed

4. Add tests for new features

5. Update CHANGELOG.md (if applicable)

### PR Guidelines

- Keep PRs focused on a single feature or fix
- Write clear PR descriptions explaining the change
- Reference related issues (e.g., "Fixes #123")
- Ensure CI checks pass
- Be responsive to code review feedback

### PR Template

When creating a PR, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Code follows project style guidelines
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] All tests passing
```

## Code Review Process

1. At least one maintainer review is required
2. All CI checks must pass
3. Address review feedback promptly
4. Once approved, a maintainer will merge your PR

## Areas for Contribution

### Good First Issues

- Documentation improvements
- Adding new markdown extensions
- Improving error messages
- Adding tests
- UI/UX enhancements

### Advanced Features

- Performance optimizations
- New template functions
- Database improvements
- Advanced search features
- Caching strategies

### Documentation

- Tutorial content
- API documentation
- Architecture diagrams
- Example configurations

## Questions?

If you have questions:

1. Check existing issues and discussions
2. Read the documentation
3. Open a new issue with the "question" label

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

Thank you for contributing! 🎉
