package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type IpAlias struct {
	Alias string `json:"alias"`
	Ip    string `json:"ip"`
}

type Output struct {
	Ip string
}

func get_my_ip() (ip string, err error) {
	res, err := http.Get("https://icanhazip.com")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("error, status code: %v\n", res.StatusCode)
		os.Exit(1)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	ip = strings.TrimSuffix(string(body), "\n")
	return ip, nil
}

func get_aliases() ([]IpAlias, bool) {
	home_dir, ok := os.LookupEnv("HOME")

	if ok == false {
		return []IpAlias{}, ok
	}
	s := []string{home_dir, ".myip.json"}
	config_file := strings.Join(s, "/")
	jsonFile, err := os.Open(config_file)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return []IpAlias{}, false
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return []IpAlias{}, false
	}

	var aliases []IpAlias
	err = json.Unmarshal(byteValue, &aliases)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return []IpAlias{}, false
	}

	return aliases, true

}
func main() {
	aliases, ok := get_aliases()
	var ip_output Output

	ip, err := get_my_ip()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// Look for alias
	if ok == true {
		for _, alias := range aliases {
			if alias.Ip == ip {
				ip_output.Ip = alias.Alias
				break
			}
		}
	}
	if ip_output.Ip == "" {
		ip_output.Ip = ip
	}

	fmt.Println(ip_output.Ip)
}
