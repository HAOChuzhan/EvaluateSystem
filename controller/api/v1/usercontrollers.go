package v1

import (
	"jim_evaluate/models/user"
	"jim_evaluate/pkg/app"


	"github.com/astaxie/beego/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func GetUserList(c *gin.Context) {
	appG := app.Gin{C: c}
	page := com.StrTo(c.DefaultPostForm("page", "1")).MustInt()
	limit := com.StrTo(c.DefaultPostForm("limit", "10")).MustInt()

	AllUserInfo := make(map[string]interface{})

	AllUserInfo = user.GetAllUserInfo(page, limit)
	appG.Response(http.StatusOK, "获取成功", AllUserInfo)

}
func UpdateUserInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	editedDate := make(map[string]interface{})

	user_id := com.StrTo(c.DefaultPostForm("user_id", "0")).MustInt()
	//user_id := c.MustGet("AuthData").(*util.Claims).User.ID
	type1 := c.DefaultPostForm("type", "")
	// if type1 == "" {
	// 	appG.Response(http.StatusBadRequest, e.BadRequest, "未填写用户类型")
	// }
	user_type := com.StrTo(type1).MustInt()
	valid.Required(user_type, "type").Message("type是必需的")
	valid.Max(user_type, 5, "type").Message("type必须小于等于5")
	valid.Min(user_type, 1, "type").Message("type必须大于等于1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, "参数有误", valid.Errors)
		return
	}
	editedDate["type"] = user_type

	if !user.UpdateUser(user_id, editedDate) {
		appG.Response(http.StatusInternalServerError, "更新数据失败", nil)
		return
	}
	appG.Response(http.StatusOK, "更新成功", nil)

}
