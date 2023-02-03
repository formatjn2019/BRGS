package web_test

import (
	"BRGS/routers"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type Button struct {
	Text string
	Url  string
}

func TestRouter(t *testing.T) {
	router := routers.InitRouter()
	router.Run(":80")
}

func TestTemplates(t *testing.T) {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html", gin.H{
			"title": "备份文件",
			// "single": "Main website",
			"data": []Button{
				{Text: "按钮1", Url: "test"},
				{Text: "按钮2", Url: "test"},
				{Text: "按钮3", Url: "test"},
				{Text: "按钮4", Url: "test"},
				{Text: "按钮5", Url: "test2"},
				{Text: "按钮6", Url: "test2"},
				{Text: "按钮7", Url: "test2"},
			},
		})
	})

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"fff": []string{"fffg", "gnn"},
		})
	})
	router.POST("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"fff": "",
		})
	})
	router.POST("/logs", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"log": strings.Repeat("logtext ffffffffffffffffffffffffffggggggggggff\n", 20),
		})
	})
	router.Run(":80")
}

type LeftAside struct {
	NavTree NavTree   `json:"navtree"`
	NavItem []NavItem `json:"navitem"`
}

type NavTree struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type NavItem struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

func Test22(t *testing.T) {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/ping", func(c *gin.Context) {
		ipAddr := c.ClientIP()
		fmt.Println(ipAddr)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/test", func(c *gin.Context) {

		lfetasides := []LeftAside{ //方法1：初始化结构体
			{NavTree: NavTree{Name: "startpages", Icon: "nav-icon fas fa-tachometer-alt"},
				NavItem: []NavItem{{Href: "/aaa/bbb", Name: "功能一"}, {Href: "/aaa/ccc", Name: "功能二"}}},
			{NavTree: NavTree{Name: "mainpages", Icon: "nav-icon fas fa-tachometer-alt"},
				NavItem: []NavItem{{Href: "/aaa/eee", Name: "功能三"}, {Href: "/aaa/fff", Name: "功能四"}}},
		}

		aaa := []LeftAside{}
		for i := 0; i <= 10; i++ { //方法2：循环给结构体赋值

			aaa = append(aaa, LeftAside{NavTree: NavTree{Name: "startpages" + strconv.Itoa(i), Icon: "nav-icon fas fa-tachometer-alt"},
				NavItem: []NavItem{{Href: "/aaa/bbb", Name: "功能一"}, {Href: "/aaa/ccc", Name: "功能二"}}})

		}
		dataM2, err := json.Marshal(aaa) //测试结构体转化为json
		if err != nil {
			fmt.Printf("序列化错误 err = %v\n", err)
		}
		fmt.Println(string(dataM2))
		c.HTML(http.StatusOK, "test.tmpl", gin.H{
			"lists":   lfetasides, //传递给前端的值
			"jsonres": string(dataM2),
			"tt":      aaa,
		})
	})

	r.Run(":8062")

}
