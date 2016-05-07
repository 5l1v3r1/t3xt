package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

type Config struct {
	AssetDir     string
	PasswordHash string
	DbPath       string

	// Parameters for whichlang.
	ClassifierPath string
	ClassifierType string

	ConfigPath string `json:"-"`
}

func GetConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return inputConfig(path)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	config.ConfigPath = path
	return &config, nil
}

func (c *Config) CheckPass(pass string) bool {
	return c.PasswordHash == hashPassword(pass)
}

func (c *Config) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.ConfigPath, data, 0755)
}

func inputConfig(path string) (*Config, error) {
	fmt.Print("Enter password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}

	prompts := []string{"asset path", "db path", "whichlang classifier path",
		"whichlang classifier type"}
	results := make([]string, len(prompts))

	for i, prompt := range prompts {
		fmt.Printf("Enter %s:", prompt)
		results[i], err = readLine()
		if err != nil {
			return nil, err
		}
	}

	c := &Config{
		AssetDir:       results[0],
		PasswordHash:   hashPassword(string(pass)),
		DbPath:         results[1],
		ClassifierPath: results[2],
		ClassifierType: results[3],
		ConfigPath:     path,
	}
	return c, c.Save()
}

func readLine() (string, error) {
	var res bytes.Buffer
	for {
		b := make([]byte, 1)
		if _, err := os.Stdin.Read(b); err != nil {
			return "", err
		}
		if b[0] == '\n' {
			break
		}
		res.WriteByte(b[0])
	}
	return res.String(), nil
}

func hashPassword(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
