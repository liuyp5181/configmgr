package data

import "time"

const DBName = "test"

const (
	ConfigTable = "config"
	RecordTable = "record"
)

type Config struct {
	Id         int       `gorm:"column:id;primary_key" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Content    string    `gorm:"column:content" json:"content"`
	Status     int       `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

type Record struct {
	Id     int    `gorm:"column:id;primary_key" json:"id"`
	UserID string `gorm:"column:user_id" json:"user_id"`
	//TraceID    string    `gorm:"column:trace_id" json:"trace_id"`
	//Method     string    `gorm:"column:method" json:"method"`
	Request    string    `gorm:"column:request" json:"request"`
	Response   string    `gorm:"column:response" json:"response"`
	Error      string    `gorm:"column:error" json:"error"`
	Status     int       `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}
