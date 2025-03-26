package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MiracleOS-Team/libxdg-go/notificationDaemon"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Function to create a single notification box
func createNotification(notification *notificationDaemon.Notification, nDaemon *notificationDaemon.Daemon) *gtk.Box {
	notificationBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 15)
	sc, _ := notificationBox.GetStyleContext()
	sc.AddClass("ntf_main_div")

	// Create top bar with app icon, title, and close button
	ntfTopBar, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)
	sc, _ = ntfTopBar.GetStyleContext()
	sc.AddClass("ntf_top_bar")

	ntfTopBarText, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)
	sc, _ = ntfTopBarText.GetStyleContext()
	sc.AddClass("nf_topbar_text")

	ntfTopBarImage, _ := gtk.ImageNewFromIconName(notification.AppIcon, gtk.ICON_SIZE_LARGE_TOOLBAR)
	ntfTopBarTextLabel, _ := gtk.LabelNew(notification.AppName)

	ntfTopBarDeleteButton, _ := gtk.ButtonNewWithLabel("✖")
	sc, _ = ntfTopBarDeleteButton.GetStyleContext()
	sc.AddClass("button")

	ntfTopBarDeleteButton.Connect("clicked", func() {
		nDaemon.CloseNotificationAsUser(notification.ID)
	})

	ntfTopBarText.PackStart(ntfTopBarImage, false, false, 0)
	ntfTopBarText.PackStart(ntfTopBarTextLabel, false, false, 0)
	ntfTopBar.PackStart(ntfTopBarText, false, false, 0)
	ntfTopBar.PackEnd(ntfTopBarDeleteButton, false, false, 0)
	notificationBox.PackStart(ntfTopBar, false, false, 0)

	// Notification content
	notificationContent, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	sc, _ = notificationContent.GetStyleContext()
	sc.AddClass("ntf_text_contents")

	if notification.Summary != "" {
		notificationSummary, _ := gtk.LabelNew(notification.Summary)
		notificationSummary.SetXAlign(0)
		sc, _ = notificationSummary.GetStyleContext()
		sc.AddClass("h2")
		notificationContent.PackStart(notificationSummary, false, false, 0)
	}

	if notification.Body != "" {
		notificationBody, _ := gtk.LabelNew(notification.Body)
		notificationBody.SetXAlign(0)
		notificationContent.PackStart(notificationBody, false, false, 0)
	}

	hours, minutes, _ := notification.Timestamp.Clock()
	timeLabel, _ := gtk.LabelNew(fmt.Sprintf("%d:%02d", hours, minutes))
	timeLabel.SetXAlign(1)
	sc, _ = timeLabel.GetStyleContext()
	sc.AddClass("h4")
	notificationContent.PackEnd(timeLabel, false, false, 0)

	notificationBox.PackStart(notificationContent, false, false, 0)

	return notificationBox
}

// Function to create the title bar of the notification panel
func createNotificationBarTitle(nDaemon *notificationDaemon.Daemon) *gtk.Box {
	tBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	title, _ := gtk.LabelNew(fmt.Sprintf("%d Notifications", len(nDaemon.Notifications)))
	sc, _ := title.GetStyleContext()
	sc.AddClass("h1")

	// Auto-update notification count
	glib.TimeoutAdd(uint(1000), func() bool {
		title.SetText(fmt.Sprintf("%d Notifications", len(nDaemon.Notifications)))
		return true
	})

	closeAllButton, _ := gtk.ButtonNewWithLabel("Clear all")
	sc, _ = closeAllButton.GetStyleContext()
	sc.AddClass("button")

	closeAllButton.Connect("clicked", func() {
		for _, elem := range nDaemon.Notifications {
			nDaemon.CloseNotificationAsUser(elem.ID)
		}
	})

	tBox.PackStart(title, false, false, 0)
	tBox.PackEnd(closeAllButton, false, false, 0)
	return tBox
}

// Function to create the notification panel
func createNotificationBar(nDaemon *notificationDaemon.Daemon) *gtk.Window {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("Notification Bar")
	win.SetDecorated(false)
	win.SetResizable(false)

	// Setup window as a top-layer shell (like KDE Plasma/GNOME Shell)
	layershell.InitForWindow(win)
	layershell.SetNamespace(win, "miracleos")
	layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_TOP)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_TOP, 10)

	mBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	mBox.PackStart(createNotificationBarTitle(nDaemon), false, false, 0)

	// Populate notifications
	for _, nt := range nDaemon.Notifications {
		mBox.PackStart(createNotification(&nt, nDaemon), false, false, 0)
	}

	win.Add(mBox)
	win.ShowAll()
	return win
}

// Function to initialize and listen for notifications
def listenNotifications() *notificationDaemon.Daemon {
	daemon := notificationDaemon.NewDaemon(notificationDaemon.Config{
		Capabilities: []string{"body", "actions", "actions-ions", "icon-static"},
	})
	if err := daemon.Start(); err != nil {
		log.Fatalf("Failed to start notification daemon: %v", err)
	}

	return daemon
}
