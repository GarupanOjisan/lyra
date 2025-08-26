package httpx

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error
type Middleware func(HandlerFunc) HandlerFunc

type Router struct {
    mux         *http.ServeMux
    middlewares []Middleware
}

func NewRouter() *Router { return &Router{mux: http.NewServeMux()} }

func (r *Router) Use(mw ...Middleware) { r.middlewares = append(r.middlewares, mw...) }

func (r *Router) Handle(method, path string, h HandlerFunc) {
    wrapped := h
    for i := len(r.middlewares) - 1; i >= 0; i-- {
        wrapped = r.middlewares[i](wrapped)
    }
    r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
        if req.Method != method { http.NotFound(w, req); return }
        if err := wrapped(w, req); err != nil {
            WriteProblem(w, req, FromError(err))
        }
    })
}

func (r *Router) GET(path string, h HandlerFunc)  { r.Handle(http.MethodGet, path, h) }
func (r *Router) POST(path string, h HandlerFunc) { r.Handle(http.MethodPost, path, h) }
func (r *Router) PUT(path string, h HandlerFunc)  { r.Handle(http.MethodPut, path, h) }
func (r *Router) DELETE(path string, h HandlerFunc){ r.Handle(http.MethodDelete, path, h) }

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) { r.mux.ServeHTTP(w, req) }

type Problem struct {
    Type   string `json:"type"`
    Title  string `json:"title"`
    Status int    `json:"status"`
    Detail string `json:"detail,omitempty"`
}

func (p Problem) Error() string { return p.Title }

func NewProblem(status int, title string) Problem {
    return Problem{Type: "about:blank", Title: title, Status: status}
}

func FromError(err error) Problem {
    if p, ok := err.(Problem); ok { return p }
    return NewProblem(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func WriteProblem(w http.ResponseWriter, _ *http.Request, p Problem) {
    w.Header().Set("Content-Type", "application/problem+json")
    w.WriteHeader(p.Status)
    _ = json.NewEncoder(w).Encode(p)
}

func OK(w http.ResponseWriter, v any) error {
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(v)
}

func Created(w http.ResponseWriter, location string) error {
    w.Header().Set("Location", location)
    w.WriteHeader(http.StatusCreated)
    return nil
}

func BindJSON(r *http.Request, dst any) error {
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(dst); err != nil {
        return NewProblem(http.StatusBadRequest, "invalid request body")
    }
    return nil
}

func Recover() Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) (err error) {
            defer func() {
                if rec := recover(); rec != nil {
                    log.Printf("panic: %v", rec)
                    err = NewProblem(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
                }
            }()
            return next(w, r)
        }
    }
}

func Logger() Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) error {
            start := time.Now()
            err := next(w, r)
            log.Printf("%s %s %s err=%v", r.Method, r.URL.Path, time.Since(start), err)
            return err
        }
    }
}

type Server struct {
    srv *http.Server
}

func NewServer(addr string, r *Router) *Server {
    return &Server{srv: &http.Server{Addr: addr, Handler: r}}
}

func (s *Server) Start(ctx context.Context) error { go func() { _ = s.srv.ListenAndServe() }(); return nil }
func (s *Server) Stop(ctx context.Context) error  { return s.srv.Shutdown(ctx) }
