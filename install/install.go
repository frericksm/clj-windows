// The main package of the installer

package main

import (
	"archive/tar"
	//	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	//	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files

//see https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

func untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			//			f.Close()
		}
	}
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
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// accepts one optional commandline  arg which is the name of the directory in which the clojure dir will be created.
//
func main() {

	version := "1.10.1.489"
	lib_name := fmt.Sprintf("clojure-tools-%s.jar", version)

	prefix_dir := os.Getenv("localappdata") + "/Programs"
	if len(os.Args) == 2 {
		prefix_dir = os.Args[1]
	}

	install_dir := prefix_dir + "/clojure"
	lib_dir := install_dir + "/lib"
	lib_exec := lib_dir + "/libexec"
	bin_dir := install_dir + "/bin"
	os.MkdirAll(install_dir, 0700)
	os.MkdirAll(lib_dir, 0700)
	os.MkdirAll(lib_exec, 0700)
	os.MkdirAll(bin_dir, 0700)

	//Installing libs into $clojure_lib_dir
	copy("deps.edn", lib_dir+"/deps.edn")
	copy("example-deps.edn", lib_dir+"/example-deps.edn")
	copy(lib_name, lib_exec+"/"+lib_name)
	copy("clojure.exe", bin_dir+"/clojure.exe")
	copy("clj.exe", bin_dir+"/clj.exe")

	// add bin_dir to path

	cmd_args := []string{"-ExecutionPolicy", "remotesigned", "-File", "setpath.ps1", "-install_dir", install_dir}

	cmd := exec.Command("Powershell.exe", cmd_args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}
