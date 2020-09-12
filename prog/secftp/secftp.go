package secftp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const RemoteDir = "/home/vagrant"

type FtpSession struct {
	signer   ssh.Signer
	username string
	address  string
	client   *sftp.Client
}

func NewFtpSession(pvtKey, username, host string, port int) (*FtpSession, error) {

	pvtKeyFile, err := ioutil.ReadFile(pvtKey)

	signer, err := ssh.ParsePrivateKey(pvtKeyFile)
	if err != nil {
		panic("ParsePrivateKey err")
	}

	return &FtpSession{
		signer:   signer,
		username: username,
		address:  fmt.Sprintf("%s:%d", host, port),
	}, nil
}

func (ftp *FtpSession) ensureConnection() error {

	if ftp.client != nil {
		_, err := ftp.client.Stat("/dev")
		if err == nil {
			return nil
		}
	}

	sshConfig := &ssh.ClientConfig{
		User: ftp.username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(ftp.signer)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey() // should we check the fingerprint?

	client, err := ssh.Dial("tcp", ftp.address, sshConfig)
	if err != nil {
		panic("dial failed: " + err.Error())
	}

	c, err := sftp.NewClient(client, sftp.MaxPacket(1<<15))
	if err != nil {
		panic("sftp NewClient failed")
	}

	ftp.client = c

	return nil
}

func (ftp *FtpSession) OpenFile(file string) (*sftp.File, error) {

	if err := ftp.ensureConnection(); err != nil {
		return nil, err
	}

	return ftp.client.Open(file)
}

func (ftp *FtpSession) CopyToRemote(localDir string, filename string) error {
	if err := ftp.ensureConnection(); err != nil {
		return fmt.Errorf("erro ao garantir conexão: %s", err.Error())
	}

	f := fmt.Sprintf("%s/%s", RemoteDir, filename)

	destFile, err := ftp.client.Create(f)
	if err != nil {
		return fmt.Errorf("Erro ao criar arquivo [%s] na máquina remota: %s", f, err.Error())
	}
	defer destFile.Close()

	srcFile, err := os.Open(fmt.Sprintf("%s/%s", localDir, filename))
	if err != nil {
		return fmt.Errorf("Erro ao abrir arquivo na máquina local: %s", err.Error())
	}
	defer srcFile.Close()

	n, err := io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("erro ao copiar arquivo: %s", err.Error())
	}

	log.Printf("Total de bytes copiados: %d\n", n)

	return nil
}

func (ftp *FtpSession) CopyFromRemote(filename string) error {
	if err := ftp.ensureConnection(); err != nil {
		return fmt.Errorf("erro ao garantir conexão: %s", err.Error())
	}

	srcFile, err := ftp.client.Open(fmt.Sprintf("%s/%s", RemoteDir, filename))
	if err != nil {
		return fmt.Errorf("Erro ao abrir arquivo na máquina remota: %s", err.Error())
	}
	defer srcFile.Close()

	destFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Erro ao criar arquivo na máquina local: %s", err.Error())
	}
	defer destFile.Close()

	n, err := io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("erro ao copiar arquivo: %s", err.Error())
	}

	log.Printf("Total de bytes copiados: %d\n", n)

	return nil
}

func (ftp *FtpSession) Close() {
	if ftp.client != nil {
		ftp.client.Close()
	}
}
