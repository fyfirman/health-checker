package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

func PrintReader(r io.Reader) {
	bodyBytes, err := ioutil.ReadAll(r)

	if err != nil {
		log.Print(err)
		return
	}

	PrintByte(bodyBytes)
}

func PrintByte(bodyBytes []byte) {
	var err error

	if len(bodyBytes) > 0 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
			log.Printf("JSON parse error: %v", err)
		}
		log.Println(prettyJSON.String())
	} else {
		log.Printf("Body: No Body Supplied\n")
	}
}
