# Blog/Docs Site - Implementation Plan

## Overview

A static site generator + reactions API server in Go. Minimal dependencies, SEO-focused, black & white design with book-like typography.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      CLI Binary                         │
├─────────────────────────────────────────────────────────┤
│  build     → Parse content, generate static HTML        │
│  serve     → Static files + Reactions API + Auth        │
│  dev       → build + serve + hot reload (fsnotify + WS) │
└─────────────────────────────────────────────────────────┘
```

**Build time**: Markdown → HTML with syntax highlighting, TOC extraction, directive parsing
**Runtime**: Serve static files + `/api/reactions` + Google OAuth

---

## Directory Structure

```
site/
├── cmd/
│   └── site/
│       └── main.go              # CLI entry point
├── internal/
│   ├── build/
│   │   ├── build.go             # Orchestrates build process
│   │   ├── markdown.go          # Markdown + directive parsing
│   │   ├── toc.go               # Table of contents extraction
│   │   └── highlight.go         # Syntax highlighting
│   ├── server/
│   │   ├── server.go            # HTTP server
│   │   ├── reactions.go         # Reactions API handlers
│   │   └── auth.go              # Google OAuth
│   ├── models/
│   │   ├── post.go              # Post struct
│   │   ├── collection.go        # Collection struct (Series/Topics)
│   │   └── reaction.go          # Reaction struct
│   └── db/
│       └── sqlite.go            # SQLite setup + queries
├── templates/
│   ├── base.html                # Base layout
│   ├── home.html                # Homepage
│   ├── collection.html          # Collection listing (Series/Topics)
│   ├── profile.html             # Profile/About page
│   └── post.html                # Post page with TOC + sidebar
├── static/
│   ├── css/
│   │   └── style.css            # Single CSS file
│   ├── js/
│   │   └── main.js              # Reactions + hot reload (minimal)
│   ├── images/                  # All images
│   └── files/                   # PDFs and other files
├── content/
│   ├── {collection-slug}/
│   │   ├── _metadata.yml        # Collection config
│   │   ├── banner.jpg           # Optional: for blog series
│   │   └── {post-slug}.md       # Posts
│   └── profile.md               # Profile/About page
├── dist/                        # Build output
├── data/
│   └── sqlite.db                # SQLite database (reactions + search)
└── go.mod
```

---

## Dependencies (Minimal)

| Package | Purpose |
|---------|---------|
| `github.com/yuin/goldmark` | Markdown parsing |
| `github.com/alecthomas/chroma/v2` | Syntax highlighting (180+ languages) |
| `modernc.org/sqlite` | Pure Go SQLite (no CGO) |
| `golang.org/x/oauth2` | Google OAuth |
| `github.com/fsnotify/fsnotify` | File watching for dev mode |
| `gopkg.in/yaml.v3` | YAML frontmatter + metadata |
| `github.com/zmtcreative/gm-alert-callouts` | GitHub/Obsidian-style callouts/alerts |
| `github.com/13rac1/goldmark-embed` | YouTube video embeds |

**Total: 8 external dependencies** (plus their transitive deps)

---

## Content Format

### Collection Metadata (`_metadata.yml`)

```yaml
name: "Getting Started"
description: "Learn the basics"
type: docs  # "docs" or "blog"
icon: "📚"  # Optional: emoji or image path (docs only)
banner: "banner.jpg"  # Optional: banner image filename (blog only)
```

**Collection Types:**
- `docs` (default): Structured documentation (called "Topics" in UI), posts ordered by `order` field, prev/next navigation
- `blog`: Chronological blog posts (called "Series" in UI), ordered by date (newest first), date prominently displayed

### Post Frontmatter

```yaml
---
title: "Introduction"
description: "A brief introduction to the topic"
date: 2024-01-15
updated: 2024-02-01      # optional, shows "Updated on" if present
draft: false             # optional, excluded from build if true
order: 1                 # optional, for ordering within topic
---
```

### Custom Directives

**Callouts** (GitHub/Obsidian compatible using gm-alert-callouts)
```md
> [!NOTE]
> This is informational content with **markdown** support.

> [!WARNING]
> Be careful with this.

> [!TIP]
> A helpful tip.

> [!CAUTION]
> Critical warning.

> [!IMPORTANT]
> Important information.

> [!TIP]-
> This callout is foldable (closed by default).
```

Supported alert types: NOTE, TIP, IMPORTANT, WARNING, CAUTION

**YouTube Embeds** (using goldmark-embed)
```md
![Video Title](https://www.youtube.com/watch?v=VIDEO_ID)
```

Uses standard markdown image syntax - automatically converts YouTube URLs to embedded iframe.

**PDF Viewer** (custom goldmark extension)
```md
:::pdf{src="/files/document.pdf"}
```

**Asides** (custom goldmark extension)
```md
::: aside
Content that appears to the right on desktop, inline on mobile.
Perfect for biographical info, historical context, or fun facts.
:::
```

On desktop (>900px): Floats to the right (240px wide)
On mobile (≤900px): Displays inline as regular block

**Code Blocks** (standard fenced code, enhanced with Chroma)
````md
```go
func main() {
    fmt.Println("Hello")
}
```
````

---

## Design System

### Typography

```css
:root {
  /* Fonts - Book-like readability */
  --font-body: "Source Serif 4", Georgia, serif;
  --font-ui: "Inter", system-ui, sans-serif;
  --font-mono: "JetBrains Mono", monospace;

  /* Scale */
  --text-sm: 0.875rem;
  --text-base: 1.125rem;    /* 18px - optimal reading */
  --text-lg: 1.25rem;
  --text-xl: 1.5rem;
  --text-2xl: 2rem;
  --text-3xl: 2.5rem;

  /* Spacing */
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-6: 1.5rem;
  --space-8: 2rem;
  --space-12: 3rem;
  --space-16: 4rem;

  /* Colors - Black & White */
  --color-bg: #ffffff;
  --color-text: #1a1a1a;
  --color-text-muted: #666666;
  --color-border: #e5e5e5;
  --color-code-bg: #f5f5f5;

  /* Layout */
  --content-width: 680px;
  --toc-width: 240px;
}
```

### Layout (Post Page)

```
┌──────────────────────────────────────────────────────────┐
│  Header: Site Title                              [Auth]  │
├────────────┬─────────────────────────────────────────────┤
│            │                                             │
│  TOC       │  Post Content                               │
│  (sticky)  │  - Title                                    │
│            │  - Meta (date, updated)                     │
│  • Intro   │  - Body                                     │
│  • Setup   │  - Reactions                                │
│  • Usage   │                                             │
│            │                                             │
├────────────┴─────────────────────────────────────────────┤
│  Footer                                                  │
└──────────────────────────────────────────────────────────┘
```

**Mobile**: TOC collapses into hamburger menu or top accordion.

---

## SEO

Each page includes:

```html
<title>{Post Title} | {Topic Name} | {Site Name}</title>
<meta name="description" content="{post.description}">
<link rel="canonical" href="{full_url}">

<!-- Open Graph -->
<meta property="og:title" content="{title}">
<meta property="og:description" content="{description}">
<meta property="og:type" content="article">
<meta property="og:url" content="{url}">

<!-- Article metadata -->
<meta property="article:published_time" content="{date}">
<meta property="article:modified_time" content="{updated}">

<!-- JSON-LD structured data -->
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "{title}",
  "datePublished": "{date}",
  "dateModified": "{updated}",
  "description": "{description}"
}
</script>
```

Additional:
- Semantic HTML (`<article>`, `<nav>`, `<main>`, `<aside>`)
- Proper heading hierarchy (single `<h1>` per page)
- `sitemap.xml` generated at build time
- Clean URLs: `/{topic}/{post}` (no `.html`)

---

## API Endpoints

### Reactions

```
GET  /api/reactions?post={slug}
     → { "👍": 5, "❤️": 3, "🎉": 1 }

POST /api/reactions
     Body: { "post": "topic/post-slug", "emoji": "👍" }
     → Requires auth, toggles reaction (add/remove)

GET  /api/reactions/user?post={slug}
     → ["👍"] (user's reactions, requires auth)
```

### Auth

```
GET  /auth/google          → Redirect to Google OAuth
GET  /auth/google/callback → Handle OAuth callback, set session cookie
GET  /auth/logout          → Clear session
GET  /api/me               → { "email": "...", "name": "..." } or 401
```

### Social Redirects

Vanity URLs that redirect to configured social profiles (301 Permanent Redirect):

```
GET  /github    → Redirects to profile.github URL from site.yml
GET  /linkedin  → Redirects to profile.linkedin URL from site.yml
GET  /email     → Redirects to mailto:{profile.email} from site.yml
```

These endpoints return 404 if the corresponding profile field is not configured in `site.yml`.

**Available Emojis**: 👍 ❤️ 😂 💡 😢 (fixed set, configurable)

---

## Database Schema

```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,           -- Google user ID
    email TEXT NOT NULL,
    name TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL REFERENCES users(id),
    post_slug TEXT NOT NULL,       -- e.g., "getting-started/intro"
    emoji TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, post_slug, emoji)
);

CREATE INDEX idx_reactions_post ON reactions(post_slug);

-- Full-text search index (FTS5)
CREATE VIRTUAL TABLE search_index USING fts5(
    post_slug,
    title,
    description,
    content_text,
    type,           -- "blog" or "docs"
    collection_name
);
```

### Search API

```
GET /api/search?q=query&type=all|blog|docs
    → Returns search results with snippets
    → JSON: [{
        title: "...",
        url: "/collection/post",
        type: "blog" | "docs",
        collection: "Collection Name",
        snippets: ["...highlighted <mark>text</mark>..."]
    }]
```

---

## CLI Commands

```bash
# Build static site to dist/
site build

# Build + serve with hot reload (dev mode, port 3000)
site dev

# Serve production (static files + API)
site serve --port 3000

# Options
--content ./content    # Content directory
--output ./dist        # Output directory
--base-url https://... # For canonical URLs
--port 3000           # Server port (default: 3000 for dev)
```

---

## Hot Reload (Dev Mode)

1. `fsnotify` watches `content/` and `templates/`
2. On change: rebuild affected files
3. WebSocket connection from browser
4. Server sends "reload" message
5. Browser refreshes

```js
// Injected in dev mode only
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = () => location.reload();
```

---

## Implementation Order

### Phase 1: Core Build
1. Project setup, go.mod, directory structure
2. Content parsing (YAML frontmatter, markdown)
3. Directive parser (callouts, pdf)
4. Syntax highlighting integration
5. TOC extraction
6. HTML templates
7. Build command

### Phase 2: Server
8. Static file server
9. Clean URL routing
10. Dev mode with hot reload

### Phase 3: Reactions
11. SQLite setup
12. Google OAuth
13. Reactions API
14. Frontend JS for reactions

### Phase 4: Polish
15. SEO (sitemap, structured data)
16. Mobile responsive TOC
17. Error pages (404)
18. Production hardening

---

## Configuration

`site.yml` in project root:

```yaml
title: "My Site"
description: "Site description for SEO"
base_url: "https://example.com"  # Production base URL
dev_base_url: "http://localhost:3000"  # Local development base URL (optional)

# Google OAuth (for reactions)
# Create a .env file with:
#   GOOGLE_CLIENT_ID=your-client-id
#   GOOGLE_CLIENT_SECRET=your-client-secret

# Profile
profile:
  photo: "/images/profile.jpg"
  bio: "Short bio for home page and profile"
  github: "https://github.com/username"
  linkedin: "https://www.linkedin.com/in/username"
  email: "email@example.com"

# Referrals
referrals:
  - name: "Service Name"
    description: "Why I recommend this"
    url: "https://referral-link.com"
    image: "/images/referral.png"  # optional
```

**Environment Variables:**

Create a `.env` file in the project root (see `.env.example`):
```bash
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here
```

The `.env` file is automatically loaded by the application (not committed to git).

**Available Reaction Emojis:** 👍 ❤️ 😂 💡 😢

---

## Decisions

1. **Collection ordering on homepage**: By most recent post in each collection
2. **Post ordering within collection**:
   - Blog (Series): By date (newest first)
   - Docs (Topics): By `order` field, then title
3. **Homepage**: Profile section at top, then grid with Series/Topics, plus referrals section
4. **Assets**: Images in `static/images/` or in content directories (for banners/icons)
5. **Terminology**: Backend uses "Collection", UI shows "Series" (blog) or "Topics" (docs)
6. **Search**: Server-side using SQLite FTS5, searches full post content with snippet extraction

---

## New Features (2024)

### Profile & About Page
- Profile section on homepage with photo, bio, social links (GitHub, LinkedIn, Email)
- Dedicated `/profile` page with full content from `content/profile.md`
- Configurable via `site.yml` profile section
- Social redirect vanity URLs: `/github`, `/linkedin`, `/email` (301 redirects to configured profiles)

### Search
- Full-text search across all content using SQLite FTS5
- Search API endpoint: `/api/search?q=...&type=all|blog|docs`
- Client-side UI with overlay dropdown, filters, debouncing
- Returns highlighted snippets with `<mark>` tags

### Collections Enhancement
- Icons for docs topics (emoji or image)
- Banners for blog series (image file in collection directory)
- Terminology: "Series" for blog, "Topics" for docs

### UI Enhancements
- Collapsible sidebar on desktop (persists state in localStorage)
- Date sorting on collection pages (newest/oldest/updated)
- Pastel color scheme for callout boxes
- Referrals section on home page

### Social Previews
- Open Graph meta tags for rich previews on Facebook, LinkedIn, Slack
- Twitter Card meta tags for enhanced Twitter/X previews
- Automatic image selection:
  - Blog posts: Use collection banner if available
  - Docs: Use collection icon or default site image
  - Fallback to site-wide default image
- Dynamic titles and descriptions from post/page metadata

---

## Notes

- No JavaScript frameworks - vanilla JS only for search, reactions, sidebar, hot reload
- CSS is single file with minification and content hashing
- Fonts loaded from Google Fonts (or self-hosted for privacy)
- All HTML valid, accessible (ARIA labels, skip links, focus states)
- Hot reload with build timing instrumentation for development

---

## OAuth + Comments Feature (Completed)

### Phase 1: GitHub OAuth

#### 1.1 Environment Variables
- Add `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` to `.env.example`

#### 1.2 Backend Changes (`internal/server/auth.go`)
- Add GitHub OAuth config function (similar to `getOAuthConfig()`)
- Add `/auth/github` endpoint to initiate OAuth flow
- Add `/auth/github/callback` to handle callback
- GitHub API returns: `id`, `login`, `name`, `email`, `avatar_url`

#### 1.3 Route Registration (`server.go`)
- Register new GitHub auth routes

---

### Phase 2: Login Popup

#### 2.1 API Endpoint
- Add `GET /api/auth/providers` that returns available providers based on env vars:
  ```json
  {
    "providers": [
      {"id": "google", "name": "Google"},
      {"id": "github", "name": "GitHub"}
    ]
  }
  ```

#### 2.2 HTML/CSS
- Add login modal to templates (reusable partial)
- Provider icons (SVG inline or CSS)
- Simple, clean design

#### 2.3 JavaScript
- Replace direct redirect (`/auth/google?redirect=...`) with modal
- Modal shows available providers fetched from `/api/auth/providers`
- Clicking a provider redirects to `/auth/{provider}?redirect=...`
- Modal can be triggered from anywhere (reactions, comments, etc.)

---

### Phase 3: Comments System

#### 3.1 Database Schema (`internal/db/sqlite.go`)
```sql
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL REFERENCES users(id),
    post_slug TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_slug);
CREATE INDEX IF NOT EXISTS idx_comments_created ON comments(post_slug, created_at);
```

#### 3.2 User Model Update
- Add `avatar_url` field to users table (for displaying in comments)
- Update OAuth callbacks to store avatar URL

#### 3.3 DB Methods
- `CreateComment(userID, postSlug, content string) (*Comment, error)`
- `GetComments(postSlug string) ([]CommentWithUser, error)` - includes user name/avatar
- `UpdateComment(id int, userID, content string) error` - only own comments
- `DeleteComment(id int, userID string) error` - only own comments

#### 3.4 API Endpoints (`internal/server/comments.go`)
- `GET /api/comments?post=<slug>` - Get all comments for a post
- `POST /api/comments` - Create comment (requires auth)
- `PUT /api/comments/{id}` - Edit own comment (requires auth)
- `DELETE /api/comments/{id}` - Delete own comment (requires auth)

#### 3.5 Markdown Rendering
- Server-side rendering using existing goldmark dependency
- Sanitize HTML output to prevent XSS
- Return rendered HTML in API response

#### 3.6 Frontend (`templates/partials/comments.html` + `main.js`)
- Comments section below reactions
- Textarea for writing comments (with markdown preview?)
- Display comments with user avatar, name, timestamp
- Edit/Delete buttons for own comments
- "Sign in to comment" prompt for logged-out users (triggers login modal)

---

### File Changes Summary

| File | Changes |
|------|---------|
| `.env.example` | Add GitHub OAuth vars |
| `internal/server/auth.go` | Add GitHub OAuth |
| `internal/server/server.go` | Register new routes |
| `internal/server/comments.go` | New file - comments API |
| `internal/db/sqlite.go` | Comments table + methods, avatar field |
| `internal/models/models.go` | Comment model |
| `templates/partials/login-modal.html` | New file - login popup |
| `templates/partials/comments.html` | New file - comments section |
| `templates/partials/post-footer.html` | Include comments partial |
| `static/js/main.js` | Login modal + comments logic |
| `static/css/style.css` | Styles for modal + comments |

---

### Decisions

1. **Avatar display**: Yes - store and display user avatars from OAuth providers
2. **Markdown preview**: Live preview while typing (best UX)
3. **Comment ordering**: Newest first
