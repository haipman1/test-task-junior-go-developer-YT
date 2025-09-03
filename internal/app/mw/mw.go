package mw

import (
	"log"

	"github.com/labstack/echo/v4"
)

const roleAdmin = "admin"
const roleNone = ""

func RoleCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Request().Header.Get("User-Role")
		UserAgent := c.Request().Header.Get("User-Agent")

		ip := c.RealIP()

		switch role {
		case roleAdmin:
			log.Printf("admin logged with ip: %s!, with UserAgent: %s", ip, UserAgent)
		case roleNone:
			log.Printf("No User-Role logged with ip: %s, with UserAgent: %s", ip, UserAgent)
		}

		err := next(c)
		if err != nil {
			return err
		}
		return nil
	}
}
