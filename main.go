package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	command := flag.String("cmd", "make tailwind", "Command name and arguments")
	folder := flag.String("dir", "views", "Folder for watching")
	extension := flag.String("ext", ".html", "File extension")

	flag.Parse()

	commands := strings.Split(*command, " ")
	folders := []string{*folder}

	entry, err := os.ReadDir(*folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range entry {
		if f.IsDir() {
			folders = append(folders, path.Join(*folder, f.Name()))
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if path.Ext(event.Name) == *extension {
						log.Println("modified file:", event.Name)
						cmd := exec.Command(commands[0], commands[1:]...)
						stdout, err := cmd.CombinedOutput()
						if err != nil {
							log.Println(err.Error())
						}
						log.Println(string(stdout))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, folder := range folders {
		err = watcher.Add(folder)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("watching", folder)
	}

	log.Printf("gowatch running with %s file extension\n", *extension)

	<-make(chan struct{})
}
