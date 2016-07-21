# Based on http://www.appveyor.com/docs/services-databases

[reflection.assembly]::LoadWithPartialName("Microsoft.SqlServer.Smo") | Out-Null
[reflection.assembly]::LoadWithPartialName("Microsoft.SqlServer.SqlWmiManagement") | Out-Null

$instanceNames = "SQL2008R2SP2", "SQL2012SP1", "SQL2014", "SQL2016"
$wmi = New-Object ('Microsoft.SqlServer.Management.Smo.Wmi.ManagedComputer')

Foreach ($instanceName in $instanceNames) {
    $uri = "ManagedComputer[@Name='$env:COMPUTERNAME']/ServerInstance[@Name='$instanceName']/ServerProtocol[@Name='Tcp']"
    $tcp = $wmi.GetSmoObject($uri)
    $tcp.IsEnabled = $true
    $tcp.Alter()
}

Set-Service SQLBrowser -StartupType Manual
Start-Service SQLBrowser
