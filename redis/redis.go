/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    redis
 * @Date:    2021/7/19 10:18 上午
 * @package: RDB
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/recovers"

	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jageros/hawox/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	zredis "github.com/zeromicro/go-zero/core/stores/redis"
)

var RDB = &Redis{}
var NotReadyErr = errors.New("RedisNotReady")
var DefaultAddr = "127.0.0.1:6379"
var DefaultClusterAddrs = "127.0.0.1:7001;127.0.0.1:7002;127.0.0.1:7003;127.0.0.1:7004;127.0.0.1:7005;127.0.0.1:7006"

var cacheConf = cache.ClusterConf{{Weight: 100, RedisConf: zredis.RedisConf{Host: DefaultAddr}}}

func init() {
	RDB = defaultRedis()
}

func Initialize(ctx contextx.Context, opfs ...func(rdb RdbConfig)) {

	for _, opf := range opfs {
		opf(RDB)
	}

	addrs := strings.Split(RDB.Addrs, ";")

	if len(addrs) > 1 {
		RDB.cluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Username: RDB.Username,
			Password: RDB.Password,
		})
		cacheConf = cache.ClusterConf{}
		for _, addr := range addrs {
			cacheConf = append(cacheConf, cache.NodeConf{
				Weight: 100,
				RedisConf: zredis.RedisConf{
					Host: addr,
					Type: "cluster",
					Pass: RDB.Password,
				},
			})
		}
	} else {
		RDB.cli = redis.NewClient(&redis.Options{
			Addr:     addrs[0],
			Username: RDB.Username,
			Password: RDB.Password,
		})
		cacheConf = cache.ClusterConf{
			{
				Weight: 100,
				RedisConf: zredis.RedisConf{
					Host: addrs[0],
					Pass: RDB.Password,
				},
			},
		}
	}

	RDB.ctx = ctx

	ctx.Go(func(ctx context.Context) error {
		tk := time.NewTicker(RDB.PingTime)
		defer tk.Stop()
		for {
			select {
			case <-ctx.Done():
				err := RDB.Close()
				if err != nil {
					logx.Infof("Redis Close err: %v", err)
				} else {
					logx.Infof("Redis Close successful!")
				}
				return ctx.Err()
			case <-tk.C:
				RDB.Ping(ctx)
			}
		}
	})
}

type RdbConfig interface {
	SetAddrs(addrs string)
	SetUsername(username string)
	SetPassword(password string)
	SetDB(db int)
	SetWaitTimeout(waitTime time.Duration)
	SetPingTime(pingTime time.Duration)
}

type Redis struct {
	Addrs       string // 连接地址
	Username    string // 用户名
	Password    string // 密码
	DB          int    // 数据库名
	WaitTimeout time.Duration
	PingTime    time.Duration
	ctx         contextx.Context
	cli         *redis.Client
	cluster     *redis.ClusterClient
}

func defaultRedis() *Redis {
	return &Redis{
		Addrs:       DefaultAddr,
		DB:          0,
		WaitTimeout: time.Second * 5,
		PingTime:    time.Second * 30,
	}
}

func CacheConf() cache.ClusterConf {
	return cacheConf
}

func (rd *Redis) SetAddrs(addrs string)                 { rd.Addrs = addrs }
func (rd *Redis) SetUsername(username string)           { rd.Username = username }
func (rd *Redis) SetPassword(password string)           { rd.Password = password }
func (rd *Redis) SetDB(db int)                          { rd.DB = db }
func (rd *Redis) SetWaitTimeout(waitTime time.Duration) { rd.WaitTimeout = waitTime }
func (rd *Redis) SetPingTime(pingTime time.Duration)    { rd.PingTime = pingTime }

func (rd *Redis) Close() error {
	var err error
	if rd.cli != nil {
		err = rd.cli.Close()
	}
	if rd.cluster != nil {
		err = rd.cluster.Close()
	}
	return err
}

func (rd *Redis) Ping(ctx context.Context) *redis.StatusCmd {
	var status *redis.StatusCmd
	if rd.cli != nil {
		status = rd.cli.Ping(ctx)
	}
	if rd.cluster != nil {
		status = rd.cluster.Ping(ctx)
	}
	return status
}

func (rd *Redis) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	if rd.cluster != nil {
		return rd.cluster.Do(ctx, args...)
	}
	if rd.cli != nil {
		return rd.cli.Do(ctx, args...)
	}
	var cmd = redis.NewCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) Get(ctx context.Context, key string) *redis.StringCmd {
	if rd.cluster != nil {
		return rd.cluster.Get(ctx, key)
	}
	if rd.cli != nil {
		return rd.cli.Get(ctx, key)
	}
	var cmd = redis.NewStringCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if rd.cluster != nil {
		return rd.cluster.Set(ctx, key, value, expiration)
	}
	if rd.cli != nil {
		return rd.cli.Set(ctx, key, value, expiration)
	}
	var cmd = redis.NewStatusCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) HMSet(ctx context.Context, key string, value ...interface{}) *redis.BoolCmd {
	if rd.cluster != nil {
		return rd.cluster.HMSet(ctx, key, value...)
	}
	if rd.cli != nil {
		return rd.cli.HMSet(ctx, key, value...)
	}
	var cmd = redis.NewBoolCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	if rd.cluster != nil {
		return rd.cluster.HGetAll(ctx, key)
	}
	if rd.cli != nil {
		return rd.cli.HGetAll(ctx, key)
	}
	var cmd = redis.NewStringStringMapCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	if rd.cluster != nil {
		return rd.cluster.SetNX(ctx, key, value, expiration)
	}
	if rd.cli != nil {
		return rd.cli.SetNX(ctx, key, value, expiration)
	}
	var cmd = redis.NewBoolCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if rd.cluster != nil {
		return rd.cluster.Del(ctx, keys...)
	}
	if rd.cli != nil {
		return rd.cli.Del(ctx, keys...)
	}
	var cmd = redis.NewIntCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) Incr(ctx context.Context, key string) *redis.IntCmd {
	if rd.cluster != nil {
		return rd.cluster.Incr(ctx, key)
	}
	if rd.cli != nil {
		return rd.cli.Del(ctx, key)
	}
	var cmd = redis.NewIntCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	if rd.cluster != nil {
		return rd.cluster.Publish(ctx, channel, message)
	}
	if rd.cli != nil {
		return rd.cli.Publish(ctx, channel, message)
	}
	var cmd = redis.NewIntCmd(ctx)
	cmd.SetErr(NotReadyErr)
	return cmd
}

func (rd *Redis) Subscribe(ctx context.Context, channels ...string) (*redis.PubSub, error) {
	if rd.cluster != nil {
		return rd.cluster.Subscribe(ctx, channels...), nil
	}
	if rd.cli != nil {
		return rd.cli.Subscribe(ctx, channels...), nil
	}
	return nil, NotReadyErr
}

// ======== get & set =========

func Get(key string) *redis.StringCmd {
	return RDB.Get(RDB.ctx, key)
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(RDB.ctx, key, value, expiration).Err()
}

func GetString(key string) (string, error) {
	return RDB.Get(RDB.ctx, key).Result()
}

func SetString(key, value string) error {
	return RDB.Set(RDB.ctx, key, value, 0).Err()
}

func GetInt(key string) (int, error) {
	return RDB.Get(RDB.ctx, key).Int()
}

func SetInt(key string, value int) error {
	return RDB.Set(RDB.ctx, key, value, 0).Err()
}

func GetInt64(key string) (int64, error) {
	return RDB.Get(RDB.ctx, key).Int64()
}

func GetUint64(key string) (uint64, error) {
	return RDB.Get(RDB.ctx, key).Uint64()
}

// =========== set ===============

func AddMembersToSet(key string, values ...interface{}) (interface{}, error) {
	var cmds = []interface{}{"SADD", key}
	cmds = append(cmds, values...)
	return RDB.Do(RDB.ctx, cmds...).Result()
}

func GetAllMembersFromSet(key string) (interface{}, error) {
	return RDB.Do(RDB.ctx, "SMEMBERS", key).Result()
}

func DelMembersInSet(key string, values ...interface{}) error {
	var cmds = []interface{}{"SREM", key}
	cmds = append(cmds, values...)
	return RDB.Do(RDB.ctx, cmds...).Err()
}

func MembersCountOfSet(key string) (int64, error) {
	return RDB.Do(RDB.ctx, "SCARD", key).Int64()
}

// =========== Hash ===============

type Encoder interface {
	Marshal() (map[string]string, error)
	Unmarshal(value map[string]string) error
}

func map2fields(v map[string]string) []interface{} {
	var values []interface{}
	for key, value := range v {
		values = append(values, key, value)
	}
	return values
}

func SetCache(key string, v Encoder) error {
	m, err := v.Marshal()
	if err != nil {
		return err
	}
	values := map2fields(m)
	return RDB.HMSet(RDB.ctx, key, values...).Err()
}

func GetCache(key string, v Encoder) error {
	result, err := RDB.HGetAll(RDB.ctx, key).Result()
	if err != nil {
		return err
	}
	return v.Unmarshal(result)
}

func LockExec(key string, f func(key string)) error {
	ctx, cancel := context.WithTimeout(RDB.ctx, RDB.WaitTimeout)
	defer cancel()
	lockKey := key + "-lock"
	var ok bool
	for !ok {
		select {
		case <-ctx.Done():
			errMsg := fmt.Sprintf("%s; key=%s has lock", ctx.Err().Error(), key)
			return errors.New(errMsg)
		default:
			ok = RDB.SetNX(RDB.ctx, lockKey, 1, RDB.WaitTimeout).Val()
		}
	}

	f(key)

	err := RDB.Del(RDB.ctx, lockKey).Err()
	for err != nil {
		select {
		case <-ctx.Done():
			errMsg := fmt.Sprintf("%s; key=%s Del err=%v", ctx.Err().Error(), key, err)
			return errors.New(errMsg)
		default:
			err = RDB.Del(RDB.ctx, lockKey).Err()
		}
	}
	return err
}

func Do(cmds ...interface{}) (interface{}, error) {
	return RDB.Do(RDB.ctx, cmds...).Result()
}

func Del(key string) error {
	return RDB.Del(RDB.ctx, key).Err()
}

func Incr(key string) (int64, error) {
	return RDB.Incr(RDB.ctx, key).Result()
}

func Publish(channel string, message interface{}) error {
	return RDB.Publish(RDB.ctx, channel, message).Err()
}

func Subscribe(f func(channel, msg string), channel ...string) error {
	sb, err := RDB.Subscribe(RDB.ctx, channel...)
	if err != nil {
		return err
	}
	message, err := sb.ReceiveMessage(RDB.ctx)
	if err != nil {
		return err
	}
	return recovers.CatchPanic(func() error {
		f(message.Channel, message.Payload)
		return nil
	})
}
