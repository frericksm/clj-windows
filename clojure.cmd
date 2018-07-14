
rem Arguments handling in a batch file
rem http://www.rgagnon.com/gp/gp-0009.html
rem https://ss64.com/nt/syntax.html
rem https://www.script-example.com/themen/cmd_Batch_Befehle.php

rem Version = 1.9.0.381

rem replace PREFIX at install time by current install_dir 
set install_dir= PREFIX
set install_dir= c:\3333
set tools_cp=%install_dir%\clojure\lib\libexec\clojure-tools-1.9.0.381.jar
  java -classpath %tools_cp% clojure.main 
