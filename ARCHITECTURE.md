# Architecture Documentation

This document provides a detailed technical overview of the matthiasbrat.com static site generator and server.

## Table of Contents

- [System Overview](#system-overview)
- [Build System](#build-system)
- [Server Architecture](#server-architecture)
- [Content Pipeline](#content-pipeline)
- [Database Schema](#database-schema)
- [Template System](#template-system)
- [Asset Processing](#asset-processing)
- [Authentication & Authorization](#authentication--authorization)
- [Performance Optimizations](#performance-optimizations)

## System Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Interface                         │
│                  (cmd/site/main.go)                         │
└────────────┬──────────────┬─────────────────┬───────────────┘
             │              │                 │
             │              │                 │
    ┌────────▼──────┐  ┌───▼──────┐  ┌──────▼────────┐
    │  build mode   │  │ dev mode │  │  serve mode   │
    │  (one-shot)   │  │ (watch)  │  │ (production)  │
    └────────┬──────┘  └───┬──────┘  └──────┬────────┘
             │              │                 │
             └──────────────┼─────────────────┘
                            │
             ┌──────────────▼──────────────┐
             │      Build System           │
             │  (internal/build/)          │
             │                             │
             │  ┌─────────────────────┐   │
             │  │ Content Loader      │   │
             │  │ Markdown Renderer   │   │
             │  │ Template Engine     │   │
             │  │ Asset Processor     │   │
             │  │ Search Indexer      │   │
             │  │ OG Image Generator  │   │
             │  └─────────────────────┘   │
             └──────────────┬──────────────┘
                            │
                   ┌────────▼─────────┐
                   │  Static Output   │
                   │    (dist/)       │
                   └────────┬─────────┘
                            │
                   ┌────────▼─────────┐
                   │   HTTP Server    │
                   │ (internal/server)│
                   │                  │
                   │  - Static Files  │
                   │  - API Endpoints │
                   │  - OAuth         │
                   │  - WebSocket     │
                   └──────────────────┘
```

### Design Principles

1. **Simplicity**: Minimal dependencies, straightforward code
2. **Performance**: Fast builds, optimized output, concurrent processing
3. **Flexibility**: Extensible markdown, customizable templates
4. **Developer Experience**: Hot reload, clear error messages, intuitive CLI
5. **Production Ready**: Docker support, proper error handling, graceful shutdown

## Build System

### Build Process Flow

```
1. Load Configuration (site.yml)
   ↓
2. Load Content (markdown files + frontmatter)
   ↓
3. Parse Collections & Posts
   ↓
4. Process Static Assets (minify, hash)
   ↓
5. Generate OG Images (concurrent)
   ↓
6. Render Templates
   ↓
7. Index Content for Search
   ↓
8. Generate Sitemap
   ↓
9. Write Output Files
```

### Key Components

#### Content Loader (`internal/build/content/loader.go`)

Responsible for:
- Discovering markdown files in `content/` directory
- Parsing YAML frontmatter
- Organizing posts into collections
- Handling nested collection hierarchies

**Algorithm:**
1. Scan directory recursively for `_metadata.yml` files
2. For each collection:
   - Load metadata (name, description, type, icon, order)
   - Find all `.md` files in the directory
   - Parse frontmatter and content for each post
   - Sort posts by date or order
3. Build collection tree (parent-child relationships)

#### Markdown Renderer (`internal/build/markdown/renderer.go`)

Built on top of [goldmark](https://github.com/yuin/goldmark) with custom extensions:

**Extensions:**
- **GitHub Flavored Markdown**: Tables, strikethrough, task lists
- **Syntax Highlighting**: Via Chroma, supports 180+ languages
- **Callouts**: `> [!NOTE]`, `> [!TIP]`, `> [!WARNING]`, etc.
- **Asides**: `::: aside` blocks
- **YouTube Embeds**: Auto-detect YouTube URLs in image syntax
- **PDF Embeds**: Embed PDFs with viewer

**Rendering Pipeline:**
1. Parse frontmatter (YAML)
2. Extract and process code blocks for syntax highlighting
3. Apply markdown extensions
4. Convert to HTML
5. Post-process for custom elements

#### Asset Processor (`internal/build/assets/processor.go`)

Handles static assets with:
- **Minification**: CSS and JavaScript minification via tdewolff/minify
- **Content Hashing**: Generate hashed filenames for cache busting (e.g., `style.css` → `style.a1b2c3d4.css`)
- **File Copying**: Copy images, fonts, and other assets
- **Source Maps**: Preserve debugging information

**Hash Algorithm:**
```go
hash := sha256.Sum256(content)
hashedName := filename + "." + hex(hash[:8]) + extension
```

#### OG Image Generator (`internal/build/og/generator.go`)

Generates Open Graph images for social sharing:

**Process:**
1. Load custom fonts (SourceSerif4)
2. Create 1200x630 canvas
3. Draw gradient background
4. Render post title (word-wrapped)
5. Add site name and author photo
6. Save as PNG to `/og/` directory

**Optimizations:**
- Concurrent generation using goroutines
- Semaphore to limit concurrent image operations
- Reuse font contexts

#### Search Indexer (`internal/build/search/indexer.go`)

SQLite FTS5 (Full-Text Search) implementation:

**Indexed Fields:**
- Post title (weight: 10)
- Post description (weight: 5)
- Post content (weight: 1)
- Collection name (weight: 3)

**Features:**
- Stemming and ranking
- Phrase matching
- Prefix matching for autocomplete
- Relevance scoring

## Server Architecture

### Server Modes

#### Development Mode (`dev`)

- **Auto-rebuild**: Watches file changes and rebuilds
- **Hot reload**: WebSocket-based live reload
- **Ephemeral user**: Creates temporary user for testing reactions
- **Debug logging**: Detailed logs for development

**File Watching:**
```go
Watched directories:
  - content/    → triggers full rebuild
  - templates/  → triggers full rebuild
  - static/     → copies changed files only
```

**WebSocket Protocol:**
```json
{
  "type": "reload",
  "message": "Files updated"
}
```

#### Production Mode (`serve`)

- **Static file serving**: Efficient file serving with caching headers
- **API endpoints**: Reactions, comments, search
- **OAuth authentication**: Google OAuth 2.0
- **Database persistence**: SQLite for user data
- **Graceful shutdown**: Clean database closure

### HTTP Routes

```
Static Routes:
  GET  /                    → index.html
  GET  /blog                → blog/index.html
  GET  /docs                → docs/index.html
  GET  /profile             → profile.html
  GET  /{collection}/{post} → {collection}/{post}/index.html
  GET  /static/*            → static files

API Routes:
  GET    /api/search                          → Full-text search
  GET    /api/posts/:slug/reactions           → Get reactions
  POST   /api/posts/:slug/reactions           → Add reaction (auth)
  DELETE /api/posts/:slug/reactions/:emoji    → Remove reaction (auth)
  GET    /api/posts/:slug/comments            → Get comments
  POST   /api/posts/:slug/comments            → Add comment (auth)

Auth Routes:
  GET  /auth/login     → Redirect to Google OAuth
  GET  /auth/callback  → OAuth callback
  GET  /auth/logout    → Logout and clear session

Dev Routes (dev mode only):
  GET  /ws             → WebSocket for hot reload
```

### Request Flow

```
Client Request
   ↓
HTTP Server (Server.ServeHTTP)
   ↓
Route Matching
   ↓
┌──────────────┬────────────────┬──────────────┐
│              │                │              │
Static File    API Endpoint     Auth Flow
   ↓              ↓                ↓
Serve from     Parse Request    OAuth2
dist/          Check Auth       Redirect
               Query DB              ↓
               Return JSON      Set Session
```

## Content Pipeline

### Content Structure

```
content/
├── blog/
│   ├── _metadata.yml          # Collection metadata
│   ├── post1.md               # Individual post
│   ├── series1/
│   │   ├── _metadata.yml      # Nested collection
│   │   ├── part1.md
│   │   └── part2.md
│   └── series2/
│       ├── _metadata.yml
│       └── post.md
└── docs/
    ├── getting-started/
    │   ├── _metadata.yml
    │   ├── installation.md
    │   └── configuration.md
    └── api-reference/
        ├── _metadata.yml
        └── endpoints.md
```

### Frontmatter Schema

```yaml
---
title: "Post Title"           # Required: Display title
description: "SEO description" # Required: Meta description
date: 2025-01-01              # Required: Publication date (YYYY-MM-DD)
updated: 2025-01-15           # Optional: Last updated date
draft: false                  # Optional: Hide from production
order: 1                      # Optional: Custom ordering
slug: "custom-slug"           # Optional: Override URL slug
---
```

### Collection Types

**Topic** (`type: topic`):
- Used for documentation
- Ordered by `order` field or alphabetically
- Displays as documentation index
- Example: Getting Started, API Reference

**Series** (`type: series`):
- Used for blog post series
- Ordered by date (newest first)
- Displays as blog listing with pagination
- Example: Go Basics, DevOps Diary

### URL Structure

```
/                                    # Homepage
/blog                                # Main blog listing
/blog/page/2                         # Paginated blog
/blog/post-slug                      # Individual blog post
/docs                                # Documentation index
/docs/getting-started                # Collection page
/docs/getting-started/installation   # Documentation page
/profile                             # Profile page
/referrals                           # Referrals page
```

## Database Schema

### SQLite Tables

#### `posts` (FTS5 Virtual Table)
```sql
CREATE VIRTUAL TABLE posts USING fts5(
  slug,           -- Post URL slug
  title,          -- Post title
  description,    -- Post description
  content,        -- Full post content (markdown)
  collection      -- Collection name
);
```

#### `users`
```sql
CREATE TABLE users (
  id TEXT PRIMARY KEY,        -- OAuth provider ID
  email TEXT NOT NULL,
  name TEXT NOT NULL,
  avatar_url TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### `reactions`
```sql
CREATE TABLE reactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  post_slug TEXT NOT NULL,
  user_id TEXT NOT NULL,
  emoji TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  UNIQUE(post_slug, user_id, emoji)
);
```

#### `comments`
```sql
CREATE TABLE comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  post_slug TEXT NOT NULL,
  user_id TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Indexes

```sql
CREATE INDEX idx_reactions_post ON reactions(post_slug);
CREATE INDEX idx_reactions_user ON reactions(user_id);
CREATE INDEX idx_comments_post ON comments(post_slug);
CREATE INDEX idx_comments_user ON comments(user_id);
```

## Template System

### Template Hierarchy

```
templates/
├── base.html              # Base layout (header, footer, scripts)
├── home.html              # Homepage
├── blog.html              # Blog listing with pagination
├── collection.html        # Collection/topic listing
├── post.html              # Individual post
├── profile.html           # Profile page
├── docs.html              # Documentation index
├── referrals.html         # Referrals page
└── partials/
    ├── sidebar-toc.html   # Table of contents sidebar
    ├── sidebar-blog.html  # Blog sidebar
    ├── post-header.html   # Post metadata header
    └── post-pager.html    # Previous/next navigation
```

### Template Composition

Each page template:
1. Defines its own template block
2. Is composed with `base.html`
3. Can include partials
4. Has access to template functions

**Example:**
```html
{{ define "post.html" }}
<article>
  {{ template "post-header.html" . }}
  <div class="content">
    {{ .Post.Content | safeHTML }}
  </div>
  {{ template "post-pager.html" . }}
</article>
{{ end }}
```

### Template Functions

Custom functions available in templates:

- `safeHTML` - Render HTML without escaping
- `safeCSS` - Render CSS without escaping
- `criticalCSS` - Inline critical CSS
- `asset` - Get hashed asset path (e.g., `/css/style.a1b2.css`)
- `hasPrefix` - String prefix checking

## Asset Processing

### Asset Pipeline

```
Static Assets (static/)
   ↓
1. Read files
   ↓
2. Minify (CSS, JS)
   ↓
3. Generate SHA-256 hash
   ↓
4. Create hashed filename
   ↓
5. Copy to output (dist/)
   ↓
6. Update asset map
   ↓
Templates use {{ asset "path" }}
   ↓
Hashed paths in HTML
```

### Cache Strategy

**Development:**
- No hashing
- No minification
- Cache-Control: no-cache

**Production:**
- Content hashing (e.g., `style.a1b2c3d4.css`)
- Minification enabled
- Cache-Control: max-age=31536000 (1 year)
- Immutable assets

### CSS Organization

```
static/css/
├── critical.css    # Above-fold styles (inlined)
├── style.css       # Main stylesheet
└── syntax.css      # Syntax highlighting themes
```

**Critical CSS** is inlined in `<head>` for faster First Contentful Paint.

## Authentication & Authorization

### OAuth Flow

```
1. User clicks "Login"
   ↓
2. Redirect to /auth/login
   ↓
3. Redirect to Google OAuth
   ↓
4. User authorizes
   ↓
5. Google redirects to /auth/callback?code=...
   ↓
6. Exchange code for access token
   ↓
7. Fetch user profile from Google
   ↓
8. Create or update user in database
   ↓
9. Set session cookie
   ↓
10. Redirect to original page
```

### Session Management

**Cookie:**
- Name: `session_token`
- HttpOnly: true
- Secure: true (in production)
- SameSite: Lax
- Max-Age: 30 days

**Session Store:**
- In-memory map: `sessionID → userID`
- Cleaned up on logout or server restart
- No persistent sessions (stateless)

### Authorization

**Public endpoints:**
- GET /api/posts/:slug/reactions
- GET /api/posts/:slug/comments
- GET /api/search

**Authenticated endpoints:**
- POST /api/posts/:slug/reactions
- DELETE /api/posts/:slug/reactions/:emoji
- POST /api/posts/:slug/comments

**Authorization Check:**
```go
func (s *Server) requireAuth(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := s.getUserFromSession(r)
        if user == nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        handler(w, r)
    }
}
```

## Performance Optimizations

### Build Performance

1. **Concurrent OG Image Generation**: Uses goroutines with semaphore
2. **Incremental Builds** (dev mode): Only rebuilds changed files
3. **Template Caching**: Parse templates once
4. **Efficient File Walking**: Single pass through content directory

**Benchmarks:**
- Small site (10 posts): ~100ms
- Medium site (50 posts): ~300ms
- Large site (200 posts): ~1s

### Runtime Performance

1. **Static File Serving**: Zero processing, direct file serving
2. **Database Indexing**: Proper indexes on frequently queried columns
3. **Asset Hashing**: Long-term caching for static assets
4. **Minification**: Reduced file sizes (~30% smaller CSS/JS)

### Memory Management

1. **Streaming**: Large files are streamed, not loaded into memory
2. **Connection Pooling**: SQLite connection pool
3. **Goroutine Limits**: Semaphores prevent unlimited goroutine spawning

### Docker Optimizations

1. **Multi-stage build**: Separate build and runtime stages
2. **Distroless base**: Minimal runtime image (~20MB)
3. **Layer caching**: Dependencies cached separately from code
4. **Static binary**: CGO_ENABLED=0 for portability

## Monitoring & Observability

### Logging

Structured logging with levels:
- INFO: Normal operations
- WARNING: Recoverable issues
- ERROR: Failures requiring attention

**Example logs:**
```
2025-01-01 12:00:00 INFO Starting server on :3000
2025-01-01 12:00:01 INFO Initial build complete in 234ms
2025-01-01 12:00:05 INFO File changed: content/blog/post.md
2025-01-01 12:00:05 INFO Rebuild complete in 156ms
```

### Error Handling

1. **User-Facing Errors**: Clear error messages in UI
2. **Internal Errors**: Logged with stack traces
3. **Graceful Degradation**: Fallbacks for non-critical features

### Graceful Shutdown

```go
// Listen for SIGINT/SIGTERM
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// On signal:
1. Stop accepting new connections
2. Finish in-flight requests (with timeout)
3. Close database connections
4. Clean up ephemeral data (dev mode)
5. Exit
```

## Future Enhancements

### Planned Features

1. **Incremental Builds**: Only rebuild changed pages
2. **Image Optimization**: Automatic image resizing and format conversion
3. **RSS Feeds**: Generate RSS/Atom feeds for blog
4. **i18n Support**: Multi-language content
5. **Plugin System**: Load external markdown extensions
6. **Admin UI**: Web-based content management
7. **Analytics**: Privacy-focused analytics
8. **CDN Integration**: Automatic asset upload to CDN

### Performance Goals

- Build times < 50ms for incremental changes
- First Contentful Paint < 500ms
- Lighthouse score 100 across all metrics
- Binary size < 10MB

## References

- [goldmark](https://github.com/yuin/goldmark) - Markdown parser
- [chroma](https://github.com/alecthomas/chroma) - Syntax highlighter
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) - Pure Go SQLite
- [gg](https://github.com/fogleman/gg) - 2D graphics library
- [OAuth 2.0](https://oauth.net/2/) - Authentication protocol
