package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

var extToDir = map[string]string{
	// Images
	".jpg":  "Images",
	".jpeg": "Images",
	".png":  "Images",
	".gif":  "Images",
	".bmp":  "Images",
	".tiff": "Images",
	".webp": "Images",

	// Documents
	".pdf":  "Documents",
	".doc":  "Documents",
	".docx": "Documents",
	".txt":  "Documents",
	".rtf":  "Documents",
	".odt":  "Documents",
	".xls":  "Documents",
	".xlsx": "Documents",
	".ppt":  "Documents",
	".pptx": "Documents",

	// Archives
	".zip": "Archives",
	".rar": "Archives",
	".7z":  "Archives",
	".tar": "Archives",
	".gz":  "Archives",
	".bz2": "Archives",

	// Executables/Installers
	".exe": "Executables",
	".msi": "Executables",
	".dmg": "Executables",
	".pkg": "Executables",
	".deb": "Executables",
	".rpm": "Executables",

	// Videos
	".mp4": "Videos",
	".avi": "Videos",
	".mkv": "Videos",
	".mov": "Videos",
	".wmv": "Videos",
	".flv": "Videos",

	// Music/Audio
	".mp3":  "Music",
	".wav":  "Music",
	".flac": "Music",
	".aac":  "Music",
	".ogg":  "Music",
	".wma":  "Music",

	// Web/Programming
	".html": "Web",
	".htm":  "Web",
	".css":  "Web",
	".js":   "Web",
	".php":  "Web",
	".py":   "Web",
	".go":   "Web",

	// Others (catch-all for unmapped extensions)
	".iso":     "Others",
	".torrent": "Others",
	".log":     "Others",
}

var downloadsDir string

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	downloadsDir = filepath.Join(homeDir, "Downloads")

	beeep.Notify("Download Organizer", "The file organizer is now running in the background.", "")
	sortExistingFiles(downloadsDir)
	go directoryChanges(downloadsDir)
	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetTitle("Download Organizer")
	sortItem := systray.AddMenuItem("Sort", "Sort Downloads")
	quitItem := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				systray.Quit()
				return
			case <-sortItem.ClickedCh:
				sortExistingFiles(downloadsDir)
			}
		}
	}()
}

func onExit() {
	fmt.Println("Exiting...")
}

func directoryChanges(downloadsDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(downloadsDir)
	if err != nil {
		fmt.Println("Error watching directory:", err)
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			fmt.Println("File event:", event)
			handleFileEvent(event, downloadsDir)
		case err := <-watcher.Errors:
			fmt.Println("Watcher error:", err)
			return
		}
	}
}

func handleFileEvent(event fsnotify.Event, downloadsDir string) {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return
	}
	fileExt := strings.ToLower(filepath.Ext(event.Name))
	targetFolderName, hasMapping := extToDir[fileExt]
	if !hasMapping {
		fmt.Printf("No folder for extension %s, skipping %s\n", fileExt, event.Name)
		return
	}
	targetFolderPath := filepath.Join(downloadsDir, targetFolderName)
	err := os.MkdirAll(targetFolderPath, 0755)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}
	fileName := filepath.Base(event.Name)
	targetFilePath := filepath.Join(targetFolderPath, fileName)
	err = os.Rename(event.Name, targetFilePath)
	if err != nil {
		fmt.Println("Error moving file:", err)
	} else {
		fmt.Printf("Moved %s to %s\n", event.Name, targetFilePath)
	}
}

func sortExistingFiles(downloadsDir string) {
	entries, err := os.ReadDir(downloadsDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(downloadsDir, entry.Name())
		fileExt := strings.ToLower(filepath.Ext(entry.Name()))
		targetFolderName, hasMapping := extToDir[fileExt]
		if hasMapping {
			targetFolderPath := filepath.Join(downloadsDir, targetFolderName)
			err := os.MkdirAll(targetFolderPath, 0755)
			if err != nil {
				fmt.Println("Error creating folder:", err)
				continue
			}
			targetFilePath := filepath.Join(targetFolderPath, entry.Name())
			err = os.Rename(filePath, targetFilePath)
			if err != nil {
				fmt.Println("Error moving existing file:", err)
			} else {
				fmt.Printf("Sorted existing file %s to %s\n", filePath, targetFilePath)
			}
		}
	}
}
