package main

import (
	"database/sql"
	"fmt"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "127.0.0.1"
	database = "testdbdll"
	user     = "im"
	password = "vie5dfs4bhers"
)

type Member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()

	router.GET("/gettest", func(c *gin.Context) {
		var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
		db, err := sql.Open("mysql", connectionString)
		checkError(err)
		defer db.Close()

		err = db.Ping()
		checkError(err)

		rows, err := db.Query("SELECT * FROM member")
		checkError(err)
		defer rows.Close()

		result := []Member{}

		for rows.Next() {
			var member Member
			err := rows.Scan(&member.ID, &member.Name, &member.Phone)
			checkError(err)
			result = append(result, member)
		}

		c.JSON(200, result)
	})

	router.GET("/gettest/:id", func(c *gin.Context) {
		var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
		db, err := sql.Open("mysql", connectionString)
		checkError(err)
		defer db.Close()

		err = db.Ping()
		checkError(err)

		var row Member
		err = db.QueryRow("SELECT * FROM `member` WHERE `id` = ?", c.Param("id")).Scan(&row.ID, &row.Name, &row.Phone)
		checkError(err)

		c.JSON(200, row)
	})

	router.POST("/gettest", func(c *gin.Context) {
		var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
		db, err := sql.Open("mysql", connectionString)
		checkError(err)
		defer db.Close()

		err = db.Ping()
		checkError(err)

		data := make(map[string]interface{})
		err = c.BindJSON(&data)
		checkError(err)

		messages := []string{}
		if data["id"] == nil {
			messages = append(messages, "id 為必填項目")
		}

		if data["name"] == nil {
			messages = append(messages, "名稱為必填項目")
		}

		if len(messages) != 0 {
			msg := strings.Join(messages, ",")
			log.Printf(msg)
			c.JSON(202, gin.H{
				"status": "失敗",
				"msg":    msg,
			})
			return
		}

		var row Member
		err = db.QueryRow("SELECT * FROM `member` WHERE `id` = ?", data["id"]).Scan(&row.ID, &row.Name, &row.Phone)
		checkError(err)

		_, err = db.Exec("UPDATE `member` SET name = ? WHERE id = ?", data["name"], data["id"])
		checkError(err)

		c.JSON(200, gin.H{
			"status": "成功",
			"msg":    "已將客戶[" + row.Name + "]名稱改為[" + data["name"].(string) + "]",
		})
	})

	router.Run(":9176")
}
