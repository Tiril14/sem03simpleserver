package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/Tiril14/funtemps/conv"
	"github.com/Tiril14/is105sem03/mycrypt"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			var x string
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke

					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))

					if strings.HasPrefix(string(dekryptertMelding), "Kjevik") {
						fields := strings.Split(string(dekryptertMelding), ";")
						if len(fields) >= 4 {
							celsius, err := strconv.ParseFloat(fields[3], 64)
							if err != nil {
								log.Println(err)
								continue
							}
							fahrenheit := conv.CelsiusToFahrenheit(celsius)
							x = fmt.Sprintf("%s;%s;%s;%.1f", fields[0], fields[1], fields[2], fahrenheit)
							if err != nil {
								log.Println(err)
								return
							}

						} else {
							log.Println("invalid input:", string(dekryptertMelding))
						}
					} else {
						x = string(dekryptertMelding)
					}

					msg := string(dekryptertMelding)
					switch msg {

					case "ping":
						svar := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(string(svar)))

					default:
						svar := mycrypt.Krypter([]rune(x), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(string(svar)))

					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
