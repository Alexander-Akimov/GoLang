package main

// generate certificate on windows 8
// 	http://www.faqforge.com/windows/use-openssl-on-windows/
// 		https://code.google.com/archive/p/openssl-for-windows/downloads

// 	1) Generating key
// 	genrsa -out key.pem
// 	genrsa -des3 -out server.key 4096
// 	genrsa -out server.key 2048

// 	2) Certificate request
// 		req -new -key key.pem -out cert.pem -config C:\openssl.cnf
// 		req -out cert.csr -key server.key -new -sha256 -config C:\openssl.cnf

// 	Unable to load config info from c:openssl/ssl/openssl.cnf
// 		http://stackoverflow.com/questions/14459078/unable-to-load-config-info-from-usr-local-ssl-openssl-cnf
// 	country name two letter code
// 		https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2

// 	3) sign the certificate request
// 		req -x509 -days 365 -key key.pem -in cert.pem -out certificate.pem -config C:\openssl.cnf
// 		req -x509 -days 365 -key server.key -in cert.csr -out certificate.crt -config C:\openssl.cnf
// 		req -x509 -days 365 -key www.mydomain.key -in www.mydomain.com.sha256.csr -out certificate.crt -config C:\openssl.cnf

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"

	"github.com/gorilla/mux"
)

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

type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	GUID       string
}

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
	routes.HandleFunc("/api/pages", APIPage)                      ///.Methods("GET").Schemes("https")
	routes.HandleFunc("/api/pages/{guid:[0-9a-zA\\-]+}", APIPage) //.Methods("GET").Schemes("https")
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/", RedirIndex)
	routes.HandleFunc("/home", ServeIndex)
	http.Handle("/", routes)
	//http.ListenAndServe(PORT, nil)

	err = http.ListenAndServeTLS(SSLPORT, "certificate.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	// certificates, err := tls.LoadX509KeyPair("certificate.crt", "server.key")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// tlsConf := &tls.Config{Certificates: []tls.Certificate{certificates}}

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
		fmt.Fprintln(w, err.Error)
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
	err := database.QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid=?", pageGUID).
		Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
		// log.Println("Couldn't get page:", pageGUID)
		log.Println(err.Error())
	}

	// html := `<html><head><title>` + thisPage.Title +
	// 	`</title></head><body><h1>` + thisPage.Title + `</h1><div>` +
	// 	thisPage.Content + `</div></body></html>`
	// fmt.Fprintln(w, html)
	t, _ := template.ParseFiles("templates/blog.html") //need to handle errors if file
	//isn't accassible
	t.Execute(w, thisPage) //handle errors if referencing struct values that
}

func (p Page) TruncatedText() template.HTML {
	chars := 0
	for i, _ := range p.Content {
		chars++
		if chars > 150 {
			return template.HTML(p.RawContent[:i] + ` ...`)
		}
	}
	return template.HTML(p.RawContent)
}
