package link

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cutter-url-go/internal/models"
	"cutter-url-go/internal/vo"
)

type Repository interface {
	Get(shortUrl vo.ShortURI) (*models.ShortLink, error)
	Insert(link *models.ShortLink) error
}

type linkRepository struct {
	c *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	ctx := context.Background()
	c := db.Collection("links")
	_, err := c.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "short_uri", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}
	return &linkRepository{c: c}
}

func (r linkRepository) Get(shortUrl vo.ShortURI) (*models.ShortLink, error) {
	ctx := context.Background()
	var link models.ShortLink
	l := r.c.FindOne(ctx, bson.M{"short_uri": shortUrl})
	if err := l.Decode(&link); err != nil {
		return nil, err
	}

	return &link, nil
}

func (r linkRepository) Insert(link *models.ShortLink) error {
	ctx := context.Background()
	_, err := r.c.InsertOne(ctx, link)
	if err != nil {
		return err
	}

	return nil
}
