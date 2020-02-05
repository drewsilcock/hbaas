package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/drewsilcock/hbaas-server/model"
	"github.com/drewsilcock/hbaas-server/seeddata"
)

// Use separate type for unmarshalling CSV as we need to modify the unmarshaling behaviour for the time field.
type Date struct {
	time.Time
}

func (date *Date) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2 January 2006", csv)
	return err
}

type SeedPerson struct {
	Name      string `csv:"name"`
	BirthDate Date   `csv:"birth_date"`
}

func init() {
	rootCmd.AddCommand(seedDatabaseCmd)
}

var seedDatabaseCmd = &cobra.Command{
	Use:   "seed-database",
	Short: "Seed person database.",
	Long:  "Seed database with some famous people's names and birthdays.",
	Run:   seedDatabase,
}

func seedDatabase(cmd *cobra.Command, args []string) {
	peopleCsvBytes, err := seeddata.Asset("people.csv")
	if err != nil {
		log.Fatal("Unable to read seed data:", err)
	}
	peopleCsv := string(peopleCsvBytes)

	var seedPeople []*SeedPerson
	if err := gocsv.UnmarshalString(peopleCsv, &seedPeople); err != nil {
		log.Fatal("Unable to unmarshal seed people from CSV:", err)
	}

	db, err := gorm.Open("postgres", viper.Get("postgres_url"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("Seeding people...")
	for _, seedPerson := range seedPeople {
		person := &model.Person{Name: seedPerson.Name, BirthDate: seedPerson.BirthDate.Time}
		db.FirstOrCreate(person, person)
	}
	fmt.Println(fmt.Sprintf("Finished seeding %d people", len(seedPeople)))
}
