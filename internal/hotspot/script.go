package hotspot

// Gli input utente passano solo tramite variabili d'ambiente (WIROUTER_*),
// mai concatenati nello script, per evitare injection.
const scriptPrelude = `
$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
Add-Type -AssemblyName System.Runtime.WindowsRuntime
[void][Windows.Networking.Connectivity.NetworkInformation, Windows.Networking.Connectivity, ContentType = WindowsRuntime]
[void][Windows.Networking.NetworkOperators.NetworkOperatorTetheringManager, Windows.Networking.NetworkOperators, ContentType = WindowsRuntime]
[void][Windows.Networking.NetworkOperators.NetworkOperatorTetheringAccessPointConfiguration, Windows.Networking.NetworkOperators, ContentType = WindowsRuntime]

$asTask = [System.WindowsRuntimeSystemExtensions].GetMethods() |
    Where-Object { $_.Name -eq 'AsTask' -and $_.GetParameters().Count -eq 1 }

function Wait-Operation($operation, $resultType) {
    $method = $asTask | Where-Object { $_.GetParameters()[0].ParameterType.Name -like 'IAsyncOperation*' } | Select-Object -First 1
    $task = $method.MakeGenericMethod($resultType).Invoke($null, @($operation))
    $task.Wait(-1) | Out-Null
    $task.Result
}

function Wait-Action($action) {
    $method = $asTask | Where-Object { $_.GetParameters()[0].ParameterType.Name -eq 'IAsyncAction' } | Select-Object -First 1
    $task = $method.Invoke($null, @($action))
    $task.Wait(-1) | Out-Null
}

function Get-Source {
    $name = $env:WIROUTER_SOURCE
    if ($name) {
        $profile = [Windows.Networking.Connectivity.NetworkInformation]::GetConnectionProfiles() |
            Where-Object { $_.ProfileName -eq $name } | Select-Object -First 1
    } else {
        $profile = [Windows.Networking.Connectivity.NetworkInformation]::GetInternetConnectionProfile()
    }
    if (-not $profile) { throw 'Nessuna connessione Internet disponibile da condividere' }
    $profile
}

function Get-Manager {
    $manager = [Windows.Networking.NetworkOperators.NetworkOperatorTetheringManager]::CreateFromConnectionProfile((Get-Source))
    $manager
}

function Assert-Success($result) {
    if ($result.Status -ne 'Success') {
        throw ('{0} {1}' -f $result.Status, $result.AdditionalErrorMessage).Trim()
    }
}
`

const scriptSources = `
[Windows.Networking.Connectivity.NetworkInformation]::GetConnectionProfiles() |
    Where-Object { $_.GetNetworkConnectivityLevel() -eq 'InternetAccess' } |
    ForEach-Object { $_.ProfileName }
`

const scriptStart = `
$manager = Get-Manager
if ($manager.TetheringOperationalState -eq 'On') { throw 'Hotspot già attivo' }
$configuration = New-Object Windows.Networking.NetworkOperators.NetworkOperatorTetheringAccessPointConfiguration
$configuration.Ssid = $env:WIROUTER_SSID
$configuration.Passphrase = $env:WIROUTER_PASSWORD
Wait-Action ($manager.ConfigureAccessPointAsync($configuration))
Assert-Success (Wait-Operation ($manager.StartTetheringAsync()) ([Windows.Networking.NetworkOperators.NetworkOperatorTetheringOperationResult]))
`

const scriptStop = `
$manager = Get-Manager
Assert-Success (Wait-Operation ($manager.StopTetheringAsync()) ([Windows.Networking.NetworkOperators.NetworkOperatorTetheringOperationResult]))
`

const scriptStatus = `
$manager = Get-Manager
'state=' + $manager.TetheringOperationalState
'ssid=' + $manager.GetCurrentAccessPointConfiguration().Ssid
'clients=' + $manager.ClientCount
`

func buildScript(body string) string {
	return scriptPrelude + "\ntry {\n" + body + "\n} catch {\n[Console]::Error.WriteLine($_.Exception.Message)\nexit 1\n}\n"
}
