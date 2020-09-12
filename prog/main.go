package main

import (
	"log"
	"prog/secftp"
)

const (
	PrivateKey = "../ssh-keys/exemplo"
	Host       = "192.168.50.11"
	User       = "vagrant"
	Port       = 22
	FilesDir   = "./files"
	Filename   = "progit.pdf"
)

func main() {
	f, err := secftp.NewFtpSession(PrivateKey, User, Host, Port)
	if err != nil {
		log.Fatalf("Erro ao iniciar sessão SFTP: %s\n", err.Error())
	}
	defer f.Close()

	err = f.CopyToRemote(FilesDir, Filename)
	if err != nil {
		log.Fatalf("Erro ao copiar arquivo para a máquina remota: %s", err.Error())
	}

	err = f.CopyFromRemote(Filename)
	if err != nil {
		log.Println("Erro ao copiar arquivo para a máquina local: ", err.Error())
	}
}
