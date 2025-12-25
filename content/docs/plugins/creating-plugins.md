---
title: Creating Plugins
description: Build your own custom plugins
order: 2
---

# Creating Plugins

Learn how to build custom plugins for your specific needs.

## Plugin Structure

Create a new Go package for your plugin:

```
plugins/
└── myplugin/
    ├── plugin.go
    └── plugin_test.go
```

## Basic Implementation

Here's a minimal plugin implementation:

```go
package myplugin

import "site/internal/models"

type MyPlugin struct {
    config Config
}

type Config struct {
    Option1 string `yaml:"option1"`
    Option2 bool   `yaml:"option2"`
}

func New() *MyPlugin {
    return &MyPlugin{}
}

func (p *MyPlugin) Name() string {
    return "myplugin"
}

func (p *MyPlugin) Init(cfg Config) error {
    p.config = cfg
    return nil
}

func (p *MyPlugin) BeforeBuild(site *models.Site) error {
    // Runs before content is processed
    return nil
}

func (p *MyPlugin) AfterBuild(site *models.Site) error {
    // Runs after all pages are generated
    return nil
}
```

## Content Transformers

Transform content during the build:

```go
func (p *MyPlugin) Transform(content string) string {
    // Modify content before rendering
    return strings.ReplaceAll(content, "{{current_year}}",
        time.Now().Format("2006"))
}
```

## Testing Your Plugin

Write tests for your plugin:

```go
func TestMyPlugin_Transform(t *testing.T) {
    p := New()
    result := p.Transform("Copyright {{current_year}}")

    if !strings.Contains(result, "2024") {
        t.Error("Expected year substitution")
    }
}
```

> [!TIP]
> Start with a simple feature and iterate. Complex plugins are easier to build incrementally.
