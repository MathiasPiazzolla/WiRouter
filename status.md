# WiRouter - Progetto Go Windows Hotspot

## Obiettivo

Creare un programma in Go per Windows che permetta di:

1. Elencare le schede di rete
2. Identificare quale scheda ha accesso a Internet
3. Scegliere una scheda Wi-Fi
4. Creare e gestire un hotspot
5. Condividere la connessione tramite Windows
6. Aggiungere successivamente una GUI

Attualmente il progetto è in fase terminale.

---

# Ambiente

## Modulo Go

```go
module github.com/MathiasPiazzolla/WiRouter

go 1.26.5
```

## Dipendenze

Installata:

```bash
go get golang.org/x/sys/windows
```

---

# Struttura attuale

```
WiRouter/
│
├── go.mod
│
├── model/
│   └── adapter.go
│
├── adapters/
│   ├── windows.go
│   ├── types.go
│   └── mac.go
│
└── Terminal/
    └── main.go
```

---

# Concetti Go utilizzati

## Package

Tutti i file nella stessa cartella con lo stesso package appartengono allo stesso modulo.

Esempio:

```go
package adapters
```

---

## Visibilità

In Go la visibilità dipende dalla prima lettera.

Pubblico:

```go
type Adapter struct{}

func List() {}
```

Privato:

```go
type adapter struct{}

func list() {}
```

Maiuscola = esportato fuori dal package.

Minuscola = solo interno.

---

# Modello Adapter

File:

```
model/adapter.go
```

Contiene la struttura dati della scheda di rete.

```go
package model

type Adapter struct {
	Name        string
	Description string
	Type        string
	MAC         string
	IP          string
	Gateway     string
	HasInternet bool
}
```

---

# Windows API utilizzata

Viene usata:

```
iphlpapi.dll
```

tramite:

```
GetAdaptersAddresses()
```

Questa API Windows permette di ottenere informazioni sulle interfacce di rete:

- nome
- descrizione hardware
- MAC
- IP
- gateway
- tipo interfaccia
- stato

---

# Caricamento API Windows

```go
var (
	iphlpapi = windows.NewLazySystemDLL("iphlpapi.dll")

	getAdaptersAddresses = iphlpapi.NewProc(
		"GetAdaptersAddresses",
	)
)
```

---

# Struttura Windows letta

La struttura C di Windows viene ricreata in Go:

```go
type ipAdapterAddresses struct {

	Length uint32
	IfIndex uint32

	Next uintptr

	AdapterName uintptr

	FirstUnicastAddress uintptr
	FirstAnycastAddress uintptr
	FirstMulticastAddress uintptr
	FirstDnsServerAddress uintptr

	DnsSuffix uintptr
	Description uintptr
	FriendlyName uintptr

	PhysicalAddress [8]byte
	PhysicalAddressLength uint32

	Flags uint32

	Mtu uint32
	IfType uint32
}
```

---

# Problema risolto

Problema iniziale:

```
Tipo: Unknown (1500)
```

Il valore `1500` non era il tipo di interfaccia.

Era il valore:

```
MTU
```

La struttura Windows era quindi disallineata.

Dopo la correzione:

Output corretto:

```
----------------
Nome: Ethernet
Descrizione: Intel(R) Ethernet Connection (17) I219-LM
Tipo: Ethernet

----------------
Nome: Wi-Fi
Descrizione: Intel(R) Wi-Fi 6E AX211 160MHz
Tipo: Wi-Fi
```

---

# Schede trovate

Esempio:

```
Ethernet
Intel(R) Ethernet Connection (17) I219-LM

Wi-Fi
Intel(R) Wi-Fi 6E AX211 160MHz

Connessione alla rete locale (LAN)* 3
Microsoft Wi-Fi Direct Virtual Adapter

Connessione alla rete locale (LAN)* 4
Microsoft Wi-Fi Direct Virtual Adapter
```

I Microsoft Wi-Fi Direct Virtual Adapter saranno probabilmente utili per la gestione hotspot.

---

# MAC Address

Aggiunto il campo:

```go
MAC string
```

Funzione:

```go
func formatMAC(mac [8]byte, length uint32) string
```

Converte:

```
AA BB CC DD EE FF
```

in:

```
AA:BB:CC:DD:EE:FF
```

---

# Prossimi step

## 1. Recuperare indirizzi IP

Usare:

```
FirstUnicastAddress
```

Obiettivo:

```
Wi-Fi
IP: 192.168.1.50
```

---

## 2. Recuperare gateway

Usare:

```
FirstGatewayAddress
```

Obiettivo:

```
Ethernet
Gateway: 192.168.1.1
```

---

## 3. Identificare la connessione Internet

Prima versione:

```go
Gateway != ""
```

significa:

```
probabile connessione Internet
```

Successivamente:

usare la tabella di routing Windows:

```
0.0.0.0/0
```

per sapere quale interfaccia viene usata realmente.

---

# Architettura finale prevista

```
             Internet
                 |
                 |
        Ethernet / Wi-Fi
                 |
                 |
       Windows Internet Sharing
                 |
                 |
             Hotspot Wi-Fi
                 |
                 |
          Dispositivi collegati
```

---

# Prossimo punto di sviluppo

Aggiungere a `Adapter`:

- IP
- Gateway
- stato connessione

per distinguere automaticamente:

```
WAN (Internet)
        |
        |
HOTSPOT (condivisione)
```


# WiRouter - Stato progetto

## Obiettivo

Creare un programma in Go per Windows che permetta di trasformare il PC in un hotspot/router software:

- rilevare le schede di rete
- capire quale scheda ha Internet
- scegliere una scheda Wi-Fi
- condividere la connessione
- creare un hotspot
- aggiungere in futuro una GUI

---

# Stato attuale

## ✅ Ambiente configurato

- Go installato e funzionante
- Modulo Go creato:

```go
module github.com/MathiasPiazzolla/WiRouter