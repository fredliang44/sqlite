package tests_test

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils/tests"
)

type (
	Toy           = tests.Toy
	Pet           = tests.Pet
	User          = tests.User
	Language      = tests.Language
	Company       = tests.Company
	Account       = tests.Account
	Coupon        = tests.Coupon
	CouponProduct = tests.CouponProduct
	Order         = tests.Order
)

var (
	AssertEqual    = tests.AssertEqual
	AssertObjEqual = tests.AssertObjEqual
	Now            = tests.Now
)

var DB *gorm.DB

func init() {
	var err error
	if DB, err = OpenTestConnection(); err != nil {
		log.Printf("failed to connect database, got error %v", err)
		os.Exit(1)
	} else {
		sqlDB, err := DB.DB()
		if err == nil {
			err = sqlDB.Ping()
		}

		if err != nil {
			log.Printf("failed to connect database, got error %v", err)
		}

		if DB.Dialector.Name() == "sqlite" {
			DB.Exec("PRAGMA foreign_keys = ON")
		}
		RunMigrations()
	}
}

func OpenTestConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	// if err == nil {
	// 	err = db.Exec(fmt.Sprintf("PRAGMA synchronous = %s;", "NORMAL")).Error
	// }
	// if err == nil {
	// 	err = db.Exec(fmt.Sprintf("PRAGMA locking_mode = %s;", "NORMAL")).Error
	// }
	// if err == nil {
	// 	err = db.Exec(fmt.Sprintf("PRAGMA busy_timeout = %d;", 5000)).Error
	// }
	// if err == nil {
	// 	err = db.Exec(fmt.Sprintf("PRAGMA journal_mode = %s;", "WAL")).Error
	// }

	if debug := os.Getenv("DEBUG"); debug == "true" {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else if debug == "false" {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}

	return
}

func RunMigrations() {
	var err error
	allModels := []interface{}{&User{}, &Account{}, &Pet{}, &Company{}, &Toy{}, &Language{}, &Coupon{}, &CouponProduct{}, &Order{}}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allModels), func(i, j int) { allModels[i], allModels[j] = allModels[j], allModels[i] })

	DB.Migrator().DropTable("user_friends", "user_speaks")

	if err = DB.Migrator().DropTable(allModels...); err != nil {
		log.Printf("Failed to drop table, got error %v\n", err)
		os.Exit(1)
	}

	if err = DB.AutoMigrate(allModels...); err != nil {
		log.Printf("Failed to auto migrate, but got error %v\n", err)
		os.Exit(1)
	}

	for _, m := range allModels {
		if !DB.Migrator().HasTable(m) {
			log.Printf("Failed to create table for %#v\n", m)
			os.Exit(1)
		}
	}
}