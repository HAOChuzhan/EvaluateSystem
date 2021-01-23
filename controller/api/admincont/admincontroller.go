package admincont

import (
	"jim_evaluate/models/order"
	"jim_evaluate/models/reason"
	"jim_evaluate/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

func SubmitReason(c *gin.Context) {
	appG := app.Gin{C: c}

	title := c.DefaultPostForm("title", "")
	admin_id := com.StrTo(c.DefaultPostForm("admin_id", "")).MustInt()
	order_id := com.StrTo(c.DefaultPostForm("order_id", "")).MustInt()

	if title == "" || admin_id == 0 || order_id == 0 {
		appG.Response(http.StatusBadRequest, "缺少参数", nil)
		return
	}
	order_data, err := order.GetOrderDetail(order_id)
	if err != nil {
		appG.Response(http.StatusBadRequest, "获取提交表单详情失败", err)
		return
	}
	user_id := order_data.UserID
	reason_r := reason.Reason{
		Title:   title,
		AdminId: admin_id,
		UserId:  user_id,
		OrderId: order_id,
	}
	if !reason.CreateReason(&reason_r) {
		appG.Response(http.StatusBadRequest, "提交驳回原因失败", nil)
		return
	}
	appG.Response(http.StatusOK, "提交驳回原因成功", nil)
}
func CheckReason(c *gin.Context) {
	appG := app.Gin{C: c}
	order_id := com.StrTo(c.DefaultPostForm("order_id", "")).MustInt()
	if order_id == 0 {
		appG.Response(http.StatusBadRequest, "缺少参数", nil)
		return
	}
	data, err := reason.IndexReason(order_id)
	if err != nil {
		appG.Response(http.StatusBadRequest, "获取驳回原因失败", nil)
		return
	}
	appG.Response(http.StatusOK, "获取原因成功", data)
}
