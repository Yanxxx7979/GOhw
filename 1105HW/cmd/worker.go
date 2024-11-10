package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
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
	var err error
	db, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連接到數據庫！", err)
	}
	db.AutoMigrate(&Todo{})
	fmt.Println("數據庫初始化完成！")
}

// rootCmd represents the base command when called without any subcommands
var workCmd = &cobra.Command{
	Use:   "worker",
	Short: "Todo Worker",
	Long:  `Will pop a job from the queue and process it`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("worker called")

		// 初始化 Redis 客戶端
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})

		// 初始化數據庫
		initDB()

		// 不斷從 Redis 列隊中阻塞地取出數據
		for {
			job, err := client.BRPop(context.Background(), 0, "test1029").Result()
			if err != nil {
				fmt.Println("取出失敗:", err)
				continue
			}
			fmt.Println("Job :", job)

			// 解析 job JSON 字符串為 Todo 結構
			var todo Todo
			err = json.Unmarshal([]byte(job[1]), &todo)
			if err != nil {
				log.Println("解析失敗:", err)
				continue
			}

			// 確認 db 已經被初始化
			if db == nil {
				log.Fatal("數據庫未初始化！")
			}

			// 將解析後的 todo 保存到數據庫中
			if err := db.Create(&todo).Error; err != nil {
				log.Println("插入數據庫失敗:", err)
				continue
			}

			fmt.Println("已成功將 job 插入數據庫:", todo)
		}
	},
}

func init() {
	rootCmd.AddCommand(workCmd)
}
