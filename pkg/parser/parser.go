package parser

import (
	"github.com/mmcdole/gofeed"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/repository"
	"github.com/tumarov/feeddy/pkg/telegram"
	"time"
)

type RSSReader struct {
	repo repository.FeedRepository
	cfg  *config.Config
	bot  *telegram.Bot
}

func NewRSSReader(repo repository.FeedRepository, cfg *config.Config, bot *telegram.Bot) *RSSReader {
	return &RSSReader{repo: repo, cfg: cfg, bot: bot}
}

func (r *RSSReader) Start() error {
	parser := gofeed.NewParser()

	for {
		users, err := r.repo.GetAll()
		if err != nil {
			return err
		}

		for _, user := range users {
			entries, err := r.fetchFeeds(user, parser)
			if err != nil {
				return err
			}
			r.bot.SendFeeds(user.ChatID, entries)
		}

		time.Sleep(time.Duration(r.cfg.ReaderTimeout) * time.Minute)
	}
}

func (r *RSSReader) fetchFeeds(user repository.UserFeed, parser *gofeed.Parser) ([]telegram.RSSEntry, error) {
	var entries []telegram.RSSEntry

	var feedsToUpdate []repository.Feed
	for _, feed := range user.Feeds {
		parsedFeed, err := parser.ParseURL(feed.URL)
		if err != nil {
			continue
		}

		for _, it := range parsedFeed.Items {
			if it.PublishedParsed.After(feed.LastRead) {
				entries = prependItem(entries, telegram.RSSEntry{
					Title:           it.Title,
					Link:            it.Link,
					Description:     it.Description,
					Published:       it.Published,
					PublishedParsed: *it.PublishedParsed,
				})
			} else {
				break
			}
		}
		if len(entries) > 0 {
			feed.LastRead = entries[len(entries)-1].PublishedParsed
		}
		feedsToUpdate = append(feedsToUpdate, feed)
	}

	if len(entries) > 0 {
		err := r.repo.Save(repository.UserFeed{ChatID: user.ChatID, Feeds: feedsToUpdate})
		if err != nil {
			return nil, err
		}
	}

	return entries, nil
}

func prependItem(items []telegram.RSSEntry, item telegram.RSSEntry) []telegram.RSSEntry {
	return append([]telegram.RSSEntry{item}, items...)
}
