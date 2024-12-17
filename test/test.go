// @Author TrandLiu
// @Date 2024/12/12 16:57:00
// @Desc
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"newginchat/models"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {

	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/newginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	err = db.AutoMigrate(&models.UserBasic{})
	err = db.AutoMigrate(&models.Group{})
	err = db.AutoMigrate(&models.Message{})
	err = db.AutoMigrate(&models.Contact{})
	fmt.Println("执行了")
	if err != nil {
		return
	}

	// Create
	//db.Create(&Product{Code: "D42", Price: 100})
	//user := &models.UserBasic{}
	//user.Name = "六志康"
	//db.Create(user)
	// Read
	//var product Product
	//fmt.Println(db.First(user, 1)) // 根据整型主键查找
	//db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	//
	//// Update - 将 product 的 price 更新为 200
	//db.Model(user).Update("PassWord", 123)
	//// Update - 更新多个字段
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	//db.Delete(&product, 1)
}
