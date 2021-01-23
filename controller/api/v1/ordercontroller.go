package v1

import (
	"fmt"
	"jim_evaluate/models"
	"jim_evaluate/models/order"
	"jim_evaluate/pkg/app"
	"jim_evaluate/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func CommitOrder(c *gin.Context)  {

	appG := app.Gin{C: c}
	order_string := c.DefaultPostForm("order","")
	if order_string == "" {
		appG.Response(500,"缺少表单",nil)
		return
	}
	order_info,err := models.JsonDecode2(order_string)
	if err != nil {
		appG.Response(500,"json解析失败",order_info)
	}
	score := 0.0;
	for _,v := range order_info{
		doNum := v["doNum"].(map[string]interface{})
		targetNum := v["targetNum"].(map[string]interface{})
		doNum_value := doNum["value"].(float64)
		targetNum_value := targetNum["value"].(float64)
		completeProgress := doNum_value / targetNum_value

		tagetPoints := v["tagetPoints"].(map[string]interface{})
		tagetPoints_value := tagetPoints["value"].(float64)

		getPoints := v["getPoints"].(map[string]interface{})
		getPoints_value := getPoints["value"].(float64)

		getPoints_new_value := completeProgress * tagetPoints_value
		getPoints_new_value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", getPoints_new_value), 64)
		getPoints_value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", getPoints_value), 64)
		if getPoints_value != getPoints_new_value {
			appG.Response(500,"数值不一样",nil)
		}
		score += getPoints_new_value

	}
	date_time := models.GetAtTime("2006-01")
	order_model := order.Order{
		UserID: 1,
		Content: order_string,
		State: 0,
		Date: date_time,
		Score: score,
	}
	if err := models.DB.Debug().Create(&order_model).Error; err != nil {
		logging.Info(err)
	}
	appG.Response(200,"添加成功",order_info)
	return
}
func LastMonthOrder(c *gin.Context)  {
	appG := app.Gin{C: c}
	user_id := c.DefaultPostForm("user_id","")
	if user_id == "" {
		appG.Response(500,"缺少用户id",nil)
		return
	}
	time_date := models.GetAtTime("2006-01")
	order_info := order.Order{}

	if err := models.DB.Debug().Where(map[string]interface{}{"user_id":user_id,"date":time_date}).Find(&order_info).Error ; err != nil {
		logging.Info(err)
	}
	appG.Response(200,"获取成功",order_info)
	return

}
func GetOrderInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	order_id := com.StrTo(c.DefaultPostForm("order_id", "")).MustInt()
	if order_id == 0 {
		appG.Response(http.StatusBadRequest, "缺少表单ID", nil)
		return
	}
	data, err := order.GetOrderDetail(order_id)
	if err != nil {
		appG.Response(http.StatusBadRequest, "获取订单详情失败", nil)
		return
	}
	appG.Response(http.StatusOK, "获取订单详情成功", data)

}
func GetOrderList(c *gin.Context) {
	appG := app.Gin{C: c}
	date := c.DefaultPostForm("date", "")
	page := com.StrTo(c.DefaultPostForm("page", "1")).MustInt()
	limit := com.StrTo(c.DefaultPostForm("limit", "10")).MustInt()

	if date == "" {
		appG.Response(http.StatusBadRequest, "缺少参数", nil)
		return
	}
	orderlist := order.GetOrderList(date, page, limit)
	appG.Response(http.StatusOK, "获取表单列表成功", orderlist)
}
func UpdateOrderState(c *gin.Context) {
	appG := app.Gin{C: c}
	order_id := com.StrTo(c.DefaultPostForm("order_id", "")).MustInt()
	admin_id := com.StrTo(c.DefaultPostForm("admin_id", "")).MustInt()
	state := com.StrTo(c.DefaultPostForm("state", "")).MustInt()

	if order_id == 0 || admin_id == 0 || state == 0 {
		appG.Response(http.StatusBadRequest, "缺少参数", nil)
		return
	}
	order_new := order.Order{
		AdminID: admin_id,
		State:   state,
	}
	if !order.UpdateOrder(order_id, &order_new) {
		appG.Response(http.StatusBadRequest, "更新表单状态失败", nil)
		return
	}
	appG.Response(http.StatusOK, "更新表单状态成功", nil)


}
