package go_redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RedisC struct {
	Conn redis.Conn
}

var RedisClient *RedisC

func InitRedis(serverHot string, password string, dbIndex int) error {
	setDb := redis.DialDatabase(dbIndex)
	setPassword := redis.DialPassword(password)
	conn, err := redis.Dial("tcp", serverHot, setDb, setPassword)
	if err != nil {
		return err
	}
	RedisClient = &RedisC{
		Conn: conn,
	}
	fmt.Println("redis connect successfully")
	return nil
}
func RedisClose() error {
	fmt.Println("redis close")
	err := RedisClient.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

//check exists
func (r *RedisC) ExistsValue(key string) bool {
	exists, _ := r.Conn.Do("EXISTS", key)
	flag, _ := redis.Int(exists, nil)
	fmt.Println("是否存在指定key，",key, flag == 1)
	return flag == 1
}

// set
func (r *RedisC) SetValue(key string, value interface{}) error {
	// 带国旗时间的k-v，Ex单位是s  set k v expire time
	if _, err := r.Conn.Do("SET", key, value); err != nil {
		return err
	}
	return nil
}

//get
func (r *RedisC) GetValue(key string) (interface{}, error) {
	if isExist := r.ExistsValue(key); !isExist {
		return nil, nil
	}
	v, err := r.Conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return v, nil
}

//del
func (r *RedisC) DelValue(key string) error {
	if isExist := r.ExistsValue(key); isExist {
		if _, err := r.Conn.Do("DEL", key); err != nil {
			return err
		}
	}
	return nil
}
