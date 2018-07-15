
rem Arguments handling in a batch file
rem http://www.rgagnon.com/gp/gp-0009.html
rem https://ss64.com/nt/syntax.html

rem https://www.script-example.com/themen/cmd_Batch_Befehle.php
rem http://www.robvanderwoude.com/battech_inputvalidation_commandline.php

rem Version = 1.9.0.381

rem replace PREFIX at install time by current install_dir 
set install_dir= PREFIX
set install_dir= c:\3333

rem Extract opts
set print_classpath= false
set describe= false
set verbose= false
set force= false
set repro= false
set tree= false
set pom= false
set resolve_tags= false
set help= false
jvm_opts=()
resolve_aliases=()
classpath_aliases=()
jvm_aliases=()
main_aliases=()
all_aliases=()


@ECHO OFF
:Loop
IF "%1"=="" GOTO Continue
   •
   ECHO "%~1"| FINDSTR /L /X /I """-J"""

   IF ERRORLEVEL   0 GOTO Loop
   set jvm_opts=%jvm_opts%%1
SHIFT
GOTO Loop
:Continue

set tools_cp=%install_dir%\clojure\lib\libexec\clojure-tools-1.9.0.381.jar
  java -classpath %tools_cp% clojure.main 
