package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	TemplatePath     = "./README.template.md"
	HostsPath        = "./hosts"
	READMEPath       = "./README.md"
	ReadmeZhCnPath   = "./README-ZH_CN.md"
	RecordRetries    = 3
	IndentWhitespace = "                " // 16 spaces
)

// var urls = []string{"github.com", "assets-cdn.github.com", "documentcloud.github.com", "gist.github.com", "help.github.com",
//
//	"nodeload.github.com", "raw.github.com", "status.github.com", "training.github.com", "avatars.githubusercontent.com",
//	"codeload.github.com", "favicons.githubusercontent.com", "gist.githubusercontent.com", "marketplace-screenshots.githubusercontent.com",
//	"raw.githubusercontent.com", "repository-images.githubusercontent.com", "user-images.githubusercontent.com",
//	"avatars0.githubusercontent.com", "avatars1.githubusercontent.com", "avatars2.githubusercontent.com", "avatars3.githubusercontent.com",
//	"avatars4.githubusercontent.com", "avatars5.githubusercontent.com", "avatars6.githubusercontent.com", "avatars7.githubusercontent.com",
//	"avatars8.githubusercontent.com"}
var urls = []string{"github.com"}

type HostConfig struct {
	Url string
	Ip  string
}

func main() {
	configs, err := resolveUrls()
	if err != nil {
		log.Fatal(err)
	}

	writeHosts(configs)
}

func retry(fn func() (*dns.Msg, error), n int) (*dns.Msg, error) {
	var err error
	for i := 0; i < n; i++ {
		response, e := fn()
		err = e
		if err == nil {
			return response, nil
		}
	}
	return nil, err
}

func getIP(url string) (string, error) {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(url), dns.TypeA)
	response, err := retry(func() (*dns.Msg, error) {
		in, _, err := c.Exchange(&m, "114.114.114.114:53")
		return in, err
	}, RecordRetries)

	if err != nil {
		return "", err
	}

	if len(response.Answer) == 0 || response.Rcode != dns.RcodeSuccess {
		return "", fmt.Errorf("DNS query error")
	}

	return response.Answer[0].(*dns.A).A.String(), nil
}

func getHostConfig(url string) (*HostConfig, error) {
	ip, err := getIP(url)
	if err != nil {
		return &HostConfig{Url: url, Ip: ""}, nil
	}
	return &HostConfig{Url: url, Ip: ip}, nil
}

func resolveUrls() ([]*HostConfig, error) {
	var configs []*HostConfig
	for _, url := range urls {
		config, err := getHostConfig(url)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func generateHosts(configs []*HostConfig) (string, string) {
	var sb strings.Builder
	sb.WriteString("# Generated by Github Hosts start \n\n")
	var updateTime time.Time
	for _, config := range configs {
		if config.Ip == "" {
			sb.WriteString(fmt.Sprintf("# %s update failed\n", config.Url))
		} else {
			whiteLen := len(IndentWhitespace) - len(config.Ip)
			sb.WriteString(fmt.Sprintf("%s%s %s\n", config.Ip, IndentWhitespace[:whiteLen], config.Url))
		}
	}
	updateTime = time.Now()
	sb.WriteString(fmt.Sprintf("\n# Last update: %s\n", updateTime.Format("2006-01-02 15:04:05")))
	sb.WriteString("# Please star: https://github.com/fliu2476/gh-hosts.git\n")
	sb.WriteString("# Generated by Github Hosts end")

	return sb.String(), updateTime.Format("2006-01-02 15:04:05")
}

func writeHosts(configs []*HostConfig) error {
	hostStr, updateTime := generateHosts(configs)

	template, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		return err
	}
	nextReadme := strings.ReplaceAll(string(template), "{{hosts}}", hostStr)
	nextReadme = strings.ReplaceAll(nextReadme, "{{last_update_time}}", updateTime)

	if err = ioutil.WriteFile(HostsPath, []byte(hostStr), 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(READMEPath, []byte(nextReadme), 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(ReadmeZhCnPath, []byte(nextReadme), 0644); err != nil {
		return err
	}

	return nil
}
