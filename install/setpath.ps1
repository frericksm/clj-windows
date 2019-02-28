# add the bin_dir to the env var Path


param
(
  $install_dir
)




#https://www.sapien.com/powershell/abouthelp/environment_variables/


if (!$install_dir) {
$install_dir= $env:localappdata +"\Programs\clojure"
}
Echo "install_dir is  " + $install_dir 


$bin_dir = $install_dir+ "\bin"
	 $path = [System.Environment]::GetEnvironmentVariable("Path", "User") 
	
	if ($path -split ';' -contains $bin_dir){
		}
		else {
           [System.Environment]::SetEnvironmentVariable("Path", $path + ";" + $bin_dir, "User") 
			Echo "added "+ $bin_dir +" to env-var Path" 
	}
