package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Todo 結構
type Todo struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"` // 自動遞增的 ID
	Todo string `json:"todo"`                               // Todo 項目內容
	Done bool   `json:"done"`                               // Todo 是否完成
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
	db.AutoMigrate(&Todo{}) // 確保資料表結構正確並且 ID 自動遞增
}

// 添加 Todo 項目到 Redis 队列
func addTodoQueue(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindBodyWithJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	b, err := json.Marshal(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = redisClient.LPush(context.Background(), "test1029", string(b)).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// 添加 Todo 項目到資料庫
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

	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "OK",
		"ret":  todos,
	})
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
func removeTodoHandler(c *gin.Context) {
	id := c.Query("id")

	// 刪除指定ID的Todo
	if err := db.Delete(&Todo{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到Todo。"})
		return
	}

	// 獲取刪除後的所有Todo項目
	var todos []Todo
	if err := db.Order("id asc").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取Todo列表失敗。"})
		return
	}

	// 使用交易重新排列ID
	tx := db.Begin() // 開始交易
	for i, todo := range todos {
		if err := tx.Model(&todo).Update("ID", i+1).Error; err != nil {
			tx.Rollback() // 回滾交易
			c.JSON(http.StatusInternalServerError, gin.H{"error": "重新排序ID失敗。"})
			return
		}
	}
	tx.Commit() // 提交交易

	c.JSON(http.StatusOK, gin.H{"message": "Todo已成功刪除並重新排序ID。"})
}

// 啟動服務
func Class02_HW(port string) {
	initDB()
	r := gin.Default()
	api := r.Group("/api")

	api.POST("/todo", addTodoHandler)        // 添加 Todo
	api.GET("/read", readTodosHandler)       // 讀取 Todo 列表
	api.POST("/mark", markTodoHandler)       // 改變 Todo 狀態
	api.DELETE("/remove", removeTodoHandler) // 刪除 Todo
	api.POST("queue", addTodoQueue)          // 添加 Todo 項目到 Redis 隊列

	if err := r.Run(":" + port); err != nil {
		fmt.Printf("啟動服務器失敗: %s\n", err)
	} else {
		fmt.Println("服務器運行在端口", port)
	}
}
