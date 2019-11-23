package main

import (
	"database/sql"
	"flag"
	"fmt"
	supportDelivery "github.com/CoolCodeTeam/CoolSupportBackend/supports/delivery"
	supportRepository "github.com/CoolCodeTeam/CoolSupportBackend/supports/repository"
	supportUseCase "github.com/CoolCodeTeam/CoolSupportBackend/supports/usecase"
	"github.com/gorilla/handlers"

	chatDelivery "github.com/CoolCodeTeam/CoolSupportBackend/chats/delivery"
	chatRepository "github.com/CoolCodeTeam/CoolSupportBackend/chats/repository"
	chatUseCase "github.com/CoolCodeTeam/CoolSupportBackend/chats/usecase"

	messageDelivery "github.com/CoolCodeTeam/CoolSupportBackend/messages/delivery"
	messageRepository "github.com/CoolCodeTeam/CoolSupportBackend/messages/repository"
	messageUseCase "github.com/CoolCodeTeam/CoolSupportBackend/messages/usecase"

	notificationDelivery "github.com/CoolCodeTeam/CoolSupportBackend/notifications/delivery"
	notificationUseCase "github.com/CoolCodeTeam/CoolSupportBackend/notifications/usecase"

	utils2 "github.com/CoolCodeTeam/CoolSupportBackend/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	DB_USER     = "supportdbuser"
	DB_PASSWORD = "1234"
	DB_NAME     = "supportdb"
	DB_DRIVER   = "postgres"
)

var (
	redisAddr = flag.String("addr", "redis://localhost:6379", "redis addr")
)

func main() {

	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	f, err := os.OpenFile("logs.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrusLogger.Error("Can`t open file:" + err.Error())
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stderr, f)
	logrusLogger.SetOutput(mw)

	utils := utils2.NewHandlersUtils(logrusLogger)

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open(DB_DRIVER, dbinfo)

	if err != nil {
		log.Printf("Error before started: %s", err.Error())
		return
	}
	if db == nil {
		log.Printf("Can not connect to database")
		return
	}

	redisConn := &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.DialURL(*redisAddr)
		},
	}
	if err != nil {
		log.Fatalf("cant connect to redis")
		return
	}
	defer redisConn.Close()

	defer db.Close()

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://boiling-chamber-90136.herokuapp.com", "https://boiling-chamber-90136.herokuapp.com", "http://localhost:3000"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	sessionRepository := supportRepository.NewSessionRedisStore(redisConn)

	supportsRepository := supportRepository.NewSupportDBStore(db)
	supportsUseCase := supportUseCase.NewSupportUseCase(supportsRepository, sessionRepository)
	supportsApi := supportDelivery.NewSupportHandlers(supportsUseCase, sessionRepository, utils)

	chatsUseCase := chatUseCase.NewChatsUseCase(chatRepository.NewChatsDBRepository(db))
	messagesUseCase := messageUseCase.NewMessageUseCase(messageRepository.NewMessageDbRepository(db), chatsUseCase)
	notificationsUseCase := notificationUseCase.NewNotificationUseCase()
	chatsApi := chatDelivery.NewChatHandlers(supportsUseCase, chatsUseCase, utils)
	notificationApi := notificationDelivery.NewNotificationHandlers(notificationsUseCase, supportsUseCase, utils)
	messagesApi := messageDelivery.NewMessageHandlers(messagesUseCase, supportsUseCase, notificationsUseCase, utils)

	r := mux.NewRouter()
	handler := r
	r.HandleFunc("/login", supportsApi.Login).Methods("POST")
	r.HandleFunc("/logout", supportsApi.Logout).Methods("DELETE")
	r.HandleFunc("/users", supportsApi.GetSupportBySession).Methods("GET")

	r.HandleFunc("/users/{id:[0-9]+}/chats", chatsApi.GetChatsByUser).Methods("GET")

	r.HandleFunc("/channels/{id:[0-9]+}/messages", messagesApi.SendMessage).Methods("POST")
	r.HandleFunc("/channels/{id:[0-9]+}/messages", messagesApi.GetMessagesByChatID).Methods("GET")

	r.HandleFunc("/chats/{id:[0-9]+}/notifications", notificationApi.HandleNewSupportWSConnection)

	r.HandleFunc("/chats/{id:[0-9]+}/messages", messagesApi.SendMessage).Methods("POST").
		HeadersRegexp("Content-Type", "application/(text|json)")
	r.HandleFunc("/chats/{id:[0-9]+}/messages", messagesApi.GetMessagesByChatID).Methods("GET")
	log.Println("Server started")

	err = http.ListenAndServe(":8081", corsMiddleware(handler))
	if err != nil {
		logrusLogger.Error(err)
		return
	}
}
