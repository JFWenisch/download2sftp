package main

import (
    "fmt"
    "log"
    "net/http"
	"os"
	"io"
	"io/ioutil"
	"encoding/json"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	
)


type Todo struct {
    Url      string    `json:"url"`
    Filename string      `json:"filename"`
}

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Successfully procesed!")
    fmt.Println("Endpoint Hit: homePage")

	//Start JSON decoding
	 var todo Todo
    body, errUnmarshal := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if errUnmarshal != nil {
        panic(errUnmarshal)
    }
    if errUnmarshal := r.Body.Close(); errUnmarshal != nil {
        panic(errUnmarshal)
    }
    if errUnmarshal := json.Unmarshal(body, &todo); errUnmarshal != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if errUnmarshal := json.NewEncoder(w).Encode(errUnmarshal); errUnmarshal != nil {
            panic(errUnmarshal)
        }
    }
 
    
	    fmt.Println("Fetching from " +todo.Url);
		fmt.Println("Downloading as " +todo.Filename);
	// END JSON decoding


	//START DOWNLOAD
	errDownload := DownloadFile(todo.Filename, todo.Url)
	if errDownload != nil {
		panic(errDownload)
	}
	//END DOWNLOAD
	fmt.Println("Sucessfully downloaded: " + todo.Url)

	user   := os.Getenv("D2SFTP_USER")
	pass   := os.Getenv("D2SFTP_PASS")
	remote := os.Getenv("D2SFTP_REMOTE")
	port   := os.Getenv("D2SFTP_PORT")

	

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		 HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		
	}

	// connect
	fmt.Println("Connecting to " +remote+":"+port +" as "+user);
	conn, err := ssh.Dial("tcp", remote+":"+port, config)
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
	dstFile, err := client.Create("./"+todo.Filename)
	if err != nil {
	   log.Fatal(err)
	}
	defer dstFile.Close()

	// create source file
	srcFile, err := os.Open(todo.Filename)
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