# 🛰️ EcoSystem Dispatcher

Il **Dispatcher** è il gateway d'ingresso centrale dell'architettura EcoSystem. Si occupa di ricevere le richieste dalle applicazioni client (Web/Mobile), validarle, autenticarle e smistarle ai microservizi interni protetti.

## 🛠️ Stack Tecnologico
- **Linguaggio:** Go 1.22+
- **Interfaccia Esterna:** HTTP/REST (JSON)
- **Comunicazione Interna:** gRPC (Protocol Buffers)
- **Cloud Provider:** Google Cloud Platform (GCP)
- **Runtime:** Cloud Run (Dockerized)

## 🏗️ Architettura di Rete
Il Dispatcher opera come un ponte sicuro:
1. **Ingresso:** Esposto via HTTPS per Web App e App Mobile.
2. **Processo:** Valida i JWT tramite GCP Identity Platform e recupera configurazioni da Secret Manager.
3. **Uscita:** Inoltra le richieste ai microservizi in una **VPC chiusa** utilizzando gRPC per massime performance e latenza minima.



## 🔒 Sicurezza e Cloud Native
- **GCP Secret Manager:** Nessuna chiave o password è salvata nel codice o in file .env locali.
- **IAM Service Accounts:** Il servizio utilizza permessi granulari per invocare i microservizi interni.
- **Observability:** Integrato con Cloud Logging e Cloud Trace per il monitoraggio dei flussi.

## 🚀 Guida Rapida (Sviluppo Locale)

### Prerequisiti
- Go installato
- Google Cloud SDK (gcloud) configurato
- Accesso al progetto GCP di riferimento

### Installazione
```bash
git clone [https://github.com/tuo-org/dispatcher.git](https://github.com/tuo-org/dispatcher.git)
cd dispatcher
go mod download
```

### Eseguzione
```bash
go run cmd/server/main.go
```

## 📂 Struttura del Progetto
- ```/cmd```: Entry point dell'applicazione.
- ```/internal/api```: Gestione degli endpoint HTTP esterni.
- ```/internal/grpcclient```: Logica di comunicazione verso i microservizi privati.
- ```/internal/config```: Integrazione nativa con GCP Secret Manager.
- ```/proto```: Definizione dei contratti gRPC.

```plaintext
/dispatcher
├── cmd/
│   └── server/          # Main entry point (main.go)
├── internal/
│   ├── api/             # Handlers HTTP (Gin/Echo)
│   ├── grpcclient/      # Client per parlare con i microservizi interni
│   ├── config/          # Caricamento segreti da GCP Secret Manager
│   └── service/         # Logica di business del dispatcher
├── proto/               # File .proto per gRPC
├── scripts/             # Script di deploy/build
├── README.md
├── go.mod
└── Dockerfile
```