package stncsession
/*
import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetSession session set eder
func SetSession(key string, value string, c *gin.Context) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

//GetSession session get eder *--- çalışmıyor
func GetSession(key string, c *gin.Context) string {
	session := sessions.Default(c)
	if session != nil {
		return session.Get(key).(string)
	} else {
		return ""
	}
}

*/