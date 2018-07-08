# define parameters

param (
 [string]
    [alias("p")]
    [string] $prefix_dir= $env:localappdata +"\Programs" 

   
    )
# Start

# define vars for download and extracting
$url = "https://download.clojure.org/install/clojure-tools-1.9.0.381.tar.gz"
$clj_tools_archive_file_tar_gz = "$PSScriptRoot\clojure-tools-1.9.0.381.tar.gz"
$clj_tools_archive_file_tar= "clojure-tools-1.9.0.381.tar"

$install_dir=$prefix_dir + "\clojure"
$lib_dir=$install_dir + "/lib"
$lib_exec=$lib_dir+ "/libexec"
$bin_dir="$prefix_dir/bin"
$man_dir="$prefix_dir/share/man/man1"
$clojure_lib_dir="$lib_dir/clojure"


function do-usage {
echo "Installs the Clojure command line tools."
echo "Usage:"
  echo "win-install-1.9.0.381.ps1[-p|--prefix <dir>]"
  exit 1
}
function cleanup {
# remove  old download files and temp dirs
Remove-Item $clj_tools_archive_file_tar_gz
Remove-Item $clj_tools_archive_file_tar
Remove-Item .\clojure-tools -Force -Recurse

}

echo downloading $url  ...
$start_time = Get-Date

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($url, $clj_tools_archive_file_tar_gz)

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"



$env:Path = $env:Path + ";C:\Program Files\7-zip"

7z.exe  x $clj_tools_archive_file_tar_gz

7z.exe  x "clojure-tools-1.9.0.381.tar"


# cleanup install dir
Get-Item  $install_dir  | Remove-Item  -Force -Recurse

# prepare install dir 
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


cleanup
#echo "Removing download"
#Remove-Item $clj_tools_archive_file_tar_gz


echo "Use clj -h for help."
