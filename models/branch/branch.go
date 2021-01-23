package branch

import (
	"encoding/json"
	"jim_evaluate/models"
	"jim_evaluate/pkg/logging"
	"time"
)

const base_url = "https://qyapi.weixin.qq.com/"
const token_url = "cgi-bin/gettoken?"
const dep_list = "cgi-bin/department/list?"
const user_cgi = "cgi-bin/user/getuserinfo?"

type Branch struct {
	models.Model
	ID     int64  `json:"-"`
	CorpId      string  `json:"corp_id"`
	AgentId       string  `json:"agent_id"`
	Secret   string  `json:"secret"`
	AccessToken     string  `json:"access_token"`
	ExpireTime       string  `json:"expire_time"`
	Del          int     `json:"del"`
	DelTime       string  `json:"del_time"`
	UpdateTime       string  `json:"update_time"`
	AddTime       string  `json:"add_time"`
}

func GetAppInfo(branch_id string) (branch_info Branch, err error) {
	if err := models.DB.Debug().Where(map[string]interface{}{"id": branch_id, "del": 0}).Find(&branch_info).Error; err != nil {
		logging.Info(err)
	}
	access_token := branch_info.AccessToken
	expire_time_string := branch_info.ExpireTime
	if access_token == "" {
		token_info, _ := UpdateToken(branch_info.CorpId, branch_info.Secret)
		branch_info.AccessToken = token_info["access_token"].(string)
		branch_info.ExpireTime = token_info["expire_time"].(string)
	}
	expire_time_format, _ := time.Parse("2006-01-02 15:04:05", expire_time_string)
	expire_time := expire_time_format.Unix()

	if expire_time < time.Now().Unix() {
		token_info, _ := UpdateToken(branch_info.CorpId, branch_info.Secret)
		branch_info.AccessToken = token_info["access_token"].(string)
		branch_info.ExpireTime = token_info["expire_time"].(string)
	}

	flag := checkToken(access_token)
	if !flag {
		token_info, _ := UpdateToken(branch_info.CorpId, branch_info.Secret)
		branch_info.AccessToken = token_info["access_token"].(string)
		branch_info.ExpireTime = token_info["expire_time"].(string)
	}

	return
}

func UpdateToken(corp_id string, secret string) (arr map[string]interface{}, err error) {
	arr = make(map[string]interface{})
	url := base_url + token_url + "corpid=" + corp_id + "&corpsecret=" + secret
	var m map[string]interface{}
	data, err := models.GetData(url)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &m); err != nil {
		return
	}
	branch := Branch{}
	//access_token
	branch.AccessToken = m["access_token"].(string)
	//时间
	expire_in := int64(m["expires_in"].(float64))

	time_now := models.GetNowUnix()
	expire_time := expire_in+time_now
	branch.ExpireTime = models.GetTime(expire_time)
	branch.UpdateTime = models.GetNowTime()

	//gengxin
	if err := models.DB.Debug().Model(Branch{}).Where(map[string]interface{}{"corp_id": corp_id, "del": 0}).Update(&branch).Error; err != nil {
		logging.Info(err)
	}
	arr["access_token"] = m["access_token"].(string)
	arr["expire_time"] = branch.ExpireTime
	return
}
func checkToken(token string) bool {
	url := base_url + dep_list + "access_token=" + token

	var m map[string]interface{}
	data, err := models.GetData(url)
	if err = json.Unmarshal([]byte(data), &m); err != nil {
		return false
	}
	errcode := m["errcode"].(float64)

	if errcode == 40014 {
		return false
	}
	return true
}
func GetUserInfo(branch_id int,code string) (map[string]interface{},error) {

	branch_info,err := GetAppInfo(string(branch_id))
	if err != nil {
		return nil,err
	}

	url := base_url+user_cgi+"access_token="+branch_info.AccessToken+"&code="+code

	user_data,err := models.GetData(url)
	if err != nil {
		return nil,err
	}

	user_info,err := models.JsonDecode(user_data)
	if err != nil {
		return nil,err
	}
	if user_info["errcode"] == 42001 {
		token_result,err := UpdateToken(branch_info.CorpId,branch_info.Secret)
		if err != nil {
			return nil,err
		}
		url = base_url+user_cgi+"access_token="+token_result["access_token"].(string)+"&code="+code
		user_data,err = models.GetData(url)
		if err != nil {
			return nil,err
		}
	}
	var m map[string] interface{}
	m = make(map[string] interface{})

 	m["user_info"] = user_info
	m["access_token"] = branch_info.AccessToken

	return m,nil

}
