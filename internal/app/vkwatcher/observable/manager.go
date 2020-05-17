package observable

import (
	"context"
	"fmt"
	"time"
)

type Manager struct {
	provider *Provider
	iterator int
	queue    chan *Observable
}

func (m *Manager) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stop manager")
				return
			default:
				data := m.provider.GetData()
				for _, ob := range data {
					m.queue <- &ob
				}
				time.Sleep(10 * time.Minute)
			}
		}
	}()
}

func NewManager(provider *Provider) *Manager {
	return &Manager{
		provider: provider,
		iterator: 0,
		queue:    make(chan *Observable, 10),
	}
}

func (m *Manager) GetObservable() <-chan *Observable {
	return m.queue
}
