package main

import (
	"flag"
	uuid2 "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/memberlist"
	"golang.org/x/net/websocket"
	"gossip-chat/chat"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type TemplateHttp struct {
	HttpPort string
}

func main() {

	hostname, _ := os.Hostname()
	uuid, _ := uuid2.GenerateUUID()
	name := flag.String("name", hostname+"_"+uuid, "name in the cluster")
	addr := flag.String("addr", "0.0.0.0", "ip to bind")
	port := flag.Int("gossip_port", 7946, "port to bind")
	httpPort := flag.Int("http_port", 8080, "port to bind")
	httpAddr := flag.String("http_addr", "0.0.0.0", "ip to bind")
	clusterAddr := flag.String("cluster", "localhost:7946", "cluster to connect to separated by comma ex : localhost:7946,localhost:7947")

	flag.Parse()

	messages := make(chan string, 1000)
	defaultConfig := memberlist.DefaultLocalConfig()
	defaultConfig.Delegate = &chat.Delegate{Messages: &messages}
	defaultConfig.Name = *name
	defaultConfig.BindAddr = *addr
	defaultConfig.BindPort = *port
	defaultConfig.AdvertisePort = *port
	chat.Gossip(defaultConfig, clusterAddr)
	chatStruct:=chat.InitChat(&messages)
	http.Handle("/chatws", websocket.Handler(chatStruct.Loop))
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	httpPortStr := strconv.Itoa(*httpPort)
	templateHttp := TemplateHttp{HttpPort: httpPortStr}
	http.HandleFunc("/", templateHttp.serveTemplate)

	httpBindAddr := *httpAddr + ":" + httpPortStr
	log.Println("bind http Addr: ", httpBindAddr)
	if err := http.ListenAndServe(httpBindAddr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
func (templateHttp TemplateHttp) serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("template", "index.html")
	tmpl := template.Must(template.ParseFiles(lp))
	err := tmpl.ExecuteTemplate(w, "js", "ws://localhost:"+templateHttp.HttpPort+"/chatws")
	if err != nil {
		log.Fatalln("Error executing template :", err)

	}
}
