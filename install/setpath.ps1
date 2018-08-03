# add the bin_dir to the env var Path

#https://www.sapien.com/powershell/abouthelp/environment_variables/

	
$install_dir= $env:localappdata +"\Programs\clojure"
$bin_dir = $install_dir+ "\bin"
	 $path = [System.Environment]::GetEnvironmentVariable("Path", "User") 
	
	if ($path -split ';' -contains $bin_dir){
		}
		else {
           [System.Environment]::SetEnvironmentVariable("Path", $path + ";" + $bin_dir, "User") 
			Echo "added "+ $bin_dir +" to env-var Path" 
	}
