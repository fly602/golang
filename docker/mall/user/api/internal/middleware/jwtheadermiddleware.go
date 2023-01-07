package middleware

import (
	"encoding/json"
	"fmt"
	"go-community/docker/mall/user/api/internal/config"
	jwtx "go-community/docker/mall/user/common/jwt"
	"net/http"
	"time"
)

type JwtheaderMiddleware struct {
	Config config.Config
}

func NewJwtheaderMiddleware(c config.Config) *JwtheaderMiddleware {
	return &JwtheaderMiddleware{
		Config: c,
	}
}

func (m *JwtheaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
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
		// Passthrough to next handler if need
		next(w, r)
	}
}
