package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danifx/dispatcher/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {

	// 1. Setup Iniziale
	ctx := context.Background()
	projectID := "eco-system-464513" // Il tuo ID progetto verificato

	// 2. Caricamento Configurazione dal Secret Manager
	// Durante i test locali, se non trova il segreto, gestiamo l'errore o usiamo default
	cfg := config.LoadConfig(ctx, projectID)
	log.Printf("🚀 Dispatcher inizializzato per l'ambiente: %s", cfg.Environment)

	// 3. Configurazione Server HTTP (Gin)
	// Usiamo il modo Release in produzione per performance migliori
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware di base
	r.Use(gin.Recovery()) // Recupera da eventuali panic
	r.Use(gin.Logger())   // Logga le richieste in arrivo

	// --- ROTTE ---

	// Health Check per Cloud Run (Fondamentale)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "active",
			"service": "eco-dispatcher",
		})
	})

	// Endpoint Principale: Ricezione richieste esterne
	r.POST("/v1/dispatch", func(c *gin.Context) {
		handleIncomingRequest(c, cfg)
	})

	// 4. Avvio Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default per Cloud Run e test locali
	}

	log.Printf("📡 Dispatcher in ascolto sulla porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Errore critico all'avvio: %v", err)
	}
}

// Struttura della richiesta che ci aspettiamo dal Web/App
type DispatchRequest struct {
	Action  string      `json:"action" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"`
}

func handleIncomingRequest(c *gin.Context, cfg *config.Config) {
	var req DispatchRequest

	// Validazione JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload non valido: " + err.Error()})
		return
	}

	log.Printf("📦 Ricevuta azione: %s", req.Action)

	// QUI ANDRÀ LA LOGICA gRPC NEI PROSSIMI STEP
	// Per ora rispondiamo con un ACK (Acknowledge)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Richiesta ricevuta e in fase di smistamento",
		"action":  req.Action,
	})
}
