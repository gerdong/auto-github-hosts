package domain

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gerdong/auto-github-hosts/config"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	base_url          = "https://ipaddress.com/search/"
	header_user_agent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebkit/737.36(KHTML, like Gecke) Chrome/52.0.2743.82 Safari/537.36"
	header_host       = "ipaddress.com"
)

func UpdateHosts() map[string]string {
	m := make(map[string]string)
	for _, host := range config.Hosts {
		ip, err := getIP(host)
		if err != nil || ip == "" {
			continue
		}
		m[host] = ip
	}
	return m
}

// 使用 resty 访问 ipaddress.com，获得host的ip
func getIP(site string) (string, error) {
	url := "https://www.ipaddress.com/search/" + site
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request: ", err)
		return "", err
	}
	req.Header.Set("User-Agent", header_user_agent)
	req.Header.Set("Host", header_host)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request: ", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error parsing HTTP response: ", err)
		return "", err
	}

	htmlStr := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Println("Error parsing HTML: ", err)
		return "", err
	}

	var trueip string
	ipRe := regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	result := doc.Find(".comma-separated")
	if result != nil {
		resultStr := result.Text()
		if len(resultStr) > 0 {
			trueip = ipRe.FindString(resultStr)
		}
	}

	return trueip, nil
}

func UpdateHostsFile(m map[string]string) error {
	hostsFile := getHostsFile()
	file, err := os.OpenFile(hostsFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening hosts file: ", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "# auto-github-hosts") {
			break
		}
		lines = append(lines, line)
	}

	file.Seek(0, 0)
	file.Truncate(0)

	for _, line := range lines {
		fmt.Fprintln(file, line)
	}

	fmt.Fprintln(file, "# auto-github-hosts")
	for host, ip := range m {
		exists := false
		scanner = bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, host) {
				exists = true
				if !strings.Contains(line, ip) {
					line = strings.Split(line, " ")[0] + " " + ip
				}
			}
			fmt.Fprintln(file, line)
		}
		if !exists {
			fmt.Fprintf(file, "%s %s\n", ip, host)
		}
	}

	return nil
}

func getHostsFile() string {
	var hostsFile string
	switch runtime.GOOS {
	case "windows":
		hostsFile = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "darwin", "linux":
		hostsFile = "/etc/hosts"
	default:
		log.Fatal("unsupported operating system: %s\n", runtime.GOOS)
	}
	return hostsFile
}
