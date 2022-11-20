package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/groovy-sky/azuredoh/v2/pkg/table"
)

func main() {
	var aztable table.AzureTable

	connStr, exist := os.LookupEnv("AzureWebJobsStorage")

	if !exist {
		fmt.Println("[ERR] Couldn't obtain connection string")
		return
	}
	tableName := "table3"

	err := aztable.Init(connStr, tableName)

	if err != nil {
		panic(err)
	}
	//aztable.SetEntry("test.com")
	//_, blocked := aztable.GetEntry("aaa.test.com")

	//fmt.Println(blocked)

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
		fmt.Println(domain[:len(domain)-1])

		aztable.SetEntry(domain[:len(domain)-1])
	}
}
