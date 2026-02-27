package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const urlsFilePath = "links.txt"

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("Установите TELEGRAM_TOKEN")
	}

	urls, err := readURLsFromFile(urlsFilePath)
	if err != nil {
		log.Fatalf("Не удалось прочитать %s: %v", urlsFilePath, err)
	}
	if len(urls) == 0 {
		log.Fatalf("Файл %s пустой", urlsFilePath)
	}

	chatID, err := strconv.ParseInt(strings.TrimSpace(os.Getenv("CHAT_ID")), 10, 64)
	if err != nil {
		log.Fatal("Установите корректный CHAT_ID")
	}

	interval := 30 * time.Second
	if raw := strings.TrimSpace(os.Getenv("CHECK_INTERVAL_SECONDS")); raw != "" {
		seconds, parseErr := strconv.Atoi(raw)
		if parseErr != nil || seconds <= 0 {
			log.Fatal("CHECK_INTERVAL_SECONDS должен быть положительным числом")
		}
		interval = time.Duration(seconds) * time.Second
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Ошибка отправки:", err)
			}
		}
	}()

	go func() {
		currentURLs := urls
		for {
			latestURLs, err := readURLsFromFile(urlsFilePath)
			if err != nil {
				log.Printf("Не удалось обновить %s: %v (использую предыдущий список)", urlsFilePath, err)
			} else if len(latestURLs) == 0 {
				log.Printf("Файл %s пустой (использую предыдущий список)", urlsFilePath)
			} else {
				currentURLs = latestURLs
			}

			for _, u := range currentURLs {
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
			time.Sleep(interval)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}

func readURLsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	urls := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
