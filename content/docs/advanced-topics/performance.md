---
title: Performance Optimization
description: Tips for optimizing your site's performance
order: 2
---

# Performance Optimization

Make your site fast and efficient with these optimization techniques.

## Build-Time Optimizations

### Code Splitting

Split your content into smaller chunks:

```go
func Build(cfg Config) error {
    // Process content in parallel
    var wg sync.WaitGroup
    for _, topic := range topics {
        wg.Add(1)
        go func(t *Topic) {
            defer wg.Done()
            processTopicContent(t)
        }(topic)
    }
    wg.Wait()
    return nil
}
```

### Asset Minification

Minify CSS and JavaScript for production:

```bash
# Minify CSS
npx csso static/css/style.css -o dist/css/style.min.css

# Minify JavaScript
npx terser static/js/main.js -o dist/js/main.min.js
```

## Runtime Optimizations

### Lazy Loading Images

Use native lazy loading for images below the fold:

```html
<img src="image.jpg" loading="lazy" alt="Description">
```

### Caching Headers

Configure proper cache headers in your server:

```go
func serveStatic(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "public, max-age=31536000")
    http.ServeFile(w, r, path)
}
```

> [!NOTE]
> Static assets should be cached aggressively since they include content hashes.
