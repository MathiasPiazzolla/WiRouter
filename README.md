# WiRouter

> Un'applicazione Windows sviluppata in Go che trasforma un PC in un router Wi-Fi intelligente.

---

# 📖 Descrizione

WiRouter è un software che permette di condividere una connessione Internet (Ethernet o Wi-Fi) creando un hotspot Wi-Fi gestibile tramite un'interfaccia moderna.

L'obiettivo non è semplicemente avviare l'Hotspot di Windows, ma offrire un pannello di controllo completo con monitoraggio, statistiche e funzionalità avanzate.

---

# 🎯 Obiettivi

- Trasformare un PC Windows in un hotspot Wi-Fi.
- Gestire la condivisione della connessione Internet.
- Monitorare i dispositivi collegati.
- Visualizzare statistiche di rete in tempo reale.
- Offrire un'interfaccia moderna e intuitiva.
- Essere leggero, veloce e open source.

---

# 🛠 Stack Tecnologico

## Backend

- Go
- Windows API
- WinRT (`NetworkOperatorTetheringManager`, via PowerShell)
- Goroutine
- JSON
- REST API (opzionale)

## Frontend

- Wails
- HTML
- CSS
- JavaScript

---

# 📂 Struttura del Progetto

```text
WiRouter
│
├── cmd/
│   └── wirouter/
│       └── main.go
│
├── internal/
│   ├── hotspot/
│   ├── network/
│   ├── clients/
│   ├── monitor/
│   ├── windows/
│   ├── settings/
│   ├── firewall/
│   ├── qr/
│   └── api/
│
├── frontend/
│
├── assets/
│
├── docs/
│
├── scripts/
│
├── go.mod
└── README.md
```

---

# 🚀 Funzionalità MVP

## Gestione Hotspot

- Avvio hotspot
- Arresto hotspot
- Cambio SSID e password (solo riavviando il programma o reinstallando l'avvio automatico, vedi [guida](docs/GUIDA_UTENTE.md))
- Stato hotspot

---

## Condivisione Internet

- Ethernet → Wi-Fi
- Wi-Fi → Wi-Fi (se supportato dall'hardware)
- Selezione automatica della connessione
- Cambio rete condivisa

---

## Gestione Client

Per ogni dispositivo collegato:

- Nome host
- Indirizzo IP
- MAC Address
- Ora di connessione
- Tempo di connessione
- Produttore del dispositivo (OUI Lookup)

---

## Dashboard

Visualizzazione di:

- Numero client
- Upload
- Download
- Velocità in tempo reale
- Stato hotspot
- Stato connessione Internet

---

# 📈 Monitoraggio

Statistiche in tempo reale:

- Upload totale
- Download totale
- Upload istantaneo
- Download istantaneo
- Utilizzo rete
- Storico

---

# ⭐ Versione 2

## QR Code

Generazione automatica del QR Code della rete Wi-Fi.

---

## Notifiche

- Nuovo dispositivo collegato
- Dispositivo disconnesso
- Errore hotspot
- Assenza connessione Internet

---

## Gestione Dispositivi

- Blocca dispositivo
- Lista nera
- Lista bianca
- Rinomina dispositivo

---

## Impostazioni

- Avvio automatico con Windows
- Avvio automatico hotspot
- Salvataggio configurazione
- Backup impostazioni

---

# 🌟 Versione Pro

- API REST
- Dashboard Web
- Multi lingua
- Tema chiaro/scuro
- Plugin
- Aggiornamenti automatici
- Esportazione log

---

# 🎨 Interfaccia

## Dashboard

```text
-----------------------------------------------------

                 WiRouter

-----------------------------------------------------

Hotspot

● ATTIVO

SSID

Casa

Password

********

Condivisione

Ethernet

Client collegati

5

Upload

32 Mbps

Download

180 Mbps

[ AVVIA ]

[ ARRESTA ]

-----------------------------------------------------
```

---

## Client

```text
Nome          IP              MAC                Stato

iPhone        192.168.137.2   XX:XX:XX           Online

Laptop        192.168.137.3   XX:XX:XX           Online

TV            192.168.137.4   XX:XX:XX           Online
```

---

# 📅 Roadmap

## Fase 1

- [x] Creazione progetto
- [ ] Setup Wails
- [ ] Gestione configurazione
- [x] Individuazione schede di rete

---

## Fase 2

- [x] Attivazione hotspot
- [x] Disattivazione hotspot
- [ ] Gestione SSID (impostato solo all'avvio, non modificabile)
- [ ] Gestione password (impostata solo all'avvio, non modificabile)
- [x] Selezione della connessione Internet condivisa

---

## Fase 3

- [ ] Elenco dispositivi collegati
- [ ] Dashboard
- [ ] Monitor traffico
- [ ] Grafici

---

## Fase 4

- [ ] QR Code
- [ ] Notifiche
- [ ] Impostazioni avanzate
- [ ] Log

---

## Fase 5

- [ ] API REST
- [ ] Dashboard Web
- [ ] Sistema plugin
- [ ] Auto Update

---

# 🔧 Architettura

```text
                Internet
                    │
          Ethernet / Wi-Fi
                    │
           ┌────────────────┐
           │    WiRouter    │
           │                │
           │ Hotspot        │
           │ Condivisione   │
           │ Monitor        │
           │ Dashboard      │
           │ API            │
           └────────────────┘
                    │
          Wi-Fi generato
                    │
        ┌───────────┴───────────┐
        │           │           │
    Smartphone   Laptop       Smart TV
```

---

# 📦 Possibili Moduli

- hotspot
- network
- windows
- settings
- monitor
- clients
- firewall
- api
- qr
- updater
- logger

---

# 🔒 Permessi

L'hotspot viene gestito tramite l'API WinRT `NetworkOperatorTetheringManager`
(la stessa delle Impostazioni di Windows): **non servono privilegi di
amministratore**. Policy aziendali o di sistema possono comunque limitare
l'Hotspot di Windows.

---

# 🖥️ Demo terminale

La demo corrente permette di:

- elencare e selezionare la connessione Internet da condividere;
- configurare SSID e password;
- avviare l'hotspot Windows;
- arrestare l'hotspot con `CTRL+C`.

Il programma accetta anche tre opzioni per l'avvio automatico (vedi sotto):
`--autostart`, `--install-autostart`, `--uninstall-autostart`.

Da PowerShell (utente standard):

```powershell
go build -o wirouter.exe ./cmd/wirouter
.\wirouter.exe
```

Requisiti: Windows 10/11 con Windows PowerShell 5.1 e una connessione Internet attiva.

Per l'uso passo per passo, cambio di SSID o password e disattivazione, vedi la
[guida d'uso](docs/GUIDA_UTENTE.md).

## Avvio automatico al login

Il programma può partire da solo all'accesso dell'utente tramite un'attività
dell'Utilità di pianificazione (nessun privilegio di amministratore).

```powershell
go build -o wirouter.exe ./cmd/wirouter
go build -ldflags "-H=windowsgui" -o wirouter-autostart.exe ./cmd/wirouter
.\wirouter.exe --install-autostart     # chiede connessione, SSID e password
.\wirouter.exe --uninstall-autostart   # rimuove l'attività
```

- `wirouter-autostart.exe` è la versione senza console usata dall'attività: deve stare
  nella stessa cartella di `wirouter.exe` al momento dell'installazione.
- All'accesso l'hotspot parte dopo 30 secondi. Se la connessione non è ancora
  disponibile, il programma attende fino a 5 minuti. Gli errori compaiono in una finestra.
- Al logout o allo spegnimento il programma dovrebbe fermare l'hotspot. Il comportamento
  non è ancora verificato sul PC: non si sa quanto Windows lascia al processo per terminare.
- Non c'è ancora un fallback sul Wi-Fi: se la connessione via cavo non è disponibile,
  l'hotspot non parte.
- Per cambiare dati, rieseguire `--install-autostart`: sovrascrive l'attività.
- Per disattivare l'avvio automatico, vedi [DISATTIVARE_AVVIO_AUTOMATICO.md](docs/DISATTIVARE_AVVIO_AUTOMATICO.md).
- **Attenzione**: la password è salvata in chiaro negli argomenti dell'attività
  e nella riga di comando del processo. Chi ha accesso al PC può leggerla.

---

# 🎯 Obiettivo Finale

Creare un software moderno, gratuito e open source che rappresenti una valida alternativa all'Hotspot integrato di Windows, offrendo un controllo completo sulla condivisione della connessione, sul monitoraggio dei dispositivi e sulle statistiche di rete.

---

# 📄 Licenza

MIT License (provvisoria)

---

# 👨‍💻 Autore

Progetto personale sviluppato in Go con l'obiettivo di creare un hotspot Wi-Fi avanzato per Windows.