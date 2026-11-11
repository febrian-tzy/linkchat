package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/linkchat/api/internal/config"
	"github.com/linkchat/api/internal/database"
	"github.com/linkchat/api/internal/handler"
	"github.com/linkchat/api/internal/middleware"
	"github.com/linkchat/api/internal/repository"
	"github.com/linkchat/api/internal/service"
	"github.com/linkchat/api/internal/session"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to MySQL")

	// ── Repositories ─────────────────────────────────────────────
	ownerRepo        := repository.NewOwnerRepository(db)
	authRepo         := repository.NewAuthRepository(db)
	contactRepo      := repository.NewContactRepository(db)
	inboxRepo        := repository.NewInboxRepository(db)
	conversationRepo := repository.NewConversationRepository(db)
	broadcastRepo    := repository.NewBroadcastRepository(db)
	webhookRepo      := repository.NewWebhookRepository(db)
	apiKeyRepo       := repository.NewAPIKeyRepository(db)
	automationRepo   := repository.NewAutomationRepository(db)
	userRepo         := repository.NewUserRepository(db)

	// ── Services ─────────────────────────────────────────────────
	sessionStore := session.New(cfg.SessionSecret)
	ownerStore   := session.NewOwnerStore(cfg.SessionSecret)
	authSvc      := service.NewAuthService(authRepo)

	// ── Handlers ─────────────────────────────────────────────────
	ownerHandler        := handler.NewOwnerHandler(ownerRepo, ownerStore)
	authHandler         := handler.NewAuthHandler(authSvc, sessionStore)
	contactHandler      := handler.NewContactHandler(contactRepo)
	inboxHandler        := handler.NewInboxHandler(inboxRepo)
	conversationHandler := handler.NewConversationHandler(conversationRepo)
	broadcastHandler    := handler.NewBroadcastHandler(broadcastRepo)
	webhookHandler      := handler.NewWebhookHandler(webhookRepo)
	apiKeyHandler       := handler.NewAPIKeyHandler(apiKeyRepo)
	automationHandler   := handler.NewAutomationHandler(automationRepo)
	userHandler         := handler.NewUserHandler(userRepo, authRepo)

	// ── Router ───────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.CleanPath)
	r.Use(corsMiddleware)

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ── Public routes ────────────────────────────────────────────
	// Tenant auth
	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login",    authHandler.Login)
	r.Post("/api/v1/auth/logout",   authHandler.Logout)

	// Owner auth (separate system)
	r.Post("/api/v1/owner/auth/login",  ownerHandler.Login)
	r.Post("/api/v1/owner/auth/logout", ownerHandler.Logout)

	// ── Protected: Tenant routes ─────────────────────────────────
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(sessionStore))

		// Auth (tenant)
		r.Get("/api/v1/auth/me", authHandler.Me)

		// Contacts
		r.Get("/api/v1/contacts",       contactHandler.List)
		r.Post("/api/v1/contacts",      contactHandler.Create)
		r.Get("/api/v1/contacts/{id}",  contactHandler.Get)
		r.Put("/api/v1/contacts/{id}",  contactHandler.Update)
		r.Delete("/api/v1/contacts/{id}", contactHandler.Delete)

		// Inboxes
		r.Get("/api/v1/inboxes",        inboxHandler.List)
		r.Post("/api/v1/inboxes",       inboxHandler.Create)
		r.Get("/api/v1/inboxes/{id}",   inboxHandler.Get)
		r.Delete("/api/v1/inboxes/{id}", inboxHandler.Delete)

		// Conversations
		r.Get("/api/v1/conversations",                conversationHandler.List)
		r.Post("/api/v1/conversations",               conversationHandler.Create)
		r.Get("/api/v1/conversations/{id}",           conversationHandler.Get)
		r.Patch("/api/v1/conversations/{id}",         conversationHandler.UpdateStatus)
		r.Get("/api/v1/conversations/{id}/messages",  conversationHandler.ListMessages)
		r.Post("/api/v1/conversations/{id}/messages", conversationHandler.SendMessage)

		// Broadcasts
		r.Get("/api/v1/broadcasts",         broadcastHandler.List)
		r.Post("/api/v1/broadcasts",        broadcastHandler.Create)
		r.Get("/api/v1/broadcasts/{id}",    broadcastHandler.Get)
		r.Delete("/api/v1/broadcasts/{id}", broadcastHandler.Delete)

		// Webhooks
		r.Get("/api/v1/webhooks",           webhookHandler.List)
		r.Post("/api/v1/webhooks",          webhookHandler.Create)
		r.Delete("/api/v1/webhooks/{id}",   webhookHandler.Delete)

		// API Keys
		r.Get("/api/v1/api-keys",           apiKeyHandler.List)
		r.Post("/api/v1/api-keys",          apiKeyHandler.Create)
		r.Delete("/api/v1/api-keys/{id}",   apiKeyHandler.Delete)

		// Automations
		r.Get("/api/v1/automations",               automationHandler.List)
		r.Post("/api/v1/automations",              automationHandler.Create)
		r.Get("/api/v1/automations/{id}",          automationHandler.Get)
		r.Patch("/api/v1/automations/{id}/toggle", automationHandler.Toggle)
		r.Delete("/api/v1/automations/{id}",       automationHandler.Delete)

		// Users / Team
		r.Get("/api/v1/users",                  userHandler.List)
		r.Post("/api/v1/users/invite",          userHandler.Invite)
		r.Put("/api/v1/users/profile",          userHandler.UpdateProfile)
		r.Post("/api/v1/users/change-password", userHandler.ChangePassword)
		r.Patch("/api/v1/users/{id}/toggle",    userHandler.ToggleActive)
		r.Delete("/api/v1/users/{id}",          userHandler.Delete)
	})

	// ── Protected: Owner routes (completely separate from tenant auth) ──
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOwner(ownerStore))
		r.Get("/api/v1/owner/auth/me", ownerHandler.Me)
		r.Get("/api/v1/owner/tenants", ownerHandler.ListTenants)
		r.Get("/api/v1/owner/stats",   ownerHandler.Stats)
	})

	addr := fmt.Sprintf(":%s", cfg.APIPort)
	log.Printf("LinkChat API listening on %s", addr)
	if err = http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:3000"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
