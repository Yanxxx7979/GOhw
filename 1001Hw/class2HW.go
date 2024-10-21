package main

import (
	"fmt"
	"github.com/gin-gonic/gin"

	// "encoding/json"

)
type Todo struct {
	Todo string `json:"To do"`
	Done bool `json:"Done"`
}
var globalTodos []Todo

func class02_HW() {
	r := gin.Default()
	api := r.Group("/api")
	//添加todo list ->api1: http://localhost:8008/api/todo?todo=
	api.GET("/todo", func(c *gin.Context) {
		todo := c.Query("todo")

		if todo != "" {
			globalTodos = append(globalTodos, Todo{Todo: todo, Done: false})			
		}
		c.JSON(200, gin.H{
			"todos": globalTodos,
		})
	})
	//印出 pending/done/all 的list -> api2: http://localhost:8008/api/read?state= pending / done / all
	api.GET("/read", func(c *gin.Context) {
		state := c.Query("state")
		var filteredTodos []Todo
		if state == "done" {
			for _, todo := range globalTodos {
				if todo.Done { // 只添加 Done 为 true 的待办事项
					filteredTodos = append(filteredTodos, todo)
					
				}
			}
		} else if state == "pending" {
			for _, todo := range globalTodos {
				if !todo.Done { // 只添加 Done 为 false 的待办事项
					filteredTodos = append(filteredTodos, todo)
				}
			}
		} else if state == "all" {
			filteredTodos = globalTodos
			c.JSON(200, gin.H{
				"All list": filteredTodos,
			})
        }
		
	})
	//改變todo狀態 -> api3 http://localhost:8008/api/mark?todo=Running&state=done
	api.GET("/mark", func(c *gin.Context){
	todoName := c.Query("todo") 
	state := c.Query("state") 
		for i, todo := range globalTodos {
			if todo.Todo == todoName {
				if state == "done" {
					globalTodos[i].Done = true 
				} else if state == "pending" {
					globalTodos[i].Done = false 
				}
				break 
			}
		}
		c.JSON(200, gin.H{
			"todos": globalTodos,
		})
	
	})
	//刪除todo項目 -> api4：http://localhost:8008/api/remove?task=Task1
	api.GET("/remove", func(c *gin.Context) {
		task := c.Query("task") 
		if task == "" {
			c.JSON(400, gin.H{"error": "Task parameter is required."})
			return
		}
		var updatedTodos []Todo
		for _, todo := range globalTodos {
			if todo.Todo != task {
				updatedTodos = append(updatedTodos, todo) // 保留未被移除的
			}
		}
		globalTodos = updatedTodos 

		c.JSON(200, gin.H{
			"todos": globalTodos,
		})
	})


	err := r.Run(":8008")
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	} else {
		fmt.Println("Server is running on port 8008")
	}
}
