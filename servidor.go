package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

// Proceso ...
type Proceso struct {
	ID     int
	Count  int
	Status bool
}

func initProceso() []Proceso {
	s := make([]Proceso, 0)
	for i := 0; i < 5; i++ {
		aux := Proceso{ID: i, Count: 0, Status: true}
		s = append(s, aux)
	}
	return s
}

// ProcesoIncrementa ...
func ProcesoIncrementa(c chan Proceso, p *Proceso) {
	for {
		switch p.Status {
		case true:
			c <- *p
			p.Count++
			time.Sleep(time.Millisecond * 500)
		case false:
			return
		}
	}
}

func printer(c chan Proceso) {
	cont := 0
	Pr := 5
	for {
		p := <-c
		Pr = 5
		if p.Status == true {
			fmt.Printf("Id:%d : %d \n", p.ID, p.Count)
		}
		cont++
		for i := 0; i < len(proces); i++ {
			if proces[i].Status == false {
				Pr--
			}
		}
		if cont == Pr {
			cont = 0
			fmt.Println("********")
		}
	}
}

var proces []Proceso
var imp = make(chan Proceso)

func servidor() {
	proces = initProceso()
	for i := int64(0); i < 5; i++ {
		go ProcesoIncrementa(imp, &proces[i])
	}
	go printer(imp)

	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Sprintln(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleClient(c)
	}
}
func handleClient(c net.Conn) {
	var p Proceso

	err := gob.NewDecoder(c).Decode(&p)

	if err != nil {
		fmt.Println("3:", p)
		fmt.Println(err)
		return
	} else {
		switch p.Count {
		case 0:
			go mandarClinte(c)
		default:
			go continuarS(&p)
		}
	}
}

func mandarClinte(c net.Conn) {
	var enviar int
	for i := int(0); i < len(proces); i++ {
		if proces[i].Status {
			proces[i].Status = false
			enviar = i
			break
		}
	}
	err := gob.NewEncoder(c).Encode(proces[enviar])
	if err != nil {
		fmt.Println(err)
		proces[enviar].Status = true
		go ProcesoIncrementa(imp, &proces[enviar])
	}
	return
}

func continuarS(p *Proceso) {
	proces[p.ID].Status = true
	proces[p.ID].Count = p.Count
	go ProcesoIncrementa(imp, &proces[p.ID])
	return
}

func main() {
	go servidor()
	var input string
	fmt.Scanln(&input)
}
