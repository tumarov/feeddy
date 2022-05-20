package boltdb

import (
	"encoding/json"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/repository"
	bolt "go.etcd.io/bbolt"
	"strconv"
	"time"
)

const FeedsBucket = "Feeds"

type FeedRepository struct {
	db  *bolt.DB
	cfg *config.Config
}

func NewFeedRepository(db *bolt.DB, cfg *config.Config) (*FeedRepository, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(FeedsBucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &FeedRepository{db: db, cfg: cfg}, nil
}

func (r *FeedRepository) Create(chatID int64) error {
	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FeedsBucket))
		feeds, err := json.Marshal([]repository.Feed{})
		if err != nil {
			return err
		}
		err = b.Put(intToBytes(chatID), feeds)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (r *FeedRepository) Save(userFeed repository.UserFeed) error {
	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FeedsBucket))
		feeds, err := json.Marshal(userFeed.Feeds)
		if err != nil {
			return err
		}
		err = b.Put(intToBytes(userFeed.ChatID), feeds)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (r *FeedRepository) Add(userFeed repository.UserFeed, feedURL string) error {
	for _, feed := range userFeed.Feeds {
		if feed.URL == feedURL {
			return nil
		}
	}

	feeds := append(userFeed.Feeds, repository.Feed{URL: feedURL, LastRead: time.Now()})

	if err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FeedsBucket))
		feeds, err := json.Marshal(feeds)
		if err != nil {
			return err
		}
		err = b.Put(intToBytes(userFeed.ChatID), feeds)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (r *FeedRepository) Get(chatID int64) (*repository.UserFeed, error) {
	var found bool
	result := &repository.UserFeed{}

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FeedsBucket))
		v := b.Get(intToBytes(chatID))
		if v == nil {
			return nil
		}

		found = true

		err := json.Unmarshal(v, &result.Feeds)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	result.ChatID = chatID

	return result, nil
}

func (r *FeedRepository) Remove(userFeed repository.UserFeed, feedURLToRemove string) error {
	for i, feed := range userFeed.Feeds {
		if feed.URL == feedURLToRemove {
			newFeeds := append(userFeed.Feeds[:i], userFeed.Feeds[i+1:]...)
			if err := r.db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(FeedsBucket))
				feeds, err := json.Marshal(newFeeds)
				if err != nil {
					return err
				}
				err = b.Put(intToBytes(userFeed.ChatID), feeds)
				return err
			}); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (r *FeedRepository) GetAll() ([]repository.UserFeed, error) {
	var results []repository.UserFeed

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FeedsBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			chatID, err := strconv.ParseInt(string(k), 10, 64)
			if err != nil {
				return err
			}

			result := repository.UserFeed{}
			result.ChatID = chatID

			err = json.Unmarshal(v, &result.Feeds)
			if err != nil {
				return err
			}

			results = append(results, result)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
