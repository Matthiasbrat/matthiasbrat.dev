package search

import (
	"fmt"

	"site/internal/build/markdown"
	"site/internal/models"
)

// DB defines the database operations needed for search indexing
// This is defined here (consumer-side) following idiomatic Go
type DB interface {
	ClearSearchIndex() error
	IndexPost(slug, collectionSlug, title, description, content, postType, url, date string) error
}

// Indexer handles search index creation
type Indexer struct {
	db DB
}

// NewIndexer creates a new search indexer
func NewIndexer(db DB) *Indexer {
	return &Indexer{db: db}
}

// IndexAll indexes all collections and posts for full-text search
func (idx *Indexer) IndexAll(collections []*models.Collection) error {
	if idx.db == nil {
		return nil
	}

	// Clear existing index
	if err := idx.db.ClearSearchIndex(); err != nil {
		return fmt.Errorf("failed to clear search index: %w", err)
	}

	// Index all posts
	for _, collection := range collections {
		postType := "blog"
		if collection.IsTopic() {
			postType = "docs"
		}

		for _, post := range collection.Posts {
			// Strip HTML tags from content for indexing
			plainContent := markdown.StripHTML(post.Content)

			// Format date
			dateStr := ""
			if !post.Date.IsZero() {
				dateStr = post.Date.Format("2006-01-02")
			}

			err := idx.db.IndexPost(
				post.Slug,
				collection.Slug,
				post.Title,
				post.Description,
				plainContent,
				postType,
				post.URL,
				dateStr,
			)
			if err != nil {
				return fmt.Errorf("failed to index post %s: %w", post.Slug, err)
			}
		}
	}

	return nil
}
