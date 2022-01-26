package main

import (
  "encoding/json"
	"log"
	"net/http"
	"os/exec"
)

var REPO_PATH = "/Users/larsbackman/Documents/Sandbox/test"
var EXPECTED_ID = "asd123"
var EXPECTED_SECRET = "SUPERPOJKEN"

type webhook struct {
    Action string
    Repository struct {
        ID string
        FullName string
    }
}

type response struct {
    Msg string
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request received")
    id, secret, provided := r.BasicAuth()

    // Allways respond as json
    w.Header().Set("Content-Type", "application/json")

    if !provided || id != EXPECTED_ID || secret != EXPECTED_SECRET {
        log.Printf("Request unauthenticated!")

        msg, _ := json.Marshal(response{
            Msg: "Unauthenticated",
        })

        http.Error(w, string(msg), http.StatusForbidden)
        return
    }

    err := triggerPull()

    if err != nil {
        msg, _ := json.Marshal(response{
            Msg: "Git pull failed. check server logs",
        })

        http.Error(w, string(msg), http.StatusInternalServerError)
        return
    }

    msg, _ := json.Marshal(response{
        Msg: "Git pull successful",
    })

    http.Error(w, string(msg), http.StatusOK)
    return
}

func triggerPull() (error) {
    cmd := exec.Command("git", "pull", "-r")
    // Sett CWD
    cmd.Dir = REPO_PATH
    log.Printf("Running command and waiting for it to finish...")
    out, err := cmd.Output()
    log.Printf("Command finished with error: %v", err)
    log.Printf("Output: %s", out)
    return err
}

func main() {
    log.Println("server started")
    http.HandleFunc("/webhook", handleWebhook)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

