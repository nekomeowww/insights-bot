package tgbot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	"github.com/nekomeowww/insights-bot/pkg/types/redis"
	"github.com/redis/rueidis"
)

func (b *Bot) couldCountRateLimitFor(key string, rate int64, perDuration time.Duration) (int64, time.Duration, bool, error) {
	if perDuration <= 0 {
		return 0, 0, true, nil
	}

	getCmd := b.rueidisClient.
		B().
		Get().
		Key(key).
		Build()

	ttlCmd := b.rueidisClient.
		B().
		Ttl().
		Key(key).
		Build()

	countedRate, err := b.rueidisClient.Do(context.TODO(), getCmd).AsInt64()
	if err != nil && !rueidis.IsRedisNil(err) {
		return 0, 0, false, err
	}

	ttl, err := b.rueidisClient.Do(context.TODO(), ttlCmd).AsInt64()
	if err != nil && !rueidis.IsRedisNil(err) {
		return countedRate, 0, false, err
	}

	ttlSeconds := time.Duration(ttl) * time.Second
	if countedRate >= rate {
		return countedRate, ttlSeconds, false, nil
	}

	countedRate++
	setCmd := b.rueidisClient.
		B().
		Set().
		Key(key).
		Value(fmt.Sprintf("%d", countedRate)).
		ExSeconds(int64(perDuration / time.Second)).
		Build()

	err = b.rueidisClient.Do(context.TODO(), setCmd).Error()
	if err != nil {
		return countedRate, ttlSeconds, false, err
	}

	return countedRate, ttlSeconds, true, nil
}

func (b *Bot) RateLimitForCommand(chatID int64, command string, rate int64, perDuration time.Duration) (int64, time.Duration, bool, error) {
	return b.couldCountRateLimitFor(redis.CommandRateLimitLock2.Format(command, bot.FromPlatformTelegram.String(), strconv.FormatInt(chatID, 10)), rate, perDuration)
}
