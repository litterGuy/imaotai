package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"imaotai/config"
	"imaotai/db"
	"imaotai/models"
	"imaotai/reqfunc"
	"log"
)

// RefreshData 刷新数据库
func RefreshData(configs *config.Config) error {
	//mtversion
	err := db.Gormdb.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.MtVersion{}).Error
	if err != nil {
		return err
	}
	version, err := reqfunc.GetMTVersion()
	if err != nil {
		return err
	}
	err = db.Gormdb.Create(&models.MtVersion{Version: version}).Error
	if err != nil {
		return err
	}
	// 产品信息插入
	err = db.Gormdb.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Session{}).Error
	if err != nil {
		return err
	}
	err = db.Gormdb.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.SessionItem{}).Error
	if err != nil {
		return err
	}
	session, err := reqfunc.GetCurrentSessionId()
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("获取产品信息为空")
	}
	err = db.Gormdb.Create(&models.Session{SessionId: session.SessionID}).Error
	if err != nil {
		return err
	}
	items := make([]*models.SessionItem, 0)
	for _, item := range session.ItemList {
		items = append(items, &models.SessionItem{
			SessionID: session.SessionID,
			Content:   item.Content,
			ItemCode:  item.ItemCode,
			JumpURL:   item.JumpURL,
			Picture:   item.Picture,
			PictureV2: item.PictureV2,
			Title:     item.Title,
		})
	}
	err = db.Gormdb.CreateInBatches(items, 100).Error
	if err != nil {
		return err
	}

	// 门店投放数量
	err = db.Gormdb.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.ShopItem{}).Error
	if err != nil {
		return err
	}
	shopitemdblist := make([]*models.ShopItem, 0)
	for _, item := range items {
		// 多账号下相同省份时,跳过处理
		provinceMap := make(map[string]string)

		for _, account := range configs.Account {

			if _, ok := provinceMap[account.Province]; ok {
				continue
			}
			provinceMap[account.Province] = account.Province

			shopitems, err := reqfunc.GetShopsByProvince(account.Province, item.ItemCode, fmt.Sprintf("%v", session.SessionID))
			if err != nil {
				log.Println(fmt.Sprintf("省市:%s,code:%s,sessionId:%v.erris : %s", account.Province, item.ItemCode, session.SessionID, err.Error()))
				return err
			}
			for _, shopitem := range shopitems {
				shopitemdblist = append(shopitemdblist, &models.ShopItem{
					ShopID:              shopitem.ShopID,
					Count:               shopitem.Count,
					MaxReserveCount:     shopitem.MaxReserveCount,
					DefaultReserveCount: shopitem.DefaultReserveCount,
					ItemID:              shopitem.ItemID,
					Inventory:           shopitem.Inventory,
					OwnerName:           shopitem.OwnerName,
					Province:            account.Province,
				})
			}
		}
	}
	err = db.Gormdb.CreateInBatches(shopitemdblist, 100).Error
	if err != nil {
		return err
	}
	// 门店详情
	err = db.Gormdb.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.ShopBean{}).Error
	if err != nil {
		return err
	}
	shopmap, err := reqfunc.GetShopList()
	if err != nil {
		return err
	}
	shoplist := make([]*models.ShopBean, 0)
	for _, v := range shopmap {
		shoplist = append(shoplist, &models.ShopBean{
			Address:       v.Address,
			City:          v.City,
			CityName:      v.CityName,
			District:      v.District,
			DistrictName:  v.DistrictName,
			FullAddress:   v.FullAddress,
			Lat:           v.Lat,
			Layaway:       v.Layaway,
			Lng:           v.Lng,
			Name:          v.Name,
			OpenEndTime:   v.OpenEndTime,
			OpenStartTime: v.OpenStartTime,
			Province:      v.Province,
			ProvinceName:  v.ProvinceName,
			ShopID:        v.ShopID,
			TenantName:    v.TenantName,
		})
	}
	err = db.Gormdb.CreateInBatches(shoplist, 100).Error
	if err != nil {
		return err
	}
	return nil
}
