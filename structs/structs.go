package structs

type Cat struct {
	Id   int    `JSON:"id"`
	Name string `JSON:"name"`
}

type Manufacturers struct {
	Id      int    `JSON:"id"`
	Name    string `JSON:"name"`
	Country string `JSON:"country"`
}
type Models struct {
	Id             int    `JSON:"id"`
	Name           string `JSON:"name"`
	ManufacturerId int    `JSON:"manufacturerId"`
	CategoryId     int    `JSON:"categoryId"`
	Year           int    `JSON:"year"`
	Specifications Specifications
	Image          string `JSON:"image"`
}
type Specifications struct {
	Engine       string `JSON:"engine"`
	Horsepower   int    `JSON:"horsepower"`
	Transmission string `JSON:"transmission"`
	Drivetrain   string `JSON:"drivetrain"`
}

// mistake should have made a struct like SbarVal2. made this for sidebar values to sort out duplicates.
// should have made it to a split of structs instead struct of splits
type SbarVal struct {
	ManuN   []string
	ManuC   []string
	Cat     []string
	ModName []string
	Drive   []string
	Year    []int
	Engine  []string
	Hp      []int
	Trans   []string
}
type SbarVal2 struct {
	ManuN   string
	ManuC   string
	Cat     string
	ModName string
	Drive   string
	Year    int
	Engine  string
	Hp      int
	Trans   string
}
