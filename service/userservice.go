package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	"math/rand"
	"newginchat/models"
	"newginchat/utils"
	"strconv"
)

// GetUserList
// @Summary 获取所有用户
// @Tags        用户模块
// @Success      200  {string}   json{"code","message}
// @Router       /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0, //0成功 -1失败
		"message": "获取所有用户",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags   用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success      200  {string}   json{"code","message}
// @Router       /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	salt := fmt.Sprintf("%06d", rand.Int31())
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "用户名已被注册",
			"data":    data,
		})
		return
	}
	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "两次输入密码不一致",
			"data":    data,
		})
		return
	}
	//user.PassWord = password
	user.Salt = salt
	user.PassWord = utils.MakePassword(password, salt)
	fmt.Println(user.PassWord)
	models.CreatUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功 -1失败
		"message": "增加用户成功",
		"data":    data,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags   用户模块
// @param id query string false "id"
// @Success      200  {string}   json{"code","message}
// @Router       /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	data := models.FindUserById(user.ID)
	if data.ID == 0 {
		c.JSON(200, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "用户不存在，删除失败",
			"data":    data,
		})
		return
	}

	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功 -1失败
		"message": "删除用户成功",
		"data":    data,
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags   用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param email formData string false "email"
// @param phone formData string false "phone"
// @Success      200  {string}   json{"code","message}
// @Router       /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("update:", user)
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, //0成功 -1失败
			"message": "修改用户成功",
			"data":    user,
		})
	}

}

// LoginUser
// @Summary 用户登录
// @Tags        用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success      200  {string}   json{"code","message}
// @Router       /user/loginUser [post]
func LoginUser(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "用户不存在",
			"data":    data,
		})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //0成功 -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, //0成功 -1失败
		"message": "登录用户成功",
		"data":    data,
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {

	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(ws)
	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {

	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2012-01-02 12:04:04")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}

}
