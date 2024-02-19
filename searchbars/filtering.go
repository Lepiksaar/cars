package searchbars

import (
	"cars/structs"
	"fmt"
	"strconv"
)

func FilterSearch(search structs.SbarVal2, car2 []structs.Models) []structs.Models {
	var filteredCars []structs.Models
	manListID := []int{}
	manu := ManElement()
	category := CatElement()
	catID := 99999
	// because we have info in three different structs we need to go over them
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
	for _, v := range car2 {
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
		switch manufacturer.Id {
		case 1:
			manufacturer.Flag = "https://en.wikipedia.org/wiki/Flag_of_Japan#/media/File:Flag_of_Japan.svg"
		case 2:
		case 3:
		case 4:
		case 5:
		case 6:
		case 7:
		case 8:
		case 9:
		case 10:
		}
	}
	return manufacturer
}
