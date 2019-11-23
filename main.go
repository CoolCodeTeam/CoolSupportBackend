package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/CoolCodeTeam/CoolSupportBackend/chats/repository"
	supportDelivery "github.com/CoolCodeTeam/CoolSupportBackend/supports/delivery"
	supportRepository "github.com/CoolCodeTeam/CoolSupportBackend/supports/repository"
	supportUseCase "github.com/CoolCodeTeam/CoolSupportBackend/supports/usecase"

	chatDelivery "github.com/CoolCodeTeam/CoolSupportBackend/chats/delivery"
	chatRepository "github.com/CoolCodeTeam/CoolSupportBackend/chats/repository"
	chatUseCase "github.com/CoolCodeTeam/CoolSupportBackend/chats/usecase"

	messagesDelivery "github.com/CoolCodeTeam/CoolSupportBackend/messages/delivery"
	messagesRepository "github.com/CoolCodeTeam/CoolSupportBackend/messages/repository"
	messagesUseCase "github.com/CoolCodeTeam/CoolSupportBackend/messages/usecase"

	utils2 "github.com/CoolCodeTeam/CoolSupportBackend/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/kabukky/httpscerts"
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
	sessionRepository := supportRepository.NewSessionRedisStore(redisConn)

	supportsRepository := supportRepository.NewSupportDBStore(db)
	supportsUseCase := supportUseCase.NewSupportUseCase(supportsRepository, sessionRepository)
	usersApi := supportDelivery.NewSupportHandlers(supportsUseCase, sessionRepository, utils)

	chatsUseCase := useCase.NewChatsUseCase(repository.NewChatsDBRepository(db), supportsRepository)
	messagesUseCase := useCase.NewMessageUseCase(repository.NewMessageDbRepository(db), chatsUseCase)
	notificationsUseCase := useCase.NewNotificationUseCase()
	chatsApi := delivery.NewChatHandlers(supportsUseCase, sessionRepository, chatsUseCase, utils)
	notificationApi := delivery.NewNotificationHandlers(supportsUseCase, sessionRepository, chatsApi.Chats, notificationsUseCase, utils)
	messagesApi := delivery.NewMessageHandlers(messagesUseCase, supportsUseCase, sessionRepository, notificationsUseCase, utils)
	middlewares := middleware.HandlersMiddlwares{
		Sessions: sessionRepository,
		Logger:   logrusLogger,
	}

	r := mux.NewRouter()
	handler := middlewares.PanicMiddleware(middlewares.LogMiddleware(r, logrusLogger))
	r.HandleFunc("/users", usersApi.SignUp).Methods("POST")
	r.HandleFunc("/login", usersApi.Login).Methods("POST")
	r.Handle("/users/{id:[0-9]+}", middlewares.AuthMiddleware(usersApi.EditProfile)).Methods("PUT")
	r.Handle("/logout", middlewares.AuthMiddleware(usersApi.Logout)).Methods("DELETE")
	r.Handle("/photos", middlewares.AuthMiddleware(usersApi.SavePhoto)).Methods("POST")
	r.Handle("/photos/{id:[0-9]+}", middlewares.AuthMiddleware(usersApi.GetPhoto)).Methods("GET")
	r.Handle("/users/{id:[0-9]+}", middlewares.AuthMiddleware(usersApi.GetUser)).Methods("GET")
	r.Handle("/users/{name:[((a-z)|(A-Z))0-9_-]+}", middlewares.AuthMiddleware(usersApi.FindUsers)).Methods("GET")
	r.HandleFunc("/users", usersApi.GetUserBySession).Methods("GET") //TODO:Добавить в API

	r.HandleFunc("/chats", chatsApi.PostChat).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}/chats", chatsApi.GetChatsByUser).Methods("GET")
	r.Handle("/chats/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.GetChatById)).Methods("GET")
	r.Handle("/chats/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.RemoveChat)).Methods("DELETE")

	r.Handle("/channels/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.GetChannelById)).Methods("GET")
	r.Handle("/channels/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.EditChannel)).Methods("PUT")
	r.Handle("/channels/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.RemoveChannel)).Methods("DELETE")
	r.Handle("/channels/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.SendMessage)).Methods("POST")
	r.Handle("/channels/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.GetMessagesByChatID)).Methods("GET")
	r.Handle("/channels/{id:[0-9]+}/messages", middlewares.AuthMiddleware(chatsApi.RemoveChannel)).Methods("DELETE")
	//TODO: r.Handle("/channels/{id:[0-9]+}/members", middlewares.AuthMiddleware(chatsApi.LogoutFromChannel)).Methods("DELETE")
	r.Handle("/workspaces/{id:[0-9]+}/channels", middlewares.AuthMiddleware(chatsApi.PostChannel)).Methods("POST")

	r.Handle("/workspaces/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.GetWorkspaceById)).Methods("GET")
	r.Handle("/workspaces/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.EditWorkspace)).Methods("PUT")
	//TODO: r.Handle("/workspaces/{id:[0-9]+}/members", middlewares.AuthMiddleware(chatsApi.LogoutFromWorkspace)).Methods("DELETE")
	r.Handle("/workspaces/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.RemoveWorkspace)).Methods("DELETE")
	r.Handle("/workspaces", middlewares.AuthMiddleware(chatsApi.PostWorkspace)).Methods("POST")
	r.Handle("/chats/{id:[0-9]+}/notifications", middlewares.AuthMiddleware(notificationApi.HandleNewWSConnection))

	r.Handle("/chats/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.SendMessage)).Methods("POST").
		HeadersRegexp("Content-Type", "application/(text|json)")
	r.Handle("/chats/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.GetMessagesByChatID)).Methods("GET")
	r.Handle("/messages/{text:[((a-z)|(A-Z))0-9_-]+}", middlewares.AuthMiddleware(messagesApi.FindMessages)).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}", middlewares.AuthMiddleware(messagesApi.DeleteMessage)).Methods("DELETE")
	r.Handle("/messages/{id:[0-9]+}", middlewares.AuthMiddleware(messagesApi.EditMessage)).Methods("PUT")
	log.Println("Server started")
	genetateSSL()

	//err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", corsMiddleware(handler))
	//if err != nil {
	//	logrus.Errorf("Can not listen https, error: %v", err.Error())
	//}

	err = http.ListenAndServe(":8080", corsMiddleware(handler))
	if err != nil {
		logrusLogger.Error(err)
		return
	}
}

func genetateSSL() {
	// Проверяем, доступен ли cert файл.
	err := httpscerts.Check("cert.pem", "key.pem")
	// Если он недоступен, то генерируем новый.
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "95.163.209.195:8080")
		if err != nil {
			logrus.Fatal("Ошибка: Не можем сгенерировать https сертификат.")
		}
	}
}
