package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goServ5/repository/postgres"
)

var Route *gin.Engine = gin.Default()

func InitRoutes() {

	Route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	Route.POST("/peoples", GetPeoples)
	Route.POST("/peoples/:id", GetPeoplesById)
	Route.POST("/peoples/add", PostPeoples)
	Route.PUT("/peoples", ModifyPeoples)
	Route.DELETE("/peoples/:id", DeletePeoplesById)

	Route.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	err := Route.Run("localhost:8888")
	if err != nil {
		fmt.Println(err)
		return
	}
}

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
	postgres.GetAll(ctx)
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
	postgres.GetById(ctx)
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
	postgres.AddPeople(ctx)
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
	postgres.ModifyOnePeople(ctx)
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
	postgres.DeleteOnePeopleById(ctx)
}
