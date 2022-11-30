package postgresql

import (
	"fmt"

	"github.com/vesicash/auth-ms/internal/config"
	"github.com/vesicash/auth-ms/utility"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Databases struct {
	Admin         *gorm.DB
	Auth          *gorm.DB
	Notifications *gorm.DB
	Payment       *gorm.DB
	Reminder      *gorm.DB
	Subscription  *gorm.DB
	Transaction   *gorm.DB
	Verification  *gorm.DB
	Cron          *gorm.DB
}

var DB Databases

// Connection gets connection of mysqlDB database
func Connection() Databases {
	return DB
}

func ConnectToDatabases() Databases {
	logger := utility.NewLogger()
	dbsCV := config.GetConfig().Databases
	databases := Databases{}
	utility.LogAndPrint("connecting to databases")
	databases.Admin = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.ADMIN_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Auth = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.AUTH_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Notifications = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.NOTIFICATIONS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Payment = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.PAYMENT_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Reminder = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.REMINDERS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Subscription = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.SUBSCRIPTIONS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Transaction = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.TRANSACTIONS_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Verification = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.VERIFICATION_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)
	databases.Cron = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.CRON_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)

	utility.LogAndPrint("connected to databases")

	fmt.Println("connected to db")
	DB = databases
	return DB
}

func connectToDb(host, user, password, dbname, port, sslmode, timezone string, logger *utility.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", host, user, password, dbname, port, sslmode, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utility.LogAndPrint(fmt.Sprintf("connection to %v db failed with: %v", dbname, err))
		panic(err)

	}
	utility.LogAndPrint(fmt.Sprintf("connected to %v db", dbname))
	return db
}

// _ = db.AutoMigrate(MigrationModels()...)
func MigrationModels(db *gorm.DB) []interface{} {
	return []interface{}{
		// model.User{},
	}
}

func MigrateModels(db *gorm.DB, models []interface{}) {
	_ = db.AutoMigrate(models...)
}
