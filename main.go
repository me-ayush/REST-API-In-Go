package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type studentInfo struct {
	Sid    string `json:"sid,omitempty"`
	Name   string `json:"name,omitempty"`
	Course string `json:"course,omitempty"`
}

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/go?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	// _ = db
	ss := []studentInfo{}
	s := studentInfo{}
	rows, err := db.Query("select * from studentinfo")
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&s.Sid, &s.Name, &s.Course)
			ss = append(ss, s)
		}
		json.NewEncoder(w).Encode(ss)
		// fmt.Fprintf(w, "GET REQUEST: ")
	}
}
func addStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer db.Close()
	s := studentInfo{}
	json.NewDecoder(r.Body).Decode(&s)
	sid, _ := strconv.Atoi(s.Sid)
	query := "insert into studentinfo(sid, name, course) values (?, ?, ?)"
	res, err := db.Exec(query, sid, s.Name, s.Course)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err = res.LastInsertId()
		if err != nil {
			json.NewEncoder(w).Encode("{error: Record not inserted}")
		} else {
			// json.NewEncoder(w).Encode(s)
			json.NewEncoder(w).Encode("Ok")
		}
	}
	// fmt.Fprintf(w, "ADD REQUEST")
}
func updateStudents(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "UPDATE REQUEST")
}
func deleteStudents(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DELETE REQUEST")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students", addStudents).Methods("POST")
	r.HandleFunc("/students/{sid}", updateStudents).Methods("PUT")
	r.HandleFunc("/students/{sid}", deleteStudents).Methods("DELTE")
	http.ListenAndServe(":3000", r)
}
