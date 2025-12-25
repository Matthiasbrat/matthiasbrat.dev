---
title: "Creating Content"
description: "Learn how to write posts and documentation"
order: 4
---

This guide covers how to create content for your site.

## Content Structure

Content is organized into topics, each containing multiple posts:

```
content/
├── topic-slug/
│   ├── _metadata.yml
│   ├── post-one.md
│   └── post-two.md
└── another-topic/
    ├── _metadata.yml
    └── some-post.md
```

## Topic Metadata

Each topic needs a `_metadata.yml` file:

```yaml
name: "Topic Name"
description: "A description of this topic"
type: docs  # or "blog"
```

### Topic Types

- **docs**: Ordered documentation (like this guide)
- **blog**: Chronological posts (newest first)

## Post Frontmatter

Each markdown file starts with YAML frontmatter:

```yaml
---
title: "Post Title"
description: "Post description for SEO"
date: 2024-01-15    # Required for blog posts
order: 1            # Required for docs posts
draft: false        # Optional, defaults to false
---
```

## Markdown Features

### Basic Formatting

Standard markdown is supported:

- **Bold** and *italic* text
- [Links](https://example.com)
- Lists (ordered and unordered)
- Code blocks with syntax highlighting

### Custom Directives

#### Callouts

```md
> [!NOTE]
> Informational callout

> [!WARNING]
> Warning callout

> [!TIP]
> Helpful tip

> [!CAUTION]
> Critical warning
```

#### PDF Embeds

```md
:::pdf{src="/files/document.pdf"}
:::
```

## Next Steps

Now that you know how to create content, learn how to deploy your site.
