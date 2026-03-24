package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

var colorSuccess = color.New(color.FgGreen)

func ConnectDB() *sql.DB {
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

func Create_db() {
	db :=ConnectDB()
	defer db.Close()

	create := ` CREATE TABLE IF NOT EXISTS Timer (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		create_at DATE DEFAULT (DATE('now')),
        timer TEXT,
		name TEXT
    );`

	_,err := db.Exec(create)

	if err !=nil{
		log.Fatal(err)
	}
}

func Insert_data(h int,m int,s int,task string) string {
	db := ConnectDB()

	defer db.Close()

	timer := fmt.Sprintf("%02d:%02d:%02d", h, m, s)

	_,err :=db.Exec("INSERT INTO Timer(timer,name) VALUES(?,?)",timer,task)

	if err !=nil{
		log.Fatal(err)
		
	}

	return  "\ninsert timer data in db successfully \n Command 'funk timer --his' "
}

func Show_history()  {
	db := ConnectDB()
	defer db.Close()

	row,err := db.Query("SELECT id, DATE(create_at), timer, name FROM Timer ORDER BY id DESC")

	if err !=nil{
		log.Fatal(err)
	}

	defer row.Close()

	table :=simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter,Text: "ID"},
			{Align: simpletable.AlignCenter,Text: "CREATE DATE"},
			{Align: simpletable.AlignCenter,Text: "TIMER"},
			{Align: simpletable.AlignCenter,Text: "NAME"},
		},	
	}

	for row.Next(){
		var id int
		var created string
		var timer string
		var name string
		row.Scan(&id,&created,&timer,&name)
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", id)},
			{Text: created},
			{Align: simpletable.AlignCenter, Text: timer},
			{Align: simpletable.AlignCenter, Text: name},
		}
		table.Body.Cells = append(table.Body.Cells,r)
	}
	table.SetStyle(simpletable.StyleDefault)

	fmt.Println(table.String())
}

func  Delete_Record(index int)  {
	db := ConnectDB()
	defer db.Close()

	record,err:=db.Exec("DELETE FROM Timer WHERE id = ?",index)

	if err !=nil{
		log.Fatal(err)
	}
	row,err  := record.RowsAffected()

	if err !=nil{
		log.Fatal(err)
	}

	if row ==0 {
		fmt.Println("no record found or record already delete")
		return
	}

	colorSuccess = color.New(color.FgRed)

	colorSuccess.Println("record delete successfully")
}

func  Delete_All_Record()  {
	db := ConnectDB()
	defer db.Close()

	record,err:=db.Exec("DELETE FROM Timer")

	if err !=nil{
		log.Fatal(err)
	}
	row,err  := record.RowsAffected()

	if err !=nil{
		log.Fatal(err)
	}

	if row ==0 {
		fmt.Println("no record found or record already delete")
		return
	}

	colorSuccess = color.New(color.FgRed)
	colorSuccess.Println("All Record delete successfully")
}
