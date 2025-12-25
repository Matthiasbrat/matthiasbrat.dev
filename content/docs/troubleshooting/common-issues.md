---
title: Common Issues
description: Solutions to frequently encountered problems
order: 1
---

# Common Issues

Quick solutions to problems you might encounter.

## Build Errors

### "Template not found"

**Problem:** The build fails with a template not found error.

**Solution:** Ensure all required templates exist in the `templates/` directory:

```bash
ls templates/
# Should show: base.html home.html post.html topic.html
```

### "Content directory is empty"

**Problem:** No pages are generated despite having content files.

**Solution:** Check that your content structure is correct:

```
content/
└── topic-name/
    ├── _metadata.yml  # Required!
    └── post.md
```

> [!WARNING]
> The `_metadata.yml` file is required for each topic directory.

## Runtime Errors

### "Port already in use"

**Problem:** The dev server won't start because the port is occupied.

**Solution:** Use a different port or kill the existing process:

```bash
# Use different port
./site dev --port 3001

# Or find and kill existing process
lsof -i :3000
kill -9 <PID>
```

### "Hot reload not working"

**Problem:** Changes aren't reflected without manual refresh.

**Solution:** Ensure WebSocket connection is established:

1. Check browser console for WS errors
2. Verify dev mode is enabled
3. Try hard refresh (Ctrl+Shift+R)

## Content Issues

### "Syntax highlighting not working"

**Problem:** Code blocks appear without colors.

**Solution:** Ensure the language is specified:

```markdown
\`\`\`javascript
// This will be highlighted
const x = 1;
\`\`\`
```

### "Images not loading"

**Problem:** Images show as broken links.

**Solution:** Place images in `static/` and reference with absolute paths:

```markdown
![Alt text](/images/photo.jpg)
```
