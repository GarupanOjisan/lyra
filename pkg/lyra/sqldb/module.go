package sqldb

import (
    "context"
    "database/sql"
    "errors"
    "os"
    "time"

    "github.com/garupanojisan/lyra/pkg/lyra/di"
    "github.com/garupanojisan/lyra/pkg/lyra/tx"
)

type Module struct {
    Driver string
    DSN    string
    DB     *sql.DB
}

func (m *Module) Name() string { return "sqldb" }

func (m *Module) Boot(c *di.Container) error {
    if m.Driver == "" { m.Driver = os.Getenv("LYRA_DB_DRIVER") }
    if m.DSN == "" { m.DSN = os.Getenv("LYRA_DB_DSN") }
    if m.Driver == "" || m.DSN == "" {
        // Not configured; do not register DB-backed Tx manager.
        return nil
    }
    db, err := sql.Open(m.Driver, m.DSN)
    if err != nil { return err }
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(time.Hour)
    if err := db.Ping(); err != nil { return err }
    m.DB = db
    di.Provide(c, db)
    di.Provide[tx.Manager](c, &sqlManager{db: db})
    return nil
}

func (m *Module) Start(ctx context.Context) error { return nil }
func (m *Module) Stop(ctx context.Context) error  { if m.DB != nil { return m.DB.Close() }; return nil }

func ModuleFromEnv() *Module { return &Module{} }

type sqlManager struct{ db *sql.DB }

func (m *sqlManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
    if m.db == nil { return errors.New("sqldb: not configured") }
    tx, err := m.db.BeginTx(ctx, nil)
    if err != nil { return err }
    ctx = WithTx(ctx, tx)
    if err := fn(ctx); err != nil {
        _ = tx.Rollback()
        return err
    }
    return tx.Commit()
}

