package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/AuruTeam/libxdg-go/desktopFiles"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// createAppGroup создает группу приложений
func createAppGroup(apps []desktopFiles.DesktopFile) *gtk.Box {
	group, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	for _, app := range apps {
		buttonBox, _ := gtk.ButtonNew()
		appBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
		appBox.GetStyleContext().AddClass("mm_applist_app")

		// Загрузка иконки приложения
		if pixbuf, err := gdk.PixbufNewFromFileAtScale(app.Icon, 16, 16, true); err == nil {
			icon, _ := gtk.ImageNewFromPixbuf(pixbuf)
			appBox.PackStart(icon, false, false, 5)
		}

		label, _ := gtk.LabelNew(app.Name)
		appBox.PackStart(label, false, false, 5)
		buttonBox.Add(appBox)
		buttonBox.Connect("clicked", func() {
			fmt.Println("Clicked on", app.Name)
			go desktopFiles.ExecuteDesktopFile(app, []string{}, "")
		})
		group.PackStart(buttonBox, false, false, 5)
	}
	return group
}

// createAppList создает список установленных приложений
func createAppList() *gtk.ScrolledWindow {
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	apps, _ := desktopFiles.ListAllApplications()
	validApps := make([]desktopFiles.DesktopFile, 0)

	// Отфильтровываем скрытые приложения
	for _, app := range apps {
		if !app.NoDisplay {
			validApps = append(validApps, app)
		}
	}

	// Группировка приложений по первой букве имени
	categories := make(map[string][]desktopFiles.DesktopFile)
	for _, app := range validApps {
		category := string([]rune(app.Name)[0]) // Получаем первую букву имени
		categories[category] = append(categories[category], app)
	}

	// Сортировка категорий
	sortedCategories := make([]string, 0, len(categories))
	for cat := range categories {
		sortedCategories = append(sortedCategories, cat)
	}
	sort.Strings(sortedCategories)

	// Создание интерфейса для категорий
	for _, category := range sortedCategories {
		label, _ := gtk.LabelNew(fmt.Sprintf("<b>%s</b>", category))
		label.SetUseMarkup(true)
		label.SetXAlign(0)
		vbox.PackStart(label, false, false, 5)

		// Сортировка приложений в категории
		sortedApps := categories[category]
		sort.Slice(sortedApps, func(i, j int) bool {
			return sortedApps[i].Name < sortedApps[j].Name
		})

		vbox.PackStart(createAppGroup(sortedApps), false, false, 5)
	}

	scroll.Add(vbox)
	return scroll
}

// createMainMenu создает главное окно меню
func createMainMenu() *gtk.Window {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("Main Menu")
	win.SetDefaultSize(600, 600)
	win.SetDecorated(false)
	win.SetResizable(false)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)

	// Настройка LayerShell
	layershell.InitForWindow(win)
	layershell.SetNamespace(win, "miracleos")
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, true)
	layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_OVERLAY)

	// Определение монитора
	if disp, err := gdk.DisplayGetDefault(); err == nil {
		if mon, err := disp.GetMonitor(0); err == nil {
			layershell.SetMonitor(win, mon)
		}
	}

	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	mainBox.GetStyleContext().AddClass("mm_menu_m2")

	// Верхняя панель с поиском
	topBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	topBox.GetStyleContext().AddClass("mm_toppart")

	searchEntry, _ := gtk.EntryNew()
	searchEntry.SetPlaceholderText("Search Anything")
	searchEntry.GetStyleContext().AddClass("mos-input")
	topBox.PackStart(searchEntry, true, true, 5)

	// Основная часть с вкладками
	contentBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	appList := createAppList()
	appList.SetSizeRequest(300, 600)

	fastApps := createPlaceholder("Most Used")
	fastApps.SetSizeRequest(300, 600)

	otherTab := createPlaceholder("Other")
	otherTab.SetSizeRequest(300, 600)

	contentBox.PackStart(appList, false, false, 10)
	contentBox.PackStart(fastApps, false, false, 10)
	contentBox.PackStart(otherTab, false, false, 10)

	mainBox.PackStart(topBox, false, false, 10)
	mainBox.PackStart(contentBox, true, true, 10)

	// Нижняя панель с пользователем и кнопками питания
	bottomBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	bottomBox.PackStart(createUserInfo(), false, false, 10)
	bottomBox.PackEnd(createPowerButtons(), false, false, 10)
	bottomBox.GetStyleContext().AddClass("mm_bottompart")

	mainBox.PackStart(bottomBox, false, false, 10)

	win.Add(mainBox)
	return win
}

func main() {
	gtk.Init(nil)
	win := createMainMenu()
	win.Connect("destroy", gtk.MainQuit)
	win.ShowAll()
	gtk.Main()
}
