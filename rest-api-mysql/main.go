//go mod init example/rest-api-mysql

package main
import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
  "net/http" 
)

//Defining our models (structs)
type Post struct {
  ID string `json:"id"`
  Title string `json:"title"`
}

var db *sql.DB
var err error
func main() {
   //Creating the database object 
  //db, err = sql.Open("mysql", "<user>:<password>@tcp(127.0.0.1:3306)/<dbname>")
  db ,err = sql.Open("mysql", "root:@/test_go")
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()

  // Initialise the router 
  router := mux.NewRouter()

    //Creating our endpoints
  router.HandleFunc("/posts", getPosts).Methods("GET") 
  router.HandleFunc("/posts", createPost).Methods("POST") 
  router.HandleFunc("/posts/{id}", getPost ).Methods("GET") 
  router.HandleFunc("/posts/{id}", updatePost).Methods("PUT") 
  router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

  //run our server on port 8000
  http.ListenAndServe(":8000", router)
}
    //getPosts()
    func getPosts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var posts []Post
    result, err := db.Query("SELECT id, title from posts")
    if err != nil {
        panic(err.Error())
    }
    defer result.Close()
    for result.Next() {
        var post Post
        err := result.Scan(&post.ID, &post.Title)
        if err != nil {
        panic(err.Error())
        }
        posts = append(posts, post)
    }
    json.NewEncoder(w).Encode(posts)
    }

    //createPost()
    func createPost(w http.ResponseWriter, r *http.Request) {
    stmt, err := db.Prepare("INSERT INTO posts(title) VALUES(?)")
    if err != nil {
        panic(err.Error())
    }
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err.Error())
    }
    keyVal := make(map[string]string)
    json.Unmarshal(body, &keyVal)
    title := keyVal["title"]
    _, err = stmt.Exec(title)
    if err != nil {
        panic(err.Error())
    }
    fmt.Fprintf(w, "New post was created")
    
    }   

    //getPost()
    func getPost(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    result, err := db.Query("SELECT id, title FROM posts WHERE id = ?", params["id"])
    if err != nil {
        panic(err.Error())
    }
    defer result.Close()
    var post Post
    for result.Next() {
        err := result.Scan(&post.ID, &post.Title)
        if err != nil {
        panic(err.Error())
        }
    }
    json.NewEncoder(w).Encode(post)
    }

    //updatePost()
    func updatePost(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    stmt, err := db.Prepare("UPDATE posts SET title = ? WHERE id = ?")
    if err != nil {
        panic(err.Error())
    }
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err.Error())
    }
    keyVal := make(map[string]string)
    json.Unmarshal(body, &keyVal)
    newTitle := keyVal["title"]
    _, err = stmt.Exec(newTitle, params["id"])
    if err != nil {
        panic(err.Error())
    }
    fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
    }

    //deletePost()
    func deletePost(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    stmt, err := db.Prepare("DELETE FROM posts WHERE id = ?")
    if err != nil {
        panic(err.Error())
    }
    _, err = stmt.Exec(params["id"])
    if err != nil {
        panic(err.Error())
    }
    fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
    }