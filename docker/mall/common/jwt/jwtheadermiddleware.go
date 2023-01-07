package jwtx

import (
	"encoding/json"
	"fmt"
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

func (m *JwtheaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		exp := m.Auth.AccessExpire
		now := time.Now().Unix()
		uid, _ := r.Context().Value("uid").(json.Number).Int64()
		fmt.Printf("exp=%+v,now=%+v,uid=%+v", exp, now, uid)
		accessToken, err := GetToken(m.Auth.AccessSecret, now, exp, uid)
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
