package server

import (
	"log"
	"net"
)

func SetIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Interfaces error:", err)
		return ""
	}
	for _, iface := range interfaces {
		if iface.Name == "Wi-Fi" || iface.Name == "Беспроводная сеть" {
			addrs, err := iface.Addrs()
			if err != nil {
				log.Println("Wi-Fi search error:", err)
				return ""
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					log.Print("Сервер запущен на адресе http://", ip, ":5051/")
					return ip
				}
			}
		}
	}
	log.Println("Не найден необходимый адрес, сервер не будет запущен.")
	return ""
}
