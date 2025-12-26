# matthiasbrat.com

A lightweight static site generator and server written in Go, designed for blogs and documentation sites with minimal dependencies and fast build times.

## Features

- **Static Site Generation**: Build fast, optimized static sites from Markdown content
- **Development Server**: Hot-reload development server with live updates
- **Production Server**: Serve static sites with reactions API
- **Full Markdown Support**: GitHub Flavored Markdown with extensions
- **Syntax Highlighting**: Support for 180+ programming languages via Chroma
- **Custom Extensions**:
  - Callouts and alerts (Tip, Warning, Note, etc.)
  - YouTube video embeds
  - PDF document embeds
- **Full-Text Search**: SQLite-powered search indexing
- **Emoji Reactions**: Google OAuth-based reactions system for blog posts
- **SEO Optimized**: Automatic sitemap generation, Open Graph images, and structured data
- **Responsive Design**: Mobile-first responsive templates
- **Zero JavaScript Frameworks**: Vanilla JavaScript only for minimal overhead

## Quick Start

### Prerequisites

- Go 1.24 or later
- Git (optional, for cloning)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Matthiasbrat/matthiasbrat.com.git
cd matthiasbrat.com
```

2. Install dependencies:
```bash
go mod download
```

3. Build the binary:
```bash
go build -o site ./cmd/site
```

## Usage

The site binary provides three main commands: `build`, `dev`, and `serve`.

### Build Mode

Generate a static site to the output directory:

```bash
./site build [options]
```

**Options:**
- `-content` - Content directory (default: `content`)
- `-output` - Output directory (default: `dist`)
- `-base-url` - Base URL for canonical links (defaults to `site.yml`)

**Example:**
```bash
./site build -content ./content -output ./dist -base-url https://matthiasbrat.com
```

### Development Mode

Run a development server with hot-reload:

```bash
./site dev [options]
```

**Options:**
- `-port` - Port to serve on (default: `3000`)
- `-content` - Content directory (default: `content`)
- `-output` - Output directory (default: `dist`)
- `-base-url` - Base URL (defaults to `site.yml` or `http://localhost:<port>`)

**Example:**
```bash
./site dev -port 3000
```

The development server will:
- Watch for changes in `content/`, `templates/`, and `static/` directories
- Automatically rebuild when files change
- Broadcast live reload events via WebSocket
- Create an ephemeral test user for reactions (cleaned up on shutdown)

### Production Server Mode

Serve pre-built static files with reactions API:

```bash
./site serve [options]
```

**Options:**
- `-port` - Port to serve on (default: `8080`)
- `-output` - Output directory (default: `dist`)
- `-base-url` - Base URL for production (defaults to `site.yml`)

**Example:**
```bash
./site serve -port 8080
```

## Configuration

### Site Configuration (`site.yml`)

Create a `site.yml` file in the root directory:

```yaml
title: "Your Site Name"
description: "Your site description"
base_url: "https://yoursite.com"
dev_base_url: "http://localhost:3000"
default_social_image: "/me.jpg"

# Profile configuration
profile:
  photo: "/me.jpg"
  bio: "Your bio"
  github: "https://github.com/yourusername"
  linkedin: "https://www.linkedin.com/in/yourusername"
  email: "you@example.com"

# Referrals (optional)
referrals:
  - name: "Person Name"
    photo: "https://example.com/photo.jpg"
    github: "https://github.com/username"
    linkedin: "https://linkedin.com/in/username"
    website: "https://example.com"
```

### Google OAuth (for reactions)

For the reactions feature, create a `.env` file:

```bash
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
```

You'll need to:
1. Create a Google Cloud project
2. Enable the Google+ API
3. Create OAuth 2.0 credentials
4. Set authorized redirect URIs to `http://localhost:3000/auth/callback` (dev) and your production URL

## Project Structure

```
.
├── cmd/
│   └── site/           # Main application entry point
│       └── main.go
├── content/            # Markdown content
│   ├── blog/           # Blog posts
│   │   ├── _metadata.yml
│   │   └── *.md
│   └── docs/           # Documentation
│       ├── _metadata.yml
│       └── */
├── internal/           # Internal packages
│   ├── build/          # Build system
│   │   ├── assets/     # Asset processing
│   │   ├── content/    # Content loading
│   │   ├── markdown/   # Markdown rendering
│   │   ├── og/         # Open Graph image generation
│   │   └── search/     # Search indexing
│   ├── db/             # SQLite database
│   ├── models/         # Data models
│   └── server/         # HTTP server
├── static/             # Static assets
│   ├── css/
│   ├── fonts/
│   ├── js/
│   └── images/
├── templates/          # HTML templates
│   ├── base.html
│   ├── partials/
│   └── *.html
├── Dockerfile          # Docker configuration
├── go.mod              # Go dependencies
└── site.yml            # Site configuration
```

## Content Organization

### Collections

Content is organized into **collections** (topics or series). Each collection has a `_metadata.yml` file:

```yaml
name: "Collection Name"
description: "Collection description"
slug: "collection-slug"
type: "topic"  # or "series"
icon: "📚"     # Optional emoji or path
order: 1       # Display order
```

**Collection Types:**
- `topic`: Documentation-style collection (e.g., API Reference, Getting Started)
- `series`: Blog series with chronological posts

### Blog Posts

Create markdown files in `content/blog/`:

```markdown
---
title: "Post Title"
description: "Post description for SEO"
date: 2025-01-01
---

Your content here...
```

**Frontmatter Fields:**
- `title` (required): Post title
- `description` (required): SEO description
- `date` (required): Publication date (YYYY-MM-DD)
- `updated` (optional): Last updated date
- `draft` (optional): Set to `true` to hide from production

### Nested Collections

You can organize posts into subcollections:

```
content/blog/
├── _metadata.yml           # Main blog collection
├── web-development/
│   ├── _metadata.yml       # Web dev subcollection
│   ├── rest-api-design.md
│   └── typescript-tips.md
└── golang-fundamentals/
    ├── _metadata.yml
    └── getting-started.md
```

## Markdown Extensions

### Callouts

```markdown
> [!NOTE]
> This is a note callout

> [!TIP]
> This is a tip callout

> [!WARNING]
> This is a warning callout

> [!IMPORTANT]
> This is an important callout
```

### Asides

```markdown
::: aside
This content will appear in a styled aside box
:::
```

### YouTube Embeds

```markdown
![YouTube Video](https://www.youtube.com/watch?v=VIDEO_ID)
```

### PDF Embeds

```markdown
![PDF Document](/path/to/document.pdf)
```

## Development

### Building from Source

```bash
# Build binary
go build -o site ./cmd/site

# Run tests
go test ./...

# Run with race detection
go run -race ./cmd/site dev
```

### Hot Reload

The development server watches these directories:
- `content/` - Markdown content
- `templates/` - HTML templates
- `static/` - Static assets

Changes trigger automatic rebuilds and browser refresh via WebSocket.

## Deployment

### GitHub Pages (Recommended)

This project includes automated GitHub Pages deployment via GitHub Actions.

**Quick Setup:**

1. Enable GitHub Pages in your repository settings:
   - Go to Settings → Pages
   - Select "GitHub Actions" as the source

2. Configure your custom domain:
   - Add your domain in Settings → Pages → Custom domain
   - Update DNS records (see [DEPLOYMENT.md](DEPLOYMENT.md) for details)

3. Push to the `main` branch - deployment happens automatically!

See [DEPLOYMENT.md](DEPLOYMENT.md) for complete deployment documentation including:
- DNS configuration
- Custom domain setup
- SSL/HTTPS configuration
- Troubleshooting guide
- Alternative deployment methods

### Docker

Build and run with Docker:

```bash
# Build image
docker build -t matthiasbrat-site .

# Run container
docker run -p 3000:3000 \
  -v $(pwd)/data:/app/data \
  -e GOOGLE_CLIENT_ID=your-id \
  -e GOOGLE_CLIENT_SECRET=your-secret \
  matthiasbrat-site
```

The Dockerfile uses multi-stage builds for minimal image size (~20MB final image).

### Traditional Deployment

1. Build static site:
```bash
./site build -base-url https://yoursite.com
```

2. Deploy `dist/` directory to:
   - Static hosting (Netlify, Vercel, GitHub Pages)
   - CDN (CloudFront, Cloudflare)
   - Object storage (S3, GCS)

3. For reactions, deploy server:
```bash
./site serve -port 8080 -base-url https://yoursite.com
```

### Environment Variables

- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `PORT` - Server port (optional, overrides `-port` flag)

## Performance

- **Build Times**: ~100-500ms for typical sites (50-100 pages)
- **Binary Size**: ~15MB (static build)
- **Docker Image**: ~20MB (distroless base)
- **Page Load**: First Contentful Paint < 1s (with CDN)
- **Lighthouse Score**: 95+ (Performance, SEO, Accessibility)

## Technology Stack

- **Language**: Go 1.24
- **Markdown**: [goldmark](https://github.com/yuin/goldmark) with extensions
- **Syntax Highlighting**: [chroma](https://github.com/alecthomas/chroma)
- **Database**: SQLite ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite))
- **OAuth**: Google OAuth 2.0
- **OG Images**: Generated with [gg](https://github.com/fogleman/gg)
- **Asset Minification**: [tdewolff/minify](https://github.com/tdewolff/minify)

## API Endpoints

When running in `serve` mode, the following API endpoints are available:

- `GET /api/search?q=query` - Full-text search
- `GET /api/posts/:slug/reactions` - Get post reactions
- `POST /api/posts/:slug/reactions` - Add reaction (requires auth)
- `DELETE /api/posts/:slug/reactions/:emoji` - Remove reaction (requires auth)
- `GET /api/posts/:slug/comments` - Get post comments
- `POST /api/posts/:slug/comments` - Add comment (requires auth)
- `GET /auth/login` - Initiate OAuth login
- `GET /auth/callback` - OAuth callback
- `GET /auth/logout` - Logout

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is available as open source. Feel free to use and modify for your own sites.

## Acknowledgments

- Built with [goldmark](https://github.com/yuin/goldmark) for markdown processing
- Syntax highlighting powered by [chroma](https://github.com/alecthomas/chroma)
- Inspired by modern static site generators like Hugo and Jekyll
