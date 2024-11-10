package todo

import (
	"fmt"
	// "gorm.io/gorm"
)

// var db *gorm.DB

// ClearDatabase 刪除資料庫中的所有 Todo 記錄並清空表格
func ClearDatabase() error {
	initDB()
	// 确保数据库连接已初始化
	if db == nil {
		return fmt.Errorf("資料庫尚未初始化")
	}

	// 删除所有 Todo 記錄
	if err := db.Exec("DELETE FROM todos").Error; err != nil {
		return fmt.Errorf("清空資料庫時發生錯誤: %v", err)
	}

	// 重置自增 ID (可选)
	if err := db.Exec("DELETE FROM sqlite_sequence WHERE name='todos'").Error; err != nil {
		return fmt.Errorf("重置自增 ID 失敗: %v", err)
	}

	if err := db.Exec("VACUUM").Error; err != nil {
		return fmt.Errorf("執行 VACUUM 操作失敗: %v", err)
	}

	fmt.Println("資料庫已成功清空")
	return nil
}
