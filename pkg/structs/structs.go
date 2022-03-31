package structs

type Search struct {
	Page    int     `json:"page" swaggertype:"integer" example:"1"`
	PerPage int     `json:"perPage" swaggertype:"integer" example:"5"`
	Filters Filters `json:"filters"`
	Sorts   Sorts   `json:"sorts"`
}

type Filters struct {
	Column string `json:"column" swaggertype:"string" example:"last_name"`
	Value  string `json:"value" swaggertype:"string" example:"Pushkin"`
}

type Sorts struct {
	Sort string `json:"sort" swaggertype:"string" example:"p.last_name"`
	Way  string `json:"way" swaggertype:"string" example:"+"`
}

type People struct {
	ID          string `json:"id" swaggertype:"string" `
	Last_name   string `json:"last_name" swaggertype:"string" example:"Kolosov"`
	First_name  string `json:"first_name" swaggertype:"string" example:"Evgenij"`
	Middle_name string `json:"middle_name" swaggertype:"string" example:"Alexandrovich"`
	Address     string `json:"address" swaggertype:"string" example:"Moscow"`
}

type PeopleToAdd struct {
	Last_name   string `json:"last_name" swaggertype:"string" example:"Kolosov"`
	First_name  string `json:"first_name" swaggertype:"string" example:"Evgenij"`
	Middle_name string `json:"middle_name" swaggertype:"string" example:"Alexandrovich"`
	Address     string `json:"address" swaggertype:"string" example:"Moscow"`
}
