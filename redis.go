/*
 * @Author       : jayj
 * @Date         : 2021-06-19 21:43:08
 * @Description  : redis db connection
 */
package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

type RConn struct {
	// Redis host
	Host string
	// Redis password
	Password string

	// target db
	DB int

	// idle timeout
	Timeout int

	MaxIdle int

	MaxActive int

	// connection will wait if not active idle left
	Wait bool
}

type RConnOption func(*RConn)

// newRedisPool 创建redis连接池
// db 要连接的数据库
func NewRedisPool(host string, options ...RConnOption) {

	rConn := &RConn{Host: host, Password: "", DB: 0, Timeout: 300, MaxIdle: 100, MaxActive: 50, Wait: false}

	for _, option := range options {
		option(rConn)
	}

	// 连接的db
	option := redis.DialDatabase(rConn.DB)

	password := redis.DialPassword(rConn.Password)

	pool = &redis.Pool{
		// 最大连接数
		MaxIdle: rConn.MaxIdle,
		// 最大活跃连接数
		MaxActive: rConn.MaxActive,
		// 在这个时间之后关闭idle
		IdleTimeout: time.Duration(rConn.Timeout) * time.Second,
		// 如果没有active就等待
		Wait: rConn.Wait,
		Dial: func() (redis.Conn, error) {
			// 1. 打开连接
			var err error

			c, err := redis.Dial("tcp", host, option, password)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("initialize redis failed, err: %s", err))
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
	}
}

// GetRedisPool get redis pool
func GetRedisPool() *redis.Pool {
	return pool
}

// WithDB target connect database
func WithDB(db int) RConnOption {
	return func(r *RConn) {
		r.DB = db
	}
}

// WithPassword redis password
func WithPassword(password string) RConnOption {
	return func(r *RConn) {
		r.Password = password
	}
}

// WithTimeout active idle will break if there is no request/response during timeout (second)
func WithTimeout(timeout int) RConnOption {
	return func(r *RConn) {
		r.Timeout = timeout
	}
}

// WithIdle maxIdle => maximum idle number, maxActive => maximum active idle number
func WithIdle(maxIdle, maxActive int) RConnOption {
	return func(r *RConn) {
		r.MaxActive = maxActive
		r.MaxIdle = maxIdle
	}
}

// WithWait request will hold if redis pool has no active idle left
func WithWait(wait bool) RConnOption {
	return func(r *RConn) {
		r.Wait = wait
	}
}
