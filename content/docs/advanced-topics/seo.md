---
title: SEO Best Practices
description: Optimize your content for search engines
order: 3
---

# SEO Best Practices

Improve your site's visibility in search engine results.

## Meta Tags

Every page should include essential meta tags:

```html
<title>Page Title | Site Name</title>
<meta name="description" content="A concise description of the page content">
<link rel="canonical" href="https://example.com/page">
```

## Structured Data

Add JSON-LD structured data for rich snippets:

```json
{
    "@context": "https://schema.org",
    "@type": "Article",
    "headline": "Article Title",
    "datePublished": "2024-01-15T09:00:00Z",
    "author": {
        "@type": "Person",
        "name": "Author Name"
    }
}
```

## Sitemap Generation

The build system automatically generates a sitemap:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>https://example.com/</loc>
    </url>
    <url>
        <loc>https://example.com/docs/getting-started</loc>
        <lastmod>2024-01-15</lastmod>
    </url>
</urlset>
```

## Performance Impact on SEO

> [!WARNING]
> Page speed is a ranking factor. Ensure your pages load quickly by following the performance optimization guide.

## Accessibility

Good accessibility improves SEO:

- Use semantic HTML elements
- Provide alt text for images
- Ensure sufficient color contrast
- Support keyboard navigation
