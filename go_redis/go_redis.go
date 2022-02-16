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
	return  nil
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
	exists,_ := r.Conn.Do("EXISTS",key)
	fmt.Println("ExistsValue",exists)
	flag,_:= redis.Int(exists,nil)
	return flag > 0
}

// set
func (r *RedisC) SetValue(key string, value interface{}) error {
	if _, err := r.Conn.Do("set", key, value); err != nil {
		return err
	}
	return nil
}
//get
func (r *RedisC) GetValue(key string) (interface{}, error) {
	if isExist := r.ExistsValue(key);!isExist{
		return nil,nil
	}
	v, err := r.Conn.Do("get", key)
	if err != nil {
		fmt.Printf("未找到key【%s】\n", key)
		return nil, err
	}
	return v, nil
}
//del
func (r *RedisC) DelValue(key string) error{
	if _,err := r.Conn.Do("del",key);err!=nil{
		return err
	}
	return nil
}

