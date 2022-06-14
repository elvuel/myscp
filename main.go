package main

import (
	"context"
	"fmt"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/integrii/flaggy"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	ltr = "ltr"
	rtl = "rtl"
)

var (
	remoteHost     string
	username       string
	passphrase     string
	privateKeyPath string

	orientation string

	localFile  string
	remoteFile string
)

func readPassword(prompt string, passphraseVar *string) {
	fmt.Print(prompt)
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	*passphraseVar = string(pass)
}

func init() {
	flaggy.DefaultParser.ShowHelpWithHFlag = false
	flaggy.SetName("MySCP")
	flaggy.SetDescription("A simple SCP tool, that can help me out for copy files between machines, due to some policy restrictions.")

	flaggy.String(&remoteHost, "h", "host", "Remote host")
	flaggy.String(&username, "u", "user", "Username")
	flaggy.String(&orientation, "ori", "orientation", "Orientation default: ltr(Local to Remote)")
	flaggy.String(&privateKeyPath, "k", "key", "Private key path")
	flaggy.String(&localFile, "l", "local", "Local file")
	flaggy.String(&remoteFile, "r", "remote", "Remove file")
	flaggy.Parse()

	if remoteHost == "" {
		flaggy.ShowHelpAndExit("Remote host is required")
	}

	if username == "" {
		flaggy.ShowHelpAndExit("Username is required")
	}

	if orientation != ltr && orientation != rtl {
		orientation = ltr
	}

	if privateKeyPath == "" {
		flaggy.ShowHelpAndExit("Private key path is required")
	}

	if localFile == "" {
		flaggy.ShowHelpAndExit("Left file is required")
	}

	if remoteFile == "" {
		flaggy.ShowHelpAndExit("Right file is required")
	}

	readPassword("Enter password: ", &passphrase)
	fmt.Println("")
}

func main() {
	clientConfig, err := auth.PrivateKeyWithPassphrase(username, []byte(passphrase), privateKeyPath, ssh.InsecureIgnoreHostKey())
	if err != nil {
		fmt.Println("Couldn't authenticate with a password protected private key", err)
		return
	}

	client := scp.NewClient(remoteHost, &clientConfig)
	client.RemoteBinary = "/usr/bin/scp"

	err = client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	defer client.Close()

	ctx := context.Background()

	var lfile *os.File

	switch orientation {
	case ltr:
		lfile, err = os.Open(localFile)
		if err != nil {
			fmt.Println("Couldn't open local file ", err)
			return
		}
		err = client.CopyFromFile(ctx, *lfile, remoteFile, "0644")
	case rtl:
		lfile, err = os.OpenFile(localFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			fmt.Println("Couldn't create local file ", err)
			return
		}
		err = client.CopyFromRemote(ctx, lfile, remoteFile)
	}
	lfile.Close()

	if err != nil {
		os.Remove(localFile)
		fmt.Println("Error while copying file ", err)
		return
	}

	fmt.Println("File copied successfully")
}
