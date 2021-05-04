/**
 * @Author: Resynz
 * @Date: 2021/4/22 11:45
 */
package db

import "export-service/config"
import "github.com/Resynz/db-handler"

var (
	ExportHandler *db_handler.DBHandler
	CrmHandler    *db_handler.DBHandler
)

func InitDBHandler() error {
	var err error

	// todo message redis
	ExportHandler, err = db_handler.New(config.Conf.ExportMySql, nil)
	if err != nil {
		return err
	}

	CrmHandler, err = db_handler.New(config.Conf.CrmMySql, nil)
	if err != nil {
		return err
	}

	return nil
}
