# Guida d'uso di WiRouter

Questa guida spiega come usare WiRouter nelle operazioni più comuni: avviare l'hotspot, cambiare SSID o password, attivare o disattivare l'avvio automatico.

Per i dettagli tecnici vedi [LOGICA_PROGRAMMA.md](./LOGICA_PROGRAMMA.md). Per lo stato di avanzamento vedi [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md).

---

## 1. Requisiti

- Windows 10 o 11 con Windows PowerShell 5.1 (incluso di serie).
- Una connessione Internet attiva sul PC, ad esempio via cavo Ethernet.
- Un adattatore Wi-Fi che supporti l'hotspot. Se non è supportato, il programma segnala un errore.
- **Nessun privilegio di amministratore**.

## 2. Compilazione

Da PowerShell, nella cartella del progetto:

```powershell
go build -o wirouter.exe ./cmd/wirouter
go build -ldflags "-H=windowsgui" -o wirouter-autostart.exe ./cmd/wirouter
```

- `wirouter.exe` è la versione con console, da usare a mano.
- `wirouter-autostart.exe` è la versione senza console, usata dall'avvio automatico. Va compilata solo se vuoi l'avvio automatico.

## 3. Uso manuale

1. Avvia `wirouter.exe`.
2. Scegli il numero della connessione Internet da condividere.
3. Inserisci SSID (nome della rete) e password (8–63 caratteri ASCII stampabili).
4. Il programma avvia l'hotspot e resta in attesa.
5. Per fermare l'hotspot premi `CTRL+C`. Il programma lo spegne e termina.

Per cambiare SSID o password, **ferma il programma e riavvialo**: i nuovi dati vengono chiesti di nuovo.

Se l'SSID o la password non sono validi, il programma mostra il motivo e chiede di nuovo il dato.

## 4. Avvio automatico al login

Con l'avvio automatico, WiRouter parte da solo ogni volta che accedi a Windows. Usa un'attività dell'Utilità di pianificazione con il nome `WiRouter`.

### Attivare

Da PowerShell, nella cartella dove si trovano i due eseguibili:

```powershell
.\wirouter.exe --install-autostart
```

Il comando chiede connessione, SSID e password e crea l'attività. Va eseguito una sola volta, oppure ogni volta che vuoi cambiare i dati (vedi sotto).

Il comando cerca `wirouter-autostart.exe` nella stessa cartella di `wirouter.exe`. Se manca, l'installazione si ferma con un errore e non crea nulla.

### Cosa succede all'accesso

- Dopo circa 30 secondi dall'accesso, WiRouter avvia l'hotspot.
- Se la connessione Internet non è ancora pronta, attende fino a 5 minuti.
- Se non riesce, mostra una finestra di errore.

### Cambiare SSID o password

Il programma non permette di cambiare SSID o password mentre è in esecuzione. Si cambia reinstallando l'attività:

```powershell
.\wirouter.exe --install-autostart
```

Il comando sovrascrive l'attività precedente e chiede di nuovo connessione, SSID e password. Il processo già in esecuzione non cambia: per applicare i nuovi dati, fermalo come descritto in [sezione 6](#6-fermare-lhotspot-e-il-programma).

### Disattivare l'avvio automatico

Non serve ricordare i comandi: c'è una guida dedicata con i metodi da riga di comando e da interfaccia grafica: [DISATTIVARE_AVVIO_AUTOMATICO.md](./DISATTIVARE_AVVIO_AUTOMATICO.md).

In breve:

```powershell
.\wirouter.exe --uninstall-autostart
```

Il comando rimuove l'attività, quindi WiRouter non partirà più al login. **Non ferma** il processo già in esecuzione né l'hotspot: vedi la sezione 6.

## 5. Verificare lo stato

- **Hotspot**: *Impostazioni → Rete e Internet → Hotspot mobile*. Mostra se è attivo e il nome della rete.
- **Attività**: *Utilità di pianificazione* (`taskschd.msc`), cartella *Libreria Utilità di pianificazione*, attività `WiRouter`. Lo stato *Pronta* o *In esecuzione* indica che è attiva.
- **Processo**: in *Gestione attività* cerca `wirouter-autostart.exe`.

## 6. Fermare l'hotspot e il programma

| Situazione | Cosa fare |
|---|---|
| Uso manuale | Premi `CTRL+C` nella finestra di `wirouter.exe` |
| Avvio automatico | Spegni l'hotspot da *Impostazioni → Hotspot mobile*, poi chiudi `wirouter-autostart.exe` da *Gestione attività* |

Chiudere la finestra o terminare il processo con `taskkill` o `Stop-Process` **non spegne l'hotspot**. Va spento a mano da Impostazioni.

Il comportamento al logout o allo spegnimento del PC è previsto (il programma dovrebbe spegnere l'hotspot), ma **non è ancora verificato**. Controllalo sempre dopo un logout.

## 7. Domande frequenti

**Posso cambiare la password senza riavviare nulla?**
No. Serve riavviare il programma (modalità manuale) oppure reinstallare l'attività (avvio automatico). La modifica dal programma è una funzione ancora da fare.

**La password è al sicuro?**
Per l'avvio automatico **no**: è salvata in chiaro nella definizione dell'attività e compare nella riga di comando del processo. Chi ha accesso al PC la può leggere. È una scelta accettata per questo PC. La protezione con DPAPI è ferma perché la policy aziendale blocca i file necessari.

**Il programma non trova la connessione Internet.**
Controlla che la connessione via cavo sia attiva. Con l'avvio automatico, il programma attende fino a 5 minuti, poi mostra un errore. Non c'è ancora un passaggio automatico al Wi-Fi.

**Ho avviato WiRouter a mano mentre l'avvio automatico era già attivo.**
Un solo processo alla volta è permesso. Chiudi quello già in esecuzione prima di avviarne un altro.

**Ho cambiato i dati ma il vecchio hotspot è ancora attivo.**
Ferma il processo e l'hotspot (sezione 6), poi riavvia.

**Dove devono stare i file per l'avvio automatico?**
L'attività punta a `wirouter-autostart.exe` nella cartella da cui hai lanciato `--install-autostart`. Se sposti o cancelli quella cartella, l'avvio automatico smette di funzionare: reinstalla l'attività dalla nuova posizione. Sul PC i file sono in `%LOCALAPPDATA%\WiRouter`.

## 8. Limiti noti

- Nessun fallback Wi-Fi quando la connessione via cavo non è disponibile.
- Il comportamento al logout non è ancora verificato.
- La password dell'avvio automatico è in chiaro.
- L'SSID e la password non si cambiano dal programma in esecuzione.
