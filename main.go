package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	log.Println("Started server")
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/pausePlay", pausePlay)
	http.HandleFunc("/volumeUp", volumeUp)
	http.HandleFunc("/volumeDown", volumeDown)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func volumeUp(w http.ResponseWriter, r *http.Request) {
	exec.Command("/usr/bin/pulseaudio-ctl", "up", "5").Output()
}

func volumeDown(w http.ResponseWriter, r *http.Request) {
	exec.Command("/usr/bin/pulseaudio-ctl", "down", "5").Output()
}

func pausePlay(w http.ResponseWriter, r *http.Request) {
	exec.Command("/usr/bin/playerctl", "play-pause").Output()
}
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,
		`<html>
    <head>
      <script>
        function pausePlay() { fetch('http://10.170.154.42:9090/pausePlay'); }
        function volumeUp() { fetch('http://10.170.154.42:9090/volumeUp'); }
        function volumeDown() { fetch('http://10.170.154.42:9090/volumeDown'); }
      </script>
    </head>
  
    <body>
      <button onClick="volumeUp()">Volume Up</button>
      <button onClick="volumeDown()">Volume Down</button>
      <br />
      <button onClick="pausePlay()">Pause/Play</button>
    </body>
  </html>`)
}
