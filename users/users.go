package users

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool // dependency
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (service *Service) Start() {
	_,err := service.pool.Exec(context.Background(), `
DROP TABLE IF EXISTS feedback;
`)
	if err != nil {
		log.Fatal("can't remove table feedback")
	}

	_, err = service.pool.Exec(context.Background(),
		`
CREATE TABLE IF NOT EXISTS feedback (
	id UUID PRIMARY KEY not null,
	feedback_topic varchar(256) not null,
	feedback_by varchar(80) not null,
    feedback_to varchar(80) not null,
	feedback_text text not null,
	score integer default 0,
    removed BOOLEAN DEFAULT FALSE

);
`)
	if err !=nil{
		log.Fatal("can't create db feedback", err)
	}

}



func (service *Service) Add(model Feedback) (err error) {
	_, err = service.pool.Exec(context.Background(),
		`INSERT INTO feedback  (id, feedback_topic, feedback_by, feedback_to, feedback_text, score)
	VALUES ($1, $2, $3, $4, $5, $6);`,
	model.Id, model.FeedbackTopic, model.FeedbackBy, model.FeedbackTo, model.FeedbackText, model.Score)
	if err != nil {
		log.Print("cant't add feedback:", err)
	}
return err
}

func (service *Service) GetFeedbackByID(id uuid.UUID)(model Feedback, err error){
	err = service.pool.QueryRow(context.Background(),
		`SELECT id, feedback_topic, feedback_by, feedback_to, feedback_text, score 
			FROM feedback WHERE removed = false AND id = $1`, id).Scan(
				&model.Id, &model.FeedbackTopic, &model.FeedbackBy, &model.FeedbackTo, &model.FeedbackText, &model.Score)
	return
}


func (service *Service) GetAllFeedback()(models []Feedback, err error) {
	rows, err := service.pool.Query(context.Background(),
		`SELECT id, feedback_topic, feedback_by, feedback_to, feedback_text, score 
			FROM feedback WHERE removed = false`)
	model := Feedback{}
	for rows.Next() {
		err = rows.Scan(
			&model.Id, &model.FeedbackTopic, &model.FeedbackBy, &model.FeedbackTo, &model.FeedbackText, &model.Score)
		if err != nil {
			return
		}
		models = append(models, model)

	}
	return
}


	func (service *Service) DeleteById(id string)(err error){
		_, err = service.pool.Exec(context.Background(),
			`update feedback set removed = true where id = $1`, id)
		if err != nil {
			log.Print("can't Delete feedback from db ")
		}
		return
	}


func (service *Service) EditFeedbackByID(model Feedback)(err error){
	_, err = service.pool.Exec(context.Background(),
		`update feedback set feedback_topic=$1, feedback_by=$2, feedback_to=$3, feedback_text=$4, score=$5 WHERE  id = $6`,
	model.Id, model.FeedbackTopic, model.FeedbackBy, model.FeedbackTo, model.FeedbackText, model.Score)
	if err != nil {
		log.Print("cant't update feedback:", err)
	}
	return err
}




	type Feedback struct {
	Id            uuid.UUID `json:"id"`
	FeedbackTopic string    `json:"feedback_topic"`
	FeedbackBy    string    `json:"feedback_by"`
	FeedbackTo    string    `json:"feedback_to"`
	FeedbackText  string    `json:"feedback_text"`
	Score         int       `json:"Score"`
}
