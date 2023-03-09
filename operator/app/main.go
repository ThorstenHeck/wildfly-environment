
// main.go
package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "io"
    "net/http"
	"os"
    "github.com/gorilla/mux"
    "time"
    "os/exec"
    "net"
    "bytes"

)

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}

	fmt.Printf("File upload in progress: %d\n", pr.BytesRead)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	
    if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

    buff := make([]byte, 512)
    _, err = file.Read(buff)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    filetype := http.DetectContentType(buff)
    if filetype != "application/zip" {
        http.Error(w, "The provided file format is not allowed. Please upload a ZIP file", http.StatusBadRequest)
        return
    }

    err = os.MkdirAll("./upload", os.ModePerm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    f, err := os.Create(fmt.Sprintf("./upload/%s", handler.Filename))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    defer f.Close()

    pr := &Progress{
        TotalSize: handler.Size,
    }

    _, err = io.Copy(f, io.TeeReader(file, pr))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "Upload successful")
}

func checkdb() string {
    timeout := 1 * time.Second
    _, err := net.DialTimeout("tcp","10.5.0.4:22", timeout)
    if err != nil {
       dbflavour := "stack-oracle"
       return dbflavour
    }
    dbflavour := "stack-postgres"
    return dbflavour
}

func replace_host(hostgroup string) {

    input, err := ioutil.ReadFile("/app/ansible/playbooks/deploy.yml")  
    if err != nil {
            fmt.Println(err)
            os.Exit(1)
    }
    output := bytes.Replace(input, []byte("HOSTGROUP_PLACEHOLDER"), []byte(hostgroup), -1)
    if err = ioutil.WriteFile("/app/ansible/playbooks/deploy.yml", output, 0666); err != nil {  
            fmt.Println(err)
            os.Exit(1)
    }
}

func deploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

    dbbackend := checkdb()
    fmt.Fprintf(w, dbbackend, " is used as the database backend!")
    fmt.Printf(dbbackend)

    replace_host(dbbackend)

    prg := "ansible-playbook"

    arg1 := "-i"
    arg2 := "inventory/hosts"
    arg3 := "playbooks/deploy.yml"

    cmd := exec.Command(prg, arg1, arg2, arg3)
    cmd.Dir = "/app/ansible"
    cmd.Env = os.Environ()
    stdout, err := cmd.Output()

    fmt.Print(string(stdout))

    if err != nil {
        fmt.Println(err.Error())
        return
    }

}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/upload", uploadHandler).Methods("POST")
    myRouter.HandleFunc("/deploy", deploy).Methods("POST")
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    handleRequests()
}