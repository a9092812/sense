// package models

// import (
// 	"log"

// 	"gorm.io/gorm"
// )

// func RunMigrations(db *gorm.DB) {

// 	err := db.AutoMigrate(
// 		&Device{},
// 		&SHT40Data{},
// 		&LuxSensorData{},
// 		&LIS2DHData{},
// 		&SoilSensorData{},
// 		&SpeedDistanceData{},
// 		&AmmoniaSensorData{},
// 		&TempLoggerData{},
// 		&DataLoggerData{},
// 		&DeviceLatestState{},
// 	)

// 	if err != nil {
// 		log.Fatalf("Migration failed: %v", err)
// 	}

// 	log.Println("GORM automigration complete")

// 	// Timescale hypertable example
// 	db.Exec(`
// 		SELECT create_hypertable(
// 			'sht40_data',
// 			'time',
// 			if_not_exists => TRUE
// 		);
// 	`)
// }

package models

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	migrationDir := "./migrations"

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		path := filepath.Join(migrationDir, file.Name())

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = sqlDB.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("migration failed %s: %v", file.Name(), err)
		}

		fmt.Println("migration applied:", file.Name())
	}

	return nil
}
