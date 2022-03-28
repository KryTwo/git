package structs

type Search struct {
	Page    int     `json:"page"`
	PerPage int     `json:"perPage"`
	Filters Filters `json:"filters"`
	Sorts   Sorts   `json:"sorts"`
}

type Filters struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}

type Sorts struct {
	Sort string `json:"sort"`
	Way  string `json:"way"`
}

type People struct {
	ID          string `json:"id"`
	Last_name   string `json:"last_name"`
	First_name  string `json:"first_name"`
	Middle_name string `json:"middle_name"`
	Address     string `json:"address"`
}
