package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	usuarios []net.Conn
	msg []string
)

func errBasic(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func panicIfErr(err error) {
	if err != nil {
		fmt.Sprintln("Error: ", err)
		return
	}
}

func server() {
	serv, err := net.Listen("tcp",":9999")
	panicIfErr(err)
	for {
		c, err := serv.Accept()
		panicIfErr(err)
		usuarios = append(usuarios,c)
		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	var cadena string
	///conecto con el cliente y recibo la informacion
	for {
		err := gob.NewDecoder(c).Decode(&cadena)
		panicIfErr(err)
		fmt.Println(cadena)
		msg = append(msg,cadena)
		dCli := strings.Contains(cadena,"Se ha desconectado")
		if dCli == true {
			for j := 0; j < len(usuarios); j++ {
				if (c == usuarios[j]) {
					usuarios = append(usuarios[:j], usuarios[j+1:]...)
				}
			}
		}

		/// habra que checar las cadenas que hay dentro 
		for i:=0; i < len(usuarios); i++ {
			//levanto bandera
			if c != usuarios[i] {
				err := gob.NewEncoder(usuarios[i]).Encode(cadena)
				errBasic(err)
				// Checando las cadenas que estuvieron disponibles
			}
		}
	}
}

//Se necesita hacer un backup
func backUp() {
	//Para el backup es necesario un archivo
	arch, err := os.Create("Chat.txt")
	panicIfErr(err)
	defer arch.Close()
	///agarrar los datos para traspasarlos al arch
	for _,str := range msg {
		_,_ = arch.WriteString(str + "\n")
	}
	fmt.Println("Guardado correcto")
}

func main() {
	var op string
	go server()
	fmt.Println("Arrancando Servidor...")
	fmt.Println("1.- Guardar mensajes. 0.- Salir.")
	fmt.Println("Mensajes: ")
	for {
		_, _ = fmt.Scanln(&op)
		switch op {
		case "1":
			backUp()
		case "0":
			fmt.Println("\n Saliendo..")
		default:
			fmt.Println("\n Opcion invalida\n\n")
		}
	}
}