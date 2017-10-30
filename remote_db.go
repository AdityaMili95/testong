package main

import "encoding/json"
import "net/http"
import "fmt"
import "html/template"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type student struct {
	ID    string
	Name  string
	Grade int
}

type testData struct {
	ID   int
	Name string
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "u512969712_prk:'cari cari remote tv 52528'@tcp(mysql.hostinger.co.id:3306)/u512969712_af")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return db, nil
}

var data = []student{
	student{"E001", "ethan", 21},
	student{"W001", "wick", 22},
	student{"B001", "bourne", 23},
	student{"B002", "bond", 23},
}

func getDataRemote(w http.ResponseWriter, r *http.Request){
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("select Word from words")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var text string
		var err = rows.Scan(&text)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(text)
	}

}

func getData(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	//var age = 27
	rows, err := db.Query("select ID,Name from test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}

	tempt, err := template.New("template.html").Funcs(funcMap).ParseFiles("template.html") //template.ParseFiles("template.html").Funcs(funcMap)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println("Success")

	var data []testData

	for rows.Next() {
		var each = testData{}
		var err = rows.Scan(&each.ID, &each.Name)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		data = append(data, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	tempt.Execute(w, map[string]interface{}{
		"row": data,
	})
	return

}

func users(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {

		tempt, err := template.ParseFiles("template.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		tempt.Execute(w, map[string]interface{}{
			"data": data,
		})
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var id = r.FormValue("id")
		var result []byte
		var err error

		for _, each := range data {
			if each.ID == id {
				result, err = json.Marshal(each)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Write(result)
				return
			}
		}

		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/users", users)
	http.HandleFunc("/user", user)
	http.HandleFunc("/db", getData)
	http.HandleFunc("/remote", getDataRemote)

	fmt.Println("starting web server at http://localhost:8088/")
	http.ListenAndServe(":8088", nil)
}
