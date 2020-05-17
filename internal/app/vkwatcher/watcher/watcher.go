package watcher

import (
	"context"
	"fmt"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/observable"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/telegramclient"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/vkclient"
	"strings"
	"sync"
	"time"
)

type Watcher struct {
	data               []observable.Observable
	workersAmount      int
	waitGroup          sync.WaitGroup
	vkclient           *vkclient.Client
	manager            *observable.Manager
	telegramclient     *telegramclient.TelegramClient
	observableProvider *observable.Provider
}

func (w *Watcher) Start(ctx context.Context) {
	fmt.Println("start!")

	w.manager.Run(ctx)
	for i := 0; i < w.workersAmount; i++ {
		w.waitGroup.Add(1)
		go w.watch(ctx)
	}

	w.waitGroup.Wait()
}

func (w *Watcher) watch(ctx context.Context) {
	defer w.waitGroup.Done()
	defer fmt.Println("watch done")
	fmt.Println("watch")
	for {
		//Делаем задачки пока сверху не придет указ остановиться
		select {

		case ob := <-w.manager.GetObservable():
			fmt.Println("%+v\n", ob)
			w.scan(*ob)
		default:
			select {
			case <-ctx.Done():
				fmt.Println("watcher done")
				return
			default:
				//	немного подождем и дальше ждать задачи
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (w *Watcher) scan(observable observable.Observable) {
	fmt.Println("start scan")
	offset := ""
	for {
		result, err := w.vkclient.GetData(observable.Value, offset)
		if err != nil {
			return
		}
		lastScan := observable.LastScan
		for _, item := range result.Response.Items {
			if observable.LastScan >= item.Date {
				break
			}
			if item.Date > lastScan {
				lastScan = item.Date
			}

			if w.hasObservableObject(observable, item) {
				w.telegramclient.SendMessage(observable, item)
			}
		}
		observable.LastScan = lastScan
		w.observableProvider.UpdateLastScan(observable)

		offset = result.Response.NextFrom
		if offset == "" {
			fmt.Println("stop scan")
			return
		}
	}
}

func (w *Watcher) hasObservableObject(observable observable.Observable, item vkclient.Item) bool {
	for _, attachment := range item.Attachments {
		if attachment.Type == "audio" {
			if strings.ToLower(attachment.Audio.Artist) == strings.ToLower(observable.Value) {
				return true
			}
		}
	}
	return false
}

func NewWatcher(
	manager *observable.Manager,
	observableProvider *observable.Provider,
	vkclientapi *vkclient.Client,
	telegramclientapi *telegramclient.TelegramClient,
) Watcher {
	return Watcher{
		workersAmount:      10,
		manager:            manager,
		vkclient:           vkclientapi,
		telegramclient:     telegramclientapi,
		observableProvider: observableProvider,
	}
}
