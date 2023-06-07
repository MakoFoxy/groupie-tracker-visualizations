package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

const (
	HomePage     string = "assets/html/mainpage.html"
	ArtistPage   string = "assets/html/contentpage.html"
	ErrorPage    string = "assets/html/error.html"
	TemplatePage string = "assets/html/templates.html"
)

type Artists struct {
	Id             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	Locations      string   `json:"locations"`
	ConcertDates   string   `json:"concertDates"`
	Relations      string   `json:"relations"`
	DatesLocations map[string][]string
}

type RelationsStruct struct {
	Index []struct {
		Id             int64               `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type ErrorStruct struct {
	Text   string
	Status int
}

func main() {
	Handler()
}

func Handler() {
	mux := http.NewServeMux()
	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	mux.HandleFunc("/", firstHandler)
	mux.HandleFunc("/artist/", postHandler)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	fmt.Println("Server is listening...")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal("ERROR:Server not listening")
	}
}

var resultArtists []Artists

func firstHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, 404) // 404 Not Found
		return
	}
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles(HomePage, TemplatePage)
		if err != nil {
			errorHandler(w, 500) // 500 Internal Server Error
			return
		}

		FuncHomeArtists(&resultArtists)

		err = tmpl.Execute(w, resultArtists)
		if err != nil {
			errorHandler(w, 500) // 500 Internal Server Error
			return
		}
	} else {
		errorHandler(w, 405)
		return
	}
}

var contentArtists RelationsStruct

func postHandler(w http.ResponseWriter, r *http.Request) {
	// path.Base()
	url := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(url[2])
	if r.URL.Path != "/artist/"+url[2] {
		errorHandler(w, 404) // 404 Not Found
		return
	}
	if err != nil || len(url) > 3 || id > len(resultArtists) || id < 1 {
		errorHandler(w, 404) // 404 Not Found
		return
	}

	if r.Method != http.MethodGet {
		errorHandler(w, 405) // 405 Method Not Allowed
		return
	}
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles(ArtistPage, TemplatePage)
		if err != nil {
			errorHandler(w, 500) // 500 Internal Server Error
			return
		}

		FuncContentArtists(&contentArtists)
		// res := RelationsData{DataRel: contentArtists}
		// res := PeopleArtists{PeopleNumArtists: resultArtists}
		err = tmpl.Execute(w, resultArtists[id-1])
		if err != nil {
			errorHandler(w, 500) // 500 Internal Server Error
			return
		}
	}
}

func FuncHomeArtists(artists *[]Artists) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonErr := json.Unmarshal(body, &artists)
	// fmt.Println(artists)

	if jsonErr != nil {
		log.Fatal(jsonErr)
		return
	}
}

func FuncContentArtists(relations *RelationsStruct) {
	url := "https://groupietrackers.herokuapp.com/api/relation"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonErr := json.Unmarshal(body, &relations)

	if jsonErr != nil {
		log.Fatal(jsonErr)
		return
	}
	// fmt.Println(relations)

	for i := range resultArtists {
		resultArtists[i].DatesLocations = contentArtists.Index[i].DatesLocations
	}
	// fmt.Println(resultArtists)
}

func errorHandler(w http.ResponseWriter, status int) {
	Res := ErrorStruct{
		Status: status,
		Text:   http.StatusText(status),
	}
	tmpl, err := template.ParseFiles(ErrorPage, TemplatePage)
	if err != nil {
		// http.Error()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(status)
	err = tmpl.Execute(w, Res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}
