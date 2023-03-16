package data

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/go-sql-driver/mysql"

	"sns/app/post/internal/biz"
	"sns/app/post/internal/conf"
	"sns/app/post/internal/data/ent"
	"sns/app/post/internal/data/ent/migrate"
)

// EntData .
type EntData struct {
	db *ent.Client
}

// NewEntData .
func NewEntData(cfg *conf.Data, logger log.Logger) (*EntData, func(), error) {
	data := &EntData{}
	l := log.NewHelper(log.With(logger, "module", "data/ent"))

	// init client
	{
		// Ent
		drv, err := sql.Open(cfg.Ent.Driver, cfg.Ent.Dsn)
		if err != nil {
			l.Errorf("failed opening: %v", err)

			return nil, nil, err
		}

		// DB Config
		drv.DB().SetMaxOpenConns(int(cfg.Ent.PoolSize))
		drv.DB().SetMaxIdleConns(int(cfg.Ent.IdleSize))
		drv.DB().SetConnMaxIdleTime(cfg.Ent.IdleTime.AsDuration())
		drv.DB().SetConnMaxLifetime(cfg.Ent.LifeTime.AsDuration())

		opts := []ent.Option{ent.Driver(drv), ent.Log(func(a ...any) { l.Debug(a...) })}
		if cfg.Ent.Debug {
			opts = append(opts, ent.Debug())
		}
		data.db = ent.NewClient(opts...)
	}

	// automatic migration
	{
		if cfg.Ent.Migration == nil {
			cfg.Ent.Migration = &conf.Data_Ent_Migration{}
		}

		if err := data.db.Schema.Create(
			context.Background(),
			migrate.WithDropIndex(cfg.Ent.Migration.DropIndex),
			migrate.WithDropColumn(cfg.Ent.Migration.DropColumn),
			migrate.WithForeignKeys(cfg.Ent.Migration.ForeignKeys),
		); err != nil {
			l.Errorf("failed creating schema resources: %v", err)

			return nil, nil, err
		}
	}

	cleanup := func() {
		l.Warn("closing ent data resources")
		if err := data.db.Close(); err != nil {
			l.Error(err)
		}
	}

	return data, cleanup, nil
}

// ------------

// NewEntTransaction .
func NewEntTransaction(d *EntData) biz.EntTransaction {
	return d
}

type entTxKey struct{}

func (d *EntData) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx = context.WithValue(ctx, entTxKey{}, tx)
	if err = fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

// ------------

func (d *EntData) PostClient(ctx context.Context) *ent.PostClient {
	if tx, ok := ctx.Value(entTxKey{}).(*ent.Tx); ok {
		return tx.Post
	}

	return d.db.Post
}
