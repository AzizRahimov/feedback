package main

import (
	"context"
	"encoding/json"
	"feedback/users"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"

	"github.com/AzizRahimov/mux/pkg/mux"

	"github.com/google/uuid"
)

func init()  {
	f, err := os.OpenFile("file.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("This is a test log entry")
}

type MainServer struct {
	router   *mux.ExactMux
	usersSvc *users.Service
}




func NewMainServer(router *mux.ExactMux, usersSvc *users.Service) *MainServer {
	return &MainServer{router: router, usersSvc: usersSvc}
}



func (m *MainServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	m.router.ServeHTTP(writer, request)
}




func (m *MainServer) GetAllFeedbackes(w http.ResponseWriter, r *http.Request) {
	response, err := m.usersSvc.GetAllFeedback()
	if err != nil {
		log.Print("can't get all feedback",err)
		_, err = w.Write([]byte("can't get all feedback"))
		if err !=nil{
			log.Print(err)

		}
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_ = json.NewEncoder(w).Encode(response)
}

	func (m *MainServer) DeleteFeedback(w http.ResponseWriter, r *http.Request) {
		//
		id, ok := mux.FromContext(r.Context(),"id")
			if !ok{
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		log.Print(id)
		err := m.usersSvc.DeleteById(id)
		if err != nil {
			log.Print("can't delete feedback",err)
			_, err = w.Write([]byte("can't delete feedback"))
			if err != nil {
				log.Print(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}

func (m *MainServer) EditFeedback(w http.ResponseWriter, r *http.Request){
	id, ok := mux.FromContext(r.Context(),"id")
	if !ok{
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	log.Print(id)
	err := m.usersSvc.EditFeedbackByID(id)
	if err != nil {
		log.Print("can't update feedback ...",err)
	}
	_, err = w.Write([]byte("can't update feedback"))
	if err != nil {
		log.Print(err)
	}
	w.WriteHeader(http.StatusBadGateway)
	return
}






func (m *MainServer) CreateFeedback(w http.ResponseWriter, r *http.Request) {
		request := users.Feedback{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Print("can't parse request error:",err)
			_,err = w.Write([]byte("wrong request"))
			if err != nil {
				log.Print(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//разобраться

	// request validation
	if len(request.FeedbackText) < 5 || len(request.FeedbackTo) < 5 ||
		len(request.FeedbackBy) < 5 || len(request.FeedbackTopic) < 5{
		_, err = w.Write([]byte("wrong data")) // передаем байты об ошибке
		if err != nil {
			log.Print(err)
		}
		w.WriteHeader(http.StatusBadRequest )
		log.Print("wrong data from client")
		}

	request.Id = uuid.New()
	 err = m.usersSvc.Add(request)
	if err != nil {
		log.Print("creating the feedback error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("some error while creating feedback"))
		return
	}
	response, err := m.usersSvc.GetFeedbackByID(request.Id)
	 _=json.NewEncoder(w).Encode(response)
}

func main() {


	pool, err := pgxpool.Connect(
		context.Background(),
		"postgres://user:pass@192.168.99.100:5400/app",
	)
	if err != nil {
		panic(err)
	}


	usersSvc := users.NewService(pool)
	usersSvc.Start()
	router := mux.NewExactMux()
	server := NewMainServer(router, usersSvc)
	router.POST("/feedback", server.CreateFeedback)
	router.GET("/feedback/all", server.GetAllFeedbackes)
	router.DELETE("/feedback/{id}", server.DeleteFeedback)
	router.POST("/feedback/{id}", server.EditFeedback)
	fmt.Println("Server is listening")
	panic(http.ListenAndServe("0.0.0.0:9995", server))


}
