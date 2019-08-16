# Clojure installer and CLI tools for windows 

this is for installing the clojure  CLI tools on windows 

## Warning


Alex Miller points out that he is working on an official installer for windows and asks that my project clearly differentiate itself from the official installers.

see https://groups.google.com/d/msg/clojure/nw66LtqsaWU/rMv1GWEfDAAJ)

I am doing this:

THIS IS NOT THE OFFICIAL  installer for clojure and cli tools on windows

see cli https://clojure.org/guides/getting_started for the official installers
and https://dev.clojure.org/jira/browse/TDEPS-67 for progress  on the windows part






## Installation
1. Download the latest release from https://github.com/frericksm/clj-windows/releases
2. Extract it to some `<local-path>`
3. Execute `<local-path>`\windows-clojure-tools-1.10.1.469\install.exe [install_dir]
with the optional install_dir

### How install.exe works:
install.exe does the following things:

1. creates and fills the folder [install_dir] if set or %localappdata%/Programs\clojure 
	
     

    where  %localappdata% is the expansion of the environment variable LOCALAPPDATA
2. It adds the path %localappdata%/Programs\clojure\bin or install_dir\bin to the environment variable PATH in scope USER
## Deinstallation 
1. Delete the folder %localappdata%/Programs\clojure or [install_dir]
2. Remove %localappdata%/Programs\clojure\bin from  the environment variable PATH in scope USER



## Build from source
### Prerequisites
1. Install https://golang.org/
2. install maven https://maven.apache.org/

### Steps 
1. Checkout https://github.com/frericksm/clj-windows.git
2. change directory to clj-windows/clojure
3. run "GOOS=windows GOARCH=amd64 go build"
4. change directory to clj-windows/clj
5. run "GOOS=windows GOARCH=amd64 go build"
6. change directory to /clj-windows/install
7. run "GOOS=windows GOARCH=amd64 go build"
8. change directory to /clj-windows
9. run "mvn package"
10. see clj-windows/target for build results  an follow the installation task above

# Proxy settings


## Proxy settings while classpath calculation

Until https://clojure.atlassian.net/browse/TDEPS-124 is fixed in an adaptable way,
i fixed it here as described in https://github.com/frericksm/clj-windows/issues/6 

the the env vars http_proxy and  https_proxy will be used  to set the appropriate system properties to the 
java call of clojure.tools.deps.alpha.script.make-classpath

# Command line arguments

 There is a subtle difference in  parsing command line arguments between unix and windows 

while following command runs on unix 


```

clj -Sdeps '{:deps {cider/cider-nrepl {:mvn/version "0.20.0"} }}' -e '(require (quote cider-nrepl.main)) (cider-nrepl.main/init ["cider.nrepl/cider-middleware"])'
```
it will not run on windows.
You have to use double quotes instead of single quotes  an additionally you have to escape double quotes in strings with a backslash. so the modified command 

```
clj -Sdeps "{:deps {cider/cider-nrepl {:mvn/version \"0.20.0\"} }}" -e "(require (quote cider-nrepl.main)) (cider-nrepl.main/init [\"cider.nrepl/cider-middleware\"])"
```
runs on windows and *surprise*  it also runs on unix
