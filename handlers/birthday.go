package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/bearbin/go-age"
	"github.com/labstack/echo/v4"

	"github.com/drewsilcock/hbaas-server/model"
)

// The exact date here isn't important, it just specifies the format.
const dateInputForm = "Jan-02"

func configureBirthdayRouter(g *echo.Group) {
	g.GET("birthday/to-name/:name", SayHappyBirthdayToName)
	g.GET("birthday/to-person/:name", SayHappyBirthdayToPerson)
	g.GET("birthday/by-date/:date", SayHappyBirthdayByDate)
}

// SayHappyBirthdayToName godoc
// @Summary Say happy birthday to a particular name.
// @ID to-name
// @Tags Saying Happy Birthday
// @Security OAuth2Implicit
// @Description Say happy birthday to a particular name.
// @Produce text/plain
// @Param name path string true "The name of the person to say happy birthday to."
// @Success 200 {string} string "Said happy birthday to requested person."
// @Router /birthday/to-name/{name} [get]
func SayHappyBirthdayToName(c echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, fmt.Sprintf("Happy birthday %s!", name))
}

// SayHappyBirthdayToPerson godoc
// @Summary Say happy birthday to a particular person.
// @ID to-person
// @Tags Saying Happy Birthday
// @Security OAuth2Implicit
// @Description Say happy birthday to a particular name.
// @Produce text/plain
// @Param name path string true "The ID of the person to say happy birthday to."
// @Success 200 {string} string "Said happy birthday to person with specified ID."
// @Failure 404 {object} echo.HTTPError "Person not found."
// @Failure 500 {object} echo.HTTPError "Database not found."
// @Router /birthday/to-person/{name} [get]
func SayHappyBirthdayToPerson(c echo.Context) error {
	urlDecode, err := urlDecodeParam(c, "name")
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Unable to decode person's name: %s", err),
		)
	}
	name := strings.ToLower(urlDecode)

	if DB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database not found.")
	}

	person := &model.Person{}
	if DB.Debug().Where("LOWER(name) = ?", name).First(person).Error != nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("No people found with name '%s'.", name),
		)
	}

	personAge := age.Age(person.BirthDate)

	today := time.Now()
	if person.BirthDate.Month() == today.Month() && person.BirthDate.Day() == today.Day() {
		return c.String(
			http.StatusOK,
			fmt.Sprintf(
				"Happy birthday to %s, who is %d years old today!",
				name,
				personAge,
			),
		)
	}

	birthdayThisYear := time.Date(today.Year(), person.BirthDate.Month(), person.BirthDate.Day(), 0, 0, 0, 0, time.UTC)
	birthdayNextYear := birthdayThisYear.AddDate(1, 0, 0)

	var nextBirthday time.Time
	if birthdayThisYear.After(today) {
		nextBirthday = birthdayThisYear
	} else {
		nextBirthday = birthdayNextYear
	}

	daysTillNextBirthday := int(math.Ceil(nextBirthday.Sub(today).Hours() / 24))

	return c.String(
		http.StatusOK,
		fmt.Sprintf("Only %d sleeps until %s is %d years old!", daysTillNextBirthday, person.Name, personAge+1),
	)
}

// SayHappyBirthdayByDate godoc
// @Summary Say happy birthday to everyone with a birthday on a particular date.
// @ID by-date
// @Tags Saying Happy Birthday
// @Security OAuth2Implicit
// @Description Say happy birthday to all people sharing specified birthday. Birthday should be specified in the format
// @Description 'Feb-02'.
// @Produce text/plain
// @Param date path string true "The date for which to wish everyone sharing that birthday a happy birthday."
// @Success 200 {string} string "Said happy birthday all person with this birthday."
// @Failure 500 {object} echo.HTTPError "Database not found."
// @Router /birthday/by-date/{date} [get]
func SayHappyBirthdayByDate(c echo.Context) error {
	dateParam := c.Param("date")

	var date time.Time
	if dateParam == "today" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse(dateInputForm, dateParam)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Sprintf("Invalid date format. Please specify in the format '%s'", dateInputForm),
			)
		}
	}

	if DB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database not found.")
	}

	var people []model.Person

	query := DB.Model(&model.Person{})
	result := query.
		Debug().
		Where(
			"EXTRACT(MONTH FROM birth_date) = ? AND EXTRACT(DAY FROM birth_date) = ?",
			int(date.Month()),
			date.Day(),
		).
		Find(&people)

	if err := result.Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	names := make([]string, len(people))
	for i, person := range people {
		names[i] = person.Name
	}

	joined := strings.Join(names, ", ")

	return c.String(http.StatusOK, fmt.Sprintf("Happy birthday %s!", joined))
}
