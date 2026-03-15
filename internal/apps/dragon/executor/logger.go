package executor

import (
	"context"
	"time"
	"valyria-backend/internal/core/database"

	"github.com/redis/go-redis/v9"
)

type logWriter struct {
	rdb *redis.Client
	key string
	ctx context.Context
}

func newLogWriter(TaskID string) *logWriter {
	return &logWriter{
		rdb: database.Redis,
		ctx: context.Background(),
		key: "release:" + TaskID + ":logs",
	}
}

func (l *logWriter) AppendLog(line string) error {
	if err := l.rdb.RPush(l.ctx, l.key, line).Err(); err != nil {
		return err
	}
	l.rdb.Expire(l.ctx, l.key, 24*time.Hour)
	return nil
}
