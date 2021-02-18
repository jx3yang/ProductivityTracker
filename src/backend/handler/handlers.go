package handler

import db "github.com/jx3yang/ProductivityTracker/src/backend/database"

func InitHandlers(d *db.MongoDatabase) error {
	initCardCollection(d)
	initListCollection(d)
	return nil
}
