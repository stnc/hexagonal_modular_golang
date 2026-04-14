package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type Flash struct {
	Success string
	Error   string
		Type    string
	Message string
}

const (
	flashTypeKey    = "flash:type"
	flashMessageKey = "flash:message"
)

func ConsumeFlash(store *session.Store, c fiber.Ctx) Flash {
	sess, err := store.Get(c)
	if err != nil {
		return Flash{}
	}
	defer sess.Release()



	
	flash := Flash{}
	if v, ok := sess.Get("flash_success").(string); ok {
		flash.Success = v
		sess.Delete("flash_success")
	}
	if v, ok := sess.Get("flash_error").(string); ok {
		flash.Error = v
		sess.Delete("flash_error")
	}
	_ = sess.Save()
	return flash
}

func SetFlash(store *session.Store, c fiber.Ctx, success, failure string) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}
	defer sess.Release()
	if success != "" {
		sess.Set("flash_success", success)
	}
	if failure != "" {
		sess.Set("flash_error", failure)
	}
	return sess.Save()
}

//  TODO: lang var 
// bu yeni eklendi 
func PopFlash(store *session.Store, c fiber.Ctx) Flash {
	sess, err := store.Get(c)
	if err != nil {
		return Flash{}
	}
	defer sess.Release()
	ft, _ := sess.Get(flashTypeKey).(string)
	msg, _ := sess.Get(flashMessageKey).(string)
	if msg != "" {
		sess.Delete(flashTypeKey)
		sess.Delete(flashMessageKey)
		_ = sess.Save()
	}
	return Flash{Type: ft, Message: msg}
}
