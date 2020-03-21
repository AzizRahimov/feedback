package db

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool // dependency
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (service *Service) Start() {
	_, err := service.pool.Exec(context.Background(), dropExistingTableFeedback)
	if err != nil {
		log.Fatal("removing table feedback error:", err)
	}

	_, err = service.pool.Exec(context.Background(), createTableFeedback)
	if err != nil {
		log.Fatal("creating feedback error:", err)
	}

}

func (service *Service) Add(model Feedback) (err error) {
	_, err = service.pool.Exec(context.Background(), addFeedback,
		model.ID, model.FeedbackTopic, model.FeedbackBy, model.FeedbackTo, model.FeedbackText, model.Score)
	if err != nil {
		log.Print("adding feedback error:", err)
	}
	return err
}

func (service *Service) GetFeedbackByID(id uuid.UUID) (model Feedback, err error) {
	err = service.pool.QueryRow(context.Background(), getFeedbackByID, id).Scan(
		&model.ID, &model.FeedbackTopic, &model.FeedbackBy, &model.FeedbackTo, &model.FeedbackText, &model.Score)
	return
}

func (service *Service) GetAllFeedback() (models []Feedback, err error) {
	rows, err := service.pool.Query(context.Background(), getAllFeedbacks)
	model := Feedback{}
	for rows.Next() {
		err = rows.Scan(
			&model.ID, &model.FeedbackTopic, &model.FeedbackBy, &model.FeedbackTo, &model.FeedbackText, &model.Score)
		if err != nil {
			return
		}
		models = append(models, model)
	}
	return
}

func (service *Service) DeleteById(id uuid.UUID) (err error) {
	_, err = service.pool.Exec(context.Background(), deleteFeedback, id)
	if err != nil {
		log.Print("deleting feedback error:", err)
	}
	return
}

func (service *Service) EditFeedbackByID(model Feedback) (err error) {
	_, err = service.pool.Exec(context.Background(), editFeedback,
		model.ID, model.FeedbackTopic, model.FeedbackBy, model.FeedbackTo, model.FeedbackText, model.Score)
	if err != nil {
		log.Print("updating feedback error:", err)
	}
	return err
}

type Feedback struct {
	ID            uuid.UUID `json:"id"`
	FeedbackTopic string    `json:"feedback_topic"`
	FeedbackBy    string    `json:"feedback_by"`
	FeedbackTo    string    `json:"feedback_to"`
	FeedbackText  string    `json:"feedback_text"`
	Score         int       `json:"Score"`
}
