package Server

import (
	"fmt"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Achivment struct {
	*gorm.Model
	Name   string         `gorm:"type:varchar(100)"`
	Heroes pq.StringArray `gorm:"type:text[]"`
}

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {

	dsn := "host=localhost user=postgres password=364678x dbname=postgres port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	if err := DB.AutoMigrate(&Achivment{}); err != nil {
		fmt.Println("Failed to migrate db")
	}
	return DB, nil
}

func UpdateTable(db *gorm.DB, names []string) {
	for _, name := range names {
		var count int64
		db.Model(&Achivment{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			newAchivment := Achivment{
				Name:   name,
				Heroes: []string{},
			}
			db.Create(&newAchivment)
			fmt.Printf("Achievement '%s' added to the database.\n", name)

		} else {
			fmt.Printf("Achievement '%s' already exists in the database.\n", name)
		}
	}
}
