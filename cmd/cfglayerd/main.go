package main

import (
    "encoding/json"
    "flag"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/LYH2263/go-cfglayer"
)

func main() {
    addr := flag.String("addr", ":8241", "listen")
    node := flag.String("node", "local", "node id")
    flag.Parse()

    m, err := cfglayer.New(cfglayer.Options{NodeID: *node, AuditPath: "audit"})
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    web := "web"
    if _, err := os.Stat(web); err != nil {
        web = filepath.Join("..", "..", "web")
    }

    mux := http.NewServeMux()
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(web))))
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        http.ServeFile(w, r, filepath.Join(web, "index.html"))
    })
    mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, m.Health())
    })
    mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, m.Stats())
    })
    mux.HandleFunc("/api/layers", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            writeJSON(w, m.ListLayers())
        case http.MethodPost:
            body, err := io.ReadAll(r.Body)
            if err != nil {
                http.Error(w, err.Error(), 400)
                return
            }
            var layer cfglayer.Layer
            if err := json.Unmarshal(body, &layer); err != nil {
                http.Error(w, err.Error(), 400)
                return
            }
            if err := m.PushLayer(r.Context(), layer); err != nil {
                http.Error(w, err.Error(), 500)
                return
            }
            writeJSON(w, map[string]string{"status": "pushed", "id": layer.ID})
        case http.MethodDelete:
            layer, err := m.PopLayer()
            if err != nil {
                http.Error(w, err.Error(), 400)
                return
            }
            writeJSON(w, layer)
        default:
            http.Error(w, "method", http.StatusMethodNotAllowed)
        }
    })
    mux.HandleFunc("/api/merge", func(w http.ResponseWriter, r *http.Request) {
        merged, err := m.MergeStack(r.Context())
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        writeJSON(w, merged)
    })
    mux.HandleFunc("/api/resolve/", func(w http.ResponseWriter, r *http.Request) {
        key := strings.TrimPrefix(r.URL.Path, "/api/resolve/")
        if key == "" {
            http.NotFound(w, r)
            return
        }
        v, err := m.Resolve(r.Context(), key)
        if err != nil {
            http.Error(w, err.Error(), 404)
            return
        }
        writeJSON(w, map[string]string{"key": key, "value": v})
    })
    mux.HandleFunc("/api/explain/", func(w http.ResponseWriter, r *http.Request) {
        key := strings.TrimPrefix(r.URL.Path, "/api/explain/")
        if key == "" {
            http.NotFound(w, r)
            return
        }
        steps, err := m.ExplainKey(key)
        if err != nil {
            http.Error(w, err.Error(), 404)
            return
        }
        writeJSON(w, steps)
    })
    mux.HandleFunc("/api/export", func(w http.ResponseWriter, r *http.Request) {
        snap, err := m.ExportMerged()
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        writeJSON(w, snap)
    })

    log.Printf("cfglayerd %s", *addr)
    log.Fatal(http.ListenAndServe(*addr, mux))
}

func writeJSON(w http.ResponseWriter, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(v)
}
