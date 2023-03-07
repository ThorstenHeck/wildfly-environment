
// main.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "io/ioutil"
    "io"
    "net/http"
	"os"
	"path/filepath"
    "github.com/gorilla/mux"
    "time"
    "os/exec"
)

// Article - Our struct for all articles
type Article struct {
    Id      string    `json:"Id"`
    Title   string `json:"Title"`
    Desc    string `json:"desc"`
    Content string `json:"content"`
}

var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    for _, article := range Articles {
        if article.Id == key {
            json.NewEncoder(w).Encode(article)
        }
    }
}


func createNewArticle(w http.ResponseWriter, r *http.Request) {
    // get the body of our POST request
    // unmarshal this into a new Article struct
    // append this to our Articles array.    
    reqBody, _ := ioutil.ReadAll(r.Body)
    var article Article 
    json.Unmarshal(reqBody, &article)
    // update our global Articles array to include
    // our new Article
    Articles = append(Articles, article)

    json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    for index, article := range Articles {
        if article.Id == id {
            Articles = append(Articles[:index], Articles[index+1:]...)
        }
    }

}

const MAX_UPLOAD_SIZE = 1024 * 1024 * 512 // 512 MB

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

	// get a reference to the fileHeaders
	files := r.MultipartForm.File["file"]


	for _, fileHeader := range files {    

        // _, err := os.Stat("./upload/" + fileHeader.Filename)
        // if err == nil {
        //     fmt.Fprintf(w, "Upload-successful")
        //     return
        // }


        if fileHeader.Size > MAX_UPLOAD_SIZE {
            http.Error(w, fmt.Sprintf("The uploaded zip file is too big: %s. Please use an file less than 512MB in size", fileHeader.Filename), http.StatusBadRequest)
            return
        }

        file, err := fileHeader.Open()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        defer file.Close()

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

        _, err = file.Seek(0, io.SeekStart)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = os.MkdirAll("./upload", os.ModePerm)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        f, err := os.Create(fmt.Sprintf("./upload/%s", fileHeader.Filename))
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        defer f.Close()

        pr := &Progress{
            TotalSize: fileHeader.Size,
        }

        _, err = io.Copy(f, io.TeeReader(file, pr))
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        fmt.Fprintf(w, "Upload successful")
    }
}


func FindLastModifiedFileBefore(dir string, t time.Time) (path string, info os.FileInfo, err error) {
    isFirst := true
    min := 0 * time.Second
    err = filepath.Walk(dir, func(p string, i os.FileInfo, e error) error {
            if e != nil {
                    return e
            }

            if !i.IsDir() && i.ModTime().Before(t) {
                    if isFirst {
                            isFirst = false
                            path = p
                            info = i
                            min = t.Sub(i.ModTime())
                    }
                    if diff := t.Sub(i.ModTime()); diff < min {
                            path = p
                            min = diff
                            info = i
                    }
            }
            return nil
    })
    return
}

func deploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

    dir := "./upload"
    path, info, err := FindLastModifiedFileBefore(dir, time.Now())
    if err != nil {
            fmt.Println(err)
            fmt.Println(info)
            fmt.Fprintf(w, "Please go back and upload a zip!")
            return 
    }
    fmt.Println("using " + path + " to deploy")

    prg := "ansible-playbook"

    arg1 := "-i"
    arg2 := "/app/environments/DEV/inventory"
    arg3 := "/app/ansible/playbooks/deploy.yml"

    cmd := exec.Command(prg, arg1, arg2, arg3)
    cmd.Env = os.Environ()
    cmd.Env = append(cmd.Env, "ANSIBLE_ROLES_PATH=/app/ansible/roles")
    cmd.Env = append(cmd.Env, "ansible_ssh_private_key_file=/home/ansible/.ssh/id_ed25519")
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Print(string(stdout))


}




func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/articles", returnAllArticles)
    myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
    myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
    myRouter.HandleFunc("/article/{id}", returnSingleArticle)
    myRouter.HandleFunc("/upload", uploadHandler).Methods("POST")
    myRouter.HandleFunc("/deploy", deploy).Methods("POST")
    log.Fatal(http.ListenAndServe(":10000", myRouter))

}

func main() {
    Articles = []Article{
        Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
        Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
    }
    handleRequests()
}