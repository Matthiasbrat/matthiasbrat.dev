---
title: Debugging Guide
description: How to debug issues effectively
order: 2
---

# Debugging Guide

Tools and techniques for troubleshooting your site.

## Verbose Mode

Run the build with verbose output:

```bash
./site build --verbose
```

This shows:
- Files being processed
- Template rendering times
- Asset copying details

## Inspecting Generated Output

Check the generated HTML:

```bash
# View page structure
cat dist/docs/getting-started/index.html | head -50

# Search for specific content
grep -r "search-term" dist/
```

## Browser DevTools

### Network Tab

Check for failed requests:

1. Open DevTools (F12)
2. Go to Network tab
3. Reload page
4. Look for red (failed) requests

### Console Tab

Look for JavaScript errors:

```javascript
// Enable verbose logging
localStorage.setItem('debug', 'true');
```

## Common Debug Patterns

### Template Debugging

Add debug output to templates:

```html
{{/* Debug: show available data */}}
<pre>{{printf "%#v" .}}</pre>
```

### Content Debugging

Print frontmatter during build:

```go
func loadPost(path string) (*Post, error) {
    log.Printf("Loading: %s", path)
    // ...
    log.Printf("Frontmatter: %+v", fm)
}
```

> [!TIP]
> Remove debug code before deploying to production.

## Getting Help

If you're still stuck:

1. Check the issue tracker for similar problems
2. Create a minimal reproduction
3. Include relevant error messages and logs
