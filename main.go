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

// Defining global variables for templates and model lists
var temp *template.Template
var Car2 = []structs.Models{}
var Search = structs.SbarVal2{}
var CompList = []structs.Models{}

func main() {
	address := ":8081"

	// Initialize the HTTP server multiplexer
	mux := http.NewServeMux()

	//  Serve static files for the frontend and API to recieve pictures and
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))
	mux.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir("api"))))

	// Fetch initial car models data
	Car2 = searchbars.ModelsElement()

	log.Printf("Starting server on %s", address)

	// Register route handlers
	mux.HandleFunc("/", HandleFunc)
	mux.HandleFunc("/sbar", Sbar)
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
	go func() {
		if err := theServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	//the shutdown logic for inside error
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("\nReceived interrupt signal. Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := theServer.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	log.Println("Server has stopped please visit: https://www.youtube.com/watch?v=dQw4w9WgXcQ for more info")

}

// Manufacturers page handler
func ManufactPage(w http.ResponseWriter, r *http.Request) {
	temp = template.Must(template.ParseFiles("frontend/manufacturer.html"))
	manId := r.FormValue("manufacturerId")
	manufacturer := searchbars.FilterManufacturer(manId)
	temp.Execute(w, manufacturer)
}

// Comparison page handler
func Comparepage(w http.ResponseWriter, r *http.Request) {
	temp = template.Must(template.ParseFiles("frontend/comparepage.html"))
	temp.Execute(w, CompList)
}

// Frontpage handler
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

// searchbar handler
func Sbar(w http.ResponseWriter, r *http.Request) {
	log.Println("POST loading search Data...")
	cyearInt, err := parseFormValueInt(r.FormValue("cyear"))
	if err != nil {
		http.Error(w, "Invalid cyear parameter", http.StatusBadRequest)
		return
	}
	hpInt, err := parseFormValueInt(r.FormValue("horsepower"))
	if err != nil {
		http.Error(w, "Invalid Hp parameter", http.StatusBadRequest)
		return
	}

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
	data := map[string]interface{}{
		"search": searchbar,
		"car":    Car2,
	}
	temp.Execute(w, data)
}

// Helper function to parse form values into integers because r.FromValue returns strings
func parseFormValueInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

// Compare handler for adding cars to the comparison list
func Compare(w http.ResponseWriter, r *http.Request) {
	cyearInt, err := parseFormValueInt(r.FormValue("carYear"))
	if err != nil {
		http.Error(w, "Invalid cyear parameter", http.StatusBadRequest)
		return
	}
	hpInt, err := parseFormValueInt(r.FormValue("carHp"))
	if err != nil {
		http.Error(w, "Invalid Hp parameter", http.StatusBadRequest)
		return
	}

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
	data := map[string]interface{}{
		"search": searchbar,
		"car":    Car2,
	}
	temp.Execute(w, data)
}
