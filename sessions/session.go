package sessions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type SessionStore struct {
	Sessiontoken string `json:"sessiontoken"`
	Username     string `json:"username"`
}

type MakeSession interface {
	Set(string, *gin.Context) SessionStore
	Get(string, *gin.Context) interface{}
}

func SessionStart() MakeSession {
	return &SessionStore{}
}
func (store *SessionStore) Set(username string, ctx *gin.Context) SessionStore {
	sessionToken := xid.New().String()
	session := sessions.Default(ctx)
	session.Set("username", username)
	session.Set("token", sessionToken)
	session.Save()
	return SessionStore{
		Sessiontoken: sessionToken,
		Username:     username,
	}
}

func (store *SessionStore) Get(sessiontoken string, ctx *gin.Context) interface{} {
	session := sessions.Default(ctx)
	sessionToken := session.Get(sessiontoken)
	return sessionToken
}
