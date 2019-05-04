package kfnetwork

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type ServerConfiguration struct {
}

func GenerateUUID() string {
	var buffer []byte
	const choices = "abcdef0123456789"

	buffer = make([]byte, 36)

	rand.Seed(int64(time.Now().UnixNano()))

	for x := range buffer {
		if x == 8 || x == 13 || x == 18 || x == 23 {
			buffer[x] = byte('-')
		} else {
			buffer[x] = choices[rand.Intn(len(choices))]
		}
	}
	return string(buffer)
}

func LoadConfig(filename string) (ServerConfiguration, error) {
	config := ServerConfiguration{}

	bytes, e := ioutil.ReadFile(filename)

	if e != nil {
		return config, e
	}

	e = json.Unmarshal(bytes, &config)

	return config, e
}

func SaveConfig(config ServerConfiguration, filename string) error {
	file, e := os.Create(filename)
	defer file.Close()

	if e != nil {
		return e
	}

	bytes, e := json.MarshalIndent(config, "", "    ")

	if e != nil {
		return e
	}

	_, e = file.WriteString(string(bytes))

	return e
}
