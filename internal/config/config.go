package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Config definisce la struttura dei dati che carichiamo da GCP
type Config struct {
	ApiKey            string `json:"api_key"`
	InternalWorkerURL string `json:"internal_worker_url"`
	Environment       string `json:"environment"`
	ProjectID         string `json:"project_id"`
}

// LoadConfig recupera il segreto "DISPATCHER_CONFIG" da Secret Manager
func LoadConfig(ctx context.Context, projectID string) *Config {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("Errore creazione client Secret Manager: %v", err)
	}
	defer client.Close()

	// Costruiamo il path del segreto usando l'ID progetto reale
	secretPath := fmt.Sprintf("projects/%s/secrets/DISPATCHER_CONFIG/versions/latest", projectID)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	// Chiamata a GCP per ottenere il contenuto del segreto
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Printf("ATTENZIONE: Impossibile leggere il segreto %s: %v", secretPath, err)
		log.Println("Uso configurazioni di default per sviluppo locale...")
		return &Config{
			Environment: "development",
			ProjectID:   projectID,
		}
	}

	// Decodifichiamo il JSON ricevuto da GCP nella nostra struct Config
	var cfg Config
	if err := json.Unmarshal(result.Payload.Data, &cfg); err != nil {
		log.Fatalf("Errore nel parsing del JSON di configurazione: %v", err)
	}

	// Aggiungiamo l'ID progetto alla config per comodità
	cfg.ProjectID = projectID

	return &cfg
}
