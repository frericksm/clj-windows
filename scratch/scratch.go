package main

import (
	"fmt"
	"net/http"
	//	"io/ioutil"
	//	"os"
)

func main() {
	_, err := http.Get("https://jira.intern/secure/Dashboard.jspa")
	if err != nil {
		fmt.Printf("%s", err)
		//		response.Body.Close()

		//		contents, err := ioutil.ReadAll(response.Body)
		//		if err != nil {
		//			fmt.Printf("%s", err)
		//			os.Exit(1)
		//		}
		//		fmt.Printf("%s\n", string(contents))
	}
}
