package config

import db_handler "github.com/Resynz/db-handler"

type Config struct {
	Mode        string               `json:"mode"`
	AppPort     int                  `json:"app_port"`
	LogConfig   logConf              `json:"log_config"`
	ExportMySql *db_handler.DbConfig `json:"export_my_sql"`
	CrmMySql    *db_handler.DbConfig `json:"crm_my_sql"`
	QueueSize   int                  `json:"queue_size"`
	Host        string               `json:"host"`
	ResultPath  string               `json:"result_path"`
}

type logConf struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Level string `json:"level"`
}
