package action

import (
	"export-service/lib/exporter"
	db_handler "github.com/Resynz/db-handler"
)

var (
	ActionMap map[string]ActionFunc
)

type ActionFunc func(condition *db_handler.Condition) (*exporter.Exporter, error)

func InitActions() {
	ActionMap = make(map[string]ActionFunc)
	ActionMap["test"] = Test
}
