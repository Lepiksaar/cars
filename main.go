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
var Search = structs.SbarVal2{}
var CompList = []structs.Models{}

func main() {

	address := ":8081"
	// making one static folder for staic files example: css
	fileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fileHandler)
	// declaring all folders under api as static to load the picture
	fileHandler = http.StripPrefix("/api/", http.FileServer(http.Dir("api")))
	http.Handle("/api/", fileHandler)
	// declaring car2 here because i can change this struct later.
	Car2 = searchbars.ModelsElement()

	log.Printf("Starting server on %s", address)

	http.HandleFunc("/", HandleFunc)
	http.HandleFunc("/action", Action)
	http.HandleFunc("/compare", Compare)
	http.HandleFunc("/comparepage", Comparepage)
	http.HandleFunc("/manufacturer", ManufactPage)
	http.ListenAndServe(address, nil)

}
func ManufactPage(w http.ResponseWriter, r *http.Request) {
	temp = template.Must(template.ParseFiles("manufacturer.html"))

	manId := r.FormValue("manufacturerId")
	manufacturer := searchbars.FilterManufacturer(manId) //add error check
	fmt.Println(manufacturer)
	temp.Execute(w, manufacturer)
}
func Comparepage(w http.ResponseWriter, r *http.Request) {
	temp = template.Must(template.ParseFiles("comparepage.html"))

	temp.Execute(w, CompList)
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	// Only parse the template once
	temp = template.Must(template.ParseFiles("index.html"))
	searchbar := searchbars.FindSearch()

	manu := map[string]interface{}{
		"search": searchbar,
		"car":    Car2,
	}
	temp.Execute(w, manu)
}

func Action(w http.ResponseWriter, r *http.Request) { // here we both recieve data from server and send updated list back.
	log.Println("POST loading search Data...")
	var cyearInt int
	var hpInt int
	var err error

	// FormVale gives us a "string" so we need to convert it into a int
	yearStr := r.FormValue("cyear")
	if yearStr == "" {
		cyearInt = 0
	} else {
		cyearInt, err = strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid cyear parameter", http.StatusBadRequest)
			return
		}
	}
	hpStr := r.FormValue("horsepower")
	if hpStr == "" {
		hpInt = 0
	} else {
		hpInt, err = strconv.Atoi(hpStr)
		if err != nil {
			http.Error(w, "Invalid Hp parameter", http.StatusBadRequest)
			return
		}
	}
	// we load te values we recieved from the sidenavbar
	Search = structs.SbarVal2{
		ManuN:   r.FormValue("manufacturer"),
		ManuC:   r.FormValue("country"),
		Cat:     r.FormValue("category"),
		ModName: r.FormValue("cname"),
		Drive:   r.FormValue("drivetrain"),
		Year:    cyearInt,
		Engine:  r.FormValue("engine"),
		Hp:      hpInt,
		Trans:   r.FormValue("transmission"),
	}
	Car2 := searchbars.FilterSearch(Search, Car2)
	searchbar := searchbars.FindSearch()

	manu := map[string]interface{}{
		"search": searchbar,
		"car":    Car2,
	}
	temp.Execute(w, manu)

}

func Compare(w http.ResponseWriter, r *http.Request) {
	// choices := []structs.SbarVal2{}
	var cyearInt int
	var hpInt int
	var err error

	// FormVale gives us a "string" so we need to convert it into a int
	yearStr := r.FormValue("carYear")
	if yearStr == "" {
		cyearInt = 0
	} else {
		cyearInt, err = strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid cyear parameter", http.StatusBadRequest)
			return
		}
	}
	hpStr := r.FormValue("carHp")
	if hpStr == "" {
		hpInt = 0
	} else {
		hpInt, err = strconv.Atoi(hpStr)
		if err != nil {
			http.Error(w, "Invalid Hp parameter", http.StatusBadRequest)
			return
		}
	}
	// we load te values we recieved from the sidenavbar
	comp := structs.Models{
		Name:  r.FormValue("carName"),
		Year:  cyearInt,
		Image: r.FormValue("carimage"),
		Specifications: structs.Specifications{
			Engine:       r.FormValue("carEngine"),
			Horsepower:   hpInt,
			Transmission: r.FormValue("carTransmission"),
			Drivetrain:   r.FormValue("carDrivetrain"),
		},
	}
	CompList = append(CompList, comp)

	fmt.Println(comp)
	fmt.Println(CompList)

	searchbar := searchbars.FindSearch()

	manu := map[string]interface{}{
		"search": searchbar,
		"car":    Car2,
	}
	temp.Execute(w, manu)
}
