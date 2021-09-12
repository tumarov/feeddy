package mongodb

import (
	"context"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type FeedRepository struct {
	db         *mongo.Client
	cfg        *config.Config
	collection *mongo.Collection
}

func NewFeedRepository(db *mongo.Client, cfg *config.Config) *FeedRepository {
	collection := db.Database(cfg.DBName).Collection(repository.FeedsCollection)
	return &FeedRepository{db: db, cfg: cfg, collection: collection}
}

func (r *FeedRepository) Create(chatID int64) error {
	userFeeds := repository.UserFeed{ChatID: chatID, Feeds: []repository.Feed{}}
	_, err := r.collection.InsertOne(context.TODO(), userFeeds)
	if err != nil {
		return err
	}
	return nil
}

func (r *FeedRepository) Save(userFeed repository.UserFeed) error {
	feeds := userFeed.Feeds
	update := bson.D{
		{"$set", bson.D{
			{"feeds", feeds},
		}},
	}
	filter := bson.D{{"chatid", userFeed.ChatID}}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
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

	filter := bson.D{{"chatid", userFeed.ChatID}}
	update := bson.D{
		{"$set", bson.D{
			{"feeds", feeds},
		}},
	}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *FeedRepository) Get(chatID int64) (*repository.UserFeed, error) {
	filter := bson.D{{"chatid", chatID}}

	var result *repository.UserFeed
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (r *FeedRepository) Remove(userFeed repository.UserFeed, feedURLToRemove string) error {
	for i, feed := range userFeed.Feeds {
		if feed.URL == feedURLToRemove {
			newFeeds := append(userFeed.Feeds[:i], userFeed.Feeds[i+1:]...)
			filter := bson.D{{"chatid", userFeed.ChatID}}
			update := bson.D{
				{"$set", bson.D{
					{"feeds", newFeeds},
				}},
			}
			_, err := r.collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (r *FeedRepository) GetAll() ([]repository.UserFeed, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var results []repository.UserFeed

	for cursor.Next(context.TODO()) {
		var feed repository.UserFeed
		err := cursor.Decode(&feed)
		if err != nil {
			return nil, err
		}

		results = append(results, feed)
	}

	return results, nil
}
