package utils

import (
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"
)

// InitDB auto create schema
func InitDB() {
	db, err := api.ConnectDB()
	if err != nil {
		panic(err)
	}

	// init tables
	tables := []interface{}{
		models.Policy[[]byte]{},
		models.NetworkRestriction{},
		models.NetworkAccessorRelation{},
		// models.BatchTask{},
	}

	for _, v := range tables {
		if !db.Migrator().HasTable(v) {
			err := db.Migrator().CreateTable(v)
			if err != nil {
				panic(err)
			}
		}
	}
}
