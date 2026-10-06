# WiRouter - Piano di Implementazione

**NOTA**: Il piano è diviso in due percorsi:
- **DEMO Interattiva** (9-13 giorni): Programma terminal interattivo standalone
- **Versione Completa** (18-25 giorni): Con API REST, GUI Wails, e funzionalità avanzate

Questo documento descrive entrambi i percorsi.

---

## ⚡ DEMO Interattiva - Programma Terminal Standalone

**Timeline**: 9-13 giorni | **Dipendenze**: Solo stdlib

Programma Go semplice che:
1. Stampa tutte le interfacce di rete
2. Chiede all'utente quale scegliere
3. Chiede SSID e password
4. Avvia l'hotspot automaticamente
5. Rimane attivo fino a CTRL+C

Dettagli: vedi [ROADMAP_CLI_DEMO.md](ROADMAP_CLI_DEMO.md)

---

## 📋 Fase 1: Struttura di Base (FONDAMENTALI)

### 1.1 Organizzazione Progetto
- [x] Creare struttura cartelle (`model/`, `Terminal/`, `docs/`)
- [x] Configurare `go.mod`
- [x] Creare `README.md`
- [x] Organizzare documentazione in `docs/`
- [ ] Creare cartella `cmd/wirouter/` per entry point principale
- [ ] Creare cartella `internal/` per logica interna
- [ ] Aggiungere `.gitignore`

### 1.2 Modello Dati
- [x] Definire struct `Adapter` e `AdapterWin`
- [ ] Creare struct `Hotspot` (SSID, password, stato)
- [ ] Creare struct `NetworkDevice` (dispositivi collegati)
- [ ] Creare struct `TrafficStats` (statistiche traffico)
- [ ] Aggiungere enum per tipi di connessione

---

## 🔌 Fase 2: Core Networking (BACKEND)

### 2.1 Gestione Adattatori di Rete
- [x] Implementare `ListWindowsAdapters()` (in `model/win.go`)
- [ ] Aggiungere `GetAdapterDetails()` (IP, gateway, DNS)
- [ ] Filtrare adattatori per tipo (Ethernet, Wi-Fi)
- [ ] Identificare quale adattatore ha accesso a Internet
- [ ] Aggiungere error handling robusto

### 2.2 Hotspot Management
- [ ] Implementare `StartHotspot(ssid, password, adapter)`
- [ ] Implementare `StopHotspot()`
- [ ] Implementare `GetHotspotStatus()` 
- [ ] Implementare `ChangeSSID(newSSID)`
- [ ] Implementare `ChangePassword(newPassword)`
- [ ] Interfacciarsi con Windows Hotspot API

### 2.3 Internet Connection Sharing (ICS)
- [ ] Implementare `ShareConnection(source, destination)` 
- [ ] Configurare ICS via COM (Ethernet → Wi-Fi)
- [ ] Configurare ICS (Wi-Fi → Wi-Fi)
- [ ] Rilevare la connessione Internet attiva
- [ ] Gestire il cambio dinamico della rete condivisa

### 2.4 Monitoraggio Client
- [ ] Implementare `GetConnectedDevices()`
- [ ] Estrarre: IP, MAC, Nome host, Ora connessione
- [ ] Aggiornare lista dispositivi in tempo reale
- [ ] Implementare polling o listener per cambamenti

### 2.5 Monitoraggio Traffico
- [ ] Implementare `GetNetworkStats()`
- [ ] Tracciare upload/download totale
- [ ] Calcolare velocità in tempo reale
- [ ] Monitoraggio per singolo dispositivo
- [ ] Utilizzare netsh o API Windows per statistiche

---

## 🎨 Fase 3: API & Backend Services (Versione Completa)

**Nota**: Questa fase è per dopo la DEMO interattiva.

### 3.1 Struttura Servizi
- [ ] Creare `Service` interface per operazioni principali
- [ ] Implementare Service per Hotspot
- [ ] Implementare Service per Network
- [ ] Implementare Service per Monitoring

### 3.2 REST API
- [ ] Setup server HTTP (net/http o Gin)
- [ ] Endpoint GET `/api/adapters` - lista adattatori
- [ ] Endpoint GET `/api/hotspot/status` - stato hotspot
- [ ] Endpoint POST `/api/hotspot/start` - avvia hotspot
- [ ] Endpoint POST `/api/hotspot/stop` - arresta hotspot
- [ ] Endpoint POST `/api/hotspot/config` - configura SSID/password
- [ ] Endpoint GET `/api/devices` - dispositivi collegati
- [ ] Endpoint GET `/api/stats` - statistiche traffico
- [ ] Aggiungere autenticazione/autorizzazione
- [ ] Aggiungere error handling HTTP

### 3.3 WebSocket (opzionale)
- [ ] Setup WebSocket per aggiornamenti real-time
- [ ] Stream statistiche traffico
- [ ] Notifiche dispositivi collegati/disconnessi

---

## 🖥️ Fase 4: Frontend (WAILS) - Versione Completa

**Nota**: Questa fase è per dopo la DEMO interattiva. Serve per aggiungere una GUI moderna.

### 4.1 Setup Wails
- [ ] Inizializzare progetto Wails
- [ ] Configurare build process
- [ ] Setup struttura cartelle frontend

### 4.2 UI Componenti Principali
- [ ] Dashboard principale
- [ ] Pannello hotspot (stato, SSID, password)
- [ ] Pannello adattatori (selezione connessione da condividere)
- [ ] Lista dispositivi collegati
- [ ] Grafici statistiche traffico
- [ ] Impostazioni generali

### 4.3 Logica Frontend
- [ ] Comunicazione con backend (HTTP/WebSocket)
- [ ] State management
- [ ] Update UI real-time
- [ ] Gestione errori e notifiche

### 4.4 Styling
- [ ] Design moderno e intuitivo
- [ ] Responsive design
- [ ] Tema chiaro/scuro (opzionale)
- [ ] Icone per configurazioni

---

## 🔒 Fase 5: Security & Robustness

### 5.1 Sicurezza
- [ ] Richiedere permessi amministratore
- [ ] Validare input (SSID, password)
- [ ] Implementare rate limiting
- [ ] Proteggere API REST
- [ ] Logging delle azioni

### 5.2 Gestione Errori
- [ ] Catturare eccezioni Windows API
- [ ] User-friendly error messages
- [ ] Retry logic per operazioni fallite
- [ ] Log dettagliati per debugging

### 5.3 Testing
- [ ] Unit test per modelli
- [ ] Unit test per funzioni core
- [ ] Integration test con Windows API
- [ ] Test UI
- [ ] Test stress/performance

---

## 📦 Fase 6: Packaging & Release

### 6.1 Build & Distribution
- [ ] Cross-compile per Windows (x86, x64)
- [ ] Creare installer MSI
- [ ] Code signing
- [ ] Release notes

### 6.2 Documentazione
- [ ] User guide
- [ ] Developer guide
- [ ] API documentation
- [ ] Troubleshooting guide

### 6.3 CI/CD
- [ ] Setup GitHub Actions
- [ ] Automated builds
- [ ] Automated tests
- [ ] Release automation

---

## 🚀 Timeline Suggerito

### 📍 MILESTONE 1: DEMO Interattiva (9-13 giorni)
Programma terminal standalone senza API/GUI:

| Fase | Durata | Priorità |
|------|--------|----------|
| Fase 1: Struttura Base | 1-2 giorni | 🔴 CRITICA |
| Fase 2: Core Hotspot Logic | 3-4 giorni | 🔴 CRITICA |
| Fase 3: ICS & Routing | 2-3 giorni | 🔴 CRITICA |
| Fase 4: Terminal Interattiva | 1-2 giorni | 🟠 ALTA |
| Fase 5: Testing & Bug Fix | 1 giorno | 🟠 ALTA |
| Fase 6: Packaging | 1 giorno | 🟡 MEDIA |
| **TOTALE DEMO** | **9-13 giorni** | |

### 🎯 MILESTONE 2: Versione Completa (18-25 giorni totali)
Aggiungere API REST, GUI Wails, e funzionalità avanzate:

| Fase | Durata | Priorità |
|------|--------|----------|
| Fase 1: Struttura Base | 1-2 giorni | 🔴 CRITICA |
| Fase 2: Core Networking | 5-7 giorni | 🔴 CRITICA |
| Fase 3: API & Services | 3-4 giorni | 🟠 ALTA |
| Fase 4: Frontend (Wails) | 4-5 giorni | 🟠 ALTA |
| Fase 5: Security & Testing | 3-4 giorni | 🟡 MEDIA |
| Fase 6: Packaging & Release | 2-3 giorni | 🟡 MEDIA |

**Tempo totale stimato: 18-25 giorni**

---

## 📝 Note Importanti

### Dipendenze Critiche

**Per DEMO Interattiva (9-13 giorni):**
- `golang.org/x/sys/windows` - Accesso API Windows
- Standard library: `os`, `fmt`, `bufio`, `strings`, `text/tabwriter`, `strconv`, `os/exec`

**Per Versione Completa (18-25 giorni):**
- `golang.org/x/sys/windows` - Accesso API Windows
- `Wails` - Frontend framework
- `net/http` o `Gin` - HTTP server

### Considerazioni Windows
- Richiedere privilegi amministratore
- Testare con antivirus attivi
- Compatibilità con diverse versioni di Windows
- Interfaccia con Windows Hotspot Service

### Rischi Potenziali
- Compatibilità hardware (non tutti i driver Wi-Fi supportano hotspot)
- Permessi di amministratore richiesti
- Instabilità dovuta a ICS
- Interferenza con software antivirus/firewall

### Prossimi Passi Immediati (DEMO Interattiva)
1. ✅ Organizzazione progetto completata
2. → Strutturare `cmd/wirouter/` e `internal/hotspot/`, `internal/network/`
3. → Implementare `ListWiFiAdapters()` con filtro per tipo Wi-Fi
4. → Implementare hotspot start/stop via `netsh wlan hostednetwork`
5. → Implementare ICS routing via netsh
6. → Creare programma interattivo in `cmd/wirouter/main.go`
7. → Build e test su Windows


