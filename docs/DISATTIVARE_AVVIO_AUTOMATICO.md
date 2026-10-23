# Disattivare l'avvio automatico di WiRouter

Questa guida spiega come rimuovere o sospendere l'avvio automatico di WiRouter
al login di Windows. Vale per l'attività pianificata `WiRouter` creata con
`--install-autostart`.

## Cosa succede quando si disattiva

Disattivare l'avvio automatico **non spegne l'hotspot** se è già acceso.
Ci sono tre elementi separati da gestire:

| Elemento | Dove si trova | Come si rimuove |
|---|---|---|
| Attività pianificata `WiRouter` | Utilità di pianificazione | Comando `--uninstall-autostart` oppure GUI |
| Processo `wirouter-autostart.exe` | In esecuzione in background | Riavvio, logout, oppure `Stop-Process` |
| Hotspot Windows | Impostazioni → Hotspot mobile | Interruttore nelle Impostazioni |

La password salvata nell'attività viene eliminata insieme all'attività.

## Metodo 1: comando (consigliato)

Dalla cartella in cui si trova `wirouter.exe`:

```powershell
.\wirouter.exe --uninstall-autostart
```

Output atteso:

```
Attività di avvio automatico rimossa.
```

Il comando non richiede privilegi di amministratore. Se l'attività non esiste,
`schtasks` restituisce un errore: significa che era già rimossa.

## Metodo 2: Utilità di pianificazione (interfaccia grafica)

1. Premere `Win + R`, digitare `taskschd.msc` e premere Invio.
2. Nel pannello sinistro aprire **Libreria Utilità di pianificazione**.
3. Trovare `WiRouter`.
4. Per **sospendere** senza cancellare: clic destro → **Disabilita**.
   Per tornare attivo: clic destro → **Abilita**.
5. Per **rimuovere** definitivamente: clic destro → **Elimina**.

## Metodo 3: riga di comando senza il programma

Per sospendere senza cancellare l'attività:

```powershell
schtasks /change /tn WiRouter /disable
```

Per riattivarla:

```powershell
schtasks /change /tn WiRouter /enable
```

Per eliminarla:

```powershell
schtasks /delete /tn WiRouter /f
```

## Fermare il processo in corso

L'avvio automatico può essere ancora in esecuzione. Rimuovere l'attività non lo
ferma. Per verificare:

```powershell
Get-Process wirouter-autostart -ErrorAction SilentlyContinue
```

Se restituisce un processo, terminarlo con il suo ID:

```powershell
Stop-Process -Id <ID>
```

Attenzione: `Stop-Process` non esegue la procedura di arresto. L'hotspot resta
acceso. Va spento a parte (vedi sotto).

## Spegnere l'hotspot

Dopo aver fermato il processo, spegnere l'hotspot da:

**Impostazioni → Rete e Internet → Hotspot mobile → disattiva**

## Verificare che sia tutto disattivato

```powershell
schtasks /query /tn WiRouter
```

Se l'attività è stata rimossa, il comando risponde con
`ERROR: Impossibile trovare il file specificato` (o messaggio equivalente).

Per controllare che il processo non sia attivo:

```powershell
Get-Process wirouter-autostart -ErrorAction SilentlyContinue
```

Nessun output significa che il processo non è in esecuzione.

## Rimuovere anche i file del programma (facoltativo)

I file vengono copiati in `%LOCALAPPDATA%\WiRouter`. Per eliminarli, dopo aver
rimosso l'attività:

```powershell
Remove-Item -Recurse -Force "$env:LOCALAPPDATA\WiRouter"
```

Il codice sorgente nel repository non viene toccato.

## Riattivare in seguito

Rieseguire `wirouter.exe --install-autostart` dalla cartella `%LOCALAPPDATA%\WiRouter`
(o da una build). Il comando chiede di nuovo connessione, SSID e password.
