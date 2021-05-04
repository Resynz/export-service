package action

import (
	"export-service/db"
	"export-service/lib/exporter"
	db_handler "github.com/Resynz/db-handler"
	"gitlab.foundingaz.cn/db-models/model-crm"
	"log"
)

func Test(condition *db_handler.Condition) (*exporter.Exporter, error) {
	var admin crm.Admin
	adminChan := db.CrmHandler.Iterate(&admin, admin.GetTableName(), condition)

	exp := &exporter.Exporter{}

	exp.Titles = []string{
		"姓名",
		"手机号",
	}

	rows := make(chan map[string]string)

	go func() {
		for c := range adminChan {
			if c == nil {
				continue
			}
			if e, ok := c.(error); ok {
				log.Printf("[Test] error:%s\n", e.Error())
				continue
			}
			v, ok := c.(*crm.Admin)
			if !ok {
				log.Printf("[Test] convert failed \n")
				continue
			}
			row := make(map[string]string)
			row["姓名"] = v.Name
			row["手机号"] = v.Phone
			rows <- row
		}
		close(rows)
	}()

	exp.Rows = rows
	return exp, nil
}
