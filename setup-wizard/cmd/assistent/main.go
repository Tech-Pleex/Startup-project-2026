// Assistenten: starter det lokale API, serverer tringuiden og åbner
// elevens browser. Al logik ligger i internal/server og er testet dér.
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/osops"
	"github.com/Tech-Pleex/Startup-project-2026/setup-wizard/internal/server"
)

func main() {
	osImpl := osops.Current()
	srv := server.New(osImpl)

	// Port 0: OS'et vælger en ledig port, så Assistenten aldrig
	// kolliderer med noget andet på elevens maskine.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("kunne ikke starte Assistenten: %v", err)
	}

	go func() {
		if err := http.Serve(ln, srv); err != nil {
			log.Fatalf("serveren stoppede uventet: %v", err)
		}
	}()

	url := "http://" + ln.Addr().String()
	log.Printf("Assistenten kører på %s", url)
	if err := osImpl.OpenURL(url); err != nil {
		log.Printf("kunne ikke åbne browseren automatisk — åbn selv %s", url)
	}

	<-srv.Quit()
}
