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

var imp = make(chan Proceso)

// ProcesoIncrementa ...
func ProcesoIncrementa(c chan Proceso, p *Proceso) {
	for {
		switch p.Status {
		case true:
			c <- *p
			p.Count++
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func printer(c chan Proceso) {
	for {
		p := <-c
		if p.Status {
			fmt.Printf("Id:%d : %d \n", p.ID, p.Count)
		}
	}
}
func pedirProceso(p *Proceso) {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewDecoder(c).Decode(&p)
	if err != nil {
		fmt.Println(err)
		return
	}
	go ProcesoIncrementa(imp, p)
	go printer(imp)
	c.Close()
	var input string
	fmt.Scanln(&input)

}

func devolver(p *Proceso) {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Close()

}

func main() {
	pr := Proceso{0, 0, true}
	pedirProceso(&pr)
	aux := pr
	pr.Status = false
	devolver(&aux)
}
