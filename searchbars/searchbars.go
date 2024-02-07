package searchbars

import (
	"cars/structs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// populating structs from json. when using interface it does not return anything concrete so that why we have to use
// extra three functions to call the populate the structs
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
func FilterSearch(search structs.SbarVal2, car2 []structs.Models) []structs.Models {
	var filteredCars []structs.Models
	manListID := []int{}
	manu := ManElement()
	category := CatElement()
	catID := 99999

	for _, v := range category {
		if search.Cat == "" {
			catID = 0
		} else if v.Name == search.Cat {
			catID = v.Id
		}
	}
	for _, v := range manu {
		if (v.Name == search.ManuN || search.ManuN == "") && (v.Country == search.ManuC || search.ManuC == "") {
			manListID = append(manListID, v.Id)
		}
	}
	fmt.Println(catID)
	for _, v := range car2 {
		//fmt.Printf(" ------>%v<--------", search.Year)
		// we compare all of the values we can directly compare against each other
		if (v.CategoryId == catID || catID == 0) && (v.Name == search.ModName || search.ModName == "") && (v.Specifications.Engine == search.Engine || search.Engine == "") && (v.Specifications.Transmission == search.Trans || search.Trans == "") && (v.Specifications.Drivetrain == search.Drive || search.Drive == "") && (v.Year == search.Year || search.Year == 0) && (v.Specifications.Horsepower == search.Hp || search.Hp == 0) {
			filteredCars = append(filteredCars, v)
		}
	}
	// second filter, of info from manufacturers struct. we do it backwards to make sure not to get out of list errors(we are changing the same list so we shorten the list)
	for i := len(filteredCars) - 1; i >= 0; i-- {
		found := false
		for _, c := range manListID {
			if filteredCars[i].ManufacturerId == c {
				found = true
				break
			}
		}
		if !found {
			filteredCars = append(filteredCars[:i], filteredCars[i+1:]...)
		}
	}

	return filteredCars
}
