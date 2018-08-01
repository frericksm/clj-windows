# define parameters

param (
 [string]
    [alias("i")]
    [string] $install_dir= $env:localappdata +"\Programs\clojure" 

   
    )
	
	
	#https://www.sapien.com/powershell/abouthelp/environment_variables/

	
	 $path = [System.Environment]::GetEnvironmentVariable("Path", "User") 
            [System.Environment]::SetEnvironmentVariable("Path", $path + ";" + $install_dir+ "\bin", "User") 