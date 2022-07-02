package api

import (
	"api/db"
	"api/model"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateHistory(c *gin.Context) {
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

	// 一日の最大出金金額を超えてないかチェックして出金履歴を登録するトランザクション開始
	tx := db.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "panic recovered"})
			return
		}
	}()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": tx.Error.Error()})
		return
	}

	// ユーザーの合計出金金額をSELECT...FOR UPDATE取得
	var sumAmount sql.NullInt64
	if err := tx.Raw("SELECT SUM(amount) FROM histories WHERE user_id = ? FOR UPDATE", req.UserID).Scan(&sumAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	// 最大出金金額を超えていたら400を返す
	if int(sumAmount.Int64)+req.Amount > model.AmountLimit {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"msg": "over the amount limit 100,000"})
		return
	}

	// 出金履歴登録
	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	// トランザクションコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, history)
}