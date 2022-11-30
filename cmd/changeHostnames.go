package main

import (
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, _ := ioutil.ReadAll(fp)
	signer, _ = ssh.ParsePrivateKey(buf)
	return
}

func makeKeyring() ssh.ClientAuth {
	signers := []ssh.Signer{}
	keys := []string{os.Getenv("HOME") + "/.ssh/id_rsa", os.Getenv("HOME") + "/.ssh/id_dsa"}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return ssh.ClientAuthKeyring(&SignerContainer{signers})
}

func main() {

	hostname := "10.168.0.1"
	config := &ssh.ClientConfig{
		User: "pojos",
		Auth: []ssh.AuthMethod{makeKeyring()},
	}

	conn, _ := ssh.Dial("tcp", hostname+":22", config)
	session, _ := conn.NewSession()
	defer session.Close()

}
