package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
)

var ip net.IP

func main() {
	var err error
	log.Println("Started server")
	ifaces, err := net.Interfaces()

	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP != nil && v.IP.To4() != nil && !v.IP.IsLoopback() && !v.IP.IsUnspecified() {
					ip = v.IP
					break
				}
			}
			// process IP address
		}
	}

	log.Println("Listening on: ", ip, ":9090")
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/pausePlay", pausePlay)
	http.HandleFunc("/volumeUp", volumeUp)
	http.HandleFunc("/volumeDown", volumeDown)
	http.HandleFunc("/test", test)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func volumeUp(w http.ResponseWriter, r *http.Request) {
	log.Println("Volume up")
	exec.Command("/usr/bin/pulseaudio-ctl", "up", "5").Output()
}

func volumeDown(w http.ResponseWriter, r *http.Request) {
	log.Println("Volume down")
	exec.Command("/usr/bin/pulseaudio-ctl", "down", "5").Output()
}

func pausePlay(w http.ResponseWriter, r *http.Request) {
	log.Println("Pause Play")
	exec.Command("/usr/bin/playerctl", "play-pause").Output()
}
func test(w http.ResponseWriter, r *http.Request) {
	log.Println("test button")
	exec.Command("xrandr", "--output", "HDMI-A-0", "--auto", "--above", "eDP").Output()
	var response, err = exec.Command("/usr/bin/cat", "/proc/acpi/button/lid/LID/state").Output()
	if err != nil {
		log.Fatal(err)
	}
	w.Write(response)
}
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,
		`<html>
    <head>
      <script>
        function pausePlay() { fetch('http://`+ip.String()+`:9090/pausePlay'); }
        function volumeUp() { fetch('http://`+ip.String()+`:9090/volumeUp'); }
        function volumeDown() { fetch('http://`+ip.String()+`:9090/volumeDown'); }
        function test() {
			let resp = fetch('http://`+ip.String()+`:9090/test')
          .then((response) => response.text())
          .then((data) => {
            console.log(data);
            document.querySelector('#TestButton').textContent = data;
        })
        }
      </script>
    </head>
  
    <body>
      <button onClick="volumeUp()">Volume Up</button>
      <button onClick="volumeDown()">Volume Down</button>
      <br />
      <button onClick="pausePlay()">Pause/Play</button>
      <button id="TestButton" onClick="test()">Test</button>
    </body>
  </html>`)
}
