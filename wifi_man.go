package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

func conn_on() bool { // проверяет наличие соединения с сетью "интернет"
	var n int
	var to = 15
	for i := 0; i < to; i++ {
		out, err := exec.Command("ping", "-c 1", "-i 1", "-W 1", "ya.ru").Output()
		if strings.Contains(string(out), "0 received") {
			n++
		}
		if err != nil {
			n += 5
		}
		if strings.Contains(string(out), "1 received") {
			time.Sleep(60 * time.Second)
			return true
		}

		if n == to {
			return false
		}
	}
	return true

}
func main() {
	path := flag.String("path", " ", "path to nmconnections file")
	flag.Parse()
	for {
		files, _ := ioutil.ReadDir(*path)
		for i := 0; i < len(files); i++ {
			var conn_name string
			var bssid string
			//  читаем файл в список
			File_content, _ := ioutil.ReadFile(*path + files[i].Name())
			List_of_file_content := strings.Split(string(File_content), "\n")
			for _, file_string := range List_of_file_content {
				// извлекаем из него мак-адрес сети
				if strings.Contains(file_string, "uuid=") {
					continue
				}
				if strings.Contains(file_string, "id=") && !strings.Contains(file_string, "bssid=") {
					conn_name = strings.Split(file_string, "id=")[1]
				}
				if strings.Contains(file_string, "bssid=") {
					bssid = strings.Split(file_string, "bssid=")[1]
				}
				if strings.Contains(file_string, "seen-bssids=") {
					bssid = strings.Split(file_string, "seen-bssids=")[1]
					bssid = strings.Split(bssid, ";")[0]
				}
				if len(conn_name) > 0 && len(bssid) > 0 {
					break
				}
			}
			// получаем список всех доступных сетей и
			//если ее  нет в списке доступных
			iw_list_scan_str, _ := exec.Command("sudo", "iwlist", "wlan0", "scan").Output()
			if !strings.Contains(string(iw_list_scan_str), bssid) {
				continue //новая итерация
			}
			// если сеть есть, повторяем ту же итерацию цикла с тем же файлом
			if conn_on() {
				i--
				continue
			} else {
				fmt.Println("No network ", time.Now())
				conn_up, _ := exec.Command("sudo", "nmcli", "con", "up", conn_name).Output()
				if strings.Contains(string(conn_up), "успешно активировано") {
					time.Sleep(60 * time.Second)
					proton_up, _ := exec.Command("protonvpn-cli", "s").Output()
					if strings.Contains(string(proton_up), "No active") {
						exec.Command("protonvpn-cli", "c", "-f").Wait()
						time.Sleep(60 * time.Second)
					}
				}
			}
		}
	}
}
