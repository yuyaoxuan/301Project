package main

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// connectSFTP establishes a connection to the SFTP server
func connectSFTP(server, username, privateKeyPath string) (*ssh.Client, *sftp.Client, error) {
	// Load SSH private key
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	// SSH Config
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// Connect to SFTP Server
	client, err := ssh.Dial("tcp", server+":22", config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %v", err)
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("failed to create SFTP client: %v", err)
	}

	return client, sftpClient, nil
}

// downloadFile downloads a file from SFTP server to local path
func downloadFile(client *sftp.Client, remotePath, localPath string) error {
	srcFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	return nil
}

// moveRemoteFile moves a file on the remote SFTP server
func moveRemoteFile(client *sftp.Client, oldPath, newPath string) error {
	return client.Rename(oldPath, newPath)
}

// ensureRemoteDir ensures a directory exists on the remote server
func ensureRemoteDir(client *sftp.Client, path string) error {
	_, err := client.Stat(path)
	if err == nil {
		// Directory exists
		return nil
	}

	// Create directory if it doesn't exist
	return client.MkdirAll(path)
}

// listRemoteDir lists files and directories in the remote path
func listRemoteDir(client *sftp.Client, path string) ([]os.FileInfo, error) {
	return client.ReadDir(path)
}