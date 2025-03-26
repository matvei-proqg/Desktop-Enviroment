package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/gdk"
	"github.com/AuruTeam/desktoplib/wallpaper"
)

const desktopFileContent = `[Desktop Entry]
Version=1.0
Type=Application
Name=AuruOS UI Config
Exec=/opt/AuruTeam/desktop/ui-config
Icon=preferences-system
Terminal=false
Categories=Settings;DesktopSettings;`

// createDesktopEntry creates a desktop icon for easy access
func createDesktopEntry() {
	usr, err := user.Current()
	if err != nil {
		log.Println("Error getting user:", err)
		return
	}
	desktopFilePath := fmt.Sprintf("%s/Desktop/auruos-ui-config.desktop", usr.HomeDir)

	file, err := os.Create(desktopFilePath)
	if err != nil {
		log.Println("Error creating desktop file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(desktopFileContent)
	if err != nil {
		log.Println("Error writing to desktop file:", err)
	}

	if err := os.Chmod(desktopFilePath, 0755); err != nil {
		log.Println("Error setting permissions for desktop file:", err)
	}
}

// loadCSS loads CSS styles
func loadCSS(path string) {
	provider, err := gtk.CssProviderNew()
	if err != nil {
		log.Println("Error creating CssProvider:", err)
		return
	}

	err = provider.LoadFromPath(path)
	if err != nil {
		log.Println("Failed to load CSS:", err)
		return
	}

	display, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Println("Failed to get display:", err)
		return
	}

	screen, err := display.GetDefaultScreen()
	if err != nil {
		log.Println("Failed to get screen:", err)
		return
	}

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

// changeFont changes the font of the interface
func changeFont(fontFamily string, fontSize int) {
	css := fmt.Sprintf(`* {
		font-family: %s;
		font-size: %dpx;
	}`, fontFamily, fontSize)

	provider, err := gtk.CssProviderNew()
	if err != nil {
		log.Println("Error creating CssProvider:", err)
		return
	}

	err = provider.LoadFromData(css)
	if err != nil {
		log.Println("Failed to apply font:", err)
		return
	}

	display, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Println("Failed to get display:", err)
		return
	}

	screen, err := display.GetDefaultScreen()
	if err != nil {
		log.Println("Failed to get screen:", err)
		return
	}

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

// setWallpaper sets the desktop wallpaper
func setWallpaper(path string) {
	err := wallpaper.SetImageWallpaper(path, "")
	if err != nil {
		log.Println("Failed to set wallpaper:", err)
	}
}

// setTransparency adjusts window transparency
func setTransparency(transparency float64) {
	css := fmt.Sprintf(`* {
		background-color: rgba(0, 0, 0, %.2f);
	}`, transparency)

	provider, err := gtk.CssProviderNew()
	if err != nil {
		log.Println("Error creating CssProvider:", err)
		return
	}

	err = provider.LoadFromData(css)
	if err != nil {
		log.Println("Failed to apply transparency:", err)
		return
	}

	display, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Println("Failed to get display:", err)
		return
	}

	screen, err := display.GetDefaultScreen()
	if err != nil {
		log.Println("Failed to get screen:", err)
		return
	}

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

// createSettingsWindow creates the settings window with controls for customization
func createSettingsWindow() {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("AuruOS UI Config")
	win.SetDefaultSize(500, 400)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	win.Add(box)

	// Theme selection buttons
	btnDark, _ := gtk.ButtonNewWithLabel("Dark Theme")
	btnLight, _ := gtk.ButtonNewWithLabel("Light Theme")
	btnCustom, _ := gtk.ButtonNewWithLabel("Custom Theme")
	btnFont, _ := gtk.ButtonNewWithLabel("Change Font")
	btnWallpaper, _ := gtk.ButtonNewWithLabel("Change Wallpaper")
	btnTransparency, _ := gtk.ButtonNewWithLabel("Adjust Transparency")

	// Connect buttons to their respective handlers
	btnDark.Connect("clicked", func() { loadCSS("/opt/AuruTeam/desktop/dark.css") })
	btnLight.Connect("clicked", func() { loadCSS("/opt/AuruTeam/desktop/light.css") })
	btnCustom.Connect("clicked", func() { loadCSS("/opt/AuruTeam/desktop/custom.css") })
	btnFont.Connect("clicked", func() { changeFont("Arial", 14) })
	btnWallpaper.Connect("clicked", func() { setWallpaper("/usr/share/backgrounds/auruos_dark_default.jpg") })
	btnTransparency.Connect("clicked", func() { setTransparency(0.7) })

	// Add buttons to the layout
	box.PackStart(btnDark, false, false, 5)
	box.PackStart(btnLight, false, false, 5)
	box.PackStart(btnCustom, false, false, 5)
	box.PackStart(btnFont, false, false, 5)
	box.PackStart(btnWallpaper, false, false, 5)
	box.PackStart(btnTransparency, false, false, 5)

	win.ShowAll()
}

func main() {
	gtk.Init(&os.Args)

	// Create desktop entry for easy access
	createDesktopEntry()

	// Launch the settings window
	createSettingsWindow()

	gtk.Main()
}
