package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"jim_evaluate/pkg/logging"
	"jim_evaluate/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	//branch_id = ""
	corpid     = "wwee4c443448294ecb"
	agentId    = "1000011"
	corpsecret = "0uo2AX0A-firYFi8kQHCujSeJ4PLlc_13_qL7_-lA2k"
)

//数据库初始化
var DB *gorm.DB

type Model struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Del        int       `json:"del"`
	AddTime    time.Time `gorm:"column:add_time" json:"add_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	DelTime    time.Time `json:"del_time"`
}

// Setup 初始化数据库实例
func Setup() {
	var err error
	DB, err = gorm.Open(setting.DatabaseSetting.Type,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Name))

	// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	// DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败，models.Setup err: %v", err)
	}

	//把回调函数注册进 GORM 的钩子里
	//创建回调
	DB.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//更新回调
	DB.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	// 设置数据表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	//开启单表模式,表名的默认值是结构体名称的小写，如果为false的话默认是小写的复数形式
	DB.SingularTable(true)
	//开启记录模式
	DB.LogMode(true)
	//设置最大空闲连接数
	DB.DB().SetMaxIdleConns(10)
	//设置最大开放连接数
	DB.DB().SetMaxOpenConns(100)
}

//updateTimeStampForCreateCallback 创建时将设置`CreatedOn`，`ModifiedOn`
func updateTimeStampForCreateCallback(scope *gorm.Scope) {

	if !scope.HasError() {
		nowTime := GetNowTime()
		if createTimeField, ok := scope.FieldByName("AddTime"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdatedTime"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback 更新时将设置`ModifiedOn`
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
    if _, ok := scope.Get("gorm:update_column"); !ok {
        scope.SetColumn("UpdatedTime", time.Now().Unix())
    }
}

//get请求
func GetData(url string) (string, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return "调用失败", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), err
}

// json加密
func JsonEncode(m map[string]interface{}) (s string, err error) {

	s_byte, err := json.Marshal(m)
	if err != nil {
		return
	}
	s = string(s_byte)

	return
}

//json解密
func JsonDecode(s string) (m map[string]interface{}, err error) {

	err = json.Unmarshal([]byte(s), &m)
	return
}

//json解密
func JsonDecode2(s string) (m []map[string]interface{}, err error) {

	err = json.Unmarshal([]byte(s), &m)
	return
}

//获取当前时间格式”Y-m-d H:i:s“
func GetNowTime() (t string) {

	time_unix := time.Now().Unix()

	t = time.Unix(time_unix, 0).Format("2006-01-02 15:04:05")

	return t

}

//根据时间戳获取时间格式
func GetTime(uni int64) (t string) {

	time_unix := time.Unix(uni, 0)
	t = time_unix.Format("2006-01-02 15:04:05")

	return t
}

//获取当前时间戳

func GetNowUnix() (t int64) {
	return time.Now().Unix()
}

//httpGET和POST方法
func httpGetJson(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func httpPostJson(url string, data map[string]interface{}) (map[string]interface{}, error) {
	xxx, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(xxx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data2 map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data2, nil
}
func SendCardMsg(ToUsers interface{}, title, description, url string) (map[string]interface{}, error) {
	btntxt := "查看详情"

	qyurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, corpsecret)
	data, err := httpGetJson(qyurl)
	if err != nil {
		logging.Info(err)
		return data, err
	}
	errcode := data["errcode"].(float64)
	if errcode != 0 {
		logging.Info("errcode: ", errcode)
		return make(map[string]interface{}), nil
	}
	access_token := data["access_token"]

	req := map[string]interface{}{
		"touser":  ToUsers,
		"msgtype": "textcard",
		"agentid": agentId,
		"textcard": map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         url,
			"btntext":     btntxt,
		},
	}

	sendurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", access_token)
	data, err = httpPostJson(sendurl, req)
	if err != nil {
		logging.Info(err)
		return nil, err
	}
	return data, nil
}

//根据格式获取当前时间
func GetAtTime(s string) (t string) {
	time_unix := time.Now().Unix()

	t = time.Unix(time_unix, 0).Format(s)

	return t
}
