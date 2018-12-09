# Clojure installer and CLI tools for windows 

this is for installing the clojure  CLI tools on windows 

## Warning


Alex Miller points out that he is working on an official installer for windows and asks that my project clearly differentiate itself from the official installers.

see https://groups.google.com/d/msg/clojure/nw66LtqsaWU/rMv1GWEfDAAJ)

I am doing this:

THIS IS NOT THE OFFICIAL  installer for clojure and cli tools on windows

see cli https://clojure.org/guides/getting_started for the official installers






## Installation
1. Download the latest release from https://github.com/frericksm/clj-windows/releases
2. Extract it to some `<local-path>`
3. Execute `<local-path>`\windows-clojure-tools-1.9.0.397\install.exe

### How install.exe works:
install.exe does the following things:

1. It downloads and extracts  https://download.clojure.org/install/clojure-tools-VERSION.tar.gz
where  VERSION is replaced by the current version  1.9.0.397 (at the time of writing)  

2. creates and fills the folder
    %localappdata%/Programs\clojure

    where  %localappdata% is the expansion of the environment variable LOCALAPPDATA
3. It adds the path %localappdata%/Programs\clojure\bin to the environment variable PATH in scope USER
## Deinstallation 
1. Delete the folder %localappdata%/Programs\clojure
2. Remove %localappdata%/Programs\clojure\bin from  the environment variable PATH in scope USER



## Build from source
### Prerequisites
1. Install https://golang.org/
2. install maven https://maven.apache.org/

### Steps 
1. Checkout https://github.com/frericksm/clj-windows.git
2. change directory to clj-windows/clojure
3. run "go build "
4. change directory to clj-windows/clj
5. run "go build "
6. change directory to /clj-windows/install
7. run "go build"
8. change directory to /clj-windows
9. run "mvn package"
10. see clj-windows/target for build results  an follow the installation task above


# Command line arguments

 There is a subtle difference in  parsing command line arguments between unix and windows 

while following command runs on unix 


```

clj -Sdeps '{:deps {cider/cider-nrepl {:mvn/version "0.18.0"} }}' -e '(require (quote cider-nrepl.main)) (cider-nrepl.main/init ["cider.nrepl/cider-middleware"])'
```
it will not run on windows.
You have to use double quotes instead of single quotes  an additionally you have to escape double quotes in strings with a backslash. so the modified command 

```
clj -Sdeps "{:deps {cider/cider-nrepl {:mvn/version \"0.18.0\"} }}" -e "(require (quote cider-nrepl.main)) (cider-nrepl.main/init [\"cider.nrepl/cider-middleware\"])"
```
runs on windows and *surprise*  it also runs on unix
