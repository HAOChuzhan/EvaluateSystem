package v1
import (
	"github.com/gin-gonic/gin"
	"jim_evaluate/models"
	paper2 "jim_evaluate/models/paper"
	"jim_evaluate/pkg/app"
	"jim_evaluate/pkg/logging"
)

func GetPaper(c *gin.Context) {
	appG := app.Gin{C: c}
	identity := c.DefaultPostForm("type","")

	if identity == "" {
		appG.Response(500,"缺少身份类型",nil)
		return
	}
	var m[]paper2.Paper
	if err := models.DB.Debug().Where(map[string]interface{}{"type": identity, "del": 0}).Find(&m).Error; err != nil {
		logging.Info(err)
	}
	appG.Response(200,"操作成功",m)
	return

}