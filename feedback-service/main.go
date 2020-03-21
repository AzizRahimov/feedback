package main

import (
	"context"
	"encoding/json"
	"feedback/feedback-service/db"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AzizRahimov/mux/pkg/mux"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func init() {
	f, err := os.OpenFile("file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("This is a test log entry")
}

type MainServer struct {
	router      *mux.ExactMux
	feedbackSvc *db.Service
}

func NewMainServer(router *mux.ExactMux, feedbackSvc *db.Service) *MainServer {
	return &MainServer{router: router, feedbackSvc: feedbackSvc}
}

func main() {

	pool, err := pgxpool.Connect(
		context.Background(),
		"postgres://user:pass@192.168.99.100:5400/app",
	)
	if err != nil {
		panic(err)
	}

	feedbackSvc := db.NewService(pool)
	feedbackSvc.Start()
	router := mux.NewExactMux()
	server := NewMainServer(router, feedbackSvc)
	router.POST("/feedback", server.CreateFeedback)
	router.GET("/feedback/all", server.GetAllFeedbackes)
	router.DELETE("/feedback/{id}", server.DeleteFeedback)
	router.POST("/feedback/{id}", server.EditFeedback)
	fmt.Println("Server is listening")
	panic(http.ListenAndServe("0.0.0.0:9995", server))
}

func (m *MainServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	m.router.ServeHTTP(writer, request)
}

func (m *MainServer) GetAllFeedbackes(w http.ResponseWriter, r *http.Request) {
	response, err := m.feedbackSvc.GetAllFeedback()
	if err != nil {
		log.Print("can't get all feedback", err)
		_, err = w.Write([]byte("can't get all feedback"))
		if err != nil {
			log.Println("sending response failed:", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("sending response failed", err)
	}
}

func (m *MainServer) DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	val, ok := mux.FromContext(r.Context(), "id")
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	id, err := uuid.Parse(val)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("wrong feedback id:" + err.Error()))
		if err != nil {
			log.Println("sending response error:", err)
		}
		return
	}
	err = m.feedbackSvc.DeleteById(id)
	if err != nil {
		log.Print("can't delete feedback", err)
		_, err = w.Write([]byte("can't delete feedback"))
		if err != nil {
			log.Print(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *MainServer) EditFeedback(w http.ResponseWriter, r *http.Request) {
	val, ok := mux.FromContext(r.Context(), "id")
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	id, err := uuid.Parse(val)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("wrong feedback id"))
		if err != nil {
			log.Println("sending response error:", err)
		}
		return
	}
	request := db.Feedback{}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("EditFeedback, parcing request error:", err)
		_, err = w.Write([]byte("wrong request"))
		if err != nil {
			log.Println("sending response error:", err)
		}
		return
	}

	request.ID = id
	err = m.feedbackSvc.EditFeedbackByID(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("editing feedback failed:", err)
		_, err = w.Write([]byte("editing feedback failed"))
		if err != nil {
			log.Println("sending response error:", err)
		}
		return
	}

	response, err := m.feedbackSvc.GetFeedbackByID(request.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("getting feedback failed:", err)
		_, err = w.Write([]byte("getting feedback failed"))
		if err != nil {
			log.Println("sending response error:", err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("sending response error:", err)
	}
}

func (m *MainServer) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	request := db.Feedback{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Print("can't parse request error:", err)
		_, err = w.Write([]byte("wrong request"))
		if err != nil {
			log.Print(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// request validation
	if len(request.FeedbackText) < 5 || len(request.FeedbackTo) < 5 ||
		len(request.FeedbackBy) < 5 || len(request.FeedbackTopic) < 5 {
		_, err = w.Write([]byte("wrong data")) // передаем байты об ошибке
		if err != nil {
			log.Print(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Print("wrong data from client")
	}

	request.ID = uuid.New()
	err = m.feedbackSvc.Add(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("creating the feedback error:", err)
		_, err = w.Write([]byte("some error while creating feedback"))
		if err != nil {
			log.Println(err)
		}
		return
	}
	response, err := m.feedbackSvc.GetFeedbackByID(request.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("getting feedback error:", err)
		_, err = w.Write([]byte("some error while getting feedback"))
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("sending response failed:", err)
	}
}
