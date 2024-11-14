package config

import (
	"addressBook/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppCtx struct {
	Logger *zap.Logger
	DB     *gorm.DB
}

func NewAppCtx() (*AppCtx, error) {
	logger, err := zap.NewProduction()
	defer logger.Sync()
	if err != nil {
		// fmt.Println("Error in Connecting with DB : ", err)
		logger.Error("error in creating a zap ", zap.Error(err))
		return nil, err
	}
	db, err := database.ConnectionDB()
	if err != nil {
		logger.Error("Error in Database Connection : ", zap.Error(err))
	}
	return &AppCtx{
		Logger: logger,
		DB:     db,
	}, nil
}
