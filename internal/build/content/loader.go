package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"site/internal/build/markdown"
	"site/internal/models"

	"gopkg.in/yaml.v3"
)

// Loader handles loading and parsing markdown content
type Loader struct {
	contentDir  string
	renderer    *markdown.Renderer
	highlighter *markdown.Highlighter
}

// NewLoader creates a new content loader
func NewLoader(contentDir string) *Loader {
	return &Loader{
		contentDir:  contentDir,
		renderer:    markdown.NewRenderer(),
		highlighter: markdown.NewHighlighter(),
	}
}

// LoadAll loads all collections from the content directory
func (l *Loader) LoadAll() ([]*models.Collection, error) {
	var collections []*models.Collection

	if err := l.scanDir(l.contentDir, "", &collections); err != nil {
		return nil, err
	}

	// Build parent-child relationships for nested collections
	l.linkChildSeries(collections)

	// Set post counts for child series
	for _, c := range collections {
		c.PostCount = len(c.Posts)
	}

	// Sort collections by most recent post
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].LatestPost.After(collections[j].LatestPost)
	})

	return collections, nil
}

// linkChildSeries links child series to their parent collection
func (l *Loader) linkChildSeries(collections []*models.Collection) {
	// Create a map for quick lookup
	collectionMap := make(map[string]*models.Collection)
	for _, c := range collections {
		collectionMap[c.Slug] = c
	}

	// Find parent-child relationships based on slug hierarchy
	for _, c := range collections {
		// Check if this collection has a parent (e.g., "blog/golang-fundamentals" -> "blog")
		if idx := strings.LastIndex(c.Slug, "/"); idx > 0 {
			parentSlug := c.Slug[:idx]
			if parent, ok := collectionMap[parentSlug]; ok {
				parent.ChildSeries = append(parent.ChildSeries, c)
			}
		}
	}

	// Sort child series by latest post date
	for _, c := range collections {
		if len(c.ChildSeries) > 0 {
			sort.Slice(c.ChildSeries, func(i, j int) bool {
				return c.ChildSeries[i].LatestPost.After(c.ChildSeries[j].LatestPost)
			})
		}
	}
}

// scanDir recursively scans a directory for collections
func (l *Loader) scanDir(dir, prefix string, collections *[]*models.Collection) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()
		fullPath := filepath.Join(dir, dirName)
		slug := dirName
		if prefix != "" {
			slug = prefix + "/" + dirName
		}

		// Check if this directory is a collection (has _metadata.yml)
		metaPath := filepath.Join(fullPath, "_metadata.yml")
		if _, err := os.Stat(metaPath); err == nil {
			// This is a collection
			collection, err := l.loadCollection(fullPath, slug)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load collection %s: %v\n", slug, err)
				continue
			}
			if len(collection.Posts) > 0 {
				*collections = append(*collections, collection)
			}
		}

		// Always scan subdirectories for nested collections
		if err := l.scanDir(fullPath, slug, collections); err != nil {
			return err
		}
	}

	return nil
}

// loadCollection loads a single collection and its posts
func (l *Loader) loadCollection(collectionPath, collectionSlug string) (*models.Collection, error) {
	// Infer type from path: anything under blog/ is a series
	defaultType := models.CollectionTypeTopic
	if strings.HasPrefix(collectionSlug, "blog") {
		defaultType = models.CollectionTypeSeries
	}

	collection := &models.Collection{Slug: collectionSlug, Type: defaultType}

	// Load metadata
	metaPath := filepath.Join(collectionPath, "_metadata.yml")
	if data, err := os.ReadFile(metaPath); err == nil {
		var meta models.CollectionMeta
		if err := yaml.Unmarshal(data, &meta); err == nil {
			collection.Name = meta.Name
			collection.Description = meta.Description
			collection.Icon = meta.Icon
			collection.Banner = meta.Banner
			// Allow explicit override via metadata
			if meta.Type == "blog" {
				collection.Type = models.CollectionTypeSeries
			} else if meta.Type == "docs" {
				collection.Type = models.CollectionTypeTopic
			}
		}
	}
	if collection.Name == "" {
		collection.Name = filepath.Base(collectionSlug)
	}

	// Load posts
	postFiles, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, err
	}

	for _, postFile := range postFiles {
		if postFile.IsDir() {
			continue
		}
		if !strings.HasSuffix(postFile.Name(), ".md") {
			continue
		}
		if strings.HasPrefix(postFile.Name(), "_") {
			continue
		}

		postPath := filepath.Join(collectionPath, postFile.Name())
		post, err := l.loadPost(postPath, collectionSlug)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load %s: %v\n", postPath, err)
			continue
		}

		if post.Draft {
			continue
		}

		collection.Posts = append(collection.Posts, post)

		if post.Date.After(collection.LatestPost) {
			collection.LatestPost = post.Date
		}
	}

	// Sort posts based on collection type
	if collection.Type == models.CollectionTypeSeries {
		sort.Slice(collection.Posts, func(i, j int) bool {
			return collection.Posts[i].Date.After(collection.Posts[j].Date)
		})
	} else {
		sort.Slice(collection.Posts, func(i, j int) bool {
			if collection.Posts[i].Order != collection.Posts[j].Order {
				return collection.Posts[i].Order < collection.Posts[j].Order
			}
			return collection.Posts[i].Title < collection.Posts[j].Title
		})
	}

	// Set up prev/next navigation
	for i, post := range collection.Posts {
		if i > 0 {
			post.PrevPost = collection.Posts[i-1]
		}
		if i < len(collection.Posts)-1 {
			post.NextPost = collection.Posts[i+1]
		}
	}

	return collection, nil
}

// loadPost reads and parses a single post
func (l *Loader) loadPost(path, collectionSlug string) (*models.Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fm, content, err := markdown.ParseFrontmatter(data)
	if err != nil {
		return nil, err
	}

	// Extract TOC before processing
	toc := markdown.ExtractTOC(content)

	// Process code blocks with syntax highlighting
	content = markdown.ProcessCodeBlocks(content, l.highlighter)

	// Render markdown to HTML
	html, err := l.renderer.Render(content)
	if err != nil {
		return nil, err
	}

	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	url := "/" + collectionSlug + "/" + slug

	post := &models.Post{
		Title:          fm.Title,
		Description:    fm.Description,
		Date:           markdown.ParseDate(fm.Date),
		Updated:        markdown.ParseDate(fm.Updated),
		Draft:          fm.Draft,
		Order:          fm.Order,
		Slug:           slug,
		CollectionSlug: collectionSlug,
		TopicSlug:      collectionSlug, // backward compatibility
		URL:            url,
		Content:        html,
		TOC:            toc,
	}

	if post.Title == "" {
		post.Title = slug
	}

	return post, nil
}
