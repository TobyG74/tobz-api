package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"

	"github.com/tobz/tobz-api/internal/apikey"
	"github.com/tobz/tobz-api/internal/auth"
	"github.com/tobz/tobz-api/internal/captcha"
	"github.com/tobz/tobz-api/internal/config"
	"github.com/tobz/tobz-api/internal/features/downloader"
	"github.com/tobz/tobz-api/internal/features/search"
	"github.com/tobz/tobz-api/internal/handlers"
	"github.com/tobz/tobz-api/internal/httpclient"
	"github.com/tobz/tobz-api/internal/middleware"
	"github.com/tobz/tobz-api/internal/response"
	"github.com/tobz/tobz-api/internal/whitelist"
)

func New(cfg *config.Config, db *gorm.DB) *fiber.App {
	tokens := auth.NewService(cfg, db)
	keys := apikey.NewService(db)
	cap := captcha.New(cfg.TurnstileSecret, cfg.CaptchaEnabled)
	oauth := auth.NewOAuthRegistry(cfg)
	safe := httpclient.NewSafe(30 * time.Second)
	wl := whitelist.NewService(db)

	h := handlers.New(cfg, db, tokens, keys, cap, oauth, safe, wl)

	app := fiber.New(fiber.Config{
		AppName:               "tobz-api v2 (Go)",
		BodyLimit:             1 << 20, // 1 MB
		ReadTimeout:           15 * time.Second,
		WriteTimeout:          20 * time.Second,
		DisableStartupMessage: cfg.IsProduction(),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			if code == fiber.StatusInternalServerError {
				return response.Internal(c)
			}
			return response.Fail(c, code, "error", err.Error())
		},
	})

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "DENY",
		ReferrerPolicy:            "no-referrer",
		CrossOriginEmbedderPolicy: "require-corp",
		HSTSMaxAge:                31536000,
		HSTSExcludeSubdomains:     false,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinOrigins(cfg.AllowedOrigins),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-API-Key",
		ExposeHeaders:    "X-Quota-Limit,X-Quota-Remaining",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	registerRoutes(app, cfg, h, tokens, keys, safe, wl)
	return app
}

func registerRoutes(app *fiber.App, cfg *config.Config, h *handlers.Handlers, tokens *auth.Service, keys *apikey.Service, safe *httpclient.Safe, wl *whitelist.Service) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.OK(c, fiber.Map{"status": "ok"})
	})

	v1 := app.Group("/api/v1")

	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return response.TooManyRequests(c, "Terlalu banyak permintaan, coba lagi nanti")
		},
	})

	authGroup := v1.Group("/auth", authLimiter)
	authGroup.Post("/register", h.Register)
	authGroup.Post("/login", h.Login)
	authGroup.Post("/refresh", h.Refresh)
	authGroup.Post("/logout", h.Logout)
	authGroup.Get("/me", middleware.RequireAuth(tokens), h.Me)
	authGroup.Post("/change-password", middleware.RequireAuth(tokens), h.ChangePassword)

	oauthGroup := authGroup.Group("/oauth")
	oauthGroup.Get("/:provider/login", h.OAuthLogin)
	oauthGroup.Get("/:provider/callback", h.OAuthCallback)

	keysGroup := v1.Group("/keys", middleware.RequireAuth(tokens))
	keysGroup.Post("/", h.CreateAPIKey)
	keysGroup.Get("/", h.ListAPIKeys)
	keysGroup.Delete("/:id", h.RevokeAPIKey)

	wlGroup := v1.Group("/whitelist", middleware.RequireAuth(tokens))
	wlGroup.Get("/", h.ListWhitelist)
	wlGroup.Post("/", h.AddWhitelist)
	wlGroup.Delete("/:id", h.RemoveWhitelist)

	svcLimiter := limiter.New(limiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: cfg.RateLimitExpiration,
		LimitReached: func(c *fiber.Ctx) error {
			return response.TooManyRequests(c, "Rate limit terlampaui")
		},
	})
	// Protected service group; each feature registers its own routes via <feature>.Register(svc, deps).
	svc := v1.Group("", svcLimiter, middleware.RequireAPIKey(keys, wl))
	downloader.Register(svc, safe)
	search.Register(svc, safe)

	app.Use(func(c *fiber.Ctx) error {
		return response.NotFound(c, "Endpoint tidak ditemukan")
	})
}

func joinOrigins(origins []string) string {
	out := ""
	for i, o := range origins {
		if i > 0 {
			out += ","
		}
		out += o
	}
	if out == "" {
		return "http://localhost:3000"
	}
	return out
}
