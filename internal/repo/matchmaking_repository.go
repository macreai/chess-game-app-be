package repo

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MatchMakingRepository interface {
	AddUserToQueue(user uint64) error
	MatchingUsers() (int64, error)
	MatchedUsers() ([]string, error)
}

type MatchMakingRepositoryImpl struct {
	Log     *logrus.Logger
	RedisDB *redis.Client
	Viper   *viper.Viper
}

func NewMatchMakingRepositoryImpl(log *logrus.Logger, redisDB *redis.Client, viper *viper.Viper) *MatchMakingRepositoryImpl {
	repo := &MatchMakingRepositoryImpl{
		Log:     log,
		RedisDB: redisDB,
		Viper:   viper,
	}

	err := repo.InitQueue()
	if err != nil {
		log.Errorf("failed to init redis data key: %+v", err)
	}

	return repo
}

func (m *MatchMakingRepositoryImpl) AddUserToQueue(userId uint64) error {
	scoreNumber, _ := m.RedisDB.ZRevRangeWithScores(context.Background(), m.getMatchRedisKey(), 0, 0).Result()

	err := m.RedisDB.ZAdd(context.Background(), m.getMatchRedisKey(), redis.Z{
		Score:  scoreNumber[0].Score + 1,
		Member: userId,
	}).Err()

	return err
}

func (m *MatchMakingRepositoryImpl) MatchingUsers() (int64, error) {
	return m.RedisDB.ZCard(context.Background(), m.getMatchRedisKey()).Result()
}

func (m *MatchMakingRepositoryImpl) MatchedUsers() ([]string, error) {
	users, err := m.RedisDB.ZRange(context.Background(), m.getMatchRedisKey(), 1, 2).Result()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (m *MatchMakingRepositoryImpl) PopUsersMatched(users []string) error {
	var interfaceUsers []interface{}
	for _, user := range users {
		interfaceUsers = append(interfaceUsers, user)
	}

	return m.RedisDB.ZRem(context.Background(), m.getMatchRedisKey(), interfaceUsers...).Err()
}

func (m *MatchMakingRepositoryImpl) InitQueue() error {
	_, err := m.RedisDB.ZAdd(context.Background(), m.getMatchRedisKey(), redis.Z{
		Score:  0,
		Member: 0,
	}).Result()

	return err
}

func (m *MatchMakingRepositoryImpl) getMatchRedisKey() string {
	return m.Viper.GetString("REDIS_MATCH_Q_KEY")
}
