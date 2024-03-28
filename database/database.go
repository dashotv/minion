package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/fae"
	"github.com/dashotv/grimoire"
)

type Connector struct {
	Jobs *grimoire.Store[*Model]
}

func New(uri, db, collection string) (*Connector, error) {
	con, err := grimoire.New[*Model](uri, db, collection)
	if err != nil {
		return nil, fae.Wrap(err, "creating job store")
	}
	grimoire.Indexes(con, &Model{})

	return &Connector{Jobs: con}, nil
}

func (c *Connector) UpdateAbandonedJobs(ctx context.Context) error {
	_, err := c.Jobs.Collection.UpdateMany(ctx, bson.M{"status": bson.M{"$in": bson.A{StatusRunning, StatusQueued}}}, bson.M{"$set": bson.M{"status": StatusFailed}, "$push": bson.M{"attempts": bson.M{"error": "minion restarted"}}})
	if err != nil {
		return fae.Errorf("querying cancelled jobs: %s", err)
	}
	return nil
}

func (c *Connector) UpdateCancelledJobs(ctx context.Context) (int64, error) {
	res, err := c.Jobs.Collection.UpdateMany(ctx, bson.M{"status": StatusCancelled}, bson.M{"$set": bson.M{"status": StatusPending}})
	if err != nil {
		return 0, fae.Errorf("querying cancelled jobs: %s", err)
	}
	return res.ModifiedCount, nil
}
