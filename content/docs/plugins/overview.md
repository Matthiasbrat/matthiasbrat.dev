---
title: Plugin Overview
description: Introduction to the plugin system
order: 1
---

# Plugin Overview

Extend your site's functionality with the plugin system.

## What Are Plugins?

Plugins are modular extensions that add new features without modifying core code. They can:

- Add new content types
- Modify the build process
- Inject custom scripts
- Transform content

## Plugin Architecture

Plugins follow a simple interface:

```go
type Plugin interface {
    Name() string
    Init(config Config) error
    BeforeBuild(site *Site) error
    AfterBuild(site *Site) error
}
```

## Installing Plugins

Add plugins to your `site.yml` configuration:

```yaml
plugins:
  - name: analytics
    config:
      tracking_id: "UA-XXXXX-Y"
  - name: search
    config:
      index_path: "search-index.json"
```

## Official Plugins

| Plugin | Description |
|--------|-------------|
| `analytics` | Google Analytics integration |
| `search` | Client-side search functionality |
| `rss` | RSS feed generation |
| `comments` | Discussion system |

> [!NOTE]
> Check the plugin registry for community-contributed plugins.
