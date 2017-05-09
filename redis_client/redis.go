package redis_client

import (
	"log"

	"server_test/redis_client/connection_pool"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

func init() {

	redis_host := beego.AppConfig.String("redis_host")
	redis_port := beego.AppConfig.String("redis_port")

	err := connection_pool.SetSize(10)
	if err != nil {
		log.Println(err)
		return
	}
	err = connection_pool.InitConnection(redis_host + ":" + redis_port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("INFO: Redis connect succeed!")
}

//
func Set(key, value string, time int) error {
	redis_connection := connection_pool.GetConnection()
	defer connection_pool.ReturnConnection(redis_connection)

	_, err := (*redis_connection).Do("set", key, value)
	if err == nil {
		if time != 0 {
			// 设置缓存过期时间
			_, err := (*redis_connection).Do("EXPIRE", key, time)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func Get(key string) (string, error) {
	redis_connection := connection_pool.GetConnection()
	defer connection_pool.ReturnConnection(redis_connection)
	ans, err := redis.String((*redis_connection).Do("get", key))

	return ans, err
}
