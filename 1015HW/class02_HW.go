package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
)

// Todo 結構
type Todo struct {
	ID   uint   `json:"id" gorm:"primaryKey"` // Todo 項目 ID
	Todo string `json:"todo"`                 // Todo 項目內容
	Done bool   `json:"done"`                 // Todo 是否完成
}

var db *gorm.DB

// 初始化數據庫
func initDB() {
	// 從環境變數獲取資料庫名稱
	dbName := os.Getenv("DB_NAME") //$env:DB_NAME="todos2.db"
	if dbName == "" {
		dbName = "todos.db" // 預設資料庫名稱
	}
	fmt.Println("連接數據庫：", dbName)
	var err error
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("無法連接到數據庫！")
	}
	println("連接數據庫成功！")
	db.AutoMigrate(&Todo{})
}

// 添加 Todo 項目
func addTodoHandler(c *gin.Context) {
    todo := c.Query("todo")
    if todo == "" {
        // 返回錯誤訊息
        c.JSON(http.StatusBadRequest, NewAppError(Error6001, nil)) // 任務內容不能為空
        return
    }

    // 獲取狀態參數
    state := c.Query("state")
    var done bool
    if state == "true" {
        done = true
    } else if state == "false" {
        done = false
    } else {
        // 如果 state 不是 "done" 或 "pending"，返回錯誤
        c.JSON(http.StatusBadRequest, NewAppError(Error6003, nil)) // 狀態無效
        return
    }

    // 新增待辦事項
    newTodo := Todo{Todo: todo, Done: done}
    if err := db.Create(&newTodo).Error; err != nil {
        // 返回創建失敗的錯誤
        c.JSON(http.StatusInternalServerError, NewAppError(Error6002, nil)) // 創建 Todo 項目失敗
        return
    }

    c.JSON(http.StatusOK, NewSuccessResponse(newTodo))
}


// 讀取 Todo 列表
func readTodosHandler(c *gin.Context) {
	var todos []Todo
	state := c.Query("state")

	if state == "done" {
		db.Where("done = ?", true).Find(&todos)
	} else if state == "pending" {
		db.Where("done = ?", false).Find(&todos)
	} else {
		db.Find(&todos)
	}

	c.JSON(http.StatusOK, NewSuccessResponse(todos))
}

// 改變 Todo 狀態並顯示所有 Todo
func markTodoHandler(c *gin.Context) {
	var todo Todo
	id := c.Query("id")

	if err := db.First(&todo, id).Error; err != nil {
		// 返回未找到的錯誤
		c.JSON(http.StatusNotFound, NewAppError(Error6005, nil)) // Todo 項目未找到
		return
	}

	state := c.Query("state")
	if state == "true" {
		todo.Done = true
	} else if state == "false" {
		todo.Done = false
	}
	if err := db.Save(&todo).Error; err != nil {
		// 返回更新失敗的錯誤
		c.JSON(http.StatusInternalServerError, NewAppError(Error6004, nil)) // 更新 Todo 項目失敗
		return
	}

	var todos []Todo
	db.Find(&todos)
	c.JSON(http.StatusOK, NewSuccessResponse(todos))
}

// 刪除 Todo 項目
// 刪除Todo項目
func removeTodoHandler(c *gin.Context) {
	id := c.Query("id")
	if err := db.Delete(&Todo{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到Todo。"})
		return
	}

	// 刪除後自動重新編排剩餘的Todo項目的ID
	var todos []Todo
	db.Find(&todos) // 獲取當前所有Todo項目

	// 重新設定ID
	for i, todo := range todos {
		todo.ID = uint(i + 1) // ID從1開始
		db.Save(&todo)        // 更新ID
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo已成功刪除。"})
}

func class02_HW() {
	initDB()
	r := gin.Default()
	api := r.Group("/api")

	api.POST("/todo", addTodoHandler)        // 添加 Todo
	api.GET("/read", readTodosHandler)       // 讀取 Todo 列表
	api.POST("/mark", markTodoHandler)       // 改變 Todo 狀態
	api.DELETE("/remove", removeTodoHandler) // 刪除 Todo

	if err := r.Run(":8082"); err != nil {
		fmt.Printf("啟動服務器失敗: %s\n", err)
	} else {
		fmt.Println("服務器運行在端口8082上")
	}
}
