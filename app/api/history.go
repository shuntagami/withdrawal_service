package api

import (
	"api/db"
	"api/model"
	"database/sql"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var m sync.Mutex

func CreateHistory(c *gin.Context) {
	m.Lock()
	defer m.Unlock()

	type request struct {
		UserID int `json:"user_id"`
		Amount int `json:"amount"`
	}
	req := request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	history := model.History{
		UserID: req.UserID,
		Amount: req.Amount,
	}

	// ユーザーの合計出金金額を取得
	var sumAmount sql.NullInt64
	if err := db.Conn.Raw("SELECT SUM(amount) FROM histories WHERE user_id = ?", req.UserID).Scan(&sumAmount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	// 最大出金金額を超えていたら400を返す
	if int(sumAmount.Int64)+req.Amount > model.AmountLimit {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "over the amount limit 100,000"})
		return
	}

	// 出金履歴登録
	if err := db.Conn.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, history)
}
