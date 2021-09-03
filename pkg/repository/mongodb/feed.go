package mongodb

import (
	"context"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	userFeeds := repository.Feed{ChatID: chatID, Feeds: []string{}}
	_, err := r.collection.InsertOne(context.TODO(), userFeeds)
	if err != nil {
		return err
	}
	return nil
}

func (r *FeedRepository) Add(feed repository.Feed, feedURL string) error {
	for _, feed := range feed.Feeds {
		if feed == feedURL {
			return nil
		}
	}

	feeds := append(feed.Feeds, feedURL)

	filter := bson.D{{"chatid", feed.ChatID}}
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

func (r *FeedRepository) Get(chatID int64) (*repository.Feed, error) {
	filter := bson.D{{"chatid", chatID}}

	var result *repository.Feed
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (r *FeedRepository) Remove(feed repository.Feed, feedURLToRemove string) error {
	for i, feedURL := range feed.Feeds {
		if feedURL == feedURLToRemove {
			newFeeds := append(feed.Feeds[:i], feed.Feeds[i+1:]...)
			filter := bson.D{{"chatid", feed.ChatID}}
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

func (r *FeedRepository) GetAll() ([]repository.Feed, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var results []repository.Feed

	for cursor.Next(context.TODO()) {
		var feed repository.Feed
		err := cursor.Decode(&feed)
		if err != nil {
			return nil, err
		}

		results = append(results, feed)
	}

	return results, nil
}
