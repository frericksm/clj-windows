
$url = "https://download.clojure.org/install/clojure-tools-1.9.0.381.tar.gz"
$output = "$PSScriptRoot\clojure-tools-1.9.0.381.tar.gz"

echo downloading $url  ...
$start_time = Get-Date

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($url, $output)
#OR
(New-Object System.Net.WebClient).DownloadFile($url, $output)

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"



$env:Path = $env:Path + ";C:\Program Files\7-zip"

7z.exe  x "clojure-tools-1.9.0.381.tar.gz"
7z.exe  x "clojure-tools-1.9.0.381.tar"$lib_dir=$install_dir + "/lib"
