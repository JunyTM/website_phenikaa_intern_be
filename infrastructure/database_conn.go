package infrastructure

import (
	"phenikaa/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openConnection() (*gorm.DB, error) {
	connectSQL := "host=" + dbHost +
		" user=" + dbUser +
		" port=" + dbPort +
		" dbname=" + dbName +
		" password=" + dbPassword +
		" sslmode=disable"
	db, err := gorm.Open(postgres.Open(connectSQL), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 1000,
		// DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		ErrLog.Printf("Not connect to database: %+v\n", err)
		return nil, err
	}

	return db, nil
}

func CloseConnection(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

// InitDatabase open connection and migrate database
func InitDatabase(allowMigrate bool) error {
	var err error
	db, err = openConnection()
	if err != nil {
		return err
	}

	if allowMigrate {
		InfoLog.Println("Migrating database...")

		db.Debug().AutoMigrate(
			&model.User{},               // Tài khoản
			&model.Role{},               // Vai trò
			&model.UserRole{},           // Phân quyền
			&model.Profile{},            // Thông tin cá nhân
			&model.InternShip{},         // Thông tin thực tập
			&model.InternJob{},          // Bài đăng tuyển dụng
			&model.InternshipEvaluate{}, // Đánh giá thực tập
			&model.Recruitment{},        // Quản lý thông tin ứng tuyển
			&model.UserForgotPassword{}, // Quản lý thông tin quên mật khẩu
		)
		InfoLog.Println("Done migrating database")
	}

	return nil
}
