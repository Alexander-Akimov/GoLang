package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"

	"github.com/gorilla/mux"
)

type Comment struct {
	Id          int
	Name        string
	Email       string
	CommentText string
}

type Page struct {
	Id         int
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	Comments   []Comment
	//Session    Session
	GUID string
}
type JSONResponse struct {
	Fields map[string]string
}

const (
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "root"
	DBPass  = "root"
	DBDbase = "cms"
	PORT    = ":8080"
	SSLPORT = ":443"
)

var database *sql.DB

func main() {
	dbConn := fmt.Sprintf("%s:%s@/%s", DBUser, DBPass, DBDbase)
	// 	fmt.Println(dbConn)
	db, err := sql.Open("mysql", dbConn)
	//db, err := sql.Open("mysql", "root:root@/cms")
	if err != nil {
		log.Println("Couldn't connect to: " + DBDbase)
		log.Println(err.Error())
	}
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/api/pages", APIPage).Methods("GET").Schemes("https")
	routes.HandleFunc("/api/pages/{guid:[0-9a-zA\\-]+}", APIPage).Methods("GET").Schemes("https")
	routes.HandleFunc("/api/comments", APICommentPost).Methods("POST")
	routes.HandleFunc("/api/comments/{id:[\\w\\d\\-]+}", APICommentPut).Methods("PUT")
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/", RedirIndex)
	routes.HandleFunc("/home", ServeIndex)
	http.Handle("/", routes)
	http.Handle("/js", http.FileServer(http.Dir("./js/")))

	//http.ListenAndServe(PORT, nil)

	// certificates, err := tls.LoadX509KeyPair("alexanderCA.pem", "alexanderCA.key")
	// log.Println(certificates)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// tlsConf := &tls.Config{Certificates: []tls.Certificate{certificates}}
	// tls.Listen("tcp",SSLPORT, tlsConf)

	err = http.ListenAndServeTLS(SSLPORT, "alexanderCA.pem", "alexanderCA.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//  tlsConf := &tls.Config{Certificates: []tls.Certificate{certificates}}

	// ln, err := tls.Listen("tcp", PORT, tlsConf)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer ln.Close()

}

func APIPage(w http.ResponseWriter, r *http.Request) {
	//log.Println("terter")
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(err)
		return
	}
	//fmt.Println(thisPage)
	APIOutput, err := json.Marshal(thisPage)
	//fmt.Println(APIOutput)
	fmt.Printf("%v\n", APIOutput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, APIOutput)
}

func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title, page_content, page_date, page_guid FROM pages ORDER BY ? DESC", "page_date")
	if err != nil {
		//fmt.Fprintln(w, err.Error())
		log.Println(err.Error())
	}
	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Pages)
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT id, page_title, page_content, page_date, page_guid FROM pages WHERE page_guid=?", pageGUID).
		Scan(&thisPage.Id, &thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
	thisPage.Content = template.HTML(thisPage.RawContent)
	//fmt.Println(thisPage)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
		// log.Println("Couldn't get page:", pageGUID)
		log.Println(err.Error())
	}
	comments, err := database.Query("SELECT id, comment_name as Name, comment_email, comment_text FROM comments WHERE page_id=?", thisPage.Id)
	if err != nil {
		log.Print(err)
	}
	for comments.Next() {
		var comment Comment
		comments.Scan(&comment.Id, &comment.Name, &comment.Email, &comment.CommentText)
		thisPage.Comments = append(thisPage.Comments, comment)
	}
	// html := "<html><head><title>" + thisPage.Title +
	// 	"</title></head><body><h1>" + thisPage.Title + "</h1><div>" + thisPage.RawContent + "</div></body></html>"
	// fmt.Fprintln(w, html)

	t, _ := template.ParseFiles("templates/blog.html") //need to handle errors if file
	//isn't accassible
	t.Execute(w, thisPage) //handle errors if referencing struct values that
	// if err != nil {
	// 	log.Println(err.Error())
	// }
}
func APICommentPut(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
	}
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")
	res, err := database.Exec("UPDATE comments SET comment_name=?,comment_email=?, comment_text=? WHERE comment_id=?", name, email, comments, id)
	fmt.Println(res)
	if err != nil {
		log.Println(err.Error())
	}
	var resp JSONResponse
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonResp)
}

func APICommentPost(w http.ResponseWriter, r *http.Request) {
	var commentAdded bool
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")
	res, err := database.Exec("INSERT INTO comments SET comment_name=?, comment_email=?, comment_text=?", name, email, comments)
	if err != nil {
		log.Println(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		commentAdded = false
	} else {
		commentAdded = true
	}
	commentAddedBool := strconv.FormatBool(commentAdded)
	var resp JSONResponse
	resp.Fields["id"] = string(id)
	resp.Fields["added"] = commentAddedBool
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonResp)
}

func (p Page) TruncatedText() template.HTML {
	chars := 0
	for i := range p.Content {
		chars++
		if chars > 150 {
			return template.HTML(p.RawContent[:i] + ` ...`)
		}
	}
	return template.HTML(p.RawContent)
}
