package order

import (
	"jim_evaluate/models"
	"jim_evaluate/pkg/logging"
)

type Order struct {
	models.Model

	UserID  int     `json:"user_id"`
	AdminID int     `json:"admin_id"`
	Content string  `json:"Content"`
	State   int     `json:"state"`
	Date    string  `json:"date"`
	Score   float64 `json:"score"`
}

func GetOrderDetail(order_id int) (Order, error) {
	orderDetail := Order{}
	if err := models.DB.Debug().Model(&Order{}).Where("id = ?", order_id).First(&orderDetail).Error; err != nil {
		logging.Info(err)
		return orderDetail, err
	}
	return orderDetail, nil

}
func GetOrderList(date string, page, limit int) (data map[string]interface{}) {
	var total int
	var orders []Order
	db := models.DB.Debug().Model(&Order{}).Where(map[string]interface{}{"date": date, "del": 0})
	if err := db.Count(&total).Error; err != nil {
		return
	}
	offset := (page - 1) * limit
	if err := db.Order("add_time DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return
	}
	data = map[string]interface{}{
		"orders":    orders,
		"total":     total,
		"page":      page,
		"pagetotal": int(total/limit) + 1,
	}
	return data
}
func CreateOrder(order *Order) bool {
	if err := models.DB.Debug().Create(&order).Error; err != nil {
		return false
	}
	return true
}
func UpdateOrder(order_id int, order *Order) bool {
	if err := models.DB.Debug().Model(&Order{}).Where(map[string]interface{}{"del": 0, "id": order_id}).Update(order).Error; err != nil {
		return false
	}
	return true
}

