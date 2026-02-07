<file_path>
go systray\readme.md
</file_path>

<edit_description>
Create README.md with compilation and usage instructions
</edit_description>

# Go Systray Download Organizer

This is a Go application that organizes files in the Downloads folder by moving them into subfolders based on file extensions. It runs in the system tray and monitors the Downloads directory for new files.

## Prerequisites

- Go 1.25.6 or later installed.

## Compilation

To compile the program, navigate to the project directory and run:

```bash
go build -o download-organizer main.go
```

This will create an executable named `download-organizer` (or `download-organizer.exe` on Windows).

To hide the console window on Windows (so the application runs in the background without a visible console), build with the following flag:

```bash
go build -ldflags="-H windowsgui" -o download-organizer.exe main.go
```

## Running

After compilation, run the executable:

On Linux/Mac:
```bash
./download-organizer
```

On Windows:
```bash
download-organizer.exe
```

The application will start in the system tray. It will sort existing files in the Downloads folder and continue monitoring for new downloads.

## Features

- Automatically moves files to folders like Images, Documents, Archives, etc., based on extensions (still inside the Downloads Folder).
- System tray menu with options to manually sort or quit.
- Notifications on startup. (Going to add more)

## Dependencies

The project uses the following Go modules:
- github.com/getlantern/systray
- github.com/fsnotify/fsnotify
- github.com/gen2brain/beeep

Dependencies are managed via `go.mod`.
