package service

import (
	"fmt"
	"imaotai/config"
	"imaotai/db"
	"imaotai/models"
	"imaotai/reqfunc"
)

// Reservation 预约
func Reservation(conf *config.Config) error {
	var mtversion *models.MtVersion
	err := db.Gormdb.First(&mtversion).Error
	if err != nil {
		return err
	}
	var session *models.Session
	err = db.Gormdb.First(&session).Error
	if err != nil {
		return err
	}
	var sessionitems []*models.SessionItem
	err = db.Gormdb.Find(&sessionitems).Error
	if err != nil {
		return err
	}
	for _, account := range conf.Account {
		for _, item := range sessionitems {
			shopid, err := GetShopId(item.ItemCode, *account)
			if err != nil {
				return err
			}
			rt, err := reqfunc.Reservation(reqfunc.UserInfo{
				UserId: account.UserId,
				Lat:    fmt.Sprintf("%v", account.Lat),
				Lng:    fmt.Sprintf("%v", account.Lng),
				Token:  account.Token,
			}, item.ItemCode, shopid, fmt.Sprintf("%v", session.SessionId), mtversion.Version)
			if err != nil {
				return err
			}
			// todo 后续添加消息通知
			println(rt)
		}
	}
	return nil
}