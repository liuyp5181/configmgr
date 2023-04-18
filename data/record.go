package data

import (
	"encoding/json"
	"github.com/liuyp5181/base/database"
	"time"
)

func InsertRecord(uid string, req interface{}, resp interface{}, inErr error) error {
	reqData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	var respData []byte
	if resp != nil {
		respData, err = json.Marshal(resp)
		if err != nil {
			return err
		}
	}

	var errData string
	if inErr != nil {
		errData = inErr.Error()
	}

	db := database.GetMysql(DBName)
	err = db.Table(RecordTable).Create(&Record{
		UserID:     uid,
		Request:    string(reqData),
		Response:   string(respData),
		Error:      errData,
		Status:     0,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}).Error
	if err != nil {
		return err
	}

	return nil
}
