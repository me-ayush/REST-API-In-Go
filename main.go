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

// Global db variable of pointer type
var db *sql.DB

// structure for storng data
type studentInfo struct {
	Sid    string `json:"sid,omitempty"`
	Name   string `json:"name,omitempty"`
	Course string `json:"course,omitempty"`
}

// Connection With MySQL Database
func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/go?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// GET Request
func getStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB() // Initializing Database
	ss := []studentInfo{}
	s := studentInfo{}

	// Executing Query
	rows, err := db.Query("select * from studentinfo")
	if err != nil {
		// If Error during executing the query
		fmt.Fprintf(w, ""+err.Error())
	} else {
		// Traversing the result and storing
		for rows.Next() {
			rows.Scan(&s.Sid, &s.Name, &s.Course)
			ss = append(ss, s)
		}
		// sending data in JSON format
		json.NewEncoder(w).Encode(ss)
	}
}

// POST Request
func addStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB() // Initializing Database
	defer db.Close()
	s := studentInfo{}

	json.NewDecoder(r.Body).Decode(&s) // Storing JSON Data in our structure
	sid, _ := strconv.Atoi(s.Sid)      // Conerting sid from string to int

	query := "insert into studentinfo(sid, name, course) values (?, ?, ?)"
	res, err := db.Exec(query, sid, s.Name, s.Course)
	if err != nil {
		// If Error during executing the query
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err = res.LastInsertId()
		// If error during insertion
		if err != nil {
			json.NewEncoder(w).Encode("{error: Record not inserted}")
		} else {
			json.NewEncoder(w).Encode("Ok")
		}
	}
}

// UPDATE Request
func updateStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB() // Initializing Database
	defer db.Close()

	s := studentInfo{}
	json.NewDecoder(r.Body).Decode(&s) // Storing JSON Data in our structure

	// Storing sid getting as parameter
	vars := mux.Vars(r)
	sid, _ := strconv.Atoi(vars["sid"])

	query := "update studentinfo set name=?, course=? where sid=?"
	res, err := db.Exec(query, s.Name, s.Course, sid)
	if err != nil {
		// If Error during executing the query
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err = res.RowsAffected()
		//If Error during updating
		if err != nil {
			json.NewEncoder(w).Encode("{error: someting went wrong}")
		} else {
			json.NewEncoder(w).Encode("ok")
		}
	}
}

// DELETE Request
func deleteStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB() // Initializing Database
	defer db.Close()

	// Storing sid getting as parameter
	vars := mux.Vars(r)
	sid, _ := strconv.Atoi(vars["sid"])

	query := "delete from studentinfo  where sid=?"
	res, err := db.Exec(query, sid)
	if err != nil {
		//If Error during executing the query
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err = res.RowsAffected()

		// If Error during Deleting
		if err != nil {
			json.NewEncoder(w).Encode("{error: someting went wrong}")
		} else {
			json.NewEncoder(w).Encode("ok")
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students", addStudents).Methods("POST")
	r.HandleFunc("/students/{sid}", updateStudents).Methods("PUT")
	r.HandleFunc("/students/{sid}", deleteStudents).Methods("DELETE")
	http.ListenAndServe(":3000", r)
}
