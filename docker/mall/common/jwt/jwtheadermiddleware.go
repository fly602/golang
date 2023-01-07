package jwtx

import (
	"encoding/json"
	"net/http"
	"time"
)

type JwtheaderMiddleware struct {
	Auth JwtAuth
}

func NewJwtheaderMiddleware(j JwtAuth) *JwtheaderMiddleware {
	return &JwtheaderMiddleware{
		Auth: j,
	}
}

// 此中间件需要放在jwt之后
func (m *JwtheaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		// 获取jwt解析后的数据
		uid, _ := r.Context().Value("uid").(json.Number).Int64()
		old, _ := r.Context().Value("uid").(json.Number).Int64()

		exp := m.Auth.AccessExpire
		now := time.Now().Unix()

		var accessToken string
		var err error
		if now-old < exp/2 {
			accessToken, err = GetToken(m.Auth.AccessSecret, now, exp, uid)
			if err != nil {
				next(w, r)
				return
			}
		} else {
			authorization := r.Header["Authorization"]
			if authorization != nil {
				accessToken = authorization[0]
			}
		}
		// Passthrough to next handler if need
		w.Header().Add("Authorization", accessToken)

		// Passthrough to next handler if need
		next(w, r)
	}
}
