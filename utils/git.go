package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	go_git_ssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"os"
)

func CloneRepo(name string) error {
	_, err := git.PlainClone(fmt.Sprintf("./%s", name), false, &git.CloneOptions{
		URL:      "git@github.com:pismo/" + name,
		Progress: os.Stdout,
		Auth:     get_ssh_key_auth(os.Getenv("HOME") + "/.ssh/id_rsa"),
	})

	if err != nil {
		return err
	}

	return nil
}

func get_ssh_key_auth(privateSshKeyFile string) transport.AuthMethod {
	var auth transport.AuthMethod
	sshKey, _ := ioutil.ReadFile(privateSshKeyFile)
	signer, _ := ssh.ParsePrivateKey([]byte(sshKey))
	auth = &go_git_ssh.PublicKeys{User: "git", Signer: signer}
	return auth
}
