package repository

type Feed struct {
	ChatID int64
	Feeds  []string
}

const FeedsCollection = "feeds"

type FeedRepository interface {
	Create(chatID int64) error
	Add(feed Feed, feedURL string) error
	Get(chatID int64) (*Feed, error)
	Remove(feed Feed, feedURLToRemove string) error
}
