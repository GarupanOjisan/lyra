package main

import (
    "context"
    "log"

    lyraapp "github.com/garupanojisan/lyra/pkg/lyra/app"
    "github.com/garupanojisan/lyra/pkg/lyra/outbox"
)

func main() {
    app := lyraapp.New()
    app.Use(outbox.ModuleDefault())

    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}

