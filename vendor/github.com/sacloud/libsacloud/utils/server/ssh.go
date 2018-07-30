package server

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/mattn/go-tty"
	"golang.org/x/crypto/ssh"
)

// SSHClientParams represents SSHClient params
type SSHClientParams struct {
	DisplayName    string
	UserName       string
	Password       string
	Host           string
	Port           int
	PrivateKeyPath string
	Quiet          bool
	Out            io.Writer
}

// TargetHost returns hostname as 'host:port'
func (p *SSHClientParams) TargetHost() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

// CreateSSHClient returns new *ssh.Client by SSHClientParams
func CreateSSHClient(params *SSHClientParams) (*ssh.Client, error) {

	// collect file info
	fileExists := false
	if params.PrivateKeyPath != "" {
		_, err := os.Stat(params.PrivateKeyPath)
		fileExists = err == nil
	}

	// collect option
	strOpt := ""
	if fileExists {
		strOpt = fmt.Sprintf(" -i %s", params.PrivateKeyPath)
	}

	if !params.Quiet {
		// output connect info
		fmt.Fprintf(params.Out, "\nconnecting server...\n\tcommand: ssh %s@%s%s\n\n", params.UserName, params.TargetHost(), strOpt)
	}
	cnf := &ssh.ClientConfig{
		User:    params.UserName,
		Timeout: 10 * time.Second,
	}

	// build auth methods
	var authMethods []ssh.AuthMethod

	// add ssh-agent
	//sshSock := os.ExpandEnv("$SSH_AUTH_SOCK")
	//if sshSock != "" {
	//	addr, _ := net.ResolveUnixAddr("unix", sshSock)
	//	agentConn, _ := net.DialUnix("unix", nil, addr)
	//	ag := agent.NewClient(agentConn)
	//	authMethods = append(authMethods, ssh.PublicKeysCallback(ag.Signers))
	//}

	// private key

	if fileExists {
		signer, err := getSigners(params.PrivateKeyPath, params.Password)
		if err != nil {
			return nil, fmt.Errorf("parse private-key(%s) is failed: %s", params.PrivateKeyPath, err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer...))
	}

	// password prompt
	authMethods = append(authMethods, ssh.PasswordCallback(func() (string, error) {
		if params.Password == "" {
			prefix := ""
			if params.DisplayName != "" {
				prefix = fmt.Sprintf("%s | ", params.DisplayName)
			}
			return pprompt(fmt.Sprintf("%spassword: ", prefix))
		}
		return params.Password, nil
	}))

	cnf.Auth = authMethods
	var conn *ssh.Client

	conn, err := ssh.Dial("tcp", params.TargetHost(), cnf)
	if err != nil {
		return nil, fmt.Errorf("connecting(%s) is failed: %s", params.TargetHost(), err)
	}

	return conn, nil
}

func pprompt(prompt string) (string, error) {
	t, err := tty.Open()
	if err != nil {
		return "", err
	}
	defer t.Close()
	fmt.Print(prompt)
	defer t.Output().WriteString("\r" + strings.Repeat(" ", runewidth.StringWidth(prompt)) + "\r")
	return t.ReadPasswordClear()
}

func getSigners(keyfile string, password string) ([]ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	b, _ := pem.Decode(buf)
	if x509.IsEncryptedPEMBlock(b) {
		pass := password
		if pass == "" {
			p, err := pprompt("pass-phrase: ")
			if err != nil {
				return nil, fmt.Errorf("ServerSsh is failed: collecting is failed: %s", err)
			}
			pass = p
		}
		buf, err = x509.DecryptPEMBlock(b, []byte(pass))
		if err != nil {
			return nil, err
		}
		pk, err := x509.ParsePKCS1PrivateKey(buf)
		if err != nil {
			return nil, err
		}
		k, err := ssh.NewSignerFromKey(pk)
		if err != nil {
			return nil, err
		}
		return []ssh.Signer{k}, nil
	}
	k, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return []ssh.Signer{k}, nil
}
