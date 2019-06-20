// win-clj project main.go
//https://stackoverflow.com/questions/13913468/how-to-start-a-process
package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	//"path/filepath"

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

func usage() {
	fmt.Println(` 	Usage: clojure [dep-opt*] [init-opt*] [main-opt] [arg*]
	       clj     [dep-opt*] [init-opt*] [main-opt] [arg*]

	The clojure script is a runner for Clojure. clj is a wrapper
	for interactive repl use. These scripts ultimately construct and
	invoke a command-line of the form:

	java [java-opt*] -cp classpath clojure.main [init-opt*] [main-opt] [arg*]

	The dep-opts are used to build the java-opts and classpath:
	 -Jopt          Pass opt through in java_opts, ex: -J-Xmx512m
	 -Oalias...     Concatenated jvm option aliases, ex: -O:mem
	 -Ralias...     Concatenated resolve-deps aliases, ex: -R:bench:1.9
	 -Calias...     Concatenated make-classpath aliases, ex: -C:dev
	 -Malias...     Concatenated main option aliases, ex: -M:test
	 -Aalias...     Concatenated aliases of any kind, ex: -A:dev:mem
 	 -Sdeps EDN     Deps data to use as the last deps file to be merged
	 -Spath         Compute classpath and echo to stdout only
	 -Scp CP        Do NOT compute or cache classpath, use this one instead
	 -Srepro        Ignore the ~/.clojure/deps.edn config file
	 -Sforce        Force recomputation of the classpath (don't use the cache)
	 -Spom          Generate (or update existing) pom.xml with deps and paths
	 -Stree         Print dependency tree
	 -Sresolve-tags Resolve git coordinate tags to shas and update deps.edn
	 -Sverbose      Print important path info to console
	 -Sdescribe     Print environment and command parsing info as data

	init-opt:
	 -i, --init path     Load a file or resource
	 -e, --eval string   Eval exprs in string; print non-nil values

	main-opt:
	 -m, --main ns-name  Call the -main function from namespace w/args
	 -r, --repl          Run a repl
	 path                Run a script from a file or resource
	 -                   Run a script from standard input
	 -h, -?, --help      Print this help message and exit

	For more info, see:
	 https://clojure.org/guides/deps_and_cli
	 https://clojure.org/reference/repl_and_main
`)

}

func main() {

	version := "1.10.0.447"

	wd, _ := os.Getwd()
	local_install_dir := wd + "/.."
	install_dir := os.Getenv("localappdata") + "/Programs/clojure"

	jarfile := fmt.Sprintf("/lib/libexec/clojure-tools-%s.jar", version)
	local_jar := local_install_dir + jarfile

	if jarfile_exists, _ := exists(local_jar); jarfile_exists {
		install_dir = local_install_dir
	}
	// fmt.Printf("local_install_dir: %s\n", local_install_dir)
	// fmt.Printf("local_jar: %s\n", local_jar)
	//fmt.Printf("install_dir: %s\n", install_dir)

	tools_cp := install_dir + jarfile
	print_classpath := false
	describe := false
	verbose := false
	force := false
	repro := false
	tree := false
	pom := false
	resolve_tags := false
	help := false
	stale := false

	jvm_opts := make([]string, 0)
	resolve_aliases := make([]string, 0)
	classpath_aliases := make([]string, 0)
	jvm_aliases := make([]string, 0)
	main_aliases := make([]string, 0)
	all_aliases := make([]string, 0)
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
		case strings.HasPrefix(arg, "-C"):
			classpath_aliases = append(classpath_aliases, arg[2:])
		case strings.HasPrefix(arg, "-O"):
			jvm_aliases = append(jvm_aliases, arg[2:])
		case strings.HasPrefix(arg, "-M"):
			main_aliases = append(main_aliases, arg[2:])
		case strings.HasPrefix(arg, "-A"):
			all_aliases = append(all_aliases, arg[2:])
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
		case strings.HasPrefix(arg, "-Sforce"):
			force = true
		case strings.HasPrefix(arg, "-Srepo"):
			repro = true
		case strings.HasPrefix(arg, "-Stree"):
			tree = true
		case strings.HasPrefix(arg, "-Spom"):
			pom = true
		case strings.HasPrefix(arg, "-Sresolve-tags"):

			resolve_tags = true
		case strings.HasPrefix(arg, "-S"):
			fmt.Println("Invalid option %s", arg)

			resolve_tags = true
		case (strings.HasPrefix(arg, "-h") || strings.HasPrefix(arg, "--help") || strings.HasPrefix(arg, "-?")):
			if len(main_aliases) == 0 && len(all_aliases) == 0 {
				help = true
			}
		default:
			additional_args = append(additional_args, arg)
		}

	}
	//	fmt.Printf("additional args: %s", additional_args)

	if help {
		usage()
		return
	}
	config_paths := make([]string, 0)

	config_dir := os.Getenv("CLJ_CONFIG")
	config_home := os.Getenv("XDG_CONFIG_HOME")
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	var config_str string

	var cmd *exec.Cmd
	var cmd_args []string

	deps_edn_file_exists = true
	if _, err := os.Stat("deps.edn"); os.IsNotExist(err) {
		deps_edn_file_exists = false
	}

	if resolve_tags {

		if !deps_edn_file_exists {
			fmt.Println("deps.edn does not exist")
			return
		}
		cmd = exec.Command("java.exe", "-Xmx256m", "-cp", tools_cp, "clojure.main", "-m",
			"clojure.tools.deps.alpha.script.resolve-tags", "--deps-file=deps.edn")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		//		fmt.Println("resolve-tags")
		//		fmt.Println(cmd_args)
		cmd.Run()
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
	if repro {
		config_paths = append(config_paths, install_dir+"/lib/deps.edn")
		config_paths = append(config_paths, "deps.edn")
	} else {
		config_paths = append(config_paths, install_dir+"/lib/deps.edn")
		config_paths = append(config_paths, config_dir+"/deps.edn")
		config_paths = append(config_paths, "deps.edn")
	}
	config_str = fmt.Sprintf(",%s", strings.Join(config_paths, ","))[1:]

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
	joined_result = append(joined_result, all_aliases...)
	joined_result = append(joined_result, jvm_aliases...)
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

	libs_file := cache_dir + "/" + ck + ".libs"
	cp_file := cache_dir + "/" + ck + ".cp"
	jvm_file := cache_dir + "/" + ck + ".jvm"
	main_file := cache_dir + "/" + ck + ".main"

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
	if force || !cp_file_exists {
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
		if len(jvm_aliases) > 0 {
			tools_args = append(tools_args, "-J"+strings.Join(jvm_aliases, ""))
		}
		if len(main_aliases) > 0 {
			tools_args = append(tools_args, "-M"+strings.Join(main_aliases, ""))
		}
		if len(all_aliases) > 0 {
			tools_args = append(tools_args, "-A"+strings.Join(all_aliases, ""))
		}
		if force_cp != "" {
			tools_args = append(tools_args, "--skip-cp")

		}

	}

	// If stale, run make-classpath to refresh cached classpath
	if stale && verbose {
		fmt.Println("Refreshing classpath")

	}
	if stale {

		make_classpath_args := []string{"-Xmx256m", "-cp", tools_cp, "clojure.main", "-m", "clojure.tools.deps.alpha.script.make-classpath", "--config-files", config_str, "--libs-file", libs_file, "--cp-file", cp_file, "--jvm-file", jvm_file, "--main-file", main_file}

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

	if pom {
		cmd_args := []string{"-Xmx256m", "-cp", tools_cp, "clojure.main", "-m", "clojure.tools.deps.alpha.script.generate-manifest", "--config-files", config_str, "--gen=pom"}
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
			":install-dir \"%s\"\n"+
			":config-dir \"%s\"\n"+
			":cache-dir \"%s\"\n"+
			":force %t\n"+
			":repro %t\n"+
			":resolve-aliases \"%s\"\n"+
			":classpath-aliases \"%s\"\n"+
			":jvm-aliases \"%s\"\n"+
			":main-aliases \"%s\"\n"+
			":all-aliases \"%s\"\n}\n", version, path_vector,
			install_dir, config_dir, cache_dir, force, repro, strings.Join(resolve_aliases, ""),
			strings.Join(classpath_aliases, ":"), strings.Join(jvm_aliases, ""),
			strings.Join(main_aliases, ":"), strings.Join(all_aliases, ""))

	} else if tree {
		cmd_args = []string{"-Xmx256m", "-cp", tools_cp, "clojure.main", "-m", "clojure.tools.deps.alpha.script.print-tree", "--libs-file", libs_file}
		cmd = exec.Command("java.exe", cmd_args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		//		fmt.Println("print-tree")
		//		fmt.Println(cmd_args)
		cmd.Run()
		return
	} else {
		var jvm_cache_opts string
		var main_cache_opts string

		if _, err := os.Stat(jvm_file); err == nil {
			content, err = ioutil.ReadFile(jvm_file)
			jvm_cache_opts = string(content)
			check(err)
		}
		if _, err := os.Stat(main_file); err == nil {
			content, err = ioutil.ReadFile(main_file)
			main_cache_opts = string(content)
			check(err)
		}

		jvm_opts_string := strings.Join(jvm_opts, " ")
		cmd_args = make([]string, 0)

		if jvm_cache_opts != "" {
			cmd_args = append(cmd_args, jvm_cache_opts)

		}

		if len(jvm_opts) > 0 {
			cmd_args = append(cmd_args, jvm_opts_string)

		}
		cmd_args = append(cmd_args, "-Dclojure.libfile="+libs_file, "-cp", cp, "clojure.main")

		if main_cache_opts != "" {
			words := strings.Fields(main_cache_opts)
			cmd_args = append(cmd_args, words...)

		}

		cmd_args = append(cmd_args, additional_args...)

		cmd = exec.Command("java.exe", cmd_args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		//		fmt.Println("clojure.main")

		//		fmt.Println(cmd_args)
		cmd.Run()

	}
}
