package config

import (
	"github.com/kholodihor/cows-shelter-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func  Connect(){
	db, err := gorm.Open(postgres.Open("postgresql://cows_shelter_xoit_user:lGqBu3AMCxKFQfK5YCuaPrE4OCv23cDL@dpg-csd0tglumphs739b9i8g-a.frankfurt-postgres.render.com/cows_shelter_xoit"))

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.User{}, 
		&models.Contact{}, 
		&models.Excursion{},
		&models.Gallery{},
		&models.News{},
		&models.Partner{},
		&models.Password{},
		&models.Pdf{},
		&models.Review{},
	)
	
	DB=db

}