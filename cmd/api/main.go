package main

import (
    "context"
    "log"
    "net/http"

    "github.com/garupanojisan/lyra/internal/usercontext"
    lyraapp "github.com/garupanojisan/lyra/pkg/lyra/app"
    "github.com/garupanojisan/lyra/pkg/lyra/di"
    "github.com/garupanojisan/lyra/pkg/lyra/httpx"
    "github.com/garupanojisan/lyra/pkg/lyra/sqldb"
    "github.com/garupanojisan/lyra/pkg/lyra/tx"
)

func main() {
    app := lyraapp.New()
    app.Use(httpx.ModuleDefault())
    // Use in-memory Tx by default, allow sqldb to override if configured via env.
    app.Use(tx.ModuleDefault())
    app.Use(sqldb.ModuleFromEnv())
    app.Use(usercontext.ModuleDefault())

    // health endpoint
    r := di.MustGet[*httpx.Router](app.Container())
    r.GET("/healthz", func(w http.ResponseWriter, r *http.Request) error {
        w.WriteHeader(http.StatusNoContent)
        return nil
    })

    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
