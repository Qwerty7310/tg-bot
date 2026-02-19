package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("Установите TELEGRAM_TOKEN")
	}

	url := []string{
		//"https://comedyconcert.ru/event/neigry-january27",
		//"https://comedyconcert.ru/event/neigry-januar27",
		//"https://comedyconcert.ru/event/neigry-janua27",
		//"https://comedyconcert.ru/event/neigry-janu27",
		//"https://comedyconcert.ru/event/neigry-jan27",
		//"https://comedyconcert.ru/event/neigry-ja27",
		//"https://comedyconcert.ru/event/neigry-j27",
		"https://comedyconcert.ru/event/neigry-february26",
		"https://comedyconcert.ru/event/neigry-februar26",
		"https://comedyconcert.ru/event/neigry-februa26",
		"https://comedyconcert.ru/event/neigry-febru26",
		"https://comedyconcert.ru/event/neigry-febr26",
		"https://comedyconcert.ru/event/neigry-feb26",
		"https://comedyconcert.ru/event/neigry-fe26",
		"https://comedyconcert.ru/event/neigry-f26",
		"https://comedyconcert.ru/event/neigry-march26",
		"https://comedyconcert.ru/event/neigry-marc26",
		"https://comedyconcert.ru/event/neigry-mar26",
		"https://comedyconcert.ru/event/neigry-ma26",
		"https://comedyconcert.ru/event/neigry-m26",
		"https://comedyconcert.ru/event/neigry-april26",
		"https://comedyconcert.ru/event/neigry-apri26",
		"https://comedyconcert.ru/event/neigry-apr26",
		"https://comedyconcert.ru/event/neigry-ap26",
		"https://comedyconcert.ru/event/neigry-a26",
		"https://comedyconcert.ru/event/neigry-may26",
		"https://comedyconcert.ru/event/neigry-ma26",
		"https://comedyconcert.ru/event/neigry-m26",
		"https://comedyconcert.ru/event/neigry-june26",
		"https://comedyconcert.ru/event/neigry-jun26",
		"https://comedyconcert.ru/event/neigry-ju26",
		"https://comedyconcert.ru/event/neigry-j26",
		"https://comedyconcert.ru/event/neigry-july26",
		"https://comedyconcert.ru/event/neigry-jul26",
		"https://comedyconcert.ru/event/neigry-ju26",
		"https://comedyconcert.ru/event/neigry-j26",
		"https://comedyconcert.ru/event/neigry-august26",
		"https://comedyconcert.ru/event/neigry-augus26",
		"https://comedyconcert.ru/event/neigry-augu26",
		"https://comedyconcert.ru/event/neigry-aug26",
		"https://comedyconcert.ru/event/neigry-au26",
		"https://comedyconcert.ru/event/neigry-a26",
		"https://comedyconcert.ru/event/neigry-september26",
		"https://comedyconcert.ru/event/neigry-septembe26",
		"https://comedyconcert.ru/event/neigry-septemb26",
		"https://comedyconcert.ru/event/neigry-septem26",
		"https://comedyconcert.ru/event/neigry-septe26",
		"https://comedyconcert.ru/event/neigry-sept26",
		"https://comedyconcert.ru/event/neigry-sep26",
		"https://comedyconcert.ru/event/neigry-se26",
		"https://comedyconcert.ru/event/neigry-s26",
		//"https://comedyconcert.ru/event/neigry-msk-dec",
		//"https://comedyconcert.ru/event/neigry-msk-dec25",
	}

	chatID := int64(1622492999)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Настраиваем получение апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обрабатываем входящие сообщения
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// простой эхо-бот
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Ошибка отправки:", err)
			}
		}
	}()

	go func() {
		for {
			for _, u := range url {
				resp, err := http.Get(u)
				if err != nil {
					log.Println("Не удалось сделать запрос:", err)
					continue
				}

				if resp.StatusCode == http.StatusOK {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("200 OK\n%s", u))
					_, err := bot.Send(msg)
					if err != nil {
						log.Println("Ошибка отправки:", err)
					}
				}

				resp.Body.Close()
			}
			time.Sleep(30 * time.Second)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
