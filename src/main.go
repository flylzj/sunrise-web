package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"model"
	"net/http"
	"path"
	"resource"
	"spider"
	"strconv"
	"strings"
	"time"
	"util"
)

func updateTimeConf(ctx *gin.Context){
	conf := model.Conf{}
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(data, &conf)
	result := resource.SetTimeConf(conf.IntervalHour, conf.IntervalMinute, conf.GoodCountInterval)
	if result{
		ctx.JSON(200, gin.H{
			"message": "ok",
			"code": 0,
		})
	}else {
		ctx.JSON(200, gin.H{
			"message": "分钟或物品间隔不能为0",
			"code": 1,
		})
	}
}

func updateEmailConf(ctx *gin.Context){
	conf := model.Conf{}
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(data, &conf)
	resource.SetEmailConf(conf.Sender, conf.SenderPwd, conf.Receiver)
	ctx.JSON(200, gin.H{
		"message": "ok",
		"code": 0,
	})
}

func getTimeConf(ctx *gin.Context){
	hour, minute, count := resource.GetTimeConf()
	ctx.JSON(200, gin.H{
		"message": "ok",
		"code": 0,
		"data": gin.H{
			"hour": hour,
			"minute": minute,
			"count": count,
		},
	})
}

func getEmailConf(ctx *gin.Context){
	sender, sender_pwd, receiver := resource.GetEmailConf()
	ctx.JSON(200, gin.H{
		"message": "ok",
		"code": 0,
		"data": gin.H{
			"sender": sender,
			"sender_pwd": sender_pwd,
			"receiver": receiver,
		},
	})
}

func getGoodHistory(ctx *gin.Context){
	abiid := ctx.Param("abiid")
	ctx.JSON(200, gin.H{
		"message": "ok",
		"code": 0,
		"data": resource.GetGoodHistory(abiid),
	})
}

func addGood(ctx *gin.Context){
	good := model.GoodBeMonitored{}
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(data, &good)
	result := resource.AddGood(good.Abiid)
	var (
		message string
		code int
	)
	if result == "添加成功"{
		message = "ok"
		code = 0
	}else {
		message = result
		code = 1
	}
	ctx.JSON(200, gin.H{
		"message": message,
		"code": code,
	})
}

func deleteGood(ctx *gin.Context){
	good := model.GoodBeMonitored{}
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(data, &good)
	result := resource.DeleteAGood(good.Abiid)
	var (
		code int
	)
	if result == "ok"{
		code = 0
	}else {
		code = 1
	}
	ctx.JSON(200, gin.H{
		"message": result,
		"code": code,
	})
}

func getBeMonitoredGoods(ctx *gin.Context){
	ctx.JSON(200, gin.H{
		"message": "ok",
		"code": 0,
		"data": resource.GetBeMonitoredGoods(),
	})
}

func getExcelFile(ctx *gin.Context){
	file, err := ctx.FormFile("file")
	if err != nil{
		ctx.JSON(200, gin.H{
			"message": "bad request",
			"code": 1,
		})
	}else {
		result := util.FilterExcel(file.Filename)
		if result != "ok"{
			ctx.JSON(200, gin.H{
				"message": result,
				"code": 1,
			})
		}else {
			util.CreatePath("excel")
			uploadFile :=  path.Join("data", "excel", strconv.Itoa(int(time.Now().Unix())) + ".xlsx")
			ctx.SaveUploadedFile(file,uploadFile)
			results := resource.AddGoodInBatches(uploadFile)
			ctx.JSON(200, gin.H{
				"message": "ok",
				"data": results,
				"code": 0,
			})
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method      //请求方法
		origin := c.Request.Header.Get("Origin")        //请求头部
		var headerKeys []string                             // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")        // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")      //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")      // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")        // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")       //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")       // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()        //  处理请求
	}
}

func StaticRouter(){
	http.Handle("/", http.FileServer(http.Dir("vue-static")))
	http.ListenAndServe(":8081",nil)
}

func StartRearEndRouter(){
	router := gin.Default()
	router.Use(Cors()) // 跨域
	router.GET("/api/time_interval", getTimeConf)
	router.POST("/api/time_interval", updateTimeConf)
	router.GET("/api/email", getEmailConf)
	router.POST("/api/email", updateEmailConf)
	router.GET("/api/history/:abiid", getGoodHistory)
	router.GET("/api/good", getBeMonitoredGoods)
	router.POST("/api/good", addGood)
	router.DELETE("/api/good", deleteGood)
	router.POST("/api/good/upload", getExcelFile)
	router.Run(":8080")
}

func main(){
	go spider.MainSpider()
	go StaticRouter()
	model.Info.Println("静态资源加载完成")
	model.Info.Println("若系统没有自动打开，请手动在浏览器输入http://127.0.0.1:8081")
	StartRearEndRouter()
}