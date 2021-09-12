package repository

import "time"

type UserFeed struct {
	ChatID int64
	Feeds  []Feed
}

type Feed struct {
	URL      string
	LastRead time.Time
}

const FeedsCollection = "feeds"

type FeedRepository interface {
	Create(chatID int64) error
	Add(userFeed UserFeed, feedURL string) error
	Get(chatID int64) (*UserFeed, error)
	Remove(userFeed UserFeed, feedURLToRemove string) error
	GetAll() ([]UserFeed, error)
	Save(userFeed UserFeed) error
}
