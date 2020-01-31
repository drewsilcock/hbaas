package handlers

import (
	"fmt"
	"net/http"

	"github.com/drewsilcock/hbaas-server/model"

	"github.com/labstack/echo/v4"
	"github.com/smallnest/gen/dbmeta"
)

func configurePersonRouter(g *echo.Group) {
	g.GET("people", GetAllPeople)
	g.GET("people/:id", GetPerson)
	g.POST("people", AddPerson)
	g.PUT("people/:id", UpdatePerson)
	g.DELETE("people/:id", DeletePerson)
}

// GetAllPeople godoc
// @Summary Get all people.
// @ID get-all-people
// @Tags Person Management
// @Description Retrieve all people present in the databas3.
// @Produce json
// @Success 200 {array} model.Person "Successfully retrieved person records."
// @Failure 400 {object} echo.HTTPError "Invalid request, indicating one of the parameters failed validation."
// @Failure 500 {object} echo.HTTPError "Database not found."
// @Router /people [get]
func GetAllPeople(c echo.Context) error {
	if DB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database not found.")
	}

	var people []model.Person

	query := DB.Model(&model.Person{})
	if err := query.Find(&people).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, people)
}

// GetPerson godoc
// @Summary Get person with specified ID.
// @ID get-person
// @Tags Person Management
// @Description Retrieve person from within database with specified ID.
// @Produce json
// @Success 200 {object} model.Person "Successfully retrieved person record."
// @Failure 404 {object} echo.HTTPError "No person found with this ID."
// @Failure 500 {object} echo.HTTPError "Database not found."
// @Param id path int true "The ID of the person."
// @Router /people/{id} [get]
func GetPerson(c echo.Context) error {
	if DB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database not found.")
	}

	id, err := readPathId(c)
	if err != nil {
		return err
	}
	person := &model.Person{}
	if DB.First(person, id).Error != nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("No person found with ID '%d'.", id),
		)
	}
	return c.JSON(http.StatusOK, person)
}

// AddPerson godoc
// @Summary Add person.
// @ID add-person
// @Tags Person Management
// @Description Add person to system.
// @Produce json
// @Success 200 {object} model.Person "Successfully created person."
// @Failure 400 {object} echo.HTTPError "Invalid person."
// @Failure 500 {object} echo.HTTPError "Database not found."
// @Param person body model.Person true "The person to add."
// @Router /people [post]
func AddPerson(c echo.Context) error {
	if DB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database not found.")
	}

	person := &model.Person{}

	if err := c.Bind(person); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := DB.Save(person).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, person)
}

// UpdatePersonareOrganisation godoc
// @Summary Update person.
// @ID update-person
// @Tags Person Management
// @Description Update person in system.
// @Produce json
// @Success 200 {object} model.Person "Successfully updated person."
// @Failure 400 {object} echo.HTTPError "Invalid updated person."
// @Failure 404 {object} echo.HTTPError "No person found with specified ID."
// @Failure 500 {object} echo.HTTPError "Database not found."
// @Param id path int true "The ID of the person to update."
// @Param person body model.Person true "The updated person."
// @Router /people/{id} [put]
func UpdatePerson(c echo.Context) error {
	if DB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database not found.")
	}

	id, err := readPathId(c)
	if err != nil {
		return err
	}

	person := &model.Person{}
	if DB.First(person, id).Error != nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("Person with ID '%d' not found.", id),
		)
	}

	updated := &model.Person{}
	if err := c.Bind(updated); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := dbmeta.Copy(person, updated); err != nil {
		return err
	}

	if err := DB.Save(person).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, person)
}

// DeletePerson godoc
// @Summary Delete person.
// @ID delete-person
// @Tags Person Management
// @Description Delete person from system.
// @Produce json
// @Success 204 "Sucessfully deleted person."
// @Failure 400 {object} echo.HTTPError "Generic internal error."
// @Failure 404 {object} echo.HTTPError "No person found with specified ID."
// @Failure 500 {object} echo.HTTPError "Database not found."
// @Param id path int true "The ID of the person to delete."
// @Router /people/{id} [delete]
func DeletePerson(c echo.Context) error {
	if DB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database not found.")
	}

	id, err := readPathId(c)
	if err != nil {
		return err
	}
	person := &model.Person{}

	if err := DB.First(person, id).Error; err != nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("Person with ID '%d' not found.", id),
		)
	}
	if err := DB.Delete(person).Error; err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
