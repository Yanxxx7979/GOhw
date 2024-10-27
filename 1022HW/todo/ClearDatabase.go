package todo

import (
	"fmt"
	// "gorm.io/gorm"
)

// var db *gorm.DB

// ClearDatabase 刪除資料庫中的所有 Todo 記錄
func ClearDatabase() error {
	initDB()
	if db == nil {
		return fmt.Errorf("資料庫尚未初始化")
	}

	// 刪除所有 Todo 記錄
	if err := db.Exec("DELETE FROM todos").Error; err != nil {
		return err
	}

	fmt.Println("資料庫已成功清空")
	return nil
}
