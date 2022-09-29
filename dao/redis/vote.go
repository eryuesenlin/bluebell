package redis

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"math"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID int64) (err error) {
	client.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()
	return err
}

func VoteForPost(userID, postID string, value float64) (err error) {
	// 1.判断投票限制
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2.跟新帖子分数
	// 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 3. 记录用户为该贴子投票的数据
	if value == 0 {
		_, err = client.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Result()
	} else {
		_, err = client.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF+postID), &redis.Z{
			Score:  value,
			Member: userID,
		}).Result()
	}
	return err
}
