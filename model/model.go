package model

import (
	"database/sql"
	"github.com/MaiaVinicius/wabot/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Project struct {
	ID       int
	Label    string
	SenderID int
	Phone    string
}
type Queue struct {
	ID        int
	SenderID  int
	Message   string
	Phone     string
	Metadata2 string
}

func dbConn() (db *sql.DB) {

	err2 := godotenv.Load()
	if err2 != nil {
		log.Fatal("Error loading .env file")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_SCHEMA")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetProjects() []Project {

	db := dbConn()

	rows, err := db.Query("SELECT p.id, p.label, s.id senderId, s.phone FROM wabot_project p INNER JOIN wabot_sender s ON s.project_id = p.id where s.status_id=1 and s.active=1 and p.active=1")

	if err != nil {
		panic(err.Error())
	}

	var projects []Project
	for rows.Next() {
		var project Project

		err = rows.Scan(&project.ID, &project.Label, &project.SenderID, &project.Phone)

		projects = append(projects, project)
	}

	return projects
}

func GetQueue(senderId int) []Queue {

	db := dbConn()

	stmt, err := db.Prepare("SELECT id, message, phone, metadata2  FROM wabot_queue where active=1 and sender_id=? and send_date<=CURDATE() LIMIT 1")

	rows, err := stmt.Query(senderId)

	if err != nil {
		panic(err.Error())
	}

	var queues []Queue
	for rows.Next() {
		var queue Queue

		err = rows.Scan(&queue.ID, &queue.Message, &queue.Phone, &queue.Metadata2)

		queues = append(queues, queue)
	}

	return queues
}

func RemoveFromQueue(queueId int) {
	db := dbConn()

	stmt, err := db.Prepare("INSERT INTO wabot_sent (id,sender_id, phone, message, price)  (SELECT id,sender_id, phone, message, price FROM wabot_queue where id=?)")

	if err != nil {
		panic(err.Error())
	}

	stmt.Exec(queueId)

	stmtDel, errDel := db.Prepare("DELETE FROM wabot_queue WHERE id=?")

	if errDel != nil {
		panic(errDel.Error())
	}

	stmtDel.Exec(queueId)
}

func InsertResponse(projectId int, phone string, id string, message string, timestamp string) {
	db := dbConn()

	stmt, err := db.Prepare("REPLACE INTO wabot_response (id, phone, message, project_id, datetime) VALUES (?,?,?,?,?)")

	if err != nil {
		panic(err.Error())
	}

	stmt.Exec(id, phone, message, projectId, timestamp)
}

func LogMessage(logType int, message string, projectId int) {
	db := dbConn()

	stmt, err := db.Prepare("INSERT INTO wabot_log (project_id, log_type_id, message) value (?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	stmt.Exec(projectId, logType, message)
}

func GetLastSent(projectId int) string {
	db := dbConn()
	stmt, err := db.Prepare("SELECT datetime FROM wabot_response WHERE project_id=? ORDER BY datetime DESC LIMIT 1")

	rows := stmt.QueryRow(projectId)

	if err != nil {
		panic(err.Error())
	}

	var response lib.Response

	rows.Scan(&response.Datetime)

	return response.Datetime
}

func AddToQueue(Phone string, message string, projectId int) {

}
