/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gogo/todo"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "call class002_HW", // api -h 會顯示更詳細的說明
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("config") // 獲取 port 的值
		if port == "" { // 如果 port 沒有設置
			port = "8082" // 設置為預設值 8081
		}
		fmt.Printf("api called, using port: %s\n", port) // 輸出使用的端口
		todo.Class02_HW(port) // 將 port 傳遞給 Class02_HW
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	apiCmd.Flags().StringP("config", "c", "8082", "port input")
}
