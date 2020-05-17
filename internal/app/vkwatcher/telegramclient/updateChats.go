package telegramclient

import (
	"context"
	"fmt"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/observable"
)

type ChatWatcher struct {
	telegramClient *TelegramClient
	provider       *observable.Provider
}

func (c *ChatWatcher) Watch(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stop chat watcher")
				return
			default:
				updates, err := c.telegramClient.GetUpdatesChan()

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				for update := range updates {
					if update.Message == nil { // ignore any non-Message Updates
						continue
					}

					ob := observable.NewMusicObservable(update.Message.Chat.UserName, "aaa", update.Message.Chat.ID)

					c.provider.Save(ob)
				}
			}
		}
	}()
}

func NewChatWatcher(client *TelegramClient, provider *observable.Provider) ChatWatcher {
	return ChatWatcher{
		telegramClient: client,
		provider:       provider,
	}
}
