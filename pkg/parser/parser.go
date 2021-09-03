package parser

import (
	"github.com/mmcdole/gofeed"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/repository"
	"github.com/tumarov/feeddy/pkg/telegram"
	"time"
)

const timeout = 5

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
			entries := fetchFeeds(user, parser)
			r.bot.SendFeeds(user.ChatID, entries)
		}

		time.Sleep(timeout * time.Minute)
	}
}

func fetchFeeds(user repository.Feed, parser *gofeed.Parser) []telegram.RSSEntry {
	var entries []telegram.RSSEntry

	for _, feedURL := range user.Feeds {
		parsedFeed, err := parser.ParseURL(feedURL)
		if err != nil {
			continue
		}

		for _, it := range parsedFeed.Items {
			if it.PublishedParsed.After(time.Now().Add(-timeout * time.Minute)) {
				entries = append(entries, telegram.RSSEntry{
					Title:       it.Title,
					Link:        it.Link,
					Description: it.Description,
					Published:   it.Published,
				})
			} else {
				break
			}
		}
	}

	return entries
}
