package main

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
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
}

type Config struct {
    Port int `json:"port"`
    Timeout int `json:"timeout"`
    Title string `json:"title"`
    Compact bool `json:"compact"`
    Categories []Category `json:"categories"`
    Services []Service `json:"services"`
    ServiceMap map[string]Service
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

    client := http.Client{
        Timeout: time.Duration(config.Timeout) * time.Millisecond,
    }

    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.File("index.html")
    })

    r.GET("/config", func(c *gin.Context) {
        c.JSON(200, config)
    })

    r.GET("/service/:key/info", func(c *gin.Context) {
        key := c.Param("key")

        service, exists := config.ServiceMap[key]
        if (!exists) {
            c.Status(404)
            return
        }

        resp, err := client.Get(service.Url)
        if err != nil || resp.StatusCode >= 400 {
            c.JSON(200, gin.H{"status": "offline"})
        } else {
            c.JSON(200, gin.H{"status": "online"})
        }
    })

    fmt.Printf("Speisekarte ist ready zu serve at http://localhost:%d\n", config.Port)
    r.Run(fmt.Sprint(":", config.Port))
}
