package lock

import (
	"context"
	"fmt"
	"miaosha-system/global"
	"strconv"
	"time"
)

// generateUniqueLockValue 根据时间戳、用户 ID 和商品 ID 生成唯一锁值
func generateUniqueLockValue(userID, goodID uint) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("%s-%d-%d", timestamp, goodID, userID)
}

// AcquireLock 尝试获取锁
func AcquireLock(ctx context.Context, key string, expiration time.Duration, userID, goodID uint) (string, bool, error) {
	value := generateUniqueLockValue(userID, goodID)
	maxRetries := 5 // 设置最大重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("用户%d第%d次获取锁\n", userID, i+1)
		set, err := global.Redis.SetNX(ctx, fmt.Sprintf("lock:%d", goodID), value, expiration).Result()
		if err == nil {
			if set {
				// 启动续期协程
				go RenewLock(ctx, key, value, expiration)
			} else {
				//fmt.Printf("用户%d第%d次获取锁失败，错误信息: %v\n", userID, i+1, err)
				// 等待一段时间后重试
				time.Sleep(time.Second * 1)
				continue
			}
			return value, set, nil
		}
	}
	return "", false, fmt.Errorf("获取锁失败，重试 %d 次后仍未成功", maxRetries)
}

// ReleaseLock 释放锁
func ReleaseLock(ctx context.Context, key, value string) (bool, error) {
	// 使用 Lua 脚本确保释放锁的原子性
	luaScript := `
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0
    end
    `
	result, err := global.Redis.Eval(ctx, luaScript, []string{key}, value).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// RenewLock 续期锁
var renewStopChan = make(chan struct{}) // 用于停止续期协程的通道

func RenewLock(ctx context.Context, key, value string, expiration time.Duration) {
	ticker := time.NewTicker(expiration / 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 使用 Lua 脚本检查并续期锁
			luaScript := `
            if redis.call("GET", KEYS[1]) == ARGV[1] then
                return redis.call("PEXPIRE", KEYS[1], ARGV[2])
            else
                return 0
            end
            `
			result, err := global.Redis.Eval(ctx, luaScript, []string{key}, value, int64(expiration.Milliseconds())).Result()
			if err != nil || result.(int64) == 0 {
				// 续期失败，停止续期协程
				return
			}
			fmt.Println("开始续期")
		case <-ctx.Done():
			return
		case <-renewStopChan:
			return
		}
	}
}

// StopRenew 停止续期协程
func StopRenew() {
	select {
	case <-renewStopChan:
		// 通道已关闭，不做任何操作
	default:
		close(renewStopChan)
	}
}
