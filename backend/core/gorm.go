package core

import (
	"gin-vue/api/utils"
	"gin-vue/global"
	"time"

	"gin-vue/modles/models"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)


func InitGorm() *gorm.DB {
	if global.Config.Mysql.Host == "" {
		global.Logger.Println("未配置mysql，取消gorm连接")
		return nil
	}
	dsn := global.Config.Mysql.Dsn()
	global.Logger.Println("DSN:", dsn)
	var mysqlLogger logger.Interface
	if global.Config.System.Env == "debug" {

		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Println("MySql连接失败...")
	}
	db.AutoMigrate(&models.User{})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour * 4)
	return db
}

func InitSqliteGorm() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("cloud_server.db"), &gorm.Config{})
	if err != nil {
		global.Logger.Println("sqlite connect failed")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Device{})
	db.AutoMigrate(&models.Bind{})
	db.AutoMigrate(&models.Token{})
	db.AutoMigrate(&models.Disk{})
	db.AutoMigrate(&models.DiskBind{})
	db.AutoMigrate(&models.Gpu{})
	db.AutoMigrate(&models.License{})
	db.AutoMigrate(&models.LicenseDate{})
	db.AutoMigrate(&models.Template{})
	InsertDefaultUser(db)
	utils.InitGpuInfo(db)
	return db
}

func InsertDefaultUser(db *gorm.DB) {
	var myUser models.User
	result := db.First(&myUser, "name=?", "vmadmin")
	if myUser.Name == "vmadmin" || result.RowsAffected > 1 {
		global.Logger.Println("default vmadmin exists")
		return
	}
	myNewUser := models.User{ID: models.NewUUID(), Name: "vmadmin", Password: "vmadmin", Role: "管理员", Status: "启用"}
	res := db.Create(&myNewUser)
	if res.Error != nil {
		global.Logger.Panicln("init user failed")
	}
	global.Logger.Println("insert default admin user successfully")
}
