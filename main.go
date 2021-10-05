package main

import (
    "fmt"
    "log"
    "net/http"
	"os"
	"io"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	
)



func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")

	fileUrl := "https://rss.cluster.wenisch.tech/klenkes/events/ics"
	err := DownloadFile("calendar.ics", fileUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded: " + fileUrl)

user   := ""
	pass   := ""
	remote := ""
	port   := ":22"

	

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		 HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		
	}

	// connect
	conn, err := ssh.Dial("tcp", remote+port, config)
	if err != nil {
	   log.Fatal(err)
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
	   log.Fatal(err)
	}
	defer client.Close()

	// create destination file
	dstFile, err := client.Create("./file.txt")
	if err != nil {
	   log.Fatal(err)
	}
	defer dstFile.Close()

	// create source file
	srcFile, err := os.Open("./file.txt")
	if err != nil {
	   log.Fatal(err)
	}

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
	   log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)
}

func handleRequests() {
    http.HandleFunc("/", homePage)
    log.Fatal(http.ListenAndServe(":10000", nil))
}
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
    handleRequests()
}