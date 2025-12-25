package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	"site/internal/models"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func New(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.init(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) init() error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT,
			avatar_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL REFERENCES users(id),
			post_slug TEXT NOT NULL,
			emoji TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, post_slug, emoji)
		);

		CREATE INDEX IF NOT EXISTS idx_reactions_post ON reactions(post_slug);

		CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id),
			expires_at DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);

		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL REFERENCES users(id),
			post_slug TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_slug);
		CREATE INDEX IF NOT EXISTS idx_comments_created ON comments(post_slug, created_at DESC);

		CREATE VIRTUAL TABLE IF NOT EXISTS search_index USING fts5(
			slug UNINDEXED,
			collection_slug UNINDEXED,
			title,
			description,
			content,
			type UNINDEXED,
			url UNINDEXED,
			date UNINDEXED
		);
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) CreateOrUpdateUser(user *models.User) error {
	_, err := db.conn.Exec(`
		INSERT INTO users (id, email, name, avatar_url, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			email = excluded.email,
			name = excluded.name,
			avatar_url = excluded.avatar_url
	`, user.ID, user.Email, user.Name, user.AvatarURL, time.Now())
	return err
}

func (db *DB) GetUser(id string) (*models.User, error) {
	var user models.User
	var avatarURL sql.NullString
	err := db.conn.QueryRow(`
		SELECT id, email, name, avatar_url, created_at FROM users WHERE id = ?
	`, id).Scan(&user.ID, &user.Email, &user.Name, &avatarURL, &user.CreatedAt)
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) CreateSession(token, userID string, expiresAt time.Time) error {
	_, err := db.conn.Exec(`
		INSERT INTO sessions (token, user_id, expires_at)
		VALUES (?, ?, ?)
	`, token, userID, expiresAt)
	return err
}

func (db *DB) GetSession(token string) (string, error) {
	var userID string
	var expiresAt time.Time
	err := db.conn.QueryRow(`
		SELECT user_id, expires_at FROM sessions WHERE token = ?
	`, token).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if time.Now().After(expiresAt) {
		db.DeleteSession(token)
		return "", nil
	}
	return userID, nil
}

func (db *DB) DeleteSession(token string) error {
	_, err := db.conn.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}

func (db *DB) CleanExpiredSessions() error {
	_, err := db.conn.Exec(`DELETE FROM sessions WHERE expires_at < ?`, time.Now())
	return err
}

func (db *DB) AddReaction(userID, postSlug, emoji string) (bool, error) {
	var exists bool
	err := db.conn.QueryRow(`
		SELECT 1 FROM reactions WHERE user_id = ? AND post_slug = ? AND emoji = ?
	`, userID, postSlug, emoji).Scan(&exists)

	if err == sql.ErrNoRows {
		_, err := db.conn.Exec(`
			INSERT INTO reactions (user_id, post_slug, emoji, created_at)
			VALUES (?, ?, ?, ?)
		`, userID, postSlug, emoji, time.Now())
		return true, err
	}
	if err != nil {
		return false, err
	}

	_, err = db.conn.Exec(`
		DELETE FROM reactions WHERE user_id = ? AND post_slug = ? AND emoji = ?
	`, userID, postSlug, emoji)
	return false, err
}

func (db *DB) GetReactionCounts(postSlug string) ([]models.ReactionCount, error) {
	rows, err := db.conn.Query(`
		SELECT emoji, COUNT(*) as count
		FROM reactions
		WHERE post_slug = ?
		GROUP BY emoji
	`, postSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []models.ReactionCount
	for rows.Next() {
		var rc models.ReactionCount
		if err := rows.Scan(&rc.Emoji, &rc.Count); err != nil {
			return nil, err
		}
		counts = append(counts, rc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range counts {
		users, err := db.getReactionUsers(postSlug, counts[i].Emoji, 3)
		if err != nil {
			return nil, err
		}
		counts[i].Users = users
	}

	return counts, nil
}

func (db *DB) getReactionUsers(postSlug, emoji string, limit int) ([]string, error) {
	rows, err := db.conn.Query(`
		SELECT u.name
		FROM reactions r
		JOIN users u ON r.user_id = u.id
		WHERE r.post_slug = ? AND r.emoji = ?
		ORDER BY r.created_at DESC
		LIMIT ?
	`, postSlug, emoji, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, rows.Err()
}

func (db *DB) GetUserReactions(userID, postSlug string) ([]string, error) {
	rows, err := db.conn.Query(`
		SELECT emoji FROM reactions WHERE user_id = ? AND post_slug = ?
	`, userID, postSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emojis []string
	for rows.Next() {
		var emoji string
		if err := rows.Scan(&emoji); err != nil {
			return nil, err
		}
		emojis = append(emojis, emoji)
	}

	return emojis, rows.Err()
}

type SearchResult struct {
	Slug           string
	CollectionSlug string
	Title          string
	Description    string
	Snippet        string
	Type           string
	URL            string
	Date           string
}

func (db *DB) ClearSearchIndex() error {
	_, err := db.conn.Exec(`DELETE FROM search_index`)
	return err
}

func (db *DB) IndexPost(slug, collectionSlug, title, description, content, postType, url, date string) error {
	_, err := db.conn.Exec(`
		INSERT INTO search_index (slug, collection_slug, title, description, content, type, url, date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, slug, collectionSlug, title, description, content, postType, url, date)
	return err
}

func (db *DB) Search(query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []SearchResult{}, nil
	}

	ftsQuery := buildFuzzyQuery(query)

	rows, err := db.conn.Query(`
		SELECT
			slug,
			collection_slug,
			title,
			description,
			snippet(search_index, 4, '<mark>', '</mark>', '...', 32) as snippet,
			type,
			url,
			date
		FROM search_index
		WHERE search_index MATCH ?
		ORDER BY rank
		LIMIT ?
	`, ftsQuery, limit)
	if err != nil {
		return []SearchResult{}, nil
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		if err := rows.Scan(&r.Slug, &r.CollectionSlug, &r.Title, &r.Description, &r.Snippet, &r.Type, &r.URL, &r.Date); err != nil {
			continue
		}
		results = append(results, r)
	}

	return results, nil
}

func buildFuzzyQuery(query string) string {
	tokens := strings.Fields(query)
	if len(tokens) == 0 {
		return ""
	}

	if len(tokens) == 1 {
		return escapeFTSToken(tokens[0]) + "*"
	}

	var parts []string
	for _, token := range tokens {
		parts = append(parts, escapeFTSToken(token)+"*")
	}
	return strings.Join(parts, " OR ")
}

func escapeFTSToken(token string) string {
	token = strings.ReplaceAll(token, `"`, `""`)
	if strings.ContainsAny(token, ` :"*`) {
		return `"` + token + `"`
	}
	return token
}

// Comment methods

func (db *DB) CreateComment(userID, postSlug, content string) (*models.Comment, error) {
	now := time.Now()
	result, err := db.conn.Exec(`
		INSERT INTO comments (user_id, post_slug, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, userID, postSlug, content, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Comment{
		ID:        id,
		UserID:    userID,
		PostSlug:  postSlug,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (db *DB) GetComments(postSlug string) ([]models.CommentWithUser, error) {
	rows, err := db.conn.Query(`
		SELECT c.id, c.user_id, c.post_slug, c.content, c.created_at, c.updated_at,
		       u.name, COALESCE(u.avatar_url, '')
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_slug = ?
		ORDER BY c.created_at DESC
	`, postSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.CommentWithUser
	for rows.Next() {
		var c models.CommentWithUser
		if err := rows.Scan(&c.ID, &c.UserID, &c.PostSlug, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.UserName, &c.UserAvatar); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, rows.Err()
}

func (db *DB) GetComment(id int64) (*models.Comment, error) {
	var c models.Comment
	err := db.conn.QueryRow(`
		SELECT id, user_id, post_slug, content, created_at, updated_at
		FROM comments WHERE id = ?
	`, id).Scan(&c.ID, &c.UserID, &c.PostSlug, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *DB) UpdateComment(id int64, userID, content string) error {
	result, err := db.conn.Exec(`
		UPDATE comments SET content = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`, content, time.Now(), id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (db *DB) DeleteComment(id int64, userID string) error {
	result, err := db.conn.Exec(`
		DELETE FROM comments WHERE id = ? AND user_id = ?
	`, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CleanupUserData removes all data for a user (comments, reactions, sessions)
func (db *DB) CleanupUserData(userID string) error {
	_, err := db.conn.Exec(`DELETE FROM comments WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}
	_, err = db.conn.Exec(`DELETE FROM reactions WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}
	_, err = db.conn.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}
	return nil
}
