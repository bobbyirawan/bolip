package bolip

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

type Dialer struct {
	client *ssh.Client
}

func (v *Dialer) Dial(address string) (net.Conn, error) {
	return v.client.Dial("tcp", address)
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

func MysqlOnTunner(user, keyPath, port, host string) {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			publicKeyFile(keyPath),
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), sshConfig)
	if client != nil {
		defer client.Close()
	}
	if err != nil {
		log.Fatal("Failed to dial. " + err.Error())
	}

	mysql.RegisterDial("mysql+tcp", (&Dialer{client}).Dial)
}
