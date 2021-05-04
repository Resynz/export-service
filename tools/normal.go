package tools

import (
	"export-service/db"
	"fmt"
	"github.com/Resynz/model-export"
	"math/rand"
	"os"
	"strings"
	"time"
)

func GenerateSn() (string, error) {
	sn := fmt.Sprintf("ET%s%s", time.Now().Format("20060102150405"), GenRandomCode(3))
	var task model_export.Task
	has, err := db.ExportHandler.GetOne(&task, task.GetTableName(), "sn", sn)
	if err != nil {
		return "", err
	}
	if has {
		return GenerateSn()
	}
	return sn, nil
}
func init() {
	rand.Seed(time.Now().Unix())
}

func GenRandomCode(length int) string {
	numDic := [9]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numDic)
	var sb strings.Builder
	for i := 0; i < length; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numDic[rand.Intn(r)])
	}
	return sb.String()
}

func FormatCsvValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, ",", "，")
	return value
}

func CheckDirExistOrCreate(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
