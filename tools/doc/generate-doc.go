package main

import "log"

func main() {
	log.Default().Printf("Generating metrics documentation")
	generateDocMetrics()
	log.Default().Printf("Generating image CRD documentation")
	generateDocImageCRD()
}
