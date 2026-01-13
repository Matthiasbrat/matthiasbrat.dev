# Deployment Guide

This document explains how to deploy the matthiasbrat.dev site to GitHub Pages.

## GitHub Pages Deployment

The site is configured to deploy automatically to GitHub Pages at https://matthiasbrat.dev using GitHub Actions.

### Workflow

**File:** `.github/workflows/deploy.yml`

This workflow uses the official GitHub Pages deployment action:
- Builds the site using Go 1.24
- Generates static files to `dist/`
- Deploys using `actions/deploy-pages@v4`
- Automatically handles CNAME and .nojekyll files

### Setup Instructions

#### Initial Setup

1. **Enable GitHub Pages in Repository Settings:**
   - Go to your repository on GitHub
   - Navigate to Settings → Pages
   - Under "Build and deployment":
     - Select "GitHub Actions" as the source

2. **Configure Custom Domain:**
   - In Settings → Pages → Custom domain
   - Enter: `matthiasbrat.dev`
   - Wait for DNS check to complete

3. **Configure DNS Records:**

   Add the following DNS records at your domain registrar:

   **A records (apex domain):**
   ```
   Type: A
   Name: @
   Value: 185.199.108.153

   Type: A
   Name: @
   Value: 185.199.109.153

   Type: A
   Name: @
   Value: 185.199.110.153

   Type: A
   Name: @
   Value: 185.199.111.153
   ```

   **AAAA records (IPv6):**
   ```
   Type: AAAA
   Name: @
   Value: 2606:50c0:8000::153

   Type: AAAA
   Name: @
   Value: 2606:50c0:8001::153

   Type: AAAA
   Name: @
   Value: 2606:50c0:8002::153

   Type: AAAA
   Name: @
   Value: 2606:50c0:8003::153
   ```

4. **Enable HTTPS:**
   - In Settings → Pages
   - Check "Enforce HTTPS"
   - GitHub will automatically provision an SSL certificate

#### Deployment Trigger

The workflow is triggered on:
- **Automatic:** Every push to the `main` branch
- **Manual:** Via the "Actions" tab → "Run workflow" button

### Build Process

The deployment workflow performs these steps:

1. **Checkout:** Clones the repository
2. **Setup Go:** Installs Go 1.24
3. **Install Dependencies:** Runs `go mod download`
4. **Build Binary:** Compiles the site generator
5. **Build Site:** Generates static files with `./site build -base-url https://matthiasbrat.dev`
6. **Add Files:** Includes CNAME and .nojekyll
7. **Deploy:** Pushes to GitHub Pages

### Deployment Output

The built site includes:
- HTML files for all pages
- CSS and JavaScript (minified and hashed)
- Static assets (images, fonts)
- Search index (SQLite database)
- Open Graph images
- Sitemap.xml
- CNAME file (for custom domain)
- .nojekyll file (disables Jekyll processing)

### Monitoring Deployments

1. **View Workflow Runs:**
   - Go to the "Actions" tab in your repository
   - Click on the latest workflow run
   - Check build logs and deployment status

2. **Check Deployment:**
   - Visit https://matthiasbrat.dev
   - Verify content is updated
   - Check browser console for errors

3. **Troubleshooting Failed Deployments:**
   - Check the Actions tab for error messages
   - Verify Go build succeeds locally
   - Ensure all required files are committed
   - Check DNS settings if domain doesn't resolve

### Local Testing Before Deployment

Before pushing to main, test the build locally:

```bash
# Build the site
go build -o site ./cmd/site
./site build -base-url https://matthiasbrat.dev

# Serve locally to test
./site serve -output dist -port 8080

# Visit http://localhost:8080
```

### Environment Variables

No environment variables are required for the GitHub Actions deployment. The build is purely static.

If you want to enable reactions/comments in production:
1. Set up Google OAuth credentials
2. Add secrets to repository settings:
   - `GOOGLE_CLIENT_ID`
   - `GOOGLE_CLIENT_SECRET`
3. Deploy the server component separately (not to GitHub Pages)

### Deployment Best Practices

1. **Test Locally First:** Always build and test locally before pushing
2. **Use Branch Protection:** Require PR reviews before merging to main
3. **Monitor Actions:** Check workflow runs for failures
4. **Cache Dependencies:** The workflow caches Go modules for faster builds
5. **Semantic Versioning:** Tag releases with version numbers

### Rollback

To rollback to a previous deployment:

1. Find the commit hash of the working version
2. Go to Actions → Deploy workflow → Run workflow
3. Select the branch/tag with the working version
4. Or revert the main branch to the working commit

### Performance Optimization

The deployment includes several optimizations:
- **Asset Hashing:** Cache-busting for CSS/JS
- **Minification:** Reduced file sizes
- **GZIP Compression:** GitHub Pages automatically compresses files
- **CDN:** GitHub's global CDN for fast delivery
- **Caching:** Long-term caching for hashed assets

### Security Considerations

1. **HTTPS:** Always use HTTPS in production
2. **Content Security Policy:** Consider adding CSP headers
3. **Subresource Integrity:** Consider adding SRI for external scripts
4. **Secrets:** Never commit API keys or secrets
5. **Branch Protection:** Protect main branch from force pushes

### Troubleshooting

#### Site not updating after push
- Check Actions tab for workflow status
- Clear browser cache
- Check if workflow is enabled
- Verify main branch protection rules

#### Custom domain not working
- Verify DNS records are correct
- Wait for DNS propagation (up to 48 hours)
- Check GitHub Pages settings
- Ensure CNAME file is in dist/

#### Build failing
- Check Go version compatibility
- Verify all dependencies in go.mod
- Test build locally
- Check workflow logs for specific errors

#### 404 errors
- Ensure .nojekyll file is present
- Check file paths in generated HTML
- Verify base URL is correct

## Alternative Deployment Options

### Netlify

```bash
# Build command
go build -o site ./cmd/site && ./site build -base-url https://matthiasbrat.dev

# Publish directory
dist
```

### Vercel

Create `vercel.json`:
```json
{
  "builds": [
    {
      "src": "cmd/site/main.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "/cmd/site/main.go"
    }
  ]
}
```

### Cloudflare Pages

- Build command: `go build -o site ./cmd/site && ./site build -base-url https://matthiasbrat.dev`
- Build output directory: `dist`

### Self-Hosted

See [ARCHITECTURE.md](ARCHITECTURE.md) for Docker deployment instructions.

## Continuous Integration

The workflow also serves as CI, ensuring:
- Code compiles successfully
- Site builds without errors
- All dependencies are available
- Content is valid

For additional CI checks, consider adding:
- Go tests: `go test ./...`
- Linting: `golangci-lint run`
- Content validation
- Link checking

## Conclusion

The GitHub Pages deployment is fully automated and requires minimal maintenance. Simply push to main, and your site will be live within 1-2 minutes.

For questions or issues, check the repository's issue tracker.
