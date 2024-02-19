package main

import (
	"cars/searchbars"
	"cars/structs"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var temp *template.Template
var Car2 = []structs.Models{}
var Search = structs.SbarVal2{}
var CompList = []structs.Models{}

func main() {

	address := ":8081"

	// declaring multiplex for server
	mux := http.NewServeMux()
	// making one static folder for staic files example: css
	fileHandler := http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend")))
	mux.Handle("/frontend/", fileHandler)
	// declaring all folders under api as static to load the picture
	fileHandler = http.StripPrefix("/api/", http.FileServer(http.Dir("api")))
	mux.Handle("/api/", fileHandler)
	// declaring car2 here because i can change this struct later.
	Car2 = searchbars.ModelsElement()

	log.Printf("Starting server on %s", address)

	mux.HandleFunc("/", HandleFunc)
	mux.HandleFunc("/action", Action)
	mux.HandleFunc("/compare", Compare)
	mux.HandleFunc("/comparepage", Comparepage)
	mux.HandleFunc("/manufacturer", ManufactPage)

	//setting up the server parameters
	theServer := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// running the server async in another channel
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := theServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-sigCh
	log.Println("\nReceived interrupt signal. Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := theServer.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	log.Println("Server has stopped please visit: https://www.youtube.com/watch?v=dQw4w9WgXcQ for more info")

}
func ManufactPage(w http.ResponseWriter, r *http.Request) {
	temp = template.Must(template.ParseFiles("frontend/manufacturer.html"))

	manId := r.FormValue("manufacturerId")
	manufacturer := searchbars.FilterManufacturer(manId)

	temp.Execute(w, manufacturer)
}
func Comparepage(w http.ResponseWriter, r *http.Request) {
	temp = template.Must(template.ParseFiles("frontend/comparepage.html"))

	temp.Execute(w, CompList)
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	// Only parse the template once
	temp = template.Must(template.ParseFiles("frontend/index.html"))
	searchbar := searchbars.FindSearch()
	// we create interface to post two different structs at the same go
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

	// FormValue gives us a "string" so we need to convert it into a int
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
	// we load te values we recieved from the compare buttons
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

	searchbar := searchbars.FindSearch()

	manu := map[string]interface{}{
		"search": searchbar,
		"car":    Car2,
	}
	temp.Execute(w, manu)
}
