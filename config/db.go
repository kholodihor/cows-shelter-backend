package config

import (
	"github.com/kholodihor/cows-shelter-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func  Connect(){
	db, err := gorm.Open(postgres.Open("postgresql://cows_shelter_user:FWlBjtlI3LHnnxHDlTFBGMz4QornIlXv@dpg-csctuadumphs7399b6q0-a.frankfurt-postgres.render.com/cows_shelter"))

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