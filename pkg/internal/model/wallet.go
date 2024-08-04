package model

import "gorm.io/gorm"

type User struct {
	ID      uint   `gorm:"primarykey"`
	UserID  uint32 `gorm:"index"`
	Balance int64  `gorm:"type:BIGINT"`
}

func UpdateBalance(uid uint32, amount int64, customDB *gorm.DB) error {

	execRes := customDB.Exec("UPDATE users SET balance = ? WHERE user_id = ?", amount, uid)
	if execRes.Error != nil {
		return execRes.Error
	}
	if execRes.RowsAffected == 0 {
		execRes = customDB.Exec("INSERT INTO users (balance, user_id) VALUES (?,?);", amount, uid)
		if execRes.Error != nil {
			return execRes.Error
		}
	}
	return nil
}
