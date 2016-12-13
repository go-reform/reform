# Based on http://www.appveyor.com/docs/services-databases

[reflection.assembly]::LoadWithPartialName("Microsoft.SqlServer.Smo") | Out-Null
[reflection.assembly]::LoadWithPartialName("Microsoft.SqlServer.SqlWmiManagement") | Out-Null

$instanceName = $env:REFORM_SQL_INSTANCE
$uri = "ManagedComputer[@Name='$env:COMPUTERNAME']/ServerInstance[@Name='$instanceName']/ServerProtocol[@Name='Tcp']"
$wmi = New-Object ('Microsoft.SqlServer.Management.Smo.Wmi.ManagedComputer')
$tcp = $wmi.GetSmoObject($uri)
$tcp.IsEnabled = $true
$tcp.Alter()

Start-Service "MSSQL`$$instanceName"

Set-Service SQLBrowser -StartupType Manual
Start-Service SQLBrowser
