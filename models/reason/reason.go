package reason

import (
	"jim_evaluate/models"
	"jim_evaluate/pkg/logging"
)

type Reason struct {
	models.Model

	Title   string `json:"title"`
	AdminId int    `json:"admin_id"`
	UserId  int    `json:"user_id"`
	OrderId int    `json:"order_id"`
}

//提交驳回原因
func CreateReason(reason *Reason) bool {
	if err := models.DB.Debug().Model(&Reason{}).Create(reason).Error; err != nil {
		logging.Info(err)
		return false
	}
	return true
}

//查看驳回原因
func IndexReason(order_id int) (data []Reason, err error) {
	if err := models.DB.Debug().Model(&Reason{}).Order("add_time DESC").Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
