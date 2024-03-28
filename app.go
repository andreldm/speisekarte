package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

type Category struct {
    Key string `json:"key"`
    Title string `json:"title"`
}

type Service struct {
    Key string `json:"key"`
    Category string `json:"category"`
    Name string `json:"name"`
    Description string `json:"description"`
    Url string `json:"url"`
    OverrideHostname bool `json:"overrideHostname"`
}

type Config struct {
    Port int `json:"-"`
    Timeout int `json:"-"`
    Title string `json:"title"`
    Compact bool `json:"compact"`
    Categories []Category `json:"categories"`
    Services []Service `json:"services"`
    ServiceMap map[string]Service `json:"-"`
}

func loadConfig() (*Config, error) {
    configPath := os.Getenv("SPEISEKARTE_CONFIG")
    if configPath == "" {
        configPath = "./config.json"
    }

    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    config := Config{
        Port: 3000,
        Timeout: 10000,
        Title: "Speisekarte",
        Compact: false,
    }

    err = json.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    config.ServiceMap = make(map[string]Service)
    for _, v := range config.Services {
        config.ServiceMap[v.Key] = v
    }

    return &config, nil
}

func main() {
    config, err := loadConfig()
    if err != nil {
        panic(err)
    }

    configMarshaled, err := json.Marshal(config)
    if err != nil {
        panic(err)
    }

    client := http.Client{
        Timeout: time.Duration(config.Timeout) * time.Millisecond,
    }

    mux := http.NewServeMux()

    mux.HandleFunc("GET /", func(rw http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(rw, r)
            return
        }

        http.ServeFile(rw, r, "index.html")
    })

    mux.HandleFunc("GET /config", func(rw http.ResponseWriter, _ *http.Request) {
        rw.Header().Add("Content-Type", "application/json")
        rw.Write(configMarshaled)
    })

    mux.HandleFunc("GET /service/{key}/info", func(rw http.ResponseWriter, r *http.Request) {
        key := r.PathValue("key")

        service, exists := config.ServiceMap[key]
        if (!exists) {
            http.NotFound(rw, r)
            return
        }

        rw.Header().Add("Content-Type", "application/json")

        resp, err := client.Get(service.Url)
        if err != nil || resp.StatusCode >= 500 {
            rw.Write([]byte("{ \"status\": \"offline\" }"))
        } else {
            rw.Write([]byte("{ \"status\": \"online\" }"))
        }
    })

    fmt.Printf("Speisekarte is running at http://localhost:%d\n", config.Port)
    http.ListenAndServe(fmt.Sprint(":", config.Port), mux)
}
