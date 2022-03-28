package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"math"
	"net/http"
)

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

func GetPeoples(ctx *gin.Context) {

	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	var json Search
	ctx.BindJSON(&json)

	var lists []People
	var list People
	page := json.Page                                  //pagination
	perPage := uint64(math.Abs(float64(json.PerPage))) //rows per page
	perPageDefault := uint64(5)                        //print to page default
	sort := json.Sorts.Sort                            //sort by column_name
	sortWay := "ASC"                                   //by default from min to max
	orderDefault := "p.id"                             //default order by p.id
	filterColumn := json.Filters.Column
	filterValue := json.Filters.Value

	if json.Sorts.Way == "-" {
		sortWay = "DESC"
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("p.id", "p.last_name", "p.first_name", "p.middle_name", "r.address").
		From("People as p").
		Join("registry as r on r.people_id = p.id")

	//pagination
	if page > 1 {
		offs := page * int(perPageDefault)
		builder = builder.Offset(uint64(offs))
	}
	if perPage > 0 {
		builder = builder.Limit(perPage)
	} else {
		builder = builder.Limit(perPageDefault)
	}

	//sorting
	if sort != "" {
		builder = builder.OrderBy(sort + " " + sortWay)
	} else {
		builder = builder.OrderBy(orderDefault + " " + sortWay)
	}

	//filreting
	if filterColumn != "" {
		builder = builder.Where(sq.Eq{filterColumn: filterValue})

	}

	req, _, err := builder.ToSql()

	if err != nil {
		fmt.Printf("%v, sql\n", err)
		fmt.Printf("%v", err)

	}
	var rows *sql.Rows
	if filterColumn != "" {
		rows, _ = db.Query(req, filterValue)
	} else {
		rows, _ = db.Query(req)
	}

	if err != nil {
		fmt.Printf("%s, rows\n", err)
		return

	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&list.ID, &list.Last_name, &list.First_name, &list.Middle_name, &list.Address)
		lists = append(lists, list)
	}

	ctx.IndentedJSON(http.StatusOK, lists)
}

func GetPeoplesById(ctx *gin.Context) {

	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	id := ctx.Param("id")

	var lists []People
	var list People

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(
		"p.id",
		"p.last_name",
		"p.first_name",
		"p.middle_name",
		"r.address").
		From("People AS p").
		Join("registry AS r ON p.id = r.people_id").
		Where("r.people_id = ?")

	sql_statement, _, err := builder.ToSql()
	if err != nil {
		fmt.Printf("%v,sql", err)
		return
	}

	rows, err := db.Query(sql_statement, id)
	if err != nil {
		fmt.Printf("%v,rows", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&list.ID, &list.Last_name, &list.First_name, &list.Middle_name, &list.Address)
		if err != nil {
			return
		}
		lists = append(lists, list)

	}
	if list.ID == "" {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "People not found"})
	} else {
		ctx.IndentedJSON(http.StatusOK, lists)
	}

}

func PostPeoples(ctx *gin.Context) {
	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	var newPeople People
	ctx.BindJSON(&newPeople)

	insertPeople := "insert into People (last_name, first_name, middle_name) VALUES ($1, $2, $3);"
	db.Query(insertPeople, newPeople.Last_name, newPeople.First_name, newPeople.Middle_name)

	insertRegistry := "insert into registry(people_id, address) values ((select max(People.id) from People),$1);"
	db.Query(insertRegistry, newPeople.Address)

	//GetPeoples(ctx)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "people is added"})
}

func ModifyPeoples(ctx *gin.Context) {
	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	var changePeopleAddress People
	ctx.BindJSON(&changePeopleAddress)

	id := changePeopleAddress.ID
	newAddress := changePeopleAddress.Address

	modifyRegistry := "UPDATE registry r SET address = $1 WHERE people_id = $2;"
	db.Query(modifyRegistry, newAddress, id)
	//GetPeoples(ctx)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Address for people_id: " + id + " successfully changed"})

}

func DeletePeoplesById(ctx *gin.Context) {
	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	id := ctx.Param("id")

	deleteRequest := "DELETE FROM People WHERE id = $1;"
	db.Query(deleteRequest, id)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "People is deleted"})
}

func main() {

	router := gin.Default()
	router.GET("/peoples/", GetPeoples)
	router.GET("/peoples/:id", GetPeoplesById)
	router.POST("/peoples/", PostPeoples)
	router.PUT("/peoples/", ModifyPeoples)
	router.DELETE("/peoples/:id", DeletePeoplesById)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run("localhost:8989")
}
