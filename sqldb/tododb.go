package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
	"github.com/alexeyco/simpletable"
)

func ConnectDB1() *sql.DB {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbPath := home + "/.funk.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InitTable() {

	db := ConnectDB1()
	defer db.Close()

	fmt.Println(("this is db file"))

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		task TEXT NOT NULL,
		status TEXT DEFAULT 'Incomplete'
	)
`)
if err != nil {
	log.Fatal(err)
}
}

func AddTask(task string) {

	db := ConnectDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO tasks(task) VALUES(?)",task)

	if err !=nil{
		log.Fatal(err)
	}
	fmt.Println("task added successfully ","'",task,"'"," \n Command 'funk todo --task'")
}

func ListTasks() {
	db := ConnectDB()
	defer db.Close()
	rows, err := db.Query(`
	SELECT id, DATE(created_at), task, status
	FROM tasks
	ORDER BY status = 'Incomplete' DESC
    `)

	if err !=nil{
		log.Fatal(err)
	}

	defer rows.Close()

	table :=simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter,Text: "ID"},
			{Align: simpletable.AlignCenter,Text: "CREATED DATE"},
			{Align: simpletable.AlignCenter,Text: "TASKS"},
			{Align: simpletable.AlignCenter,Text: "STATUS"},
		},
	}

	for rows.Next() {
		var id int
		var create string
		var task string
		var status string

		rows.Scan(&id, &create,&task,&status)
		r:=[]*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", id)},
			{Text: create},
			{Align: simpletable.AlignCenter, Text: task},
			{Align: simpletable.AlignCenter, Text: status},
		}
		table.Body.Cells =append(table.Body.Cells, r)

	}
	table.SetStyle(simpletable.StyleDefault)

	fmt.Println(table.String())

}

func Delete_task(index int)  {
	db:=ConnectDB()

	defer db.Close()

	task,err := db.Exec("DELETE FROM tasks WHERE id =?",index)

	if err !=nil{
		log.Fatal(err)
	}

	row,err := task.RowsAffected()

	if row == 0 {
		fmt.Println("no task found or already task delete")
		return
	}

	fmt.Println("task delete successfully")
}

func Update_Task(index int)  {
	db:=ConnectDB()

	defer db.Close()

	_,err := db.Exec("UPDATE tasks SET status = 'complete' WHERE id = ?",index)

	if err !=nil{
		log.Fatal(err)
	}

	log.Println("Congratulation for task completed\nchanges save in db \ncommand 'funk todo --task'")
}