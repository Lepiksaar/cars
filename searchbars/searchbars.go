package searchbars

import (
	"cars/structs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// populating structs from json. when using interface it does not return anything concrete so that why we have to use
// extra three functions to call the populate the structs. meaning we cannot call fetchData directly.
func fetchData(url string, v interface{}) {
	looking, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer looking.Body.Close()

	if looking.StatusCode == http.StatusOK {
		json.NewDecoder(looking.Body).Decode(v)
	}
}
func ManElement() []structs.Manufacturers {
	var item []structs.Manufacturers
	fetchData("http://localhost:3000/api/manufacturers", &item)
	return item
}

func CatElement() []structs.Cat {
	var item []structs.Cat
	fetchData("http://localhost:3000/api/categories", &item)
	return item
}

func ModelsElement() []structs.Models {
	var item []structs.Models
	fetchData("http://localhost:3000/api/models", &item)
	return item
}

// this function removes double values for searchbar.
func FindSearch() structs.SbarVal {
	uniqueCount := make(map[string]struct{})
	s := structs.SbarVal{}

	// method for checking uniques string values
	checkAndAppend := func(key string, slice *[]string) {
		strKey := fmt.Sprint(key)
		if _, exists := uniqueCount[strKey]; !exists {
			uniqueCount[strKey] = struct{}{}
			*slice = append(*slice, key)
		}
	}
	// method for checking uniques int values
	checkAndAppend2 := func(key int, slice *[]int) {
		strKey := fmt.Sprint(key)
		if _, exists := uniqueCount[strKey]; !exists {
			uniqueCount[strKey] = struct{}{}
			*slice = append(*slice, key)
		}
	}
	// now runnig all searchfields and removing duplicates
	models := ModelsElement()
	for _, v := range models {

		checkAndAppend2(v.Year, &s.Year)
		checkAndAppend(v.Specifications.Engine, &s.Engine)
		checkAndAppend2(v.Specifications.Horsepower, &s.Hp)
		checkAndAppend(v.Specifications.Transmission, &s.Trans)
		checkAndAppend(v.Specifications.Drivetrain, &s.Drive)

		s.ModName = append(s.ModName, v.Name)
	}

	manu := ManElement()
	for _, v := range manu {
		checkAndAppend(v.Country, &s.ManuC)
		checkAndAppend(v.Name, &s.ManuN)
	}
	category := CatElement()
	for _, v := range category {
		checkAndAppend(v.Name, &s.Cat)
	}
	return s
}
