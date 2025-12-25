---
title: "Working with Embedded Content"
description: "Learn how to embed YouTube videos and PDF documents in your posts"
date: 2025-01-11
---

This site supports embedding rich media content directly in your posts. Let's explore the available options.

## YouTube Videos

You can embed YouTube videos using standard markdown image syntax. Just use a YouTube URL as the image source.

### Example: Go Programming Introduction

![Introduction to Go Programming](https://www.youtube.com/watch?v=YS4e4q9oBaU)

The syntax is simple:

```markdown
![Video Title](https://www.youtube.com/watch?v=VIDEO_ID)
```

::: aside
**How it works**: The goldmark-embed extension automatically detects YouTube URLs and converts them to responsive iframe embeds.
:::

> [!TIP]
> Use descriptive alt text for accessibility. This text appears if the video can't load and helps screen readers.

### Another Example: Web Development

![Learn Web Development Basics](https://www.youtube.com/watch?v=dQw4w9WgXcQ)

## PDF Documents

For PDF files, use the special `:::pdf` directive:

```markdown
:::pdf{src="/files/document.pdf"}
```

> [!NOTE]
> PDF files should be placed in your `static/files/` directory.

### PDF Browser Support

Not all browsers support embedded PDFs natively. The embed provides:
- An iframe viewer for supported browsers
- A fallback download link for unsupported browsers

| Browser | PDF Support |
|---------|-------------|
| Chrome | ✅ Native |
| Firefox | ✅ Native |
| Safari | ✅ Native |
| Edge | ✅ Native |
| Mobile Safari | ⚠️ Limited |
| Mobile Chrome | ✅ Yes |

## Best Practices

### YouTube Embeds

> [!IMPORTANT]
> **Performance Tip**: Each YouTube embed loads additional JavaScript and resources. Use them sparingly to keep pages fast.

- Use specific start times when relevant: `?t=120` for 2 minutes
- Consider linking to YouTube instead of embedding for lists of videos
- Test on mobile devices - embeds are responsive but take screen space

### PDF Embeds

> [!WARNING]
> Large PDF files can slow down page loading. Optimize PDFs before embedding:
> - Compress images
> - Remove unnecessary metadata
> - Consider splitting large documents

- Keep PDFs under 5MB when possible
- Provide a direct download link as backup
- Use descriptive filenames
- Consider providing a text summary for accessibility

## Alternative: Links

Sometimes a simple link is better than an embed:

**YouTube Link Example:**
[Watch: Introduction to Go Programming](https://www.youtube.com/watch?v=YS4e4q9oBaU)

**PDF Link Example:**
[Download Documentation (PDF)](/files/document.pdf)

::: aside
**When to link vs embed**:
- **Embed** when the content is central to the post
- **Link** when providing supplementary resources
- **Link** when listing multiple videos/documents
:::

## Accessibility Considerations

### For YouTube

- Ensure videos have captions
- Provide transcripts for important content
- Don't rely solely on video to convey critical information
- Use descriptive alt text

### For PDFs

- Ensure PDFs are tagged and accessible
- Provide alternative formats when possible (HTML, plain text)
- Include a text summary of PDF contents
- Test with screen readers

## Advanced: Custom Styling

YouTube embeds use responsive CSS by default. PDFs are displayed in a bordered container.

You can customize the appearance with your own CSS:

```css
/* Customize YouTube embed container */
iframe[src*="youtube.com"] {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

/* Customize PDF container */
.pdf-embed {
    border: 2px solid #333;
    margin: 2rem 0;
}
```

## Code Example: Fetching Video Metadata

Here's how you might fetch YouTube video metadata in Go:

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type VideoInfo struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Duration    string `json:"duration"`
}

func getVideoInfo(videoID string) (*VideoInfo, error) {
    url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=YOUR_API_KEY&part=snippet,contentDetails", videoID)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Items []struct {
            Snippet struct {
                Title       string `json:"title"`
                Description string `json:"description"`
            } `json:"snippet"`
            ContentDetails struct {
                Duration string `json:"duration"`
            } `json:"contentDetails"`
        } `json:"items"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    if len(result.Items) == 0 {
        return nil, fmt.Errorf("video not found")
    }

    return &VideoInfo{
        Title:       result.Items[0].Snippet.Title,
        Description: result.Items[0].Snippet.Description,
        Duration:    result.Items[0].ContentDetails.Duration,
    }, nil
}
```

## Summary

This site supports two types of rich media embeds:

- **YouTube Videos**: Use markdown image syntax with YouTube URLs
- **PDF Documents**: Use the `:::pdf{src="path"}` directive

Both are optimized for responsive display and provide fallbacks for better compatibility.

> [!TIP]
> Check out the [Markdown Syntax Guide](/blog/markdown-syntax-guide) for more formatting options!
