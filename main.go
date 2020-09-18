package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"io/ioutil"
	"fmt"
)

type Student struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB
var err error

var students []Student

func main(){

	db, err = sql.Open("mysql", "<name>:<password>@tcp(127.0.0.1:3306)/studentsdb")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router:= mux.NewRouter()


	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/students", createStudent).Methods("POST")
	router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	http.ListenAndServe(":8080", router)

}

func getStudents(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var students []Student
	
	result, err := db.Query("SELECT id, name FROM students")
	if err != nil{
		panic(err.Error())
	}

	defer result.Close()

	for result.Next(){
		var student Student
		err := result.Scan(&student.ID, &student.Name)
		if err != nil {
			panic(err.Error())
		}
		students = append(students,student)
	}
	json.NewEncoder(w).Encode(students)
}

func createStudent(w http.ResponseWriter, r *http.Request){
	stmt, err := db.Prepare("INSERT INTO students(name) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal)
		name := keyVal["name"]

		_, err = stmt.Exec(name)
		if err != nil {
			panic(err.Error())
		}

		fmt.Fprintf(w, "New student was created")
}

func getStudent(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type", "application/json")
	params:=mux.Vars(r)

	result, err := db.Query("SELECT id, name FROM students WHERE id = ?", params["id"])
		if err != nil{
			panic(err.Error())
		}

		defer result.Close()
		var student Student
		for result.Next(){
			err := result.Scan(&student.ID, &student.Name)
			if err != nil {
				panic(err.Error())
			}
		}

	json.NewEncoder(w).Encode(student)
}
	

func updateStudent(w http.ResponseWriter, r *http.Request){
	params:=mux.Vars(r)

	//1. create a query statement
	stmt, err := db.Prepare("UPDATE students SET name = ? WHERE id = ? ")
		if err != nil {
			panic(err.Error())
		}
	//2. set up the body for request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	//decode json body and store map referenced 
	json.Unmarshal(body, &keyVal)
	newName := keyVal["name"]
	//print &keyVal vs keyVal
	
	//execute statement with the params
	_, err = stmt.Exec(newName, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Name with ID = %s was updated", params["id"])
}

func deleteStudent(w http.ResponseWriter, r *http.Request){
	
	params:=mux.Vars(r)

	//1. create statement
	stmt, err := db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "student with ID = %s was deleted", params["id"])
	
	


}