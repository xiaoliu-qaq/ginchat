// @Author TrandLiu
// @Date 2024/12/17 16:51:00
// @Desc
package models

import "gorm.io/gorm"

//人员关系

type Contact struct {
	gorm.Model
	OwnerID  uint //谁的关系
	TargetId uint //对应的谁
	Type     int  //对应的类型 0 1 3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
