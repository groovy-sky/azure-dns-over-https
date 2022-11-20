package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

func main() {

	url := "https://raw.githubusercontent.com/Ultimate-Hosts-Blacklist/Ultimate.Hosts.Blacklist/master/domains/domains0.list"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("[ERR] Couldn't reach " + url)
	}

	defer res.Body.Close()

	file := bufio.NewReader(res.Body)

	for {
		domain, err := file.ReadString('\n')

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		}

		fmt.Print(domain)
	}

}
