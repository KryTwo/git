package handlers

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"goServ5/pkg/structs"
	"goServ5/repository/postgres"
	"log"
	"math"
	"net/http"
)

func GetPeoples(ctx *gin.Context) {

	var json structs.Search
	err := ctx.BindJSON(&json)
	if err != nil {
		log.Fatal(err)
	}

	var lists []structs.People
	var list structs.People
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

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
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

	//filtering
	if filterColumn != "" {
		builder = builder.Where(squirrel.Eq{filterColumn: filterValue})

	}

	req, _, err := builder.ToSql()

	if err != nil {
		fmt.Printf("%v, sql\n", err)
		fmt.Printf("%v", err)

	}
	var rows *sql.Rows
	if filterColumn != "" {
		rows, _ = postgres.Db.Query(req, filterValue)
	} else {
		rows, _ = postgres.Db.Query(req)
	}

	if err != nil {
		fmt.Printf("%s, rows\n", err)
		return

	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}() //rows.close

	for rows.Next() {
		if err := rows.Scan(&list.ID, &list.Last_name, &list.First_name, &list.Middle_name, &list.Address); err != nil {
			log.Fatal(err)
		}
		lists = append(lists, list)
	}

	ctx.IndentedJSON(http.StatusOK, lists)
}

func GetPeoplesById(ctx *gin.Context) {
	id := ctx.Param("id")

	var lists []structs.People
	var list structs.People

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Select(
		"p.id",
		"p.last_name",
		"p.first_name",
		"p.middle_name",
		"r.address").
		From("People AS p").
		Join("registry AS r ON p.id = r.people_id").
		Where("r.people_id = ?")

	stmt, _, err := builder.ToSql()
	if err != nil {
		fmt.Printf("%v,sql", err)
		return
	}

	rows, err := postgres.Db.Query(stmt, id)
	if err != nil {
		fmt.Printf("%v,rows", err)
		return
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}() //rows.close

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
	var newPeople structs.People
	err := ctx.BindJSON(&newPeople)
	if err != nil {
		log.Fatal(err)
		return
	}

	insertPeople := "INSERT INTO people (last_name, first_name, middle_name) VALUES ($1, $2, $3);"
	func() {
		_, err := postgres.Db.Query(insertPeople, newPeople.Last_name, newPeople.First_name, newPeople.Middle_name)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	insertRegistry := "INSERT INTO registry(people_id, address) VALUES ((SELECT max(People.id) FROM People),$1);"
	func() {
		_, err := postgres.Db.Query(insertRegistry, newPeople.Address)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	//GetPeoples(ctx)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "people is added"})
}

func ModifyPeoples(ctx *gin.Context) {
	var changePeopleAddress structs.People
	err := ctx.BindJSON(&changePeopleAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	id := changePeopleAddress.ID
	newAddress := changePeopleAddress.Address

	modifyRegistry := "UPDATE registry r SET address = $1 WHERE people_id = $2;"
	func() {
		_, err := postgres.Db.Query(modifyRegistry, newAddress, id)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	//GetPeoples(ctx)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Address for people_id: " + id + " successfully changed"})

}

func DeletePeoplesById(ctx *gin.Context) {
	id := ctx.Param("id")

	deleteRequest := "DELETE FROM People WHERE id = $1;"
	func() {
		_, err := postgres.Db.Query(deleteRequest, id)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "People is deleted"})
}
