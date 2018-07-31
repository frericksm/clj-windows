package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"io"
		"log"
	"net/http"
	"os"
	"path/filepath"
)

//see https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r io.Reader) error {

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

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
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

func main() {

	url := "https://download.clojure.org/install/clojure-tools-1.9.0.381.tar.gz"
	fname := "clojure-tools-1.9.0.381.tar.gz"

	prefix_dir := os.Getenv("localappdata") + "/Programs"
	if len(os.Args) == 2 {
		prefix_dir = os.Args[1]
	}

	install_dir := prefix_dir + "/clojure"
	lib_dir := install_dir + "/lib"
	lib_exec := lib_dir + "/libexec"
	bin_dir := install_dir + "/bin"
	out, err := os.Create(fname)
	check(err)

	//	defer out.Close()

	resp, err := http.Get(url)
	check(err)

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)
	out.Close()

	// untar file
	f, err := os.Open(fname)
	check(err)
	//	defer f.Close()

	r := bufio.NewReader(f)
	err = Untar(".", r)
	check(err)
	err = f.Close()
	check(err)

	os.MkdirAll(install_dir, 0700)
	os.MkdirAll(lib_dir, 0700)
	os.MkdirAll(lib_exec, 0700)
	os.MkdirAll(bin_dir, 0700)

	//Installing libs into $clojure_lib_dir
	copy("clojure-tools/deps.edn", lib_dir+"/deps.edn")
	copy("clojure-tools/example-deps.edn", lib_dir+"/example-deps.edn")
	copy("clojure-tools/clojure-tools-1.9.0.381.jar", lib_exec+"/clojure-tools-1.9.0.381.jar")
	copy("../clojure/clojure.exe", bin_dir+"/clojure.exe")

	// delete  install files
	err = os.RemoveAll("clojure-tools")
	check(err)
	err = os.Remove(fname)
	check(err)

}
