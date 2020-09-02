package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vlamitin/pusher/internal/messages"
	"github.com/vlamitin/pusher/internal/notifier"
	"github.com/vlamitin/pusher/internal/persistence"
	"github.com/vlamitin/pusher/internal/server"
	"os"
)

type Config struct {
	serverPort  int
	appToken    string
	userToken   string
	logFilePath string
	pgHost      string
	pgPort      string
	pgUser      string
	pgPassword  string
	pgDb        string
}

func main() {
	config := getConfig()
	err := setupLogger(config.logFilePath)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := persistence.CreatePool(config.pgHost, config.pgPort, config.pgUser, config.pgPassword, config.pgDb)
	if err != nil {
		log.Fatal(fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err))
	}

	pushoverConnector := notifier.NewPushoverConnector(config.userToken, config.appToken)
	messagesRepo := messages.NewPgMessageRepo(pool)
	messagesService := messages.NewMessageService(messagesRepo, pushoverConnector)

	err = server.Start(config.serverPort,
		server.Endpoint{
			Path:        "/messages",
			Verb:        "POST",
			HandlerFunc: messages.NewMessageSenderController(&messagesService).PostMessage,
		},
		server.Endpoint{
			Path:        "/messages",
			Verb:        "GET",
			HandlerFunc: messages.NewCountsReceiverController(&messagesService).GetStatusCounts,
		},
	)

	log.Warn(fmt.Errorf("Shutting down: %v\n", err))
	os.Exit(0)
}

func getConfig() Config {
	serverPortPtr := flag.Int("port", 8080, "Port to run server on")
	appTokenPtr := flag.String("pushover_app_token", "", "App token from https://pushover.net/apps/{app_token} ")
	userTokenPtr := flag.String("pushover_user_token", "", "Your User Key from https://pushover.net")
	logFilePtr := flag.String("log_file", "", "Name of log file. If not specified - log to stdout")

	flag.Parse()

	return Config{
		serverPort:  *serverPortPtr,
		appToken:    *appTokenPtr,
		userToken:   *userTokenPtr,
		logFilePath: *logFilePtr,
		pgHost:      os.Getenv("POSTGRES_HOST"),
		pgPort:      os.Getenv("POSTGRES_PORT"),
		pgUser:      os.Getenv("POSTGRES_USER"),
		pgPassword:  os.Getenv("POSTGRES_PASSWORD"),
		pgDb:        os.Getenv("POSTGRES_DB"),
	}
}

func setupLogger(logFilePath string) error {
	if logFilePath == "" {
		log.SetOutput(os.Stdout)
	} else {
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("error when tried to open log file %s: %v\n", logFilePath, err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.SetLevel(log.InfoLevel)
	return nil
}
