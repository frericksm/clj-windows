param (
 [string]
    [alias("p")]
    [string] $prefix_dir= $env:localappdata +"\Programs" 

   
    )
# Start
function do-usage {
echo "Installs the Clojure command line tools."
echo "Usage:"
  echo "win-install-1.9.0.381.ps1[-p|--prefix <dir>]"
  exit 1
}

echo "Downloading and expanding tar"
$url = "https://download.clojure.org/install/clojure-tools-1.9.0.381.tar.gz"
$file1= "clojure-tools-1.9.0.381.tar.gz"
$file2= "clojure-tools-1.9.0.381.tar"
$output = "$PSScriptRoot\clojure-tools-1.9.0.381.tar.gz"

$install_dir=$prefix_dir + "\clojure"
$lib_dir=$install_dir + "/lib"
$lib_exec=$lib_dir+ "/libexec"
$bin_dir="$prefix_dir/bin"
$man_dir="$prefix_dir/share/man/man1"
$clojure_lib_dir="$lib_dir/clojure"


Remove-Item $file1
Remove-Item $file2
Remove-Item .\clojure-tools -Force -Recurse
Remove-Item $install_dir -Force -Recurse
echo downloading $url  ...
$start_time = Get-Date

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($url, $output)

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"



$env:Path = $env:Path + ";C:\Program Files\7-zip"

7z.exe  x "clojure-tools-1.9.0.381.tar.gz"
7z.exe  x "clojure-tools-1.9.0.381.tar"


New-Item -Path $install_dir -ItemType "directory"
New-Item -Path $lib_dir  -ItemType "directory"
New-Item -Path $lib_exec   -ItemType "directory"

echo "Installing libs into $clojure_lib_dir"
Copy-Item clojure-tools/deps.edn $lib_dir
Copy-Item clojure-tools/example-deps.edn $lib_dir
Copy-Item clojure-tools/clojure-tools-1.9.0.381.jar $lib_exec

echo "Installing clojure and clj into $bin_dir"
sed -i -e 's@PREFIX@'"$clojure_lib_dir"'@g' clojure-tools/clojure
install -Dm755 clojure-tools/clojure "$bin_dir/clojure"
install -Dm755 clojure-tools/clj "$bin_dir/clj"

echo "Installing man pages into $man_dir"
install -Dm644 clojure-tools/clojure.1 "$man_dir/clojure.1"
install -Dm644 clojure-tools/clj.1 "$man_dir/clj.1"

echo "Removing download"
rm -rf clojure-tools
rm -rf clojure-tools-1.9.0.381.tar.gz

echo "Use clj -h for help."
