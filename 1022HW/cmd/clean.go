// cmd/clean.go
package cmd

import (
	"fmt"
	"gogo/todo"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清除資料庫中的所有 Todo 項目",
	Long:  `清除資料庫中的所有 Todo 項目，這個操作會永久刪除資料庫中的所有記錄。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("開始清除資料庫中的所有 Todo 項目...")
		err := todo.ClearDatabase()
		if err != nil {
			fmt.Println("清除失敗:", err)
		} else {
			fmt.Println("資料庫已成功清除！")
		}
	},
}

func init() {
	fmt.Println("Registering clean command")
	rootCmd.AddCommand(cleanCmd) // 添加 cleanCmd
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
