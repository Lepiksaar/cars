package searchbars

import (
	"cars/structs"
	"fmt"
	"strconv"
)

// This function filteres out model structs that are not looked for.
func FilterSearch(search structs.SbarVal2, car2 []structs.Models) []structs.Models {
	var filteredCars []structs.Models
	var manListID []int
	manu := ManElement()
	category := CatElement()
	catID := 99999

	// because we have info in three different structs we need to go over them.
	// first we check if category was specified
	for _, v := range category {
		if search.Cat == "" {
			catID = 0
		} else if v.Name == search.Cat {
			catID = v.Id
		}
	}

	// we check if recieved any preferences from client
	for _, v := range manu {
		if (v.Name == search.ManuN || search.ManuN == "") && (v.Country == search.ManuC || search.ManuC == "") {
			manListID = append(manListID, v.Id)
		}
	}

	// now can go through Car2 list to add all models that are in criteria
	for _, v := range car2 {
		// we compare all of the values we can directly compare against each other
		if (v.CategoryId == catID || catID == 0) &&
			(v.Name == search.ModName || search.ModName == "") &&
			(v.Specifications.Engine == search.Engine || search.Engine == "") &&
			(v.Specifications.Transmission == search.Trans || search.Trans == "") &&
			(v.Specifications.Drivetrain == search.Drive || search.Drive == "") &&
			(v.Year == search.Year || search.Year == 0) &&
			(v.Specifications.Horsepower == search.Hp || search.Hp == 0) {
			filteredCars = append(filteredCars, v)
		}
	}

	// we have to do extra check for manufacturers, because there can be more than one id element.
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

// Looks up for required manufacturers info from manufacturers struct and adds wikipedia page to struct
func FilterManufacturer(needInt string) structs.Manufacturers {
	manufacturer := structs.Manufacturers{}
	loadstruct := ManElement()
	needInt2, err := strconv.Atoi(needInt)

	if err != nil {
		// Handle the error if the conversion failed
		fmt.Println("Error during manufacturer id conversion:", err)
	}

	for i, v := range loadstruct {
		if v.Id == needInt2 {
			manufacturer = loadstruct[i]
		}
	}
	// we dont need it, but why not :)
	switch manufacturer.Id {
	case 1:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Toyota"
	case 2:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Honda"
	case 3:
		manufacturer.Info = "https://en.wikipedia.org/wiki/BMW"
	case 4:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Audi"
	case 5:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Mercedes-Benz"
	case 6:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Ford"
	case 7:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Chevrolet"
	case 8:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Hyundai"
	case 9:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Lexus"
	case 10:
		manufacturer.Info = "https://en.wikipedia.org/wiki/Nissan"
	}
	return manufacturer
}
