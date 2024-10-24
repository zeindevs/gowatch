package main

import (
	"flag"
	"io/fs"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	command := flag.String("cmd", "make tailwind", "Command name")
	folder := flag.String("dir", "views", "Folder path")
	extension := flag.String("ext", ".html", "File extension")

	flag.Parse()

	commands := strings.Split(*command, " ")
	folders := []string{}

	if err := filepath.Walk(*folder, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			folders = append(folders, path)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
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
		if err := watcher.Add(folder); err != nil {
			log.Fatal(err)
		}
		log.Println("watching", folder)
	}

	log.Printf("gowatch running with %s file extension\n", *extension)

	<-make(chan struct{})
}
