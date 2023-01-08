package jwtx

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type JwtheaderMiddleware struct {
	Auth JwtAuth
	logx.Logger
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
		// 初始化日志
		m.Logger = logx.WithContext(r.Context())

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
			m.Infof("jwt refresh accessToken =%v", accessToken)
		} else {
			authorization := r.Header["Authorization"]
			if authorization != nil {
				accessToken = authorization[0]
			}
			m.Infof("jwt use old accessToken =%v", accessToken)
		}
		// Passthrough to next handler if need
		w.Header().Add("Authorization", accessToken)
		// Passthrough to next handler if need
		next(w, r)
	}
}
