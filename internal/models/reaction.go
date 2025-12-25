package models

import (
	"slices"
	"time"
)

type User struct {
	ID        string
	Email     string
	Name      string
	AvatarURL string
	CreatedAt time.Time
}

type Reaction struct {
	ID        int64
	UserID    string
	PostSlug  string
	Emoji     string
	CreatedAt time.Time
}

type ReactionCount struct {
	Emoji string   `json:"emoji"`
	Count int      `json:"count"`
	Users []string `json:"users,omitempty"`
}

var AllowedEmojis = []string{"👍", "❤️", "😂", "💡", "😢"}

func IsValidEmoji(emoji string) bool {
	return slices.Contains(AllowedEmojis, emoji)
}

type Comment struct {
	ID        int64
	UserID    string
	PostSlug  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentWithUser struct {
	Comment
	UserName   string
	UserAvatar string
}
