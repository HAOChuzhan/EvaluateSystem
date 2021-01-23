package user

import (
	"errors"
	"jim_evaluate/models"
	"jim_evaluate/models/branch"
	"jim_evaluate/pkg/logging"
)

const base_url = "https://qyapi.weixin.qq.com/"
const token_url = "cgi-bin/gettoken?"
const dep_list = "cgi-bin/department/list?"
const user_cgi = "cgi-bin/user/getuserinfo?"

type User struct {
	models.Model
	Title    string `json:"title"`
	Mobile   string `json:"mobile"`
	WxUserId string `json:"wx_user_id"`
	Avatar   string `json:"avatar"` // 头像
	Power    int    `json:"power"`
	Type     int    `json:"type"`
	BranchId int    `json:"branch_id"`
}

//用户登录
func UserLogin(code string, branch_id int) (m map[string]interface{}, err error) {

	var data map[string]interface{}
	data, err = branch.GetUserInfo(branch_id, code)
	if err != nil {
		return nil, err
	}
	//user_string := data["user_info"]
	//user_info,err := models.JsonDecode(user_string)
	user_info := data["user_info"].(map[string]interface{})

	user := User{}
	if err := models.DB.Debug().Where(map[string]interface{}{"wx_user_id": user_info["UserId"], "del": 0, "branch_id": branch_id}).Find(&user).Error; err != nil {
		logging.Info(err)
	}
	user_id := user.ID
	if user.ID == 0 {
		url := base_url + user_cgi + "access_token=" + data["access_token"].(string) + "&userid=" + user_info["UserId"].(string)
		user_data_json, err := models.GetData(url)
		if err != nil {
			return nil, err
		}
		user_data, err := models.JsonDecode(user_data_json)
		if err != nil {
			return nil, err
		}
		if user_data["errcode"] == 0 {
			return nil, errors.New("企业微信获取错误")
		}
		var user_new = User{
			Title:    user_data["name"].(string),
			WxUserId: user_data["userid"].(string),
			Mobile:   user_data["mobile"].(string),
			Avatar:   user_data["thumb_avatar"].(string),
			BranchId: branch_id,
		}
		if err := models.DB.Debug().Model(User{}).Create(&user_new).Error; err != nil {
			logging.Info(err)
		}
		user_id = user_new.ID
	}
	var result map[string]interface{}
	result = make(map[string]interface{})
	result["user_id"] = user_id
	result["branch_id"] = branch_id
	result["wx_user_id"] = user_info["UserId"]
	return result, nil
}

//获取所有用户信息列表
func GetAllUserInfo(page int, limit int) (data map[string]interface{}) {
	var total int
	var users []User
	db := models.DB.Debug().Model(&User{}).Where("del = ?", 0)
	if err := db.Count(&total).Error; err != nil {
		return
	}
	offset := (page - 1) * limit
	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return
	}
	data = map[string]interface{}{
		"users":     users,
		"total":     total,
		"page":      page,
		"pagetotal": int(total/limit) + 1,
	}
	return data

}

//更新用户信息
func UpdateUser(userID int, data interface{}) bool {
	if err := models.DB.Debug().Model(&User{}).Where("id = ?", userID).Update(data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

/*******************************************/
// Create 新建用户
func CreateUser(user *User) bool {
	if models.DB.NewRecord(user) {
		models.DB.Debug().Create(user)
		return !models.DB.NewRecord(user)
	}
	return false
}

// GetInfo 根据id获取用户信息
func GetInfo(id int) (user User) {
	if err := models.DB.Debug().Where(map[string]interface{}{"id": id}).Find(&user).Error; err != nil {
		logging.Info(err)
	}
	return
}

// Login 登陆
func Login(user User) bool {
	var find User
	if models.DB.Where(user).Select("id").First(&find); find.ID > 0 {
		return true
	}
	return false
}

// CreateByPassword 根据密码新建用户
func CreateByPassword(user *User) bool {
	if models.DB.NewRecord(user) {
		models.DB.Debug().Create(user)
		return !models.DB.NewRecord(user)
	}

	return false
}

// Update 更新信息
func Update(user *User, data interface{}) bool {
	if err := models.DB.Debug().Model(user).Where("id = ?", user.ID).Update(data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

// UpdateColumn 更新指定列
func UpdateColumn(user *User, column string, data interface{}) bool {
	if err := models.DB.Debug().Model(User{}).Where("id = ?", user.ID).UpdateColumn(column, data).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

// AddUnionid 添加 Unionid
func AddUnionid(userId int, unionid string) bool {
	if err := models.DB.Debug().Model(User{}).Where("id = ?", userId).UpdateColumn("unionid", unionid).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

// QueryUserByUnionid 通过unionid获取用户信息
func QueryUserByUnionid(unionid string) (user User) {
	models.DB.Debug().Model(User{}).Where("unionid = ?", unionid).First(&user)
	return
}

// QueryUserByOpenid 通过Openid获取用户信息
func QueryUserByOpenid(openid string) (user User) {
	models.DB.Debug().Model(User{}).Where("openid = ?", openid).First(&user)
	return
}

// QueryUserById 通过id获取用户信息
func QueryUserById(id int) (user User) {
	models.DB.Debug().Model(User{}).Where("id = ?", id).First(&user)
	return
}
