package db

const (
	dropExistingTableFeedback = `DROP TABLE IF EXISTS feedback;`

	createTableFeedback = `CREATE TABLE IF NOT EXISTS feedback (
		id UUID PRIMARY KEY not null,
		feedback_topic varchar(256) not null,
		feedback_by varchar(80) not null,
		feedback_to varchar(80) not null,
		feedback_text text not null,
		score integer default 0,
		removed BOOLEAN DEFAULT FALSE);`
	addFeedback = `INSERT INTO feedback  (id, feedback_topic, feedback_by, feedback_to, feedback_text, score)
		VALUES ($1, $2, $3, $4, $5, $6);`
	getFeedbackByID = `SELECT id, feedback_topic, feedback_by, feedback_to, feedback_text, score 
		FROM feedback WHERE removed = false AND id = $1`
	getAllFeedbacks = `SELECT id, feedback_topic, feedback_by, feedback_to, feedback_text, score 
		FROM feedback WHERE removed = false`
	deleteFeedback = `update feedback set removed = true where id = $1`
	editFeedback   = `update feedback set feedback_topic=$1, feedback_by=$2, feedback_to=$3, feedback_text=$4, score=$5 WHERE  id = $6`
)
