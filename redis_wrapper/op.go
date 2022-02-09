package redis_wrapper

import (
    "context"
    "github.com/bilibili/kratos/pkg/cache/redis"
    "github.com/bilibili/kratos/pkg/log"
    "strconv"
    "time"
)

const (
    DEFAULT_LOCK_WAIT_SECOND = 3
    DEFAULT_LOCK_RETRY_TIME  = 3
)

type ZSetElement struct {
    Member string
    Score  float64
}

type BatchCMD struct {
    CMD    string
    Params []interface{}
}

type Scanner func(reply interface{}, err error)

//common cmd======================================

func Expire(ctx context.Context, key string, second int64) (err error) {
    _, err = toInt64Reply(redisExecutor.Do(ctx, "EXPIRE", key, second))
    return
}

func Del(ctx context.Context, key string) (delNum int64, err error) {
    delNum, err = toInt64Reply(redisExecutor.Do(ctx, "DEL", key))
    return
}

func Keys(ctx context.Context, pattern string) (keys []string, err error) {
    return toStringListReply(redisExecutor.Do(ctx, "KEYS", pattern))
}

func Get(ctx context.Context, key string) (value string, err error) {
    value, err = toStringReply(redisExecutor.Do(ctx, "GET", key))
    return
}

func Set(ctx context.Context, key string, value string) (err error) {
    _, err = toStringReply(redisExecutor.Do(ctx, "SET", key, value))
    return
}

func SetNX(ctx context.Context, key string, value string) (isSet bool, err error) {
    var ret int64
    ret, err = toInt64Reply(redisExecutor.Do(ctx, "SET", key, value, "NX"))
    if ret == 1 {
        isSet = true
    }
    return
}

//SimpleLock 简单的redis分布式锁 当redis非单一节点时并不可靠
func SimpleLock(ctx context.Context, key string, second int64) (gotLock bool) {
    _, err := redis.String(redisExecutor.Do(ctx, "SET", key, "lock", "NX", "EX", second))
    if err != nil {
        log.Error("SimpleLock key:%s second:%d error:%s", key, second, err.Error())
        return false
    }
    return true
}

func WaitLock(ctx context.Context, key string, second int64) (gotLock bool) {
    beginTime := time.Now()
    tryTime := 0
    for {
        diffTime := time.Now().Sub(beginTime)
        if diffTime > time.Second*DEFAULT_LOCK_WAIT_SECOND {
            break
        }
        if SimpleLock(ctx, key, second) {
            return true
        }
        tryTime++
        if tryTime >= DEFAULT_LOCK_RETRY_TIME {
            break
        }
        time.Sleep(time.Millisecond * 200)
    }

    return false
}

//简单Del一下 理想情况应该用脚本对比Lock是否是自己加的再删除
func UnLock(ctx context.Context, key string) () {
    _, _ = Del(ctx, key)
}

func IncrBy(ctx context.Context, key string, delta int64) (curValue int64, err error) {
    curValue, err = toInt64Reply(redisExecutor.Do(ctx, "INCRBY", key, delta))
    return
}

//hash cmd======================================

func HGet(ctx context.Context, key string, field string) (value string, err error) {
    value, err = toStringReply(redisExecutor.Do(ctx, "HGET", key, field))
    return
}

func HGetAll(ctx context.Context, key string) (values map[string]string, err error) {
    values = map[string]string{}
    var ret []string
    if ret, err = toStringListReply(redisExecutor.Do(ctx, "HGETALL", key)); err != nil {
        return
    }
    for i := 0; i < len(ret)-1; i += 2 {
        values[ret[i]] = ret[i+1]
    }
    return
}

func HKeys(ctx context.Context, key string) (hkeys []string, err error) {
    hkeys, err = toStringListReply(redisExecutor.Do(ctx, "HKEYS", key))
    return
}

func HSet(ctx context.Context, key string, field string, value string) (err error) {
    _, err = toInt64Reply(redisExecutor.Do(ctx, "HSET", key, field, value))
    return
}

func HSetNX(ctx context.Context, key string, field string, value string) (isSet bool, err error) {
    var ret int64
    if ret, err = toInt64Reply(redisExecutor.Do(ctx, "HSETNX", key, field, value)); err != nil {
        return
    }
    if ret == 1 {
        isSet = true
    }
    return
}

func HMGet(ctx context.Context, key string, fields []string) (values map[string]string, err error) {
    values = map[string]string{}
    if len(fields) == 0 {
        return
    }
    params := []interface{}{key}
    for _, v := range fields {
        params = append(params, v)
    }
    var ret []string
    if ret, err = toStringListReply(redisExecutor.Do(ctx, "HMGET", params...)); err != nil {
        return
    }
    for i, v := range fields {
        values[v] = ret[i]
    }
    return
}

func HIncr(ctx context.Context, key string, field string, incr int64) (cur int64, err error) {
    cur, err = toInt64Reply(redisExecutor.Do(ctx, "HINCRBY", key, field, incr))
    return
}

func HDel(ctx context.Context, key string, field string) (delNum int64, err error) {
    delNum, err = toInt64Reply(redisExecutor.Do(ctx, "HDEL", key, field))
    return
}

func HLen(ctx context.Context, key string) (length int64, err error) {
    length, err = toInt64Reply(redisExecutor.Do(ctx, "HLEN", key))
    return
}

//zset cmd======================================

func ZAdd(ctx context.Context, key string, score float64, member string) (changeNum int64, err error) {
    changeNum, err = toInt64Reply(redisExecutor.Do(ctx, "ZADD", key, score, member))
    return
}

func ZIncr(ctx context.Context, key string, incr float64, member string) (cur float64, err error) {
    cur, err = toFloat64Reply(redisExecutor.Do(ctx, "ZINCRBY", key, incr, member))
    return
}

func ZRange(ctx context.Context, key string, start int64, stop int64, withScores bool) (elements []*ZSetElement, err error) {
    return zRange(ctx, key, start, stop, false, withScores)
}

func ZRevRange(ctx context.Context, key string, start int64, stop int64, withScores bool) (elements []*ZSetElement, err error) {
    return zRange(ctx, key, start, stop, true, withScores)
}

func ZRangeByScore(ctx context.Context, key string, min int64, max int64, withScores bool) (elements []*ZSetElement, err error) {
    return zRangeByScore(ctx, key, min, max, false, withScores)
}

func ZRevRangeByScore(ctx context.Context, key string, min int64, max int64, withScores bool) (elements []*ZSetElement, err error) {
    return zRangeByScore(ctx, key, min, max, true, withScores)
}


func ZScore(ctx context.Context, key string, member string) (score float64, err error) {
    score, err = toFloat64Reply(redisExecutor.Do(ctx, "ZSCORE", key, member))
    return
}

func ZRank(ctx context.Context, key string, member string) (rank int64, err error) {
    return zRank(ctx, key, member, false)
}

func ZRevRank(ctx context.Context, key string, member string) (rank int64, err error) {
    return zRank(ctx, key, member, true)
}

//list cmd======================================

func LPush(ctx context.Context, key string, values []string) (curLen int64, err error) {
    params := []interface{}{key}
    for _, v := range values {
        params = append(params, v)
    }
    curLen, err = toInt64Reply(redisExecutor.Do(ctx, "LPUSH", params...))
    return
}

func LPop(ctx context.Context, key string) (value string, err error) {
    return toStringReply(redisExecutor.Do(ctx, "LPOP", key))
}

func RPush(ctx context.Context, key string, values []string) (curLen int64, err error) {
    params := []interface{}{key}
    for _, v := range values {
        params = append(params, v)
    }
    curLen, err = toInt64Reply(redisExecutor.Do(ctx, "RPUSH", params...))
    return
}

func RPop(ctx context.Context, key string) (value string, err error) {
    return toStringReply(redisExecutor.Do(ctx, "RPOP", key))
}

func LTrim(ctx context.Context, key string, start int64, stop int64) (err error) {
    _, err = redisExecutor.Do(ctx, "LTRIM", key, start, stop)
    return
}

func LRange(ctx context.Context, key string, start int64, stop int64) (ret []string, err error) {
    ret, err = toStringListReply(redisExecutor.Do(ctx, "LRANGE", key, start, stop))
    return
}

func LSet(ctx context.Context, key string, index int64, value string) (err error) {
    _, err = toStringReply(redisExecutor.Do(ctx, "LSET", key, index, value))
    return
}

func LLen(ctx context.Context, key string) (length int64, err error) {
    length, err = toInt64Reply(redisExecutor.Do(ctx, "LLEN", key))
    return
}

func RPopLPush(ctx context.Context, keyFrom string, keyTo string) (movedElement string, err error) {
    movedElement, err = toStringReply(redisExecutor.Do(ctx, "RPOPLPUSH", keyFrom, keyTo))
    return
}

// LRem count > 0 从头到尾移除count个
// count < 0 从尾到头移除count个
// count = 0 从尾到头移除所有
func LRem(ctx context.Context, key string, count int64, value string) (removedNum int64, err error) {
    removedNum, err = toInt64Reply(redisExecutor.Do(ctx, "LREM", key, count, value))
    return
}

//set cmd======================================

func SAdd(ctx context.Context, key string, value string) (err error) {
    _, err = toInt64Reply(redisExecutor.Do(ctx, "SADD", key, value))
    return
}

//private function ======================================

func zRange(ctx context.Context, key string, start int64, stop int64, isRev bool, withScores bool) (elements []*ZSetElement, err error) {
    cmd := "ZRANGE"
    if isRev {
        cmd = "ZREVRANGE"
    }
    params := []interface{}{key, start, stop}
    if withScores {
        params = append(params, "WITHSCORES")
    }
    var values []string
    if values, err = toStringListReply(redisExecutor.Do(ctx, cmd, params...)); err != nil {
        return
    }
    if withScores {
        for i := 0; i < len(values)-1; i += 2 {
            var score float64
            if score, err = toFloat64(values[i+1]); err != nil {
                return
            }
            elements = append(elements, &ZSetElement{
                Member: values[i],
                Score:  score,
            })
        }
    } else {
        for _, v := range values {
            elements = append(elements, &ZSetElement{
                Member: v,
                Score:  0,
            })
        }
    }
    return
}

func zRangeByScore(ctx context.Context, key string, min int64, max int64, isRev bool, withScores bool)(elements []*ZSetElement, err error) {
    cmd := "ZRANGEBYSCORE"
    params := []interface{}{key, min, max}
    if min == -1 {
        params = []interface{}{key, "-inf", max}
    }
    if max == -1 {
        params = []interface{}{key, min, "+inf"}
    }
    if isRev {
        cmd = "ZREVRANGEBYSCORE"
        params = []interface{}{key, max, min}
        if min == -1 {
            params = []interface{}{key, max, "-inf"}
        }
        if max == -1 {
            params = []interface{}{key, "+inf", min}
        }
    }
    if withScores {
        params = append(params, "WITHSCORES")
    }
    var values []string
    if values, err = toStringListReply(redisExecutor.Do(ctx, cmd, params...)); err != nil {
        return
    }
    if withScores {
        for i := 0; i < len(values)-1; i += 2 {
            var score float64
            if score, err = toFloat64(values[i+1]); err != nil {
                return
            }
            elements = append(elements, &ZSetElement{
                Member: values[i],
                Score:  score,
            })
        }
    } else {
        for _, v := range values {
            elements = append(elements, &ZSetElement{
                Member: v,
                Score:  0,
            })
        }
    }
    return
}

func zRank(ctx context.Context, key string, member string, isRev bool) (rank int64, err error) {
    cmd := "ZRANK"
    if isRev {
        cmd = "ZREVRANK"
    }
    var reply interface{}
    if reply, err = redisExecutor.Do(ctx, cmd, key, member); err != nil {
        return
    }
    rank, err = redis.Int64(reply, err)
    if err == redis.ErrNil {
        err = nil
        rank = -1
    }
    return
}

func toInt64Reply(reply interface{}, err error) (ret int64, newErr error) {
    ret, newErr = redis.Int64(reply, err)
    if newErr == redis.ErrNil {
        newErr = nil
    }
    return
}

func toFloat64Reply(reply interface{}, err error) (ret float64, newErr error) {
    ret, newErr = redis.Float64(reply, err)
    if newErr == redis.ErrNil {
        newErr = nil
    }
    return
}

func toStringListReply(reply interface{}, err error) (ret []string, newErr error) {
    ret, newErr = redis.Strings(reply, err)
    if newErr == redis.ErrNil {
        newErr = nil
    }
    return
}

func toStringReply(reply interface{}, err error) (ret string, newErr error) {
    ret, newErr = redis.String(reply, err)
    if newErr == redis.ErrNil {
        newErr = nil
    }
    return
}

func toInt64(str string) (ret int64, err error) {
    ret, err = strconv.ParseInt(str, 10, 64)
    return
}

func toFloat64(str string) (ret float64, err error) {
    ret, err = strconv.ParseFloat(str, 64)
    return
}
