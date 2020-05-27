package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xcorter/vkwatcher/internal/app/infrastructure/config"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/observable"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/telegramclient"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/vkclient"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/watcher"
	. "io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("run")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.New()

	client := &http.Client{}

	db := getDB(cfg)
	provider := observable.NewProvider(db)
	manager := observable.NewManager(provider)

	vkclientapi := vkclient.NewClient(client, cfg.VkApiKey)
	telegramClient := getTelegramClient(cfg)
	w := watcher.NewWatcher(
		manager,
		provider,
		vkclientapi,
		telegramClient,
	)
	ctx, cancel := context.WithCancel(context.Background())

	chatWatcher := telegramclient.NewChatWatcher(telegramClient, provider)
	chatWatcher.Watch(ctx)

	//Регистрируем сигналы SIGINT и SIGTERM, чтобы завершить работу красиво
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		fmt.Println("Нажали crtl+c. Полылаем сингал завершения")
		cancel()
		fmt.Println("Ждем завершения вокреков")
	}()

	w.Start(ctx)
}

func getTelegramClient(cfg *config.Config) *telegramclient.TelegramClient {
	client := getTelegramHTTPClient(cfg)

	return telegramclient.NewClient(
		cfg.TelegramApiKey,
		client,
	)
}

func getTelegramHTTPClient(cfg *config.Config) *http.Client {
	//creating the proxyURL
	if cfg.UseProxy == "" {
		return &http.Client{
			//Timeout: 10 * time.Second,
		}
	}
	proxyURL, err := url2.Parse(cfg.UseProxy)
	if err != nil {
		log.Println(err)
	}
	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	return &http.Client{
		Transport: transport,
		//Timeout: 5 * time.Second,
	}
}

func getDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("sqlite3", "./observable.db")
	checkErr(err)

	if cfg.Env == "dev" {
		tableSql, err := ReadFile("./resources/schema.sql")
		checkErr(err)

		_, err = db.Exec(string(tableSql))
		if err != nil {
			panic(err)
		}
	}
	
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
