package redis_client

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

var redis_connection *redis.Conn

//redis是但线程模型，读写操作必须加锁
var redis_lock sync.Mutex

func init() {
	redis_host := beego.AppConfig.String("redis_host")
	redis_port := beego.AppConfig.String("redis_port")
	redis_index := beego.AppConfig.String("redis_index")

	conn := redis_host + ":" + redis_port
	rs, err := redis.Dial("tcp", conn)

	if err != nil {
		fmt.Println("redis connect error!")
	} else {
		fmt.Println("redis connect succceed!")
	}
	rs.Do("SELECT", redis_index)
	redis_connection = &rs
}

//
func Set(key, value string, time int) error {
	redis_lock.Lock()
	_, err := (*redis_connection).Do("set", key, value)
	defer redis_lock.Unlock()
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
	redis_lock.Lock()
	ans, err := redis.String((*redis_connection).Do("get", key))
	redis_lock.Unlock()
	return ans, err
}
