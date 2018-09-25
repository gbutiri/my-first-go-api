package main

import "fmt"
import (
    "net/http"
    "os/exec"
    "runtime"
    "log"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type Message struct {
	Name string
	Body string
	Time int64
}

type User struct {
	Id int
	Name string
	Time int
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_sample")
		
	var (
		id int
		name string
		time int
	)
	
	users := []*User{}
	
	rows, err := db.Query("SELECT id, name, time FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &time)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name, time)
		
		user := new(User)
		user.Id = id
		user.Name = name
		user.Time = time
		
		
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(users)
	if err != nil {
        panic(err)
    }
	fmt.Println(string(b))
    fmt.Fprintf(w, "%s", b)
}
func GetContractUsers(w http.ResponseWriter, r *http.Request) {
	m := Message{
		Name: "Alice", 
		Body: "Hello", 
		Time: 1294706395881547000,
	}
	b, err := json.Marshal(m)
	if err != nil {
        panic(err)
    }
    fmt.Fprintf(w, "%s", b)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	m := Message{
		Name: "Single User", 
		Body: "Welcome back!", 
		Time: 1294706395881547000,
	}
	b, err := json.Marshal(m)
	if err != nil {
        panic(err)
    }
	
	ids, ok := r.URL.Query()["id"]
	if !ok {
        fmt.Fprint(w, "Url Param 'id' is missing")
        return
    }
	key := ids[0]
	fmt.Fprint(w, "Url Param 'id' is " + string(key))
    fmt.Fprintf(w, "%s", b)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "The Docs.\n/getusers -> get all users\n/getuser?id={id} -> get one user based on ID")
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}

func main() {

	db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/golang_sample")
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

    http.HandleFunc("/", myHandler)
    http.HandleFunc("/getusers", GetUsers)
    http.HandleFunc("/getuser", GetUser)
    panic(http.ListenAndServe(":8080", nil))
}