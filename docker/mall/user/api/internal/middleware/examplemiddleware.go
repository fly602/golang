package middleware

import (
	"encoding/json"
	"fmt"
	"go-community/docker/mall/user/api/internal/config"
	jwtx "go-community/docker/mall/user/common/jwt"
	"net/http"
	"time"
)

type ExampleMiddleware struct {
	Config config.Config
}

func NewExampleMiddleware(c config.Config) *ExampleMiddleware {
	return &ExampleMiddleware{
		Config: c,
	}
}

func (m *ExampleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		exp := m.Config.Auth.AccessExpire
		now := time.Now().Unix()
		uid, _ := r.Context().Value("uid").(json.Number).Int64()
		fmt.Printf("exp=%+v,now=%+v,uid=%+v", exp, now, uid)
		accessToken, err := jwtx.GetToken(m.Config.Auth.AccessSecret, now, exp, uid)
		if err != nil {
			next(w, r)
			return
		}
		// Passthrough to next handler if need
		w.Header().Add("Authorization", accessToken)
		next(w, r)
	}
}
