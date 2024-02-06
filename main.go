package main

import (
	"cars/searchbars"
	"cars/structs"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var temp *template.Template
var Car2 = []structs.Models{}

func main() {

	address := ":8088"
	// making one static folder for staic files example: css
	fileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fileHandler)
	// declaring all folders under api as static to load the picture
	fileHandler = http.StripPrefix("/api/", http.FileServer(http.Dir("api")))
	http.Handle("/api/", fileHandler)

	Car2 = searchbars.ModelsElement()

	log.Printf("Starting server on %s", address)

	http.HandleFunc("/", HandleFunc)
	http.HandleFunc("/action", Action)
	http.ListenAndServe(address, nil)

}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	// Only parse the template once
	if temp == nil {
		temp = template.Must(template.ParseFiles("index.html"))
	}
	searchbar := searchbars.FindSearch()
	//car2 := searchbars.ModelsElement()

	manu := map[string]interface{}{
		"search": searchbar,
		"car":    Car2,
	}
	temp.Execute(w, manu)
}

func Action(w http.ResponseWriter, r *http.Request) { // here we both recieve data from server and send to it
	log.Println("POST Input Data...")

	// FormVale gives us a "string"
	cyearStr, err := strconv.Atoi(r.FormValue("cyear"))
	if err != nil {
		http.Error(w, "Invalid cyear parameter", http.StatusBadRequest)
		return
	}

	hpStr, err := strconv.Atoi(r.FormValue("horsepower"))
	if err != nil {
		http.Error(w, "Invalid Hp parameter", http.StatusBadRequest)
	}
	// we load te values we recieved from the sidenavbar
	search := structs.SbarVal2{
		ManuN:   r.FormValue("manufacturer"),
		ManuC:   r.FormValue("country"),
		Cat:     r.FormValue("category"),
		ModName: r.FormValue("cname"),
		Drive:   r.FormValue("drivetrain"),
		Year:    cyearStr,
		Engine:  r.FormValue("engine"),
		Hp:      hpStr,
		Trans:   r.FormValue("transmission"),
	}

	fmt.Println(search)
	temp.Execute(w, nil)
}
