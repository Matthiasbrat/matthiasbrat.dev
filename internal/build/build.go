package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"site/internal/build/assets"
	"site/internal/build/content"
	"site/internal/build/markdown"
	"site/internal/build/og"
	"site/internal/build/search"
	"site/internal/models"

	"gopkg.in/yaml.v3"
)

const postsPerPage = 10

// Site holds all site data
type Site struct {
	Config      Config
	Collections []*models.Collection
	AssetHashes assets.Hashes // Maps original filename to hashed filename
}

// Build generates the static site
func Build(cfg Config) error {
	// Load site config
	if cfg.SiteName == "" {
		cfg.SiteName = "Site"
	}

	// Try to load site.yml for config
	siteConfigPath := "site.yml"
	if data, err := os.ReadFile(siteConfigPath); err == nil {
		var siteCfg struct {
			Title              string             `yaml:"title"`
			Description        string             `yaml:"description"`
			BaseURL            string             `yaml:"base_url"`
			DefaultSocialImage string             `yaml:"default_social_image"`
			Profile            ProfileConfig      `yaml:"profile"`
			Referrals          []ReferralConfig   `yaml:"referrals"`
		}
		if err := yaml.Unmarshal(data, &siteCfg); err == nil {
			if siteCfg.Title != "" {
				cfg.SiteName = siteCfg.Title
			}
			if siteCfg.Description != "" {
				cfg.SiteDesc = siteCfg.Description
			}
			if siteCfg.BaseURL != "" && cfg.BaseURL == "" {
				cfg.BaseURL = siteCfg.BaseURL
			}
			if siteCfg.DefaultSocialImage != "" {
				cfg.DefaultSocialImage = siteCfg.DefaultSocialImage
			}
			cfg.Profile = siteCfg.Profile
			cfg.Referrals = siteCfg.Referrals
		}
	}

	site := &Site{
		Config: cfg,
	}

	// Clean output directory
	if err := os.RemoveAll(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to clean output directory: %w", err)
	}

	// Load content
	if err := site.loadContent(); err != nil {
		return fmt.Errorf("failed to load content: %w", err)
	}

	// Index content for search (if database provided)
	if cfg.DB != nil {
		indexer := search.NewIndexer(cfg.DB)
		if err := indexer.IndexAll(site.Collections); err != nil {
			return fmt.Errorf("failed to index content: %w", err)
		}
	}

	// Copy static files first to generate hashes
	processor := assets.NewProcessor(cfg.StaticDir, cfg.OutputDir, cfg.DevMode)
	assetHashes, err := processor.ProcessAll()
	if err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}
	site.AssetHashes = assetHashes

	// Load templates (after copying static, so asset hashes are available)
	tmpl, err := site.loadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Generate OG images concurrently
	if err := site.generateOGImages(); err != nil {
		return fmt.Errorf("failed to generate OG images: %w", err)
	}

	// Generate pages
	if err := site.generatePages(tmpl); err != nil {
		return fmt.Errorf("failed to generate pages: %w", err)
	}

	// Generate sitemap
	if err := site.generateSitemap(); err != nil {
		return fmt.Errorf("failed to generate sitemap: %w", err)
	}

	return nil
}

// loadContent reads and parses all content files
func (s *Site) loadContent() error {
	loader := content.NewLoader(s.Config.ContentDir)
	collections, err := loader.LoadAll()
	if err != nil {
		return err
	}
	s.Collections = collections
	return nil
}

func (s *Site) generateOGImages() error {
	profilePhotoPath := ""
	if s.Config.Profile.Photo != "" {
		photoURL := s.Config.Profile.Photo
		if strings.HasPrefix(photoURL, "/") {
			profilePhotoPath = filepath.Join(s.Config.StaticDir, photoURL)
		}
	}

	fontPath := filepath.Join(s.Config.StaticDir, "fonts", "SourceSerif4-Regular.ttf")
	fontBoldPath := filepath.Join(s.Config.StaticDir, "fonts", "SourceSerif4-Semibold.ttf")

	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return nil
	}
	if _, err := os.Stat(fontBoldPath); os.IsNotExist(err) {
		return nil
	}

	generator, err := og.NewGenerator(
		s.Config.OutputDir,
		profilePhotoPath,
		fontPath,
		fontBoldPath,
		s.Config.SiteName,
		s.Config.BaseURL,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU())
	errChan := make(chan error, 1)

	for _, collection := range s.Collections {
		for _, post := range collection.Posts {
			wg.Add(1)
			go func(p *models.Post, c *models.Collection) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				ogPath, err := generator.Generate(p, c)
				if err != nil {
					select {
					case errChan <- err:
					default:
					}
					return
				}
				p.OGImage = s.Config.BaseURL + ogPath
			}(post, collection)
		}
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// TemplateSet holds parsed templates for each page type
type TemplateSet struct {
	Home       *template.Template
	Docs       *template.Template
	Blog       *template.Template
	Collection *template.Template
	Post       *template.Template
	Profile    *template.Template
	Referrals  *template.Template
}

// loadTemplates loads all HTML templates
func (s *Site) loadTemplates() (*TemplateSet, error) {
	// Read critical CSS for inlining
	criticalCSSPath := filepath.Join(s.Config.StaticDir, "css", "critical.css")
	criticalCSS := ""
	if data, err := os.ReadFile(criticalCSSPath); err == nil {
		criticalCSS = string(data)
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeCSS": func(s string) template.CSS {
			return template.CSS(s)
		},
		"criticalCSS": func() template.CSS {
			return template.CSS(criticalCSS)
		},
		"asset": func(path string) string {
			// Normalize path to use forward slashes
			normalizedPath := filepath.ToSlash(path)

			// Return hashed path if available, otherwise original
			if hashed, ok := s.AssetHashes[normalizedPath]; ok {
				return "/" + hashed
			}
			return "/" + normalizedPath
		},
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
	}

	basePath := filepath.Join(s.Config.TemplateDir, "base.html")
	baseContent, err := os.ReadFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read base.html: %w", err)
	}

	// Load all partials
	partialFiles, err := filepath.Glob(filepath.Join(s.Config.TemplateDir, "partials", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob partials: %w", err)
	}

	// Helper to create a template with base + partials + page
	parseWithBase := func(pageName string) (*template.Template, error) {
		pagePath := filepath.Join(s.Config.TemplateDir, pageName)
		pageContent, err := os.ReadFile(pagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", pageName, err)
		}

		tmpl := template.New("base.html").Funcs(funcMap)

		// Parse base template
		if _, err := tmpl.Parse(string(baseContent)); err != nil {
			return nil, fmt.Errorf("failed to parse base.html: %w", err)
		}

		// Parse all partials
		for _, partialPath := range partialFiles {
			partialContent, err := os.ReadFile(partialPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read partial %s: %w", partialPath, err)
			}
			partialName := filepath.Base(partialPath)
			if _, err := tmpl.New(partialName).Parse(string(partialContent)); err != nil {
				return nil, fmt.Errorf("failed to parse partial %s: %w", partialName, err)
			}
		}

		// Parse page template
		if _, err := tmpl.New(pageName).Parse(string(pageContent)); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", pageName, err)
		}

		return tmpl, nil
	}

	ts := &TemplateSet{}

	ts.Home, err = parseWithBase("home.html")
	if err != nil {
		return nil, err
	}

	ts.Docs, err = parseWithBase("docs.html")
	if err != nil {
		return nil, err
	}

	ts.Blog, err = parseWithBase("blog.html")
	if err != nil {
		// Fallback to collection.html
		ts.Blog, _ = parseWithBase("collection.html")
	}

	ts.Collection, err = parseWithBase("collection.html")
	if err != nil {
		// Fallback to topic.html for backward compatibility
		ts.Collection, err = parseWithBase("topic.html")
		if err != nil {
			return nil, err
		}
	}

	ts.Post, err = parseWithBase("post.html")
	if err != nil {
		return nil, err
	}

	ts.Profile, err = parseWithBase("profile.html")
	if err != nil {
		return nil, err
	}

	ts.Referrals, err = parseWithBase("referrals.html")
	if err != nil {
		// Referrals page is optional
		ts.Referrals = nil
	}

	return ts, nil
}

// generatePages creates all HTML pages
func (s *Site) generatePages(ts *TemplateSet) error {
	// Generate homepage
	if err := s.generateHome(ts.Home); err != nil {
		return err
	}

	// Generate profile page
	if err := s.generateProfile(ts.Profile); err != nil {
		return err
	}

	// Generate referrals page
	if ts.Referrals != nil && len(s.Config.Referrals) > 0 {
		if err := s.generateReferrals(ts.Referrals); err != nil {
			return err
		}
	}

	// Generate docs landing page
	if err := s.generateDocsLanding(ts.Docs); err != nil {
		return err
	}

	// Generate collection pages
	for _, collection := range s.Collections {
		// Use blog template for main blog with pagination
		if collection.IsMainBlog() && ts.Blog != nil {
			if err := s.generateBlogLanding(ts.Blog, collection); err != nil {
				return err
			}
		} else {
			if err := s.generateCollection(ts.Collection, collection); err != nil {
				return err
			}
		}

		// Generate post pages
		for _, post := range collection.Posts {
			if err := s.generatePost(ts.Post, collection, post); err != nil {
				return err
			}
		}
	}

	return nil
}

// PageData holds common data for all pages
type PageData struct {
	Title           string
	Description     string
	CanonicalURL    string
	OGType          string
	OGImage         string
	DatePublished   string
	DateModified    string
	SiteName        string
	SiteDescription string
	Collections     []*models.Collection
	Year            int
	StructuredData  template.JS
	DevMode         bool
	User            *models.User
}

// getSocialImage determines the appropriate Open Graph image for a page
func (s *Site) getSocialImage(collection *models.Collection) string {
	// For blog posts/collections, prefer the banner image
	if collection != nil && collection.Banner != "" {
		return s.Config.BaseURL + "/images/" + collection.Banner
	}

	// For docs collections with icons, use the icon if it's an image
	if collection != nil && collection.Icon != "" {
		if strings.HasPrefix(collection.Icon, "/") || strings.HasPrefix(collection.Icon, "http") {
			return collection.Icon
		}
	}

	// Use profile photo if available
	if s.Config.Profile.Photo != "" {
		return s.Config.Profile.Photo
	}

	// Fall back to default social image
	if s.Config.DefaultSocialImage != "" {
		return s.Config.DefaultSocialImage
	}

	return ""
}

// generateHome creates the homepage
func (s *Site) generateHome(tmpl *template.Template) error {
	// Collect all posts from blog series and sort by date
	var allPosts []*models.Post
	for _, collection := range s.Collections {
		if collection.IsSeries() {
			allPosts = append(allPosts, collection.Posts...)
		}
	}
	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].Date.After(allPosts[j].Date)
	})

	// Take only the 5 most recent
	latestPosts := allPosts
	if len(latestPosts) > 5 {
		latestPosts = latestPosts[:5]
	}

	// Get docs collections (topics) limited to 5
	var docsCollections []*models.Collection
	for _, collection := range s.Collections {
		if collection.Type == models.CollectionTypeTopic {
			docsCollections = append(docsCollections, collection)
		}
	}
	if len(docsCollections) > 5 {
		docsCollections = docsCollections[:5]
	}

	data := struct {
		PageData
		Profile         ProfileConfig
		Referrals       []ReferralConfig
		LatestPosts     []*models.Post
		DocsCollections []*models.Collection
	}{
		PageData: PageData{
			Title:           s.Config.SiteName,
			Description:     s.Config.SiteDesc,
			CanonicalURL:    s.Config.BaseURL,
			OGType:          "website",
			OGImage:         s.getSocialImage(nil),
			SiteName:        s.Config.SiteName,
			SiteDescription: s.Config.SiteDesc,
			Collections:     s.Collections,
			Year:            time.Now().Year(),
			DevMode:         s.Config.DevMode,
		},
		Profile:         s.Config.Profile,
		Referrals:       s.Config.Referrals,
		LatestPosts:     latestPosts,
		DocsCollections: docsCollections,
	}

	return s.renderPage(tmpl, "index.html", data)
}

// generateProfile creates the /profile page
func (s *Site) generateProfile(tmpl *template.Template) error {
	// Read profile markdown file
	profilePath := filepath.Join(s.Config.ContentDir, "profile.md")
	profileData, err := os.ReadFile(profilePath)
	if err != nil {
		// Profile page is optional, skip if not found
		return nil
	}

	// Parse markdown
	fm, content, err := markdown.ParseFrontmatter(profileData)
	if err != nil {
		return fmt.Errorf("failed to parse profile.md: %w", err)
	}

	// Create renderer and highlighter
	renderer := markdown.NewRenderer()
	highlighter := markdown.NewHighlighter()

	// Process code blocks with syntax highlighting
	content = markdown.ProcessCodeBlocks(content, highlighter)

	// Render markdown to HTML
	html, err := renderer.Render(content)
	if err != nil {
		return fmt.Errorf("failed to render profile markdown: %w", err)
	}

	// Use title from frontmatter or default
	title := fm.Title
	if title == "" {
		title = "About"
	}

	data := struct {
		PageData
		Profile ProfileConfig
		Content template.HTML
	}{
		PageData: PageData{
			Title:           title + " | " + s.Config.SiteName,
			Description:     fm.Description,
			CanonicalURL:    s.Config.BaseURL + "/profile",
			OGType:          "profile",
			OGImage:         s.getSocialImage(nil),
			SiteName:        s.Config.SiteName,
			SiteDescription: s.Config.SiteDesc,
			Collections:     s.Collections,
			Year:            time.Now().Year(),
			DevMode:         s.Config.DevMode,
		},
		Profile: s.Config.Profile,
		Content: template.HTML(html),
	}

	return s.renderPage(tmpl, "profile.html", data)
}

// generateReferrals creates the /referrals page
func (s *Site) generateReferrals(tmpl *template.Template) error {
	data := struct {
		PageData
		Referrals []ReferralConfig
	}{
		PageData: PageData{
			Title:        "Referrals | " + s.Config.SiteName,
			Description:  "People I recommend and work with",
			CanonicalURL: s.Config.BaseURL + "/referrals",
			OGType:       "website",
			OGImage:      s.getSocialImage(nil),
			SiteName:     s.Config.SiteName,
			Collections:  s.Collections,
			Year:         time.Now().Year(),
			DevMode:      s.Config.DevMode,
		},
		Referrals: s.Config.Referrals,
	}

	return s.renderPage(tmpl, filepath.Join("referrals", "index.html"), data)
}

// generateDocsLanding creates the /docs landing page
func (s *Site) generateDocsLanding(tmpl *template.Template) error {
	// Filter for docs collections (topics) only
	var docsCollections []*models.Collection
	for _, collection := range s.Collections {
		if collection.Type == models.CollectionTypeTopic {
			docsCollections = append(docsCollections, collection)
		}
	}

	data := struct {
		PageData
		DocsCollections []*models.Collection
	}{
		PageData: PageData{
			Title:        "Documentation | " + s.Config.SiteName,
			Description:  "Browse all documentation and guides",
			CanonicalURL: s.Config.BaseURL + "/docs",
			OGType:       "website",
			OGImage:      s.getSocialImage(nil),
			SiteName:     s.Config.SiteName,
			Collections:  s.Collections,
			Year:         time.Now().Year(),
			DevMode:      s.Config.DevMode,
		},
		DocsCollections: docsCollections,
	}

	return s.renderPage(tmpl, filepath.Join("docs", "index.html"), data)
}

// generateBlogLanding creates paginated /blog landing pages
func (s *Site) generateBlogLanding(tmpl *template.Template, collection *models.Collection) error {
	// Get all posts from the main blog collection
	allPosts := make([]*models.Post, len(collection.Posts))
	copy(allPosts, collection.Posts)

	// Sort by date (newest first)
	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].Date.After(allPosts[j].Date)
	})

	totalPosts := len(allPosts)
	totalPages := (totalPosts + postsPerPage - 1) / postsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	// Generate a page for each
	for page := 1; page <= totalPages; page++ {
		start := (page - 1) * postsPerPage
		end := start + postsPerPage
		if end > totalPosts {
			end = totalPosts
		}

		var pagePosts []*models.Post
		if start < totalPosts {
			pagePosts = allPosts[start:end]
		}

		pageNumbers := buildPageNumbers(page, totalPages)

		data := struct {
			PageData
			Collection  *models.Collection
			Posts       []*models.Post
			CurrentPage int
			TotalPages  int
			PrevPage    int
			NextPage    int
			PageNumbers []int
			SortOrder   string
			SortQuery   string
		}{
			PageData: PageData{
				Title:        collection.Name + " | " + s.Config.SiteName,
				Description:  collection.Description,
				CanonicalURL: s.Config.BaseURL + "/blog",
				OGType:       "website",
				OGImage:      s.getSocialImage(collection),
				SiteName:     s.Config.SiteName,
				Collections:  s.Collections,
				Year:         time.Now().Year(),
				DevMode:      s.Config.DevMode,
			},
			Collection:  collection,
			Posts:       pagePosts,
			CurrentPage: page,
			TotalPages:  totalPages,
			PrevPage:    page - 1,
			NextPage:    page + 1,
			PageNumbers: pageNumbers,
			SortOrder:   "newest",
			SortQuery:   "",
		}

		var outPath string
		if page == 1 {
			outPath = filepath.Join("blog", "index.html")
		} else {
			outPath = filepath.Join("blog", "page", fmt.Sprintf("%d", page), "index.html")
		}

		if err := s.renderPage(tmpl, outPath, data); err != nil {
			return err
		}
	}

	return nil
}

// buildPageNumbers creates page numbers for pagination (returns -1 for ellipsis)
func buildPageNumbers(current, total int) []int {
	if total <= 7 {
		pages := make([]int, total)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}

	var pages []int
	pages = append(pages, 1)

	if current > 3 {
		pages = append(pages, -1)
	}

	for i := current - 1; i <= current+1; i++ {
		if i > 1 && i < total {
			pages = append(pages, i)
		}
	}

	if current < total-2 {
		pages = append(pages, -1)
	}

	if total > 1 {
		pages = append(pages, total)
	}

	return pages
}

// generateCollection creates a collection listing page
func (s *Site) generateCollection(tmpl *template.Template, collection *models.Collection) error {
	data := struct {
		PageData
		Collection *models.Collection
	}{
		PageData: PageData{
			Title:        collection.Name + " | " + s.Config.SiteName,
			Description:  collection.Description,
			CanonicalURL: s.Config.BaseURL + "/" + collection.Slug,
			OGType:       "website",
			OGImage:      s.getSocialImage(collection),
			SiteName:     s.Config.SiteName,
			Collections:  s.Collections,
			Year:         time.Now().Year(),
			DevMode:      s.Config.DevMode,
		},
		Collection: collection,
	}

	outPath := filepath.Join(collection.Slug, "index.html")
	return s.renderPage(tmpl, outPath, data)
}

// generatePost creates a post page
func (s *Site) generatePost(tmpl *template.Template, collection *models.Collection, post *models.Post) error {
	// Build structured data
	structuredData := map[string]interface{}{
		"@context":      "https://schema.org",
		"@type":         "Article",
		"headline":      post.Title,
		"description":   post.Description,
		"datePublished": post.Date.Format(time.RFC3339),
	}
	if !post.Updated.IsZero() {
		structuredData["dateModified"] = post.Updated.Format(time.RFC3339)
	}

	sdJSON, _ := json.MarshalIndent(structuredData, "", "  ")

	ogImage := post.OGImage
	if ogImage == "" {
		ogImage = s.getSocialImage(collection)
	}

	data := struct {
		PageData
		Collection *models.Collection
		Post       *models.Post
		Emojis     []string
	}{
		PageData: PageData{
			Title:          post.Title + " | " + collection.Name + " | " + s.Config.SiteName,
			Description:    post.Description,
			CanonicalURL:   s.Config.BaseURL + post.URL,
			OGType:         "article",
			OGImage:        ogImage,
			DatePublished:  post.Date.Format(time.RFC3339),
			SiteName:       s.Config.SiteName,
			Collections:    s.Collections,
			Year:           time.Now().Year(),
			StructuredData: template.JS(sdJSON),
			DevMode:        s.Config.DevMode,
		},
		Collection: collection,
		Post:       post,
		Emojis:     models.AllowedEmojis,
	}

	if !post.Updated.IsZero() {
		data.PageData.DateModified = post.Updated.Format(time.RFC3339)
	}

	outPath := filepath.Join(post.TopicSlug, post.Slug, "index.html")
	return s.renderPage(tmpl, outPath, data)
}

// renderPage renders a template to a file
func (s *Site) renderPage(tmpl *template.Template, outPath string, data interface{}) error {
	var buf bytes.Buffer

	// Execute the base template (which includes the content block)
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write output file
	fullPath := filepath.Join(s.Config.OutputDir, outPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, buf.Bytes(), 0644)
}

// generateSitemap creates sitemap.xml
func (s *Site) generateSitemap() error {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString("\n")
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	buf.WriteString("\n")

	// Homepage
	buf.WriteString(fmt.Sprintf("  <url><loc>%s</loc></url>\n", s.Config.BaseURL))

	// Collections and posts
	for _, collection := range s.Collections {
		buf.WriteString(fmt.Sprintf("  <url><loc>%s/%s</loc></url>\n", s.Config.BaseURL, collection.Slug))

		for _, post := range collection.Posts {
			lastmod := post.Date
			if !post.Updated.IsZero() {
				lastmod = post.Updated
			}
			buf.WriteString(fmt.Sprintf("  <url><loc>%s%s</loc><lastmod>%s</lastmod></url>\n",
				s.Config.BaseURL, post.URL, lastmod.Format("2006-01-02")))
		}
	}

	buf.WriteString("</urlset>\n")

	sitemapPath := filepath.Join(s.Config.OutputDir, "sitemap.xml")
	return os.WriteFile(sitemapPath, buf.Bytes(), 0644)
}
