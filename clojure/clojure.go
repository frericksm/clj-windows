// clj-win project main.go
package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"log"
	"net/url"
	"os/exec"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// returns env var http_proxy an https_proxy as java options -Dhttp.proxyHost...
func proxyargs() []string {

	proxy_args := []string{}

	http_proxy_url, _ := url.Parse(os.Getenv("http_proxy"))
	https_proxy_url, _ := url.Parse(os.Getenv("https_proxy"))
	//check env var http_proxy and transform to java option
	if http_proxy_url != nil {
		proxy_args = append(proxy_args, "-Dhttp.proxyHost="+http_proxy_url.Hostname())
		http_proxy_port := http_proxy_url.Port()
		if http_proxy_port != "" {
			proxy_args = append(proxy_args, "-Dhttp.proxyPort="+http_proxy_port)
		}
	}

	//check env var https_proxy and transform to java option
	if https_proxy_url != nil {
		proxy_args = append(proxy_args, "-Dhttps.proxyHost="+https_proxy_url.Hostname())
		https_proxy_port := https_proxy_url.Port()
		if https_proxy_port != "" {
			proxy_args = append(proxy_args, "-Dhttps.proxyPort="+https_proxy_port)
		}
	}

	//	fmt.Println("proxy_args")
	//	fmt.Println(proxy_args)

	return proxy_args
}

// copies file at path src to path dest
func copy(src string, dest string) {
	from, err := os.Open(src)
	check(err)

	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer to.Close()

	_, err = io.Copy(to, from)
	check(err)

}

func usage(version string) {
	fmt.Printf("Version: %s\n", version)	
	fmt.Println(`You use the Clojure tools ('clj' or 'clojure') to run Clojure programs
	on the JVM, e.g. to start a REPL or invoke a specific function with data.
	The Clojure tools will configure the JVM process by defining a classpath
	(of desired libraries), an execution environment (JVM options) and
	specifying a main class and args. 

	Using a deps.edn file (or files), you tell Clojure where your source code
	resides and what libraries you need. Clojure will then calculate the full
	set of required libraries and a classpath, caching expensive parts of this
	process for better performance.

	The internal steps of the Clojure tools, as well as the Clojure functions
	you intend to run, are parameterized by data structures, often maps. Shell
	command lines are not optimized for passing nested data, so instead you
	will put the data structures in your deps.edn file and refer to them on the
	command line via 'aliases' - keywords that name data structures.

	'clj' and 'clojure' differ in that 'clj' has extra support for use as a REPL
	in a terminal, and should be preferred unless you don't want that support,
	then use 'clojure'.

	Usage:
	  Start a REPL   clj     [clj-opt*] [-A:aliases] [init-opt*]
	  Exec function  clojure [clj-opt*] -X[:aliases] [a/fn] [kpath v]*
	  Run main       clojure [clj-opt*] -M[:aliases] [init-opt*] [main-opt] [arg*]
	  Prepare        clojure [clj-opt*] -P [other exec opts]

	exec-opts:
	 -A:aliases     Use aliases to modify classpath
	 -X[:aliases]   Use aliases to modify classpath or supply exec fn/args
	 -M[:aliases]   Use aliases to modify classpath or supply main opts
	 -P             Prepare deps - download libs, cache classpath, but don't exec

	clj-opts:
	 -Jopt          Pass opt through in java_opts, ex: -J-Xmx512m
	 -Sdeps EDN     Deps data to use as the last deps file to be merged
	 -Spath         Compute classpath and echo to stdout only
	 -Spom          Generate (or update) pom.xml with deps and paths
	 -Stree         Print dependency tree
	 -Scp CP        Do NOT compute or cache classpath, use this one instead
	 -Srepro        Ignore the ~/.clojure/deps.edn config file
	 -Sforce        Force recomputation of the classpath (don't use the cache)
	 -Sverbose      Print important path info to console
	 -Sdescribe     Print environment and command parsing info as data
	 -Sthreads      Set specific number of download threads
	 -Strace        Write a trace.edn file that traces deps expansion
	 --             Stop parsing dep options and pass remaining arguments to clojure.main

	init-opt:
	 -i, --init path     Load a file or resource
	 -e, --eval string   Eval exprs in string; print non-nil values
	 --report target     Report uncaught exception to "file" (default), "stderr", or "none"

	main-opt:
	 -m, --main ns-name  Call the -main function from namespace w/args
	 -r, --repl          Run a repl
	 path                Run a script from a file or resource
	 -                   Run a script from standard input
	 -h, -?, --help      Print this help message and exit

	Programs provided by :deps alias:
	 -X:deps mvn-install       Install a maven jar to the local repository cache
	 -X:deps git-resolve-tags  Resolve git coord tags to shas and update deps.edn

	For more info, see:
	 https://clojure.org/guides/deps_and_cli
	 https://clojure.org/reference/repl_and_main
`)
	
}

func main() {

	
	version := "1.10.1.763"
	exec_path, _ := os.Executable()
	wd := filepath.Dir(exec_path)
	local_install_dir := filepath.Dir(wd)
	install_dir := os.Getenv("localappdata") + "/Programs/clojure"
	
	jarfile := fmt.Sprintf("/lib/libexec/clojure-tools-%s.jar", version)
	local_jar := local_install_dir + jarfile
	
	if jarfile_exists, _ := exists(local_jar); jarfile_exists {
		install_dir = local_install_dir
	}
	tools_cp := install_dir + jarfile
	print_classpath := false
	describe := false
	threads:= ""
	verbose := false
	trace := false
	force := false
	repro := false
	tree := false
	pom := false
	help := false
	prep :=false
	stale := false

	jvm_opts := make([]string, 0)
	resolve_aliases := make([]string, 0)
	classpath_aliases := make([]string, 0)
	main_aliases := make([]string, 0)
	exec_aliases := make([]string, 0)
	repl_aliases := make([]string, 0)
	mode := "repl"
	var deps_data string
	var force_cp string
	var deps_edn_file_exists bool
	var cache_dir string
	additional_args := make([]string, 0)

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case strings.HasPrefix(arg, "-J"):
			jvm_opts = append(jvm_opts, arg[2:])
		case strings.HasPrefix(arg, "-R"):
			resolve_aliases = append(resolve_aliases, arg[2:])
			fmt.Println("-R is deprecated, use -A with repl, -M for main, or -X for exec")
		case strings.HasPrefix(arg, "-C"):
			classpath_aliases = append(classpath_aliases, arg[2:])
			fmt.Println("-C is deprecated, use -A with repl, -M for main, or -X for exec")
		case strings.HasPrefix(arg, "-O"):
			fmt.Println("-O is no longer supported, use -A with repl, -M for main, or -X for exec")
			break
		case strings.HasPrefix(arg, "-A"):
			repl_aliases = append(repl_aliases, arg[2:])
		case arg == "-M":
			mode = "main"
			break
		case strings.HasPrefix(arg, "-M"):
			mode = "main"
			main_aliases = append(main_aliases, arg[2:])
			break
		case arg == "-X":
			mode = "exec"
			break
		case strings.HasPrefix(arg, "-X"):
			mode = "exec"
			main_aliases = append(exec_aliases, arg[2:])
			break
		case arg == "-P":
			prep = true
		case strings.HasPrefix(arg, "-Sdeps"):
			i = i + 1
			deps_data = os.Args[i]
		case strings.HasPrefix(arg, "-Scp"):
			i = i + 1
			force_cp = os.Args[i]
		case strings.HasPrefix(arg, "-Spath"):
			print_classpath = true
		case strings.HasPrefix(arg, "-Sverbose"):
			verbose = true
		case strings.HasPrefix(arg, "-Sdescribe"):
			describe = true
		case strings.HasPrefix(arg, "-Sthreads"):
			i = i + 1
			threads = os.Args[i]
		case strings.HasPrefix(arg, "-Strace"):
			trace = true
		case strings.HasPrefix(arg, "-Sforce"):
			force = true
		case strings.HasPrefix(arg, "-Srepro"):
			repro = true
		case strings.HasPrefix(arg, "-Stree"):
			tree = true
		case strings.HasPrefix(arg, "-Spom"):
			pom = true
		case strings.HasPrefix(arg, "-Sresolve-tags"):
			fmt.Println("Option changed, use: clj -X:deps git-resolve-tags")
			break
		case strings.HasPrefix(arg, "-S"):
			fmt.Println("Invalid option %s", arg)
		case (strings.HasPrefix(arg, "-h") || strings.HasPrefix(arg, "--help") || strings.HasPrefix(arg, "-?")):
			if len(main_aliases) == 0 && len(repl_aliases) == 0 {
				help = true
			}
                case strings.HasPrefix(arg, "--"):
			break
		default:
			additional_args = append(additional_args, arg)
		}

	}
	//	fmt.Printf("additional args: %s", additional_args)

	if help {
		usage(version)
		return
	}
	config_paths := make([]string, 0)
	config_user := ""
	config_project := ""

	config_dir := os.Getenv("CLJ_CONFIG")
	config_home := os.Getenv("XDG_CONFIG_HOME")
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	var cmd *exec.Cmd
	var cmd_args []string

	deps_edn_file_exists = true
	if _, err := os.Stat("deps.edn"); os.IsNotExist(err) {
		deps_edn_file_exists = false
	}
	
	// Determine user config directory

	if config_dir == "" && config_home != "" {
		config_dir = config_home + "/clojure"
	} else {
		config_dir = home + "/.clojure"
	}

	if _, err := os.Stat(config_dir); os.IsNotExist(err) {

		// If user config directory does not exist, create it

		os.Mkdir(config_dir, os.ModePerm)
		copy(install_dir+"/lib/example-deps.edn", config_dir+"/deps.edn")
	}

	// Determine user cache directory
	user_cache_dir := os.Getenv("CLJ_CACHE")
	cache_home := os.Getenv("XDG_CACHE_HOME")

	if user_cache_dir == "" && cache_home != "" {
		user_cache_dir = cache_home + "/clojure"
	} else {
		user_cache_dir = config_dir + "/.cpcache"
	}

	// Chain deps.edn in config paths. repro=skip config dir
	config_project="deps.edn"
	if repro {
		config_paths = append(config_paths, install_dir+"/lib/deps.edn")
		config_paths = append(config_paths, "deps.edn")
	} else {
		config_user=config_dir + "/deps.edn"
		config_paths = append(config_paths, install_dir+"/lib/deps.edn")
		config_paths = append(config_paths, config_dir+"/deps.edn")
		config_paths = append(config_paths, "deps.edn")
	}
	// Determine whether to use user or project cache

	if deps_edn_file_exists {
		cache_dir = ".cpcache"
	} else {
		cache_dir = user_cache_dir
	}

	// Construct location of cached classpath file

	joined_result := make([]string, 0)
	joined_result = append(joined_result, resolve_aliases...)
	joined_result = append(joined_result, classpath_aliases...)
	joined_result = append(joined_result, repl_aliases...)
	joined_result = append(joined_result, exec_aliases...)
	joined_result = append(joined_result, main_aliases...)
	joined_result = append(joined_result, deps_data)

	val := strings.Join(joined_result, "")

	//	fmt.Printf("joined_result= %s", joined_result)

	for _, cpf := range config_paths {
		if _, err := os.Stat(cpf); os.IsNotExist(err) {
			val = val + "|NIL"
		} else {
			val = val + "|" + cpf
		}

	}
	//	fmt.Printf("final val= %s\n", val)

	crc32q := crc32.MakeTable(0xD5828281)
	ck := fmt.Sprintf("%08x", crc32.Checksum([]byte(val), crc32q))

	//	fmt.Printf("ck= %s\n", ck)
	libs_file := filepath.Join(cache_dir, ck + ".libs")
	cp_file := filepath.Join(cache_dir, ck + ".cp")
	jvm_file := filepath.Join(cache_dir, ck + ".jvm")
	main_file := filepath.Join(cache_dir, ck + ".main")
	basis_file := filepath.Join(cache_dir, ck + ".basis")
	// Print paths in verbose mode
	if verbose {
		fmt.Printf("version      = %s\n", version)
		fmt.Printf("install_dir  =%s\n", install_dir)
		fmt.Printf("config_dir   =%s\n", config_dir)
		fmt.Printf("config_paths =%s\n", strings.Join(config_paths, " "))
		fmt.Printf("cache_dir    = %s\n", cache_dir)
		fmt.Printf("cp_file      =%s\n", cp_file)

	}
	// Check for stale classpath file
	stale = false
	cp_file_exists := false
	if _, err := os.Stat(cp_file); err == nil {
		cp_file_exists = true
	}
	if force || trace || tree ||prep || !cp_file_exists {
		stale = true
	} else {
		cp_file_info, err := os.Stat(cp_file)
		if err != nil {
			check(err)
		}

		for _, cpf := range config_paths {

			cpf_info, err := os.Stat(cpf)
			if err != nil {

				stale = true
				break
			}
			if cpf_info.ModTime().After(cp_file_info.ModTime()) {
				stale = true
				break
			}
		}
	}
	var tools_args []string

	// Make tools args if needed
	if stale || pom {
		tools_args = make([]string, 0)

		if deps_data != "" {
			tools_args = append(tools_args, "--config-data")
			tools_args = append(tools_args, deps_data)
		}
		if len(resolve_aliases) > 0 {
			tools_args = append(tools_args, "-R"+strings.Join(resolve_aliases, ""))
		}
		if len(classpath_aliases) > 0 {
			tools_args = append(tools_args, "-C"+strings.Join(classpath_aliases, ""))
		}
		if len(main_aliases) > 0 {
			tools_args = append(tools_args, "-M"+strings.Join(main_aliases, ""))
		}
		if len(repl_aliases) > 0 {
			tools_args = append(tools_args, "-A"+strings.Join(repl_aliases, ""))
		}
		if len(exec_aliases) > 0 {
			tools_args = append(tools_args, "-X"+strings.Join(exec_aliases, ""))
		}
		if force_cp != "" {
			tools_args = append(tools_args, "--skip-cp")

		}
		if trace {
			tools_args = append(tools_args, "--trace")

		}
		if tree {
			tools_args = append(tools_args, "--tree")

		}
		if threads!= "" {
			tools_args = append(tools_args, "--threads", threads)
		}
		
	}

	// If stale, run make-classpath to refresh cached classpath
	if stale && verbose {
		fmt.Println("Refreshing classpath")

	}
	if stale {
		
		make_classpath_args := []string{ "-cp", tools_cp, "clojure.main", "-m", "clojure.tools.deps.alpha.script.make-classpath2", "--config-user", config_user,"--config-project", config_project,"--basis-file", basis_file,"--libs-file",libs_file,"--cp-file", cp_file,"--jvm-file", jvm_file,"--main-file", main_file}
		
		cmd_args = proxyargs()
		cmd_args = append(cmd_args, make_classpath_args...)
		cmd_args = append(cmd_args, tools_args...)
		cmd = exec.Command("java.exe", cmd_args...)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		//		fmt.Println("make-classpath")
		//		fmt.Println(cmd_args)
		cmd.Run()
	}

	var cp string
	var content []byte
	var err error
	if force_cp != "" {
		cp = force_cp
	} else {

		content, err = ioutil.ReadFile(cp_file)
		check(err)
		cp = string(content)
	}
	if prep {
		os.Exit(0)
	} 
	
	if pom {
		cmd_args := []string{ "-cp", tools_cp, "clojure.main", "-m", "clojure.tools.deps.alpha.script.generate-manifest2","--config-user", config_user,"--config-project", config_project, "--gen=pom"}
		cmd_args = append(cmd_args, tools_args...)
		cmd = exec.Command("java.exe", cmd_args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		//		fmt.Println("make-classpath")
		//		fmt.Println(cmd_args)
		cmd.Run()
	} else if print_classpath {
		fmt.Println(cp)
	} else if describe {
		path_vector := ""
		for _, cpf := range config_paths {
			//			fmt.Println("cpf" + cpf)
			if _, err := os.Stat(cpf); err == nil {
				path_vector = path_vector + `"` + cpf + `"` + " "
			}

		}

		fmt.Printf("\n{:version \"%s\"\n"+
			":config-files [%s]\n"+
			":config-user [%s]\n"+
			":config-project [%s]\n"+
			":install-dir \"%s\"\n"+
			":config-dir \"%s\"\n"+
			":cache-dir \"%s\"\n"+
			":force %t\n"+
			":repro %t\n"+
			":main-aliases \"%s\"\n"+
			":repl-aliases \"%s\"\n}\n", version, path_vector,config_user, config_project,
			install_dir, config_dir, cache_dir, force, repro, strings.Join(main_aliases, ":"), strings.Join(repl_aliases, ""))

	} else if tree {
		return
	} else if trace {
		fmt.Println("Wrote trace.edn")
	} else {
		var jvm_cache_opts string
		var main_cache_opts string
		
		if _, err := os.Stat(jvm_file); err == nil {
			content, err = ioutil.ReadFile(jvm_file)
			jvm_cache_opts = string(content)
			check(err)
		}
		
		
		cmd_args = make([]string, 0)
		if mode == "exec" {
			var exec_args = make([]string, 0)
			
			if len(exec_aliases) > 0 {
				exec_args= append(exec_args, "--aliases")
				exec_args= append(exec_args, exec_aliases...)
				
			} else {
				jvm_opts_string := strings.Join(jvm_opts, " ")
				
				if jvm_cache_opts != "" {
					cmd_args = append(cmd_args, jvm_cache_opts)
					
				}
				
				if len(jvm_opts) > 0 {
					cmd_args = append(cmd_args, jvm_opts_string)
					
				}
				
				
			}
			cmd_args = append(cmd_args, fmt.Sprintf("-Dclojure.basis=%s", basis_file), "-classpath", cp + string(os.PathListSeparator)+ install_dir+ "/libexec/exec.jar", "clojure.main", "-m", "clojure.run.exec")
			cmd_args = append(cmd_args, exec_args...)
			cmd = exec.Command("java.exe", cmd_args...)
			
		} else { // else mode == "exec"
			if _, err := os.Stat(main_file); err == nil {
				content, err = ioutil.ReadFile(main_file)
				main_cache_opts = string(content)
				check(err)
			}
			if len(additional_args) > 0 && mode== "repl" {
				fmt.Println("WARNING: When invoking clojure.main, use -M")
			}
			
			cmd_args = append(cmd_args, fmt.Sprintf("-Dclojure.basis=%s", basis_file), "-classpath",  cp, "clojure.main", main_cache_opts)
			cmd_args = append(cmd_args, additional_args...)
			
			cmd = exec.Command("java.exe", cmd_args...)
			
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Run()		
		}
	}
}
