package middleware

import (
	"time"

	"ims/app/facades"

	contractshttp "github.com/goravel/framework/contracts/http"
	contractssession "github.com/goravel/framework/contracts/session"
)

// StartSession initializes and saves sessions for stateful requests
func StartSession() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		// 1. Get the session cookie name from configuration
		sessionCookieName := facades.Config().GetString("session.cookie", "goravel_session")

		// 2. Retrieve session ID from cookie
		sessionID := ctx.Request().Cookie(sessionCookieName)

		// 3. Get the default session driver
		drv, err := facades.Session().Driver()
		if err != nil {
			ctx.Request().Next()
			return
		}

		// 4. Build the session store
		var sessionStore contractssession.Session
		if sessionID != "" {
			sessionStore, err = facades.Session().BuildSession(drv, sessionID)
		} else {
			sessionStore, err = facades.Session().BuildSession(drv)
		}
		if err != nil {
			ctx.Request().Next()
			return
		}

		// 5. Start the session
		sessionStore.Start()

		// 6. Set the session on the request context
		ctx.Request().SetSession(sessionStore)

		// 7. Process the request
		ctx.Request().Next()

		// 8. Save the session back to storage driver
		_ = sessionStore.Save()

		// 9. Write session ID cookie to response
		lifetime := facades.Config().GetInt("session.lifetime", 120)
		ctx.Response().Cookie(contractshttp.Cookie{
			Name:     sessionCookieName,
			Value:    sessionStore.GetID(),
			Path:     facades.Config().GetString("session.path", "/"),
			Domain:   facades.Config().GetString("session.domain", ""),
			MaxAge:   lifetime * 60,
			Expires:  time.Now().Add(time.Duration(lifetime) * time.Minute),
			Secure:   facades.Config().GetBool("session.secure", false),
			HttpOnly: facades.Config().GetBool("session.http_only", true),
			SameSite: facades.Config().GetString("session.same_site", "lax"),
		})
	}
}
