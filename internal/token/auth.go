package token

import (
	"be-ai/util"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Allow(roleList ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var allowed bool
		authStr := getTokenAuth(c)
		if authStr == "" {
			return util.SendUnauth(c)
		}

		dataJWT, err := NewJWT().Verify(authStr)
		if err != nil {
			log.Println("error verify token :", err.Error())
			return util.SendUnauth(c)
		}

		//if CheckLogout(authStr) {
		//	log.Println("already logged out")
		//	return util.SendUnauth(c)
		//}

		if len(roleList) > 0 {
			for _, val := range roleList {
				if val == dataJWT.Role {
					allowed = true
					break
				}
			}

			if !allowed {
				return util.SendUnauth(c)
			}
		}

		return c.Next()
	}
}

func GetInfoAuth(c *fiber.Ctx) *Payload {
	authStr := getTokenAuth(c)
	dataJWT, _ := NewJWT().Verify(authStr)

	return dataJWT
}

func getTokenAuth(c *fiber.Ctx) string {
	authStr := c.Get("Authorization")
	if authStr == "" {
		return ""
	}

	return authStr[7:]
}
