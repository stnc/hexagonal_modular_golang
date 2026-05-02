package middleware

import (
	"fmt"

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
	sess := session.FromContext(c)

	flash := Flash{}
	if v, ok := sess.Get("FlashSuccess").(string); ok {
		fmt.Println("flash succces")
		fmt.Println(v)
		flash.Success = v
		sess.Delete("FlashSuccess")
	}
	if v, ok := sess.Get("FlashError").(string); ok {
		fmt.Println("flash eroorr")
		fmt.Println(v)
		flash.Error = v
		sess.Delete("FlashError")
	}

	return flash
}

func SetFlash(store *session.Store, c fiber.Ctx, success, failure string) {
	sess := session.FromContext(c)

	if success != "" {
		sess.Set("FlashSuccess", success)
	}
	if failure != "" {
		sess.Set("FlashError", failure)
	}

}

//	TODO: lang var
//
// bu yeni eklendi
func PopFlash(store *session.Store, c fiber.Ctx) Flash {
	sess := session.FromContext(c)

	ft, _ := sess.Get(flashTypeKey).(string)
	msg, _ := sess.Get(flashMessageKey).(string)
	if msg != "" {
		sess.Delete(flashTypeKey)
		sess.Delete(flashMessageKey)
	}
	return Flash{Type: ft, Message: msg}
}
