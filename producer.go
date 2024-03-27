package minion

import (
	"context"
	"time"

	"github.com/dashotv/minion/database"
)

type Producer struct {
	Minion *Minion
	Queue  *Queue
	ch     chan string
}

func (p *Producer) Run(ctx context.Context) {
	p.ch = make(chan string, 1)
	p.Minion.Subscribe(func(n *Notification) {
		if n.Event == "job:created" && n.Kind == p.Queue.Name {
			p.ch <- n.JobID
		}
	})
	go p.listen(ctx)
}

func (p *Producer) listen(ctx context.Context) {
	for {
		select {
		case <-p.ch:
		case <-time.After(time.Duration(p.Queue.Interval) * time.Second):
			p.handle()
		case <-ctx.Done():
			close(p.ch)
			return
		}
	}
}

func (p *Producer) handle() {
	if p.Queue.Full() {
		return
	}

	i := p.Queue.Remaining()
	list, err := p.Minion.db.Jobs.Query().Where("client", p.Minion.Client).Where("queue", p.Queue.Name).Where("status", database.StatusPending).Asc("created_at").Limit(i).Run()
	if err != nil {
		p.Minion.Log.Errorf("querying pending jobs: %s", err)
	}

	for _, j := range list {
		j.Status = string(database.StatusQueued)
		err = p.Minion.db.Jobs.Save(j)
		if err != nil {
			p.Minion.Log.Errorf("updating job: %s", err)
		}

		p.Minion.notify("job:queued", j.ID.Hex(), j.Kind)
		p.Queue.channel <- j.ID.Hex()
	}
}
