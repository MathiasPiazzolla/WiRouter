# WiRouter - Roadmap Demo Terminal Program

## 🎯 Obiettivo
Creare un programma Go interattivo da **terminale che stampa tutte le interfacce di rete**, l'utente sceglie quale usare, e il programma **avvia automaticamente l'hotspot** con SSID e password configurabili.

---

## 📊 Fase 1: Fondamenta (Struttura Base)
**Durata: 1-2 giorni** | **Priorità: 🔴 CRITICA**

### 1.1 Riorganizzare Struttura Progetto
```
WiRouter/
├── cmd/
│   └── wirouter/
│       └── main.go           ← Entry point CLI principale
├── internal/
│   ├── hotspot/
│   │   └── hotspot.go        ← Logica hotspot (START/STOP)
│   ├── network/
│   │   └── windows.go        ← Windows API calls (adapter listing)
│   └── utils/
│       └── admin.go          ← Check privilegi admin
├── model/
│   ├── adapter.go            ← Struct dati (mantenere)
│   └── win.go                ← Windows API (mantenere)
├── Terminal/main.go          ← Backup/reference (deprecato)
└── go.mod
```

**Task:**
- [X] Creare cartella `cmd/wirouter/`
- [X] Creare cartella `internal/hotspot/`, `internal/network/`, `internal/utils/`
- [ ] Refactor `model/win.go` in `internal/network/`
- [ ] Creare `internal/utils/admin.go` per check privilegi
- [X] Implementare `cmd/wirouter/main.go` con logica interattiva
- [ ] Update `go.mod` per imports relativi

### 1.2 Setup Go Terminal Program - Approccio Interattivo
- [ ] Usare `bufio.Reader` per leggere input da stdin (no framework, no flag package)
- [ ] Nessun argomento da linea comando - tutto interattivo
- [ ] Setup error handling
- [ ] Input/output semplice con os e fmt

---

## 🔧 Fase 2: Core Hotspot Logic (Implementazione Windows API)
**Durata: 3-4 giorni** | **Priorità: 🔴 CRITICA**

### 2.1 Enumerare Tutti gli Adattatori di Rete
**File: `internal/network/windows.go`**

- [x] ✅ `ListWindowsAdapters()` - Già implementato
- [ ] Perfezionare per ottenere TUTTE le interfacce (non solo Wi-Fi)
- [ ] Aggiungere error handling robusto
- [ ] Testare output con diversi tipi di adapter

**Funzioni Finali Necessarie:**
```go
func ListAllAdapters() ([]AdapterWin, error)  // Elenca TUTTE le interfacce
func ValidateAdapter(adapter *AdapterWin) error
```

### 2.2 Implementare Hotspot Start via netsh
**File: `internal/hotspot/hotspot.go`**

Windows permette di gestire hotspot tramite `netsh wlan hostednetwork` o Windows API.
Approccio CLI via `netsh` è più semplice e stabile.

**Comandi netsh necessari:**
```powershell
# Impostare SSID e password
netsh wlan set hostednetwork mode=allow ssid="MioHotspot" key="password123"

# Avviare hotspot
netsh wlan start hostednetwork

# Fermare hotspot
netsh wlan stop hostednetwork

# Visualizzare stato
netsh wlan show hostednetwork
```

**Struct per hotspot:**
```go
type HotspotConfig struct {
    SSID     string
    Password string
    Adapter  string // es. "Wi-Fi"
}

type HotspotStatus struct {
    IsRunning   bool
    SSID        string
    ConnectedDevices int
}
```

**Funzioni da implementare:**
- [ ] `SetHotspotConfig(config HotspotConfig) error` - Imposta SSID/password
- [ ] `StartHotspot() error` - Avvia hotspot
- [ ] `StopHotspot() error` - Arresta hotspot
- [ ] `GetHotspotStatus() (HotspotStatus, error)` - Verifica stato
- [ ] `ValidatePassword(password string) error` - Valida password (min 8 chars)
- [ ] `ValidateSSID(ssid string) error` - Valida SSID

**Implementazione interna:**
```go
func execCommand(cmd string, args ...string) (string, error) {
    // Esegue comando Windows tramite cmd.exe
    // Necessita privilegi admin
}
```

### 2.3 Verificare Privilegi Amministratore
**File: `internal/utils/admin.go`**

- [ ] Creare funzione `IsAdmin() bool`
- [ ] Se non admin, visualizzare messaggio di errore chiaro
- [ ] Forse implementare richiesta elevazione privilegi (opzionale per fase 1)

---

## 📋 Fase 3: Internet Connection Sharing (ICS)
**Durata: 2-3 giorni** | **Priorità: 🟠 ALTA** (richiesto per demo funzionante)

### 3.1 Identificare Connessione Internet Attiva
**File: `internal/network/windows.go`**

- [ ] Funzione `GetInternetAdapter()` - Rileva quale adapter ha Internet
- [ ] Check gateway predefinito
- [ ] Check connectivity Google DNS (8.8.8.8)

**Funzione:**
```go
func GetInternetAdapter() (*AdapterWin, error)
```

### 3.2 Configurare ICS (Internet Connection Sharing)
**Approccio: Via netsh (semplice e stabile)**

```powershell
# Abilitare forwarding su Ethernet
netsh int ip set interface name="Ethernet" forwarding=enabled

# Disabilitare quando hotspot ferma
netsh int ip set interface name="Ethernet" forwarding=disabled
```

**Funzioni necessarie:**
- [ ] Implementare `EnableICS(sourceAdapter string) error` - Abilita forwarding
- [ ] Implementare `DisableICS(sourceAdapter string) error` - Disabilita forwarding
- [ ] Integrata nel flusso di `StartHotspot()` / `StopHotspot()`

**Nota**: Per DEMO semplice, ICS automatico è opzionale. L'utente potrebbe dover configurare manualmente se netsh ICS non funziona.

---

## 🖥️ Fase 4: Terminal Interface Interattiva
**Durata: 1-2 giorni** | **Priorità: 🟠 ALTA**

### 4.1 Flusso Interattivo del Programma
**File: `cmd/wirouter/main.go`**

Programma completamente interattivo:

```
1. Avvia programma: wirouter.exe
   ↓
2. Stampa lista interfacce di rete con numeri
   ↓
3. Chiede: "Seleziona interfaccia (es. 1): "
   ↓
4. Chiede: "Inserisci SSID hotspot: "
   ↓
5. Chiede: "Inserisci password (min 8 caratteri): "
   ↓
6. Avvia hotspot e mostra status
```

**Logica main():**
```go
func main() {
    // 0. Check privilegi admin (CRITICO!)
    if !isAdmin() {
        fmt.Println("❌ Errore: Richiesti privilegi amministratore!")
        return
    }
    
    // 1. Elenca interfacce
    adapters, err := listAllAdapters()
    if err != nil {
        fmt.Println("❌ Errore:", err)
        return
    }
    
    // 2. Mostra tabella interfacce
    printAdaptersTable(adapters)
    
    // 3. Legge scelta utente
    selectedIndex := getIntInput("Seleziona interfaccia (numero): ")
    selectedAdapter := adapters[selectedIndex]
    
    // 4. Legge SSID
    ssid := getStringInput("Inserisci SSID: ")
    
    // 5. Legge password
    password := getStringInput("Inserisci password: ")
    
    // 6. Valida input
    if err := validateSSID(ssid); err != nil { return }
    if err := validatePassword(password); err != nil { return }
    
    // 7. Avvia hotspot e ICS
    if err := startHotspotAndICS(selectedAdapter, ssid, password); err != nil {
        fmt.Println("❌ Errore:", err)
        return
    }
    
    fmt.Println("✅ Hotspot avviato! Premi CTRL+C per fermare.")
    
    // 8. Rimani in esecuzione finché non viene interrotto
    select {} // Rimane in attesa di CTRL+C
}
```

### 4.2 Funzioni Helper Necessarie
- [ ] Funzione `isAdmin() bool` - Verifica privilegi amministratore (CRITICA)
- [ ] Funzione `listAllAdapters() []AdapterWin` - Legge tutte le interfacce da Windows
- [ ] Funzione `printAdaptersTable(adapters)` - Stampa tabella con numero, nome, tipo, descrizione
- [ ] Funzione `getStringInput(prompt string) string` - Legge testo da stdin
- [ ] Funzione `getIntInput(prompt string) int` - Legge numero da stdin con validazione
- [ ] Funzione `validateIndex(index int, max int) error` - Valida indice e ritorna errore se non valido
- [ ] Funzione `validatePassword(pwd string) error` - Verifica password >= 8 caratteri
- [ ] Funzione `validateSSID(ssid string) error` - Verifica SSID non vuoto
- [ ] Funzione `startHotspotAndICS(adapter, ssid, pwd string) error` - Wrapper che combina hotspot + ICS

### 4.3 Output Formattazione
- [ ] Stampare output chiaro con fmt.Println/Printf
- [ ] Messaggi di successo (✅) e errore (❌)
- [ ] Tabella adapter in modo ordinato con text/tabwriter
- [ ] Informazioni leggibili per utente

---

## 🧪 Fase 5: Testing & Bug Fixing
**Durata: 1 giorno** | **Priorità: 🟠 ALTA**

### 5.1 Test Manuale
- [ ] Build eseguibile: `go build -o wirouter.exe ./cmd/wirouter`
- [ ] Eseguire: `wirouter.exe` (da PowerShell con privilegi admin)
- [ ] Verificare che tabella interfacce appaia correttamente
- [ ] Scegliere un'interfaccia (numero valido)
- [ ] Inserire SSID valido
- [ ] Inserire password >= 8 caratteri
- [ ] Verificare che hotspot appaia in Windows
- [ ] Verificare connessione da smartphone
- [ ] Premere CTRL+C per fermare l'hotspot
- [ ] Verificare che hotspot si sia fermato correttamente

### 5.2 Gestione Errori
- [ ] Validare che netsh sia disponibile
- [ ] Catturare output errori da Windows
- [ ] User-friendly error messages
- [ ] Logging dettagliato (opzionale)

### 5.3 Edge Cases
- [ ] Cosa succede se hotspot è già in esecuzione?
- [ ] Cosa succede se password è troppo corta?
- [ ] Cosa succede se adapter Wi-Fi non c'è?
- [ ] Cosa succede se nessun internet connection?

---

## 📦 Fase 6: Packaging per Demo
**Durata: 1 giorno** | **Priorità: 🟡 MEDIA**

### 6.1 Build Eseguibile
- [ ] Build con: `go build -o wirouter.exe ./cmd/wirouter`
- [ ] Eseguire in PowerShell (richiede privilegi admin)
- [ ] Testare su Windows 10/11

### 6.2 Documentazione Demo
- [ ] Creare file `DEMO.md` in docs/ con istruzioni
- [ ] Screenshot/example di utilizzo
- [ ] Prerequisiti (Windows 10+, privilegi admin)

---

## ⏱️ Timeline Totale

| Fase | Task | Giorni | Status |
|------|------|--------|--------|
| 1 | Struttura progetto | 1-2 | ⏳ TODO |
| 2 | Hotspot logic (netsh) | 3-4 | ⏳ TODO |
| 3 | ICS routing | 2-3 | ⏳ TODO |
| 4 | CLI interface | 1-2 | ⏳ TODO |
| 5 | Testing | 1 | ⏳ TODO |
| 6 | Packaging | 1 | ⏳ TODO |
| **TOTALE** | | **9-13 giorni** | |

---

### 🎯 **Demo Terminal Finale - Flusso Interattivo**

Programma completamente interattivo, basta eseguirlo:

```
C:\> wirouter.exe

========================================
       🌐 WiRouter - Hotspot Manager
========================================

🔍 Interfacce di rete disponibili:

┌───┬─────────────┬──────────────────┬───────────┐
│ # │ Nome        │ Descrizione      │ Tipo      │
├───┼─────────────┼──────────────────┼───────────┤
│ 1 │ Ethernet    │ Intel Ethernet   │ Ethernet  │
│ 2 │ Wi-Fi       │ Qualcomm Wi-Fi   │ Wi-Fi     │
│ 3 │ VirtualBox  │ VirtualBox Host  │ Virtual   │
└───┴─────────────┴──────────────────┴───────────┘

📌 Seleziona interfaccia per l'hotspot (numero): 2

📝 Inserisci SSID del hotspot: MioHotspot

🔐 Inserisci password (minimo 8 caratteri): password123

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Hotspot configurato con successo!
✅ Hotspot avviato!
✅ ICS abilitato: Ethernet → Wi-Fi

🎯 Il tuo hotspot "MioHotspot" è online!
Dispositivi possono connettersi a questo SSID.

Premi CTRL+C per fermare l'hotspot.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Dopo CTRL+C:**
```
⏹️  Arresto hotspot...
✅ Hotspot fermato
✅ ICS disabilitato

Grazie per aver usato WiRouter!
```

---

## ⚠️ Note Importanti

### Prerequisiti Assoluti
1. **Privilegi Amministratore** - Richiesto per netsh e hotspot
2. **Windows 10/11** - Funzionalità hostednetwork
3. **Adapter Wi-Fi compatibile** - Deve supportare hosted network

### Limitazioni Attuali (Fase CLI)
- Solo da CLI (no GUI)
- No monitoraggio real-time devices
- No statistiche traffico
- No QR code
- No Web UI

### Dipendenze Go
- `golang.org/x/sys/windows` - Windows API (già presente)
- Standard library: `os`, `os/exec`, `fmt`, `bufio`, `strings`, `text/tabwriter`, `strconv`

### Prossimi Step Dopo Demo CLI
1. Aggiungere monitoraggio dispositivi in tempo reale
2. Implementare REST API
3. Aggiungere web dashboard con Wails
4. Statistiche traffico dettagliate

