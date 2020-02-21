package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

var DB *sql.DB

type result struct {
	Success    bool   `json:"success"`
	Reason     string `json:"reason,omitempty"`
	Name       string `json:"name,omitempty"`
	LastAction int64  `json:"lastAction,omitempty"`
}

func main() {
	var err error
	DB, err = sql.Open("mysql", os.Getenv("UNIBRIDGE_DSN"))
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}
	DB.SetConnMaxLifetime(100 * time.Second)
	DB.SetMaxOpenConns(100)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "UNIBRIDGE for XJCRAFT")
	})
	r.GET("/checkpass", checkPass)
	r.Run(os.Getenv("UNIBRIDGE_LISTEN"))
}

func checkPass(c *gin.Context) {
	name := c.Query("name")
	pass := c.Query("pass")
	if name == "" {
		c.JSON(200, result{Success: false, Reason: "用户名为空"})
		return
	}
	if pass == "" {
		c.JSON(200, result{Success: false, Reason: "密码为空"})
		return
	}
	// query db
	var password string
	var lastAction time.Time
	var loginFails int
	err := DB.QueryRow("SELECT `password`,`lastAction`,`loginFails` FROM CrazyLogin_accounts WHERE `name`= ?", name).Scan(&password, &lastAction, &loginFails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(200, result{Success: false, Reason: "用户不存在"})
		} else {
			c.JSON(200, result{Success: false, Reason: "系统错误"})
		}
		return
	}
	if loginFails > 3 {
		c.JSON(200, result{Success: false, Reason: "尝试次数过多"})
		return
	}
	// check pass
	h := sha256.New()
	h.Write([]byte(name))
	h.Write([]byte(password))
	passHash := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(passHash, pass) {
		DB.Exec("UPDATE CrazyLogin_accounts SET loginFails=loginFails+1 WHERE `name`=?", name)
		c.JSON(200, result{Success: false, Reason: "密码错误"})
		return
	}
	c.JSON(200, result{Success: true, Name: name, LastAction: lastAction.Unix()})
}
