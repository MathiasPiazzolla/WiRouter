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
- [x] Creare cartella `cmd/wirouter/` per entry point principale
- [x] Creare cartella `internal/` per logica interna
- [ ] Aggiungere `.gitignore`

### 1.2 Modello Dati
- [x] Definire struct `Adapter` e `AdapterWin`
- [x] Creare struct `Hotspot` (SSID, password, stato): `model.HotspotConfig` e `model.HotspotStatus`
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
- [x] Implementare `StartHotspot(config)`
- [x] Implementare `StopHotspot(source)`
- [x] Implementare `GetHotspotStatus(source)`
- [ ] Implementare `ChangeSSID(newSSID)`
- [ ] Implementare `ChangePassword(newPassword)`
- [x] Interfacciarsi con l'Hotspot di Windows (WinRT `NetworkOperatorTetheringManager` via PowerShell)

### 2.3 Condivisione Internet
La condivisione è gestita da `NetworkOperatorTetheringManager` a partire dal profilo di connessione scelto.
- [x] Rilevare le connessioni Internet attive
- [ ] Wi-Fi → Wi-Fi (se supportato dall'hardware)
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
- [ ] Utilizzare le API Windows per statistiche

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
- [x] Funzionamento con utente standard (nessuna elevazione richiesta)
- [x] Validare input (SSID, password)
- [ ] Implementare rate limiting
- [ ] Proteggere API REST
- [ ] Logging delle azioni

### 5.2 Gestione Errori
- [ ] Catturare eccezioni Windows API
- [ ] User-friendly error messages
- [x] Retry logic per operazioni fallite (avvio automatico: 3 tentativi)
- [ ] Log dettagliati per debugging

### 5.3 Testing
- [ ] Unit test per modelli
- [x] Unit test per funzioni core (validazione, runner, parsing, argomenti dell'attività)
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
- Testare con antivirus attivi
- Compatibilità con diverse versioni di Windows
- Interfaccia con Windows Hotspot Service

### Rischi Potenziali
- Compatibilità hardware (non tutti i driver Wi-Fi supportano hotspot)
- Policy aziendali che possono bloccare l'Hotspot di Windows
- Dipendenza da Windows PowerShell 5.1 per l'accesso a WinRT
- Interferenza con software antivirus/firewall

### Prossimi Passi Immediati (DEMO Interattiva)
1. ✅ Organizzazione progetto completata
2. ✅ Struttura `cmd/wirouter/` e `internal/hotspot/`
3. ✅ Elenco connessioni Internet condivisibili
4. ✅ Hotspot start/stop/status via WinRT (PowerShell)
5. ✅ Programma interattivo in `cmd/wirouter/main.go`
6. ✅ Avvio reale su Windows: il programma parte senza errori con sorgente Ethernet. Da confermare: hotspot visibile e spegnimento
7. → Test di login, logout e riavvio dopo errore (vedi Fase 4)

---

## 🚀 Fase 4: Avvio automatico al login

**Stato**: in corso. Il codice è collegato a `main.go` e committato (`c8b9be5`). L'attività è creata sul PC e il primo test manuale è riuscito. Restano da verificare login, logout e riavvio dopo errore.

**Decisioni confermate**
- Avvio tramite Utilità di pianificazione, attività "All'accesso" con ritardo di 30 s, riavvio fino a 3 volte in caso di errore, nessun limite di durata
- Password in chiaro come argomento dell'attività pianificata (scelta 2). DPAPI non usato: `internal/secret` è escluso dalla policy aziendale. Rischio accettato: la password è leggibile nella definizione dell'attività e nella riga di comando del processo, da chi ha accesso al PC
- Nessun file di configurazione: `--install-autostart` chiede i dati e crea l'attività con gli argomenti `--source`, `--ssid`, `--password`
- All'accesso chiusura o logout: l'hotspot viene fermato
- Due eseguibili: `wirouter.exe` con console (uso interattivo e installazione dell'avvio automatico) e `wirouter-autostart.exe` senza console (`-H=windowsgui`)
- Errori all'avvio automatico mostrati con una finestra di Windows

**Già implementato**
- [x] Istanza unica con mutex `Local\WiRouter` (`internal/instance`), in modalità interattiva e automatica
- [x] Finestra di errore di Windows (`internal/notify`), mostrata in modalità `--autostart`
- [x] Generazione XML dell'attività in UTF-16 con BOM, con test (`internal/autostart`)
- [x] Ciclo di avvio: attesa della fonte, riavvio se già acceso, tentativi di avvio, arresto allo stop (`internal/runner`, test con finto hotspot)
- [x] Adattatore `hotspot.Controller` che collega `internal/hotspot` a `runner.Hotspot`
- [x] `cmd/wirouter/main.go`: flag `--autostart`, `--install-autostart`, `--uninstall-autostart`, con test del round-trip degli argomenti
- [x] Finestre PowerShell e `schtasks` nascoste (`HideWindow`, `CREATE_NO_WINDOW`)
- [x] Build dei due eseguibili verificata (`wirouter.exe` e `wirouter-autostart.exe` con `-H=windowsgui`)
- [x] Verifica che `schtasks /create /xml` accetti il file generato
- [x] Attività `WiRouter` creata sul PC con `--install-autostart`: stato Ready, trigger all'accesso, ritardo 30 s, utente `ad-ct011351`
- [x] Test manuale con `schtasks /run`: processo `wirouter-autostart` in esecuzione, nessun errore
- [x] Documentazione: [DISATTIVARE_AVVIO_AUTOMATICO.md](DISATTIVARE_AVVIO_AUTOMATICO.md) e sezione nel README
- [x] [LOGICA_PROGRAMMA.md](LOGICA_PROGRAMMA.md) aggiornato con la modalità automatica

**Da fare**
- [ ] Verificare che l'hotspot sia acceso dopo il test (non ancora confermato con un dispositivo o nelle Impostazioni)
- [ ] Test di login reale (riavvio o logout e nuovo accesso): l'hotspot deve partire dopo circa 30 s
- [ ] Test di logout: l'hotspot deve essere fermato
- [ ] Test di riavvio dopo errore (3 tentativi)
- [ ] Test con cavo Ethernet scollegato: attesa fino a 5 minuti, poi finestra di errore
- [ ] Misurare il tempo che Windows concede al programma dopo il segnale di logout, per verificare che l'arresto dell'hotspot finisca in tempo
- [ ] **In attesa di decisione**: fallback sul Wi-Fi quando la LAN non è disponibile. Da decidere il tempo di attesa della LAN prima del passaggio al Wi-Fi. Da verificare se la scheda Wi-Fi supporta la condivisione con hotspot: il test richiede il cavo di rete scollegato.
- [ ] Disattivazione dell'avvio automatico documentata, ma non ancora provata sul PC



