---
title: Custom Themes
description: Create and customize your own themes
order: 1
---

# Custom Themes

Learn how to create beautiful custom themes for your site.

## Theme Structure

A theme consists of several key components:

```
themes/
└── my-theme/
    ├── templates/
    │   ├── base.html
    │   ├── home.html
    │   └── post.html
    └── static/
        └── css/
            └── theme.css
```

## CSS Variables

The theming system uses CSS custom properties for easy customization:

```css
:root {
    --color-primary: #0066cc;
    --color-text: #1a1a1a;
    --color-background: #ffffff;
    --font-body: 'Inter', sans-serif;
    --font-heading: 'Source Serif 4', serif;
}
```

## Creating a Dark Theme

Override the default variables for dark mode:

```css
@media (prefers-color-scheme: dark) {
    :root {
        --color-text: #e5e5e5;
        --color-background: #1a1a1a;
    }
}
```

> [!TIP]
> Test your theme with both light and dark mode to ensure good contrast and readability.
