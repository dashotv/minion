package minion

import (
	"go.mongodb.org/mongo-driver/bson"
)

type stat struct {
	Id struct {
		Status string `bson:"status"`
		Queue  string `bson:"queue"`
	} `bson:"_id"`
	Count int `bson:"count"`
}

type Stats map[string]map[string]int

func (m *Minion) SubscribeStats(f func(Stats)) {
	if m.statsEntry == 0 {
		id, err := m.cron.AddFunc("* * * * * *", m.stats)
		if err != nil {
			m.Log.Errorf("error scheduling stats: %s", err)
			return
		}
		m.statsEntry = id
	}
	m.statsSubs = append(m.statsSubs, f)
}

func (m *Minion) stats() {
	if m.statsEntry == 0 {
		return
	}
	if len(m.statsSubs) == 0 {
		m.Remove(m.statsEntry)
		return
	}

	// Equivalent to the following MongoDB query:
	// db.jobs.aggregate([
	//	{ $group: { _id: {status:"$status",queue:"$queue"}, count: {$sum: 1}}},
	//  { $project: { count: 1 } }
	// ])
	cur, err := m.db.Collection.Aggregate(m.Context, bson.A{
		bson.M{"$group": bson.M{"_id": bson.M{"queue": "$queue", "status": "$status"}, "count": bson.M{"$sum": 1}}},
		bson.M{"$project": bson.M{"count": 1}},
	})
	if err != nil {
		m.Log.Errorf("error querying stats: %s", err)
		m.Remove(m.statsEntry)
	}

	results := make([]*stat, 0)
	if err = cur.All(m.Context, &results); err != nil {
		m.Log.Errorf("error decoding stats: %s", err)
		m.Remove(m.statsEntry)
	}

	stats := Stats{}
	for _, s := range results {
		if _, ok := stats["totals"]; !ok {
			stats["totals"] = make(map[string]int)
		}
		if _, ok := stats[s.Id.Queue]; !ok {
			stats[s.Id.Queue] = make(map[string]int)
		}
		stats[s.Id.Queue][s.Id.Status] = s.Count
		stats["totals"][s.Id.Status] += s.Count
	}

	for _, f := range m.statsSubs {
		f(stats)
	}
}
