package main

import (
  "encoding/json"
	"log"
	"net/http"
  "os"
	"os/exec"
  "path/filepath"

  "github.com/willroberts/minecraft-client"
)

var ID string
var SECRET string
var RCON string
var RCON_PORT string

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

    if !provided || id != ID || secret != SECRET {
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
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)

    cmd := exec.Command("git", "pull", "-r")

    // Sett CWD
    cmd.Dir = exPath
    log.Printf("Running command and waiting for it to finish...")
    out, err := cmd.Output()
    log.Printf("Command finished with error: %v", err)
    log.Printf("Output: %s", out)
    return err
}

func reloadServer() {
    client, err := minecraft.NewClient("127.0.0.1:" + RCON_PORT)

    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Send some commands.
    if err := client.Authenticate(RCON); err != nil {
        log.Fatal(err)
    }

    resp, err := client.SendCommand("whitelist reload")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf(resp.Body)
}

func main() {
    log.Println("server started")
    http.HandleFunc("/webhook", handleWebhook)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
