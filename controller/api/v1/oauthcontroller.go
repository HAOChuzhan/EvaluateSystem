package v1

import (
	"encoding/json"
	"jim_evaluate/models"
	"jim_evaluate/models/oauth"
	"jim_evaluate/models/branch"
	"jim_evaluate/models/user"
	"jim_evaluate/pkg/app"
	"jim_evaluate/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	appG := app.Gin{C: c}
	branch_id := c.DefaultPostForm("branch_id", "0")

	data, err := branch.GetAppInfo(branch_id)
	if err != nil {
		appG.Response(http.StatusBadRequest,"操作错误",data)
	}

	back_url := "https://127.0.0.1:8000/login"
	var m map[string]interface{}
	m = make(map[string]interface{})
	m["branch_id"] = branch_id

	state_byte, _ := json.Marshal(m)
	state := string(state_byte)

	scope := "snsapi_base"
	url := oauth.GetOauthUrl(data.CorpId, back_url, state, scope)

	c.Redirect(200,url)
	appG.Response(200,"操作成功",data)

}
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	state_string := c.DefaultPostForm("state", "")
	code := c.DefaultPostForm("code", "")

	state, err := models.JsonDecode(state_string)
	if err != nil {
		appG.Response(http.StatusBadRequest,"操作错误",nil)
	}
	branch_id := state["branch_id"].(int)

	login_result, err := user.UserLogin(code, branch_id)
	if err != nil {
		appG.Response(http.StatusBadRequest,"操作错误",nil)
	}
	token, err := util.GenerateToken(login_result["user_id"].(int))
	if err != nil {
		appG.Response(500,"获取token失败",nil)
	}
	appG.Response(http.StatusOK, "e.OK", token)

}
func SendMsg(c *gin.Context) {
	appG := app.Gin{C: c}

	toUsers := c.DefaultPostForm("tousers", "haochuzhan")
	title := "测试"
	description := "请在本月之前提交相关的任务表单"
	url := "www.baidu.com"
	data, err := models.SendCardMsg(toUsers, title, description, url)
	if err != nil {
		appG.Response(http.StatusBadRequest, "失败", nil)
		return
	}
	appG.Response(http.StatusOK, "获取成功", data)
}
