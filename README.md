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
- COM (Internet Connection Sharing)
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
- Cambio SSID
- Cambio password
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

- [ ] Creazione progetto
- [ ] Setup Wails
- [ ] Gestione configurazione
- [ ] Individuazione schede di rete

---

## Fase 2

- [ ] Attivazione hotspot
- [ ] Disattivazione hotspot
- [ ] Gestione SSID
- [ ] Gestione password
- [ ] Configurazione Internet Connection Sharing (ICS)

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
           │ ICS            │
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

# 🔒 Permessi Richiesti

L'applicazione richiederà privilegi di amministratore per:

- Gestione hotspot
- Configurazione ICS
- Accesso alle API di rete
- Modifica firewall (eventuale)

---

# 🎯 Obiettivo Finale

Creare un software moderno, gratuito e open source che rappresenti una valida alternativa all'Hotspot integrato di Windows, offrendo un controllo completo sulla condivisione della connessione, sul monitoraggio dei dispositivi e sulle statistiche di rete.

---

# 📄 Licenza

MIT License (provvisoria)

---

# 👨‍💻 Autore

Progetto personale sviluppato in Go con l'obiettivo di creare un hotspot Wi-Fi avanzato per Windows.