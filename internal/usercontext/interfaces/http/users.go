package httpi

import (
    "net/http"

    "github.com/garupanojisan/lyra/internal/usercontext/application/usecase"
    "github.com/garupanojisan/lyra/pkg/lyra/httpx"
)

type UsersAPI struct{ Create *usecase.CreateUserHandler }

func (a *UsersAPI) Register(r *httpx.Router) {
    r.POST("/v1/users", a.postUsers)
}

func (a *UsersAPI) postUsers(w http.ResponseWriter, r *http.Request) error {
    var in struct{ Email string `json:"email"` }
    if err := httpx.BindJSON(r, &in); err != nil { return err }
    if err := a.Create.Handle(r.Context(), usecase.CreateUser{Email: in.Email}); err != nil {
        return httpx.NewProblem(http.StatusConflict, "email taken")
    }
    return httpx.Created(w, "/v1/users/…")
}

