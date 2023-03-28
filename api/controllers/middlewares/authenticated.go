package middlewares

import (
	"context"
	"net/http"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/jwt"
	"github.com/julienschmidt/httprouter"
)

func Authenticated(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		result := jwt.VerifyToken(authorization, conf.API_TOKEN_SIGNING_KEY)
		if !result.Valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, CONTEXT_KEY_USER_ID, result.UserId)
		h(w, r.WithContext(ctx), ps)
	}
}
