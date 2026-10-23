# WiRouter - Logica del programma (versione tecnica)

Questo documento descrive il funzionamento interno di WiRouter: architettura, flusso di esecuzione, chiamate a Windows, validazione, gestione degli errori e limiti noti.

---

## 1. Obiettivo e vincoli di progetto

- **Scopo**: esporre come hotspot Wi-Fi una connessione Internet già attiva sul PC.
- **Nessun privilegio di amministratore**: la versione precedente usava `netsh wlan ... hostednetwork` e Internet Connection Sharing (ICS), che richiedono elevazione. Sull'hardware di test (Intel AX211) il driver riporta `Rete ospitata supportata: No`, quindi `hostednetwork` non è utilizzabile nemmeno da amministratore.
- **Approccio attuale**: API WinRT `NetworkOperatorTetheringManager`, la stessa usata da *Impostazioni → Rete e Internet → Hotspot mobile*. Verificata in modalità non elevata: `GetTetheringCapabilityFromConnectionProfile` restituisce `Enabled`.
- **Stack**: Go per CLI e logica; Windows PowerShell 5.1 come ponte verso l'API WinRT.

---

## 2. Architettura

| Livello | File | Responsabilità |
|---|---|---|
| CLI | [main.go](../cmd/wirouter/main.go) | I/O utente, ordine dei passi, gestione di CTRL+C |
| API di dominio | [hotspot.go](../internal/hotspot/hotspot.go) | Validazione, costruzione dell'input, esecuzione degli script, parsing dell'output |
| Ponte WinRT | [script.go](../internal/hotspot/script.go) | Script PowerShell che invoca WinRT |
| Modello | [hotspot.go (model)](../model/hotspot.go) | `HotspotConfig` e `HotspotStatus` |

```
main.go
 ├─ hotspot.ListSources()          → powershell.exe (scriptSources) → WinRT
 ├─ hotspot.StartHotspot(config)   → powershell.exe (scriptStart)   → WinRT
 ├─ attesa CTRL+C
 └─ hotspot.StopHotspot(source)    → powershell.exe (scriptStop)    → WinRT
```

Una sessione completa avvia **tre processi PowerShell separati**. Ogni chiamata è stateless: il `NetworkOperatorTetheringManager` viene ricreato ogni volta con `CreateFromConnectionProfile`. Il programma non conserva stato tra le chiamate; lo stato vero è quello di Windows.

---

## 3. Flusso di esecuzione (`main.go`)

1. `main` chiama `run(os.Stdin, os.Stdout)`. In caso di errore stampa su stderr ed esce con codice 1.
2. `hotspot.ListSources()` restituisce i nomi dei profili di connessione con `InternetAccess`. Se la lista è vuota, `run` termina con errore.
3. `readSource` chiede un numero, lo converte con `strconv.Atoi` e lo verifica nell'intervallo `[1, n]`. Riprova finché l'input non è valido.
4. `readValidated` chiede SSID e password. Riceve il validatore come parametro (`func(string) error`) e segue lo stesso schema: errore → messaggio → nuovo prompt.
5. `signal.Notify(interrupt, os.Interrupt)` viene registrato **prima** di `StartHotspot`. Il canale ha buffer 1, quindi un CTRL+C ricevuto durante l'avvio non va perso.
6. `hotspot.StartHotspot(config)` accende l'hotspot.
7. `<-interrupt` blocca il thread principale fino a CTRL+C.
8. `hotspot.StopHotspot(source)` spegne l'hotspot. `signal.Stop` è in `defer`, quindi resta attivo fino alla fine di `run`.

**Lettura dell'input** (`readLine`): usa `bufio.Reader.ReadString('\n')` e poi `TrimSpace`. Se arriva EOF con stringa vuota restituisce `input terminato`. Se arriva EOF con dati parziali accetta il valore.

---

## 4. Validazione (`hotspot.go`)

| Campo | Regola | Motivazione |
|---|---|---|
| SSID | Non vuoto dopo `TrimSpace`; al massimo 32 byte | Limite SSID dello standard 802.11 |
| Password | Da 8 a 63 byte; ogni byte nell'intervallo 32–126 (ASCII stampabile) | Passphrase WPA2-PSK (802.11i): 8–63 caratteri ASCII |

Note tecniche:

- `len` in Go conta **byte**, non caratteri. Un SSID con lettere accentate (es. `è` = 2 byte in UTF-8) raggiunge il limite prima. Il controllo è quindi quello previsto dallo standard.
- `StartHotspot` riesegue la validazione anche se il CLI l'ha già fatta (difesa in profondità).
- `ValidateSSID` non esclude i caratteri di controllo. Dal CLI non è possibile inserire un `\n` (il lettore si ferma lì), ma la regola non è esplicita. Vedi §10.

---

## 5. Il ponte PowerShell

### 5.1 Lancio del processo

`runScript` costruisce il comando tramite `executeCommand`:

```
powershell.exe -NoProfile -NonInteractive -Command <script>
```

- `-NoProfile`: non carica i profili utente, quindi l'output è sempre uguale.
- `-NonInteractive`: nessun prompt che possa bloccare il processo.
- `-Command`: lo script è passato come stringa. La *execution policy* riguarda i file `.ps1`, quindi non si applica.
- `CombinedOutput`: stdout e stderr vengono uniti. Se l'exit code è diverso da zero, il testo diventa il messaggio di errore Go.
- `buildScript` costruisce lo script finale: `scriptPrelude` + `try { body } catch { stderr; exit 1 }`.

Il motivo per cui si usa **Windows PowerShell 5.1** (`powershell.exe`) e non PowerShell 7 (`pwsh`): la proiezione dei tipi WinRT con `ContentType=WindowsRuntime` e il bridge `AsTask` sono supportati in 5.1. Non sono pensati per PowerShell 7.

### 5.2 Passaggio dei dati: variabili d'ambiente

| Variabile | Contenuto | Letta da |
|---|---|---|
| `WIROUTER_SOURCE` | Nome del profilo sorgente | `Get-Source` |
| `WIROUTER_SSID` | SSID scelto dall'utente | `scriptStart` |
| `WIROUTER_PASSWORD` | Password scelta dall'utente | `scriptStart` |

Il testo dello script **non contiene mai** i valori utente. Questo evita l'injection: un SSID come `Wi'Router"; Remove-Item *` resta un dato. Il test `TestStartHotspotPassesInputsOnlyThroughEnvironment` verifica che:

- le variabili arrivino nell'ambiente con i valori corretti;
- né l'SSID né la password compaiano nel testo dello script.

**Limite noto**: la password è visibile nell'ambiente del processo `powershell.exe` per la durata della chiamata, quindi a processi dello stesso utente che possono leggerne l'ambiente. È anche presente nella memoria del processo Go. Un'alternativa è passarla via stdin del processo figlio; non è stata implementata.

### 5.3 Prelude (`scriptPrelude`)

Eseguito prima di ogni script:

1. `$ErrorActionPreference = 'Stop'`: gli errori non terminanti diventano eccezioni, così il `catch` li intercetta.
2. `[Console]::OutputEncoding = UTF8`: l'output (nomi di profilo con caratteri non ASCII) arriva integro a Go.
3. `Add-Type -AssemblyName System.Runtime.WindowsRuntime` e la proiezione `[void][Tipo, Assembly, ContentType=WindowsRuntime]`: rendono disponibili `NetworkInformation`, `NetworkOperatorTetheringManager` e `NetworkOperatorTetheringAccessPointConfiguration`.
4. **Bridge async**. Le API WinRT restituiscono `IAsyncAction` o `IAsyncOperation<T>`, non i `Task` di .NET. Il prelude:
   - trova l'extension method `AsTask` di `WindowsRuntimeSystemExtensions`;
   - `Wait-Operation` seleziona la overload generica per `IAsyncOperation*`, la istanzia con `MakeGenericMethod($resultType)` e chiama `Wait(-1)`, restituendo `Result`;
   - `Wait-Action` fa lo stesso per `IAsyncAction`, senza risultato.

   `Wait(-1)` è **bloccante senza timeout**. Vedi §10.
5. `Get-Source`: se `WIROUTER_SOURCE` è valorizzata cerca il profilo per `ProfileName`; altrimenti usa `GetInternetConnectionProfile()`. Se non trova nulla lancia un'eccezione.
6. `Get-Manager`: `NetworkOperatorTetheringManager::CreateFromConnectionProfile(Get-Source)`.
7. `Assert-Success`: confronta `Status` con `'Success'`. Altrimenti lancia un'eccezione con lo stato e `AdditionalErrorMessage`.

### 5.4 Script per operazione

| Script | Funzione Go | Operazioni WinRT |
|---|---|---|
| `scriptSources` | `ListSources` | `GetConnectionProfiles()`, filtro su `GetNetworkConnectivityLevel() -eq 'InternetAccess'`, stampa dei `ProfileName` |
| `scriptStart` | `StartHotspot` | Legge `TetheringOperationalState` (se `On` → errore `Hotspot già attivo`); crea la configurazione con `Ssid` e `Passphrase`; `ConfigureAccessPointAsync`; `StartTetheringAsync`; `Assert-Success` |
| `scriptStop` | `StopHotspot` | `StopTetheringAsync`; `Assert-Success`. Non controlla lo stato prima |
| `scriptStatus` | `GetHotspotStatus` | `TetheringOperationalState`, `GetCurrentAccessPointConfiguration().Ssid`, `ClientCount`; stampa `state=`, `ssid=`, `clients=` |

Note:

- In `scriptStart` la configurazione deve essere applicata **prima** di `StartTetheringAsync`. Per questo le due chiamate sono in sequenza.
- Il confronto `TetheringOperationalState -eq 'On'` mette a confronto un enum con una stringa. Funziona perché PowerShell converte la stringa nel tipo dell'enum.
- `scriptStatus` esiste ma il CLI non lo chiama. Serve per diagnosi e test.

### 5.5 Parsing dell'output (`parseStatus`)

Formato dell'output: una riga `chiave=valore` per campo.

- Ogni riga è divisa al **primo** `=` con `strings.Cut`, quindi i valori possono contenere `=`.
- `state` → `IsRunning = EqualFold(value, "On")`.
- `clients` → `strconv.Atoi`; l'errore è ignorato e il valore resta 0.
- Le righe senza `=` vengono ignorate.

`ListSources` fa il parsing riga per riga con `TrimSpace`. Questo significa che un nome di profilo con spazi iniziali o finali verrebbe modificato, e poi `Get-Source` non lo troverebbe. Vedi §10.

---

## 6. Gestione degli errori

Ci sono due livelli:

1. **PowerShell**: `$ErrorActionPreference = 'Stop'` e il `catch` finale scrivono il messaggio su stderr ed escono con `exit 1`.
2. **Go**: se l'exit code non è zero, `runScript` restituisce il testo di stderr come errore. Le funzioni pubbliche lo avvolgono con un contesto (`avvio hotspot: %w`, `arresto hotspot: %w`, `lettura connessioni Internet: %w`).

Casi specifici:

- **Stato WinRT diverso da `Success`**: `Assert-Success` lancia un messaggio con l'enum `NetworkOperatorTetheringOperationStatus` (es. `Error`) e `AdditionalErrorMessage`, che spesso è vuoto.
- **`IAsyncAction` fallita** (es. `ConfigureAccessPointAsync`): `Wait` lancia una `AggregateException`. Il messaggio è generico (tipicamente `One or more errors occurred.`), quindi poco informativo.
- **Errori di validazione**: sono recuperabili. Il CLI ripropone il prompt.
- **Errori di Windows o di PowerShell**: sono terminali. Il programma esce con codice 1.

---

## 7. Gestione di CTRL+C e della chiusura

Questa è la parte più delicata del programma.

| Momento | Cosa succede |
|---|---|
| CTRL+C durante la selezione o l'inserimento | Nessuna registrazione attiva: il processo termina normalmente, l'hotspot non è ancora acceso |
| CTRL+C durante `StartHotspot` | Il segnale resta nel canale (buffer 1). Quando `StartHotspot` ritorna, `<-interrupt` si sblocca subito e il programma esegue lo stop |
| CTRL+C durante l'attesa | Nessun figlio PowerShell è in esecuzione. Go riceve il segnale e chiama `StopHotspot` |
| CTRL+C durante `StopHotspot` | Vedi sotto |
| Chiusura della finestra, `taskkill /F`, crash | Go non esegue lo stop. L'hotspot **resta attivo** |

**Punto da verificare** (comportamento standard di Windows, non ancora testato): CTRL+C viene inviato a **tutti** i processi collegati alla stessa console, compreso `powershell.exe` se è in esecuzione. Quindi un CTRL+C durante `StartHotspot` o `StopHotspot` può terminare il figlio PowerShell. La chiamata restituisce errore e lo stato dell'hotspot può essere incerto. In quel caso va controllato in *Impostazioni → Hotspot mobile*.

**Conseguenza del crash**: non c'è un file di stato per ripristinare l'hotspot. Al successivo avvio `scriptStart` trova `On` e restituisce `Hotspot già attivo`. Va spento a mano.

---

## 8. Strategia di test

- `executeCommand` è una variabile di funzione: nei test viene sostituita con una finta. Nessun test esegue PowerShell o tocca l'hotspot.
- I test coprono:
  - validazione di SSID e password, con casi limite e caratteri non ASCII;
  - passaggio degli input solo via ambiente, con un SSID di tipo injection;
  - rifiuto degli input non validi senza eseguire il comando;
  - propagazione del messaggio di errore dello script;
  - parsing di `ListSources` (incluse righe vuote e `\r\n`);
  - parsing di `parseStatus`.
- Il file [model/win.go](../model/win.go) (enumerazione adattatori con `GetAdaptersAddresses`) non fa parte del flusso del CLI. Ha test propri (`model/win_test.go`).

Verifica manuale eseguita **senza** privilegi di amministratore:

- `ListSources` → `[Ethernet]`.
- `GetHotspotStatus` → `IsRunning: false`.
- Avvio reale dell'hotspot: **non ancora eseguito**.

Comando di verifica: `go test -count=1 ./...`.

---

## 9. Cosa non fa ancora

- Non mostra i dispositivi collegati (`ClientCount` è letto solo da `scriptStatus`).
- Non mostra statistiche di traffico.
- Non ha interfaccia grafica né API web.
- Non ripristina lo stato dopo un crash.

---

## 10. Limiti noti e possibili miglioramenti

| # | Limite | Possibile miglioramento |
|---|---|---|
| 1 | `Wait(-1)` non ha timeout: una chiamata WinRT bloccata blocca il programma | Usare `Wait(timeoutMs)` e segnalare il timeout |
| 2 | La password è visibile nell'ambiente del processo figlio | Passarla via stdin del figlio |
| 3 | `StopHotspot` non controlla lo stato prima di fermare | Leggere `TetheringOperationalState` e rispondere in modo esplicito se è già `Off` |
| 4 | `ValidateSSID` non esclude i caratteri di controllo | Rifiutare i caratteri `< 32` e `127` |
| 5 | Un nome di profilo con spazi iniziali o finali viene modificato dal `TrimSpace` | Non fare trim sui nomi sorgente, oppure confrontare con trim da entrambe le parti |
| 6 | La sorgente Wi-Fi non è stata testata: il test è stato fatto solo con Ethernet | Test con sorgente Wi-Fi e verifica del supporto del driver |
| 7 | Le policy aziendali possono bloccare il tethering anche se la capability è `Enabled` | Da verificare sulla macchina del dominio |
| 8 | Non è chiaro se Windows conservi SSID e password dopo lo stop | Da verificare in *Impostazioni → Hotspot mobile* |
| 9 | Nessun ripristino dopo crash o chiusura della finestra | Salvare lo stato su file e ripristinarlo all'avvio |
| 10 | I messaggi di `AggregateException` sono generici | Estrarre `InnerException.Message` |
