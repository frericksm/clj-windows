# define parameters

param (
 [string]
    [alias("i")]
    [string] $install_dir= $env:localappdata +"\Programs\clojure" 

   
    )
	
	
	#https://www.sapien.com/powershell/abouthelp/environment_variables/

	
	$bin_dir = $install_dir+ "\bin"
	 $path = [System.Environment]::GetEnvironmentVariable("Path", "User") 
	
	if ($path -split ';' -contains $bin_dir){
		}
		else {
			Write-Output "added "+ $bin_dir +" to env-var Path" 
           [System.Environment]::SetEnvironmentVariable("Path", $path + ";" + $bin_dir, "User") 
	}
