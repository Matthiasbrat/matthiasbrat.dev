package models

import "time"

type Post struct {
	Title       string
	Description string
	Date        time.Time
	Updated     time.Time
	Draft       bool
	Order       int

	Slug           string
	CollectionSlug string
	TopicSlug      string
	URL            string
	Content        string
	TOC            []TOCItem
	RawContent     string
	OGImage        string

	PrevPost *Post
	NextPost *Post
}

type TOCItem struct {
	Level int
	ID    string
	Text  string
}

type PostFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Date        string `yaml:"date"`
	Updated     string `yaml:"updated"`
	Draft       bool   `yaml:"draft"`
	Order       int    `yaml:"order"`
}
