# WiRouter - Roadmap Demo Terminal Program

## 🎯 Obiettivo
Programma Go interattivo da **terminale**: l'utente sceglie la connessione Internet da condividere, inserisce SSID e password e il programma **avvia l'hotspot di Windows**, **senza privilegi di amministratore**.

## 🧭 Approccio tecnico
L'hotspot è gestito dall'API WinRT `Windows.Networking.NetworkOperators.NetworkOperatorTetheringManager`, la stessa usata dalle Impostazioni di Windows. L'API è invocata da Go tramite `powershell.exe` (Windows PowerShell 5.1); gli input utente passano **solo** da variabili d'ambiente `WIROUTER_*`, mai concatenati nello script. La condivisione Internet è gestita dal sistema a partire dal profilo di connessione scelto.

---

## 📊 Fase 1: Struttura

```
WiRouter/
├── cmd/wirouter/main.go        ← CLI interattiva
├── internal/hotspot/
│   ├── hotspot.go              ← API Go (validazione, start/stop/status)
│   └── script.go               ← Script PowerShell/WinRT
├── model/
│   ├── adapter.go              ← Struct dati
│   ├── hotspot.go              ← HotspotConfig / HotspotStatus
│   └── win.go                  ← Elenco adapter (GetAdaptersAddresses)
├── Terminal/deprecated.go      ← Riferimento (deprecato)
└── go.mod
```

- [x] Creare `cmd/wirouter/` e `internal/hotspot/`
- [x] `bufio.Reader` per stdin, nessun argomento da linea di comando
- [x] Error handling con messaggi chiari

---

## 🔧 Fase 2: Hotspot (WinRT)

**File: `internal/hotspot/`**

- [x] `ValidateSSID(ssid string) error` - non vuoto, max 32 byte
- [x] `ValidatePassword(password string) error` - 8-63 caratteri ASCII stampabili
- [x] `ListSources() ([]string, error)` - connessioni con accesso a Internet
- [x] `StartHotspot(config HotspotConfig) error` - configura e avvia
- [x] `StopHotspot(source string) error`
- [x] `GetHotspotStatus(source string) (HotspotStatus, error)`
- [x] Input passati tramite variabili d'ambiente (no injection)
- [x] Funzionamento senza privilegi di amministratore
- [ ] Avvio reale verificato (può essere bloccato da policy di sistema)

---

## 🖥️ Fase 3: Terminal Interface

**File: `cmd/wirouter/main.go`**

```
1. Avvia wirouter.exe
2. Elenco connessioni con accesso a Internet
3. Selezione connessione da condividere (numero)
4. SSID
5. Password (8-63 caratteri)
6. Avvio hotspot, attesa di CTRL+C
7. CTRL+C: arresto hotspot
```

- [x] Elenco numerato e selezione validata (riprova in caso di errore)
- [x] Validazione SSID e password con nuova richiesta
- [x] Arresto pulito con CTRL+C
- [ ] Gestione hotspot già attivo all'avvio

---

## 🧪 Fase 4: Test

### 4.1 Test automatici
- [x] Validazione input
- [x] Input solo via ambiente, mai nello script
- [x] Parsing stato e connessioni
- [x] Propagazione dei messaggi d'errore dello script

### 4.2 Test manuale
- [ ] `go build -o wirouter.exe ./cmd/wirouter`
- [ ] Eseguire da PowerShell come utente standard
- [ ] Verificare che l'hotspot appaia in Windows
- [ ] Connettere uno smartphone
- [ ] CTRL+C e verifica dell'arresto

### 4.3 Edge cases
- [ ] Hotspot già in esecuzione
- [ ] Nessuna connessione Internet
- [ ] Hotspot bloccato da policy
- [ ] Adapter Wi-Fi assente o non compatibile

---

## 📦 Fase 5: Packaging
- [ ] Build `wirouter.exe` e test su Windows 10/11
- [ ] `docs/DEMO.md` con istruzioni e prerequisiti

---

## ⚠️ Prerequisiti e limiti
- Windows 10/11 con Windows PowerShell 5.1
- Connessione Internet attiva e adapter Wi-Fi compatibile
- Solo CLI: nessun monitoraggio client, statistiche, QR code o Web UI

### Dipendenze Go
- `golang.org/x/sys/windows` (enumerazione adapter)
- Standard library: `os`, `os/exec`, `bufio`, `strings`, `strconv`
