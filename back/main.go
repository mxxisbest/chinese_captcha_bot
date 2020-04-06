package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

//Response 就是个回应（合规性注释
type Response struct {
	Code    int
	Message string
	URL     string
	ANS     string
}

type redisPoolConf struct {
	maxIdle        int
	maxActive      int
	maxIdleTimeout int
	host           string
	password       string
	db             int
	handleTimeout  int
}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const defaultPort int = 8002
const defaultRedisConfig = "redis:6379"
const defaultRedisPasswd = ""

var redisPool *redis.Pool
var redisPoolConfig *redisPoolConf
var redisClient redis.Conn

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"https://www.is.sy","http://www.is.sy","https://vip.is.sy","http://vip.is.sy"}
	router.Use(cors.New(config))

	port := flag.Int("port", defaultPort, "服务端口")
	conn := flag.String("conn", defaultRedisConfig, "Redis连接，格式: host:port")
	passwd := flag.String("passwd", defaultRedisPasswd, "Redis连接密码")

	flag.Parse()

	redisPoolConfig = &redisPoolConf{
		maxIdle:        1024,
		maxActive:      1024,
		maxIdleTimeout: 30,
		host:           *conn,
		password:       *passwd,
		db:             1,
		handleTimeout:  30,
	}
	initRedisPool()

	router.POST("/", func(context *gin.Context) {
		redisClient = redisPool.Get()
		defer redisClient.Close()

		res := &Response{
			Code:    1,
			Message: "",
			URL:     "",
			ANS:     "",
		}

		URL := context.PostForm("URL")

		if URL == "" {
			res.Code = 0
			res.Message = "图片地址为空"
			context.JSON(400, *res)
			return
		}

		ANS := context.PostForm("ANS")

		if ANS == "" {
			res.Code = 0
			res.Message = "答案为空"
			context.JSON(400, *res)
			return
		}

		log.Println(URL, ANS)

		res.URL = URL

		if ANS != "" {
			log.Println("开始写入", URL, ANS)
			_, _ = redisClient.Do("mset", URL, ANS)

			res.ANS = URL + ":" + ANS

			context.JSON(200, *res)
		}
	})

	router.GET("/", func(context *gin.Context) {

		URL := getchallenge()

		if URL == "" {
			context.String(http.StatusNotFound, "有点不对劲")
		} else {
			context.JSON(200, URL)
		}
	})

	router.Run(fmt.Sprintf(":%d", *port))
}

// 随机来一个题目
func getchallenge() string {
	redisClient = redisPool.Get()
	defer redisClient.Close()

	URL = redisClient.Do("RANDOMKEY")
	ANS, _ := redis.String(redisClient.Do("get", URL))

	return URL + ":" + ANS
}

func initRedisPool() {
	// 建立连接池
	redisPool = &redis.Pool{
		MaxIdle:     redisPoolConfig.maxIdle,
		MaxActive:   redisPoolConfig.maxActive,
		IdleTimeout: time.Duration(redisPoolConfig.maxIdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", redisPoolConfig.host,
				redis.DialPassword(redisPoolConfig.password),
				redis.DialDatabase(redisPoolConfig.db),
				redis.DialConnectTimeout(time.Duration(redisPoolConfig.handleTimeout)*time.Second),
				redis.DialReadTimeout(time.Duration(redisPoolConfig.handleTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(redisPoolConfig.handleTimeout)*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}
