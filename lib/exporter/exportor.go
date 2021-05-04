package exporter

import (
	"encoding/csv"
	"export-service/config"
	"export-service/tools"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	utf8_bom = "\xef\xbb\xbf"
)

type Exporter struct {
	Titles     []string
	Rows       <-chan map[string]string
	Creator    int64
	ResultPath string
	RowFlush   bool
}

func (s *Exporter) genWriter() (io.WriteCloser, error) {
	filePath := fmt.Sprintf("%s", time.Now().Format("2006/01/02"))
	fileName := fmt.Sprintf("%d_%d_%s.csv", s.Creator, time.Now().Unix(), tools.GenRandomCode(3))
	dirPath := fmt.Sprintf("%s/%s", config.Conf.ResultPath, filePath)
	if err := tools.CheckDirExistOrCreate(dirPath); err != nil {
		return nil, err
	}
	s.ResultPath = fmt.Sprintf("%s/%s", filePath, fileName)
	return os.OpenFile(fmt.Sprintf("%s/%s", dirPath, fileName), os.O_TRUNC|os.O_CREATE|os.O_RDWR, os.ModePerm)
}

func (s *Exporter) Exec() error {

	if len(s.Titles) == 0 {
		return nil
	}

	w, err := s.genWriter()
	if err != nil {
		return err
	}

	defer w.Close()

	fCsv := csv.NewWriter(w)
	defer fCsv.Flush()
	// todo 1. 写入标题
	oldTitle := s.Titles[0]
	hasBom := strings.HasPrefix(oldTitle, utf8_bom)
	if !hasBom {
		s.Titles[0] = fmt.Sprintf("%s%s", utf8_bom, oldTitle)
	}
	if err = fCsv.Write(s.Titles); err != nil {
		return err
	}
	if !hasBom {
		s.Titles[0] = oldTitle
	}

	// todo 2.写入内容
	for r := range s.Rows {
		row := make([]string, len(s.Titles))
		for i, t := range s.Titles {
			v, ok := r[t]
			_v := ""
			if ok {
				_v = tools.FormatCsvValue(v)
			}
			row[i] = _v
		}
		if err = fCsv.Write(row); err != nil {
			return err
		}
		if s.RowFlush {
			fCsv.Flush()
		}
	}
	return nil
}
