package migrate

import (
	"higo-game-bus/database"
	"higo-game-bus/model"
)

// MigrateModel 迁移数据库
func MigrateModel() {
	db := database.GetInstance()
	_ = db.AutoMigrate(&model.Game{})
	_ = db.AutoMigrate(&model.Battle{})
}
