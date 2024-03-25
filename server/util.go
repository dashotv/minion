package main

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type H map[string]interface{}

func QueryParamInt(c echo.Context, name string, def int) int {
	param := c.QueryParam(name)
	result, err := strconv.Atoi(param)
	if err != nil {
		return def
	}
	return result
}

func (r *Router) jobStats() (*Stats, error) {
	total, err := r.DB.Jobs.Query().Count()
	if err != nil {
		return nil, err
	}

	cur, err := r.DB.Jobs.Collection.Aggregate(context.Background(), bson.A{
		bson.M{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	list := []*result{}
	if err = cur.All(context.Background(), &list); err != nil {
		return nil, err
	}

	stats := &Stats{
		Total: total,
	}
	for _, raw := range list {
		switch raw.ID {
		case "pending":
			stats.Pending = raw.Count
		case "queued":
			stats.Queued = raw.Count
		case "running":
			stats.Running = raw.Count
		case "cancelled":
			stats.Cancelled = raw.Count
		case "failed":
			stats.Failed = raw.Count
		case "finished":
			stats.Finished = raw.Count
		}
	}
	return stats, nil
}

type result struct {
	ID    string `bson:"_id"`
	Count int64  `bson:"count"`
}

type Stats struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Queued    int64 `json:"queued"`
	Running   int64 `json:"running"`
	Cancelled int64 `json:"cancelled"`
	Failed    int64 `json:"failed"`
	Finished  int64 `json:"finished"`
}
