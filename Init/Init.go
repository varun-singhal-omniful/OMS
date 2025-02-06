package Init

import (
	"context"

	"github.com/varun-singhal-omniful/oms-service/database"
)

func InitializeDB(c context.Context) {
	database.ConnectMongo(c)
}
