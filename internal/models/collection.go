package models

import "time"

type CollectionType string

const (
	CollectionTypeSeries CollectionType = "series"
	CollectionTypeTopic  CollectionType = "topic"
)

type Collection struct {
	Slug        string
	Name        string
	Description string
	Type        CollectionType
	Icon        string
	Banner      string
	Posts       []*Post
	LatestPost  time.Time
	ChildSeries []*Collection // Nested series under this collection
	PostCount   int           // Total posts in this series
}

type CollectionMeta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Icon        string `yaml:"icon"`
	Banner      string `yaml:"banner"`
}

// IsSeries returns true for blog-style collections (sorted by date)
func (c *Collection) IsSeries() bool {
	return c.Type == CollectionTypeSeries
}

// IsTopic returns true for docs-style collections (sorted by order)
func (c *Collection) IsTopic() bool {
	return c.Type == CollectionTypeTopic || c.Type == ""
}

// IsBlog is an alias for IsSeries (for template readability)
func (c *Collection) IsBlog() bool {
	return c.IsSeries()
}

// IsDocs is an alias for IsTopic (for template readability)
func (c *Collection) IsDocs() bool {
	return c.IsTopic()
}

// IsMainBlog returns true if this is the main blog collection (slug == "blog")
func (c *Collection) IsMainBlog() bool {
	return c.Slug == "blog"
}
