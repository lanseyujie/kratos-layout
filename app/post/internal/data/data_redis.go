package data

import (
	"context"
	"net"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"

	"sns/app/post/internal/biz"
	"sns/app/post/internal/conf"
)

// RedisData .
type RedisData struct {
	db *redis.Client
}

// NewRedisData .
func NewRedisData(c *conf.Data, logger log.Logger, hook redis.Hook) (*RedisData, func(), error) {
	data := &RedisData{}
	pCtx := context.TODO()

	// Redis
	{
		client := redis.NewClient(&redis.Options{
			Addr:         c.Redis.Addr,
			Password:     c.Redis.Password,
			DB:           int(c.Redis.Db),
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})

		if c.Redis.Debug && hook != nil {
			client.AddHook(hook)
		}

		ctx, cancel := context.WithTimeout(pCtx, c.Redis.DialTimeout.AsDuration())
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			return nil, nil, err
		}

		data.db = client

		cleanup := func() {
			log.NewHelper(logger).Warn("closing redis data resources")
			if err := data.db.Close(); err != nil {
				log.NewHelper(logger).Errorf("Redis.Disconnect Err: %v", err)
			}
		}

		return data, cleanup, nil
	}
}

// ------------

// NewRedisTransaction .
func NewRedisTransaction(d *RedisData) biz.RedisTransaction {
	return d
}

type redisTxKey struct{}

func (data *RedisData) DB(ctx context.Context) redis.Cmdable {
	if tx, ok := ctx.Value(redisTxKey{}).(redis.Cmdable); ok {
		return tx
	}

	return data.db
}

// ExecTx performs sequential single-threaded writes, and command execution does not immediately return a result.
// It should be noted that v8 is concurrency safe, but locks were removed in v9 making it non-concurrency safe.
// https://redis.uptrace.dev/zh/guide/go-redis-pipelines.html
func (data *RedisData) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := data.db.TxPipeline()
	defer tx.Discard()

	ctx = context.WithValue(ctx, redisTxKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	_, err := tx.Exec(ctx)

	return err
}

// ------------

var _ redis.Hook = (*RedisHook)(nil)

type RedisHook struct {
	logger log.Logger
}

func NewRedisHook(log log.Logger) redis.Hook {
	return &RedisHook{logger: log}
}

func (rh *RedisHook) DialHook(hook redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := hook(ctx, network, addr)

		return conn, err
	}
}

func (rh *RedisHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		log.NewHelper(rh.logger).Debugf("starting processing: <%s>\n", cmd)
		err := hook(ctx, cmd)
		log.NewHelper(rh.logger).Debugf("finished processing: <%s>\n", cmd)

		return err
	}
}

func (rh *RedisHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		log.NewHelper(rh.logger).Debugf("pipeline starting processing: %v\n", cmds)
		err := hook(ctx, cmds)
		log.NewHelper(rh.logger).Debugf("pipeline finished processing: %v\n", cmds)

		return err
	}
}
