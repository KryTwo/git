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
	"strconv"
)

// GetPeoples godoc
// @Summary      Show all people
// @Tags Peoples
// @Description  Show people with sorting and filtering
// @Accept json
// @Produce json
// @Param input body structs.Search true "search val"
// @Success 200 {array} structs.People
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /peoples [post]
func GetPeoples(ctx *gin.Context) {

	var json structs.Search
	err := ctx.BindJSON(&json)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
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
	builder := psql.Select(
		"p.id",
		"p.last_name",
		"p.first_name",
		"p.middle_name",
		"r.address").
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
		rows, err = postgres.Db.Query(req, filterValue)
	} else {
		rows, err = postgres.Db.Query(req)
	}

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
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

	if lists == nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, lists)
}

// GetPeoplesById godoc
// @Summary      Show People By ID
// @Tags Peoples
// @Description  Show One People
// @Produce json
// @Param people_id path int true "people_id"
// @Success 200 {object} structs.People
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /peoples/{people_id} [post]
func GetPeoplesById(ctx *gin.Context) {

	// take people_id from param and check To Integer Type
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

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
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
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

// PostPeoples godoc
// @Summary      Add People
// @Tags Peoples
// @Description  Add one people
// @Accept json
// @Produce json
// @Param input body structs.PeopleToAdd true "post values"
// @Success 201 {object} structs.PeopleToAdd
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /peoples/add [post]
func PostPeoples(ctx *gin.Context) {
	var newPeople structs.PeopleToAdd
	err := ctx.BindJSON(&newPeople)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	insertPeople := "INSERT INTO people (last_name, first_name, middle_name) VALUES ($1, $2, $3);"
	func() {
		_, err := postgres.Db.Query(insertPeople, newPeople.Last_name, newPeople.First_name, newPeople.Middle_name)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	insertRegistry := "INSERT INTO registry(people_id, address) VALUES ((SELECT max(People.id) FROM People),$1);"
	func() {
		_, err := postgres.Db.Query(insertRegistry, newPeople.Address)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	ctx.IndentedJSON(http.StatusCreated, newPeople)
}

// ModifyPeoples godoc
// @Summary      Modify People
// @Tags Peoples
// @Description  Modify People
// @Accept json
// @Produce json
// @Param input body structs.People true "post values"
// @Success 200 {object} structs.People{last_name=string,first_name=string,middle_name=string}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /peoples [put]
func ModifyPeoples(ctx *gin.Context) {

	// response struct

	//into modifyPeople we pass JSON body (method [PUT])
	var modifyPeople structs.People
	err := ctx.BindJSON(&modifyPeople)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := strconv.Atoi(modifyPeople.ID); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	//use squirrel to create sql Query
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	var builder squirrel.UpdateBuilder

	// update table Registry, if address in request JSON body not nil
	if modifyPeople.Address != "" {
		builder = psql.Update("registry").
			Set("address", modifyPeople.Address).
			Where("people_id = ?")
		req, _, err := builder.ToSql()
		if err != nil {
			fmt.Println(err)
			return
		}
		func() {
			_, err := postgres.Db.Query(req, modifyPeople.Address, modifyPeople.ID)
			if err != nil {
				newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}()
	}

	// update table People apart from Registry
	builder = psql.Update("people")

	if modifyPeople.Last_name != "" {
		builder = builder.Set("last_name", modifyPeople.Last_name)
	}

	if modifyPeople.First_name != "" {
		builder = builder.Set("first_name", modifyPeople.First_name)
	}

	if modifyPeople.Middle_name != "" {
		builder = builder.Set("middle_name", modifyPeople.Middle_name)
	}

	builder = builder.Where(squirrel.Eq{"id": modifyPeople.ID})

	req, arg, err := builder.ToSql()

	// send request to sql
	func() {
		_, err := postgres.Db.Query(req, arg...)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	// result

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "changes are successful"})

}

// DeletePeoplesById godoc
// @Summary      DeletePeoplesById
// @Tags Peoples
// @Description  DeletePeoplesById
// @Produce json
// @Param people_id path int true "people_id"
// @Success 200 {string} string "people is delete"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /peoples/{people_id} [delete]
func DeletePeoplesById(ctx *gin.Context) {
	// take people_id from param and check To Integer Type
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	deleteRequest := "DELETE FROM People WHERE id = $1;"
	func() {
		_, err := postgres.Db.Query(deleteRequest, id)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}()
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "People is deleted"})
}
