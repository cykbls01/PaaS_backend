package util

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// ssh连接
func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth   []ssh.AuthMethod
		addr   string
		clientConfig *ssh.ClientConfig
		sshClient *ssh.Client
		sftpClient *sftp.Client
		err   error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		Timeout: 30 * time.Second,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

// 传入文件和地址存储文件
func setFile(data []byte, remoteDir string){
	var (
		err  error
		sftpClient *sftp.Client
	)

	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = connect("root", "rootpass", "127.0.0.1", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(remoteDir)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()
	dstFile.Write(data)

    log.Println("copy file to remote server finished!")

}

func getFile(remoteFilePath string)*sftp.File{
	var (
		err  error
		sftpClient *sftp.Client
	)

	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = connect("root", "rootpass", "127.0.0.1", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()


    log.Println("copy file from remote server finished!")
	return srcFile
}

func ReadFile(path string)[]byte{
	var data  []byte
	var err error
	if data, err = ioutil.ReadFile(path); err != nil {
		fmt.Println(err)
	}
	return data
}

func WriteFile(path string,b []byte){
	os.Mkdir(strings.Split(path,"/")[0],os.ModePerm)
	err := ioutil.WriteFile(path,b,777)
	log.Println(err)
}
//func ReadFile(path string)[]byte{
//	var data  []byte
//	dstFile := getFile(path)
//	dstFile.Read(data)
//	return data
//}
//
//
//
//func WriteFile(path string,b []byte){
//	setFile(b,path)
//}