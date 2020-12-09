package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type Persona struct {
	Nombre string
	Msg string
}

var ing = make(chan string)
var cM = make(chan Persona)

func errP(err error) {
	if err != nil {
		fmt.Sprintln("Error: ", err)
		return
	}
}

func bError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func cliente(ing chan string, cM chan Persona) {
	var msjStr string
	c,err := net.Dial("tcp",":9999")
	errP(err)
	defer c.Close()
	//Ingresando para conversar
	for {
		select {
		case msj := <-ing:
			err = gob.NewEncoder(c).Encode(msj)
			errP(err)
		case cChat := <-cM:
			msgChat := cChat.Nombre + ":" + cChat.Msg
			err = gob.NewEncoder(c).Encode(msgChat)
			errP(err)
		}
	}
	//tendra que receptar los mensajes
	le := gob.NewDecoder(c).Decode(&msjStr)
	bError(le)
	fmt.Println(msjStr)
}

func mostrar(ing chan string, cM chan Persona) {
	var msjStr string
	c,err := net.Dial("tcp",":9999")
	errP(err)
	defer c.Close()
	for {
		select {
		case msj := <-ing:
			err = gob.NewEncoder(c).Encode(msj)
			errP(err)
		case cChat := <-cM:
			msgChat := cChat.Nombre + ":" + cChat.Msg
			err = gob.NewEncoder(c).Encode(msgChat)
			errP(err)
		}
	}
	le := gob.NewDecoder(c).Decode(&msjStr)
	bError(le)		
}

func main() {
	var op string
	var usuarioNuevo Persona
	inp := bufio.NewScanner(os.Stdin)
	go cliente(ing,cM)
	fmt.Print("Nickname: ")
	inp.Scan()
	nmUs := inp.Text() //No tiene un valor
	usuarioNuevo.Nombre = nmUs
	ing <- "Esta conectado: " + nmUs
	fmt.Println("1.-  Enviar mensaje")
	fmt.Println("2.- Enviar archivo")
	fmt.Println("3.- Mostrar mensajes")
	fmt.Println("0.- Salir")
	fmt.Println("Opcion: ")
	for {
		_,_ = fmt.Scanln(&op)
		switch op {
		case "1":

			inp.Scan()
			msg := inp.Text()
			usuarioNuevo.Msg = msg
			cM <- usuarioNuevo
		// enviando mensaje
		case "2":
			fmt.Println("No implementada!")
		case "3":
			fmt.Println("Mostrando..")
			go cliente(ing,cM)
		case "0":
			ing <- "Esta desconectado: " + nmUs
			var exit string
			_, _ = fmt.Scanln(&exit)
			return
		// Default.
		default:
			fmt.Print("\n-Opcion invalida!\n\n")
		}
	}
}