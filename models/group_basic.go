// @Author TrandLiu
// @Date 2024/12/17 16:53:00
// @Desc
package models

import (
	"gorm.io/gorm"
)

// 群消息
type Group struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Type    int
	Desc    string
}

func (table *Group) TableName() string {
	return "group"
}
