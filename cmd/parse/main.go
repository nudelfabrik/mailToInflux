package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nudelfabrik/mailToInflux/dmarc"
)

func main() {
	xmlFile, err := os.Open("dmarc.xml")
	if err != nil {
		log.Fatal(err)
	}

	defer xmlFile.Close()

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Println(err)
	}

	report, err := dmarc.ParseDmarc(bytes)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(report)
}
