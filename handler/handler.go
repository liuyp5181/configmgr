package handler

import (
	"context"
	"fmt"
	"github.com/liuyp5181/base/config"
	"github.com/liuyp5181/base/database"
	"github.com/liuyp5181/base/log"
	"github.com/liuyp5181/base/service/extend"
	pb "github.com/liuyp5181/configmgr/api"
	"github.com/liuyp5181/configmgr/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm/clause"
	"time"
)

type ServiceImpl struct {
	pb.UnimplementedGreeterServer
}

func (s *ServiceImpl) GetConfig(ctx context.Context, req *pb.GetConfigReq) (resp *pb.GetConfigRes, err error) {
	fmt.Println(req, config.ServiceName)
	var cfg data.Config
	db := database.GetMysql(data.DBName)
	err = db.Table(data.ConfigTable).Where("name = ?", req.Key).Take(&cfg).Error
	if err != nil {
		log.Error("get db config err = ", err)
		return nil, status.Error(codes.DataLoss, err.Error())
	}
	resp = &pb.GetConfigRes{Val: cfg.Content}
	return resp, nil
}

func (s *ServiceImpl) SetConfig(ctx context.Context, req *pb.SetConfigReq) (resp *pb.SetConfigRes, err error) {
	e := extend.NewContext(ctx)
	defer func() {
		err := data.InsertRecord(e.GetClient("user_id"), req, resp, err)
		if err != nil {
			log.Error("InsertRecord err = %v", err)
		}
	}()

	db := database.GetMysql(data.DBName)
	// 在id冲突时，将name更新为新值
	// 等同于insert into...on duplicate key update
	// insert into user('id','name') values(?,?) on duplicate key update id=values('id'), name=values('name')
	// var user = User{}
	// db.Clauses(clause.OnConflict{
	// 	 Columns: []clause.Column{{Name: "id"}},
	//	 DoUpdates: clause.AssignmentColumns([]string{"name"}),
	// }).Create(&user)
	var cfg = data.Config{
		Name:       req.Key,
		Content:    req.Val,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err = db.Table(data.ConfigTable).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"content", "update_time"}),
	}).Omit("status").Create(&cfg).Error
	if err != nil {
		return nil, status.Error(codes.DataLoss, err.Error())
	}

	PutWatch(req.Key, []byte(req.Val))

	resp = &pb.SetConfigRes{}
	return resp, nil
}

func (s *ServiceImpl) DelConfig(ctx context.Context, req *pb.DelConfigReq) (resp *pb.DelConfigRes, err error) {
	e := extend.NewContext(ctx)
	defer func() {
		err := data.InsertRecord(e.GetClient("user_id"), req, resp, err)
		if err != nil {
			log.Error("InsertRecord err = %v", err)
		}
	}()

	fmt.Println(req)
	var cfg data.Config
	db := database.GetMysql(data.DBName)
	err = db.Table(data.ConfigTable).Where("name = ?", req.Key).Delete(&cfg).Error
	if err != nil {
		return nil, status.Error(codes.DataLoss, err.Error())
	}

	DeleteWatch(req.Key)

	resp = &pb.DelConfigRes{}
	return resp, nil
}
