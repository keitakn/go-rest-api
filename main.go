package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {

	users := Users{
		Store: map[string]*User{},
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", HealthCheck),
		rest.Get("/users", users.GetAllUsers),
		rest.Post("/users", users.PostUser),
		rest.Get("/users/:id", users.GetUser),
		rest.Put("/users/:id", users.PutUser),
		rest.Delete("/users/:id", users.DeleteUser),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type User struct {
	Id   string
	Name string
}

type Users struct {
	sync.RWMutex
	Store map[string]*User
}

func (u *Users) GetAllUsers(w rest.ResponseWriter, r *rest.Request) {
	u.RLock()
	users := make([]User, len(u.Store))
	i := 0
	for _, user := range u.Store {
		users[i] = *user
		i++
	}
	u.RUnlock()
	w.WriteJson(&users)
}

func (u *Users) GetUser(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	u.RLock()
	var user *User
	if u.Store[id] != nil {
		user = &User{}
		*user = *u.Store[id]
	}
	u.RUnlock()
	if user == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(user)
}

func (u *Users) PostUser(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Lock()
	id := fmt.Sprintf("%d", len(u.Store))
	user.Id = id
	u.Store[id] = &user
	u.Unlock()
	w.WriteJson(&user)
}

func (u *Users) PutUser(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	u.Lock()
	if u.Store[id] == nil {
		rest.NotFound(w, r)
		u.Unlock()
		return
	}
	user := User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		u.Unlock()
		return
	}
	user.Id = id
	u.Store[id] = &user
	u.Unlock()
	w.WriteJson(&user)
}

func (u *Users) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	u.Lock()
	delete(u.Store, id)
	u.Unlock()
	w.WriteHeader(http.StatusOK)
}

func HealthCheck(w rest.ResponseWriter, r *rest.Request) {
	responseId := GenerateResponseId()

	type HealthCheckResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var logger *zap.Logger
	logger = CreateLogger()
	logger.Info(
		"HealthCheck",
		zap.String("Channel", "go-rest-api"),
		zap.String("ResponseId", responseId),
	)

	// Fargateでセットした環境変数が読み込めるかテスト
	slackToken := os.Getenv("SLACK_TOKEN")

	res := HealthCheckResponse{Code: 200, Message: slackToken}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(res)
}

func CreateLogger() *zap.Logger {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	zapConfig := zap.Config{
		Level:    level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := zapConfig.Build()

	return logger
}

func GenerateResponseId() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	uu := u.String()

	return uu
}
