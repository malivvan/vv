package menu

import (
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
	"sync"
)

type MenuItem struct {
	*cui.Box
	Title    string
	SubItems []*MenuItem
	onClick  func(*MenuItem)
	sync.RWMutex
}

func NewMenuItem(title string) *MenuItem {
	return &MenuItem{
		Box:      cui.NewBox(),
		Title:    title,
		SubItems: make([]*MenuItem, 0),
	}
}

func (menuItem *MenuItem) AddItem(item *MenuItem) *MenuItem {
	menuItem.Lock()
	defer menuItem.Unlock()

	menuItem.SubItems = append(menuItem.SubItems, item)
	return menuItem
}

func (menuItem *MenuItem) SetOnClick(fn func(*MenuItem)) *MenuItem {
	menuItem.Lock()
	defer menuItem.Unlock()

	menuItem.onClick = fn
	return menuItem
}

func (menuItem *MenuItem) Draw(screen tcell.Screen) {
	if !menuItem.GetVisible() {
		return
	}

	menuItem.Box.Draw(screen)

	menuItem.Lock()
	defer menuItem.Unlock()

	x, y, _, _ := menuItem.GetInnerRect()

	cui.PrintSimple(screen, []byte(menuItem.Title), x, y)
}

type SubMenu struct {
	*cui.Box
	Items         []*MenuItem
	parent        *MenuBar
	childMenu     *SubMenu
	currentSelect int
	sync.RWMutex
}

func NewSubMenu(parent *MenuBar, items []*MenuItem) *SubMenu {
	subMenu := &SubMenu{
		Box:           cui.NewBox(),
		Items:         items,
		parent:        parent,
		currentSelect: -1,
	}
	subMenu.SetBorder(true)
	return subMenu
}

func (subMenu *SubMenu) Draw(screen tcell.Screen) {
	anySubItems := false
	maxWidth := 0
	for _, item := range subMenu.Items {
		if itemTitleLen := len(item.Title); itemTitleLen > maxWidth {
			maxWidth = itemTitleLen
		}
		if len(item.SubItems) > 0 {
			anySubItems = true
		}
	}

	rectX, rectY, _, _ := subMenu.GetRect()
	rectWid := maxWidth
	if anySubItems {
		rectWid += 1
	}
	rectHig := len(subMenu.Items)
	// +2 - add space one space for each side of rect - to fit text inside
	subMenu.SetRect(rectX, rectY, rectWid+2, rectHig+2)

	if !subMenu.GetVisible() {
		return
	}

	subMenu.Box.Draw(screen)

	subMenu.Lock()
	defer subMenu.Unlock()

	x, y, _, _ := subMenu.GetInnerRect()
	for i, item := range subMenu.Items {
		if i == subMenu.currentSelect {
			cui.Print(screen, []byte(item.Title), x, y+i, 20, 0, tcell.ColorBlue)
			if len(item.SubItems) > 0 {
				cui.Print(screen, []byte(">"), x+maxWidth, y+i, 20, 0, tcell.ColorBlue)
			}
			continue
		}
		cui.PrintSimple(screen, []byte(item.Title), x, y+i)
		if len(item.SubItems) > 0 {
			cui.PrintSimple(screen, []byte(">"), x+maxWidth, y+i)
		}
	}
	if subMenu.childMenu != nil {
		subMenu.childMenu.Draw(screen)
	}
}

func (subMenu *SubMenu) MouseHandler() func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
	return subMenu.WrapMouseHandler(func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
		if subMenu.childMenu != nil {
			consumed, capture = subMenu.childMenu.MouseHandler()(action, event, setFocus)

			if consumed {
				return
			}
		}
		rectX, rectY, rectW, _ := subMenu.Box.GetInnerRect()
		if !subMenu.Box.InRect(event.Position()) {
			// Close the menu if the user clicks outside the menu box
			if action == cui.MouseLeftClick {
				subMenu.parent.subMenu = nil
			}
			return false, nil
		}
		_, y := event.Position()
		index := y - rectY

		subMenu.currentSelect = index
		consumed = true

		if action == cui.MouseLeftClick {
			setFocus(subMenu)
			if index >= 0 && index < len(subMenu.Items) {
				handler := subMenu.Items[index].onClick
				if handler != nil {
					handler(subMenu.Items[index])
				}
				if len(subMenu.Items[index].SubItems) > 0 {
					subMenu.childMenu = NewSubMenu(subMenu.parent, subMenu.Items[index].SubItems)
					subMenu.childMenu.SetRect(rectX+rectW, y, 15, 10)
					return
				}
			}
			subMenu.parent.subMenu = nil
		}
		return
	})
}

type MenuBar struct {
	*cui.Box
	MenuItems     []*MenuItem
	subMenu       *SubMenu // sub menu if not nil will be drawn
	currentOption int
	sync.RWMutex
}

func NewMenuBar() *MenuBar {
	return &MenuBar{
		Box:       cui.NewBox(),
		MenuItems: make([]*MenuItem, 0),
	}
}

func (menuBar *MenuBar) AfterDraw() func(tcell.Screen) {
	return func(screen tcell.Screen) {
		if menuBar.subMenu != nil {
			menuBar.subMenu.Draw(screen)
		}
	}
}

func (menuBar *MenuBar) AddItem(item *MenuItem) *MenuBar {
	menuBar.Lock()
	defer menuBar.Unlock()

	menuBar.MenuItems = append(menuBar.MenuItems, item)
	return menuBar
}

func (menuBar *MenuBar) Draw(screen tcell.Screen) {
	if !menuBar.GetVisible() {
		return
	}

	menuBar.Box.Draw(screen)

	menuBar.Lock()
	defer menuBar.Unlock()

	x, y, width, _ := menuBar.GetInnerRect()

	for i := 0; i < width; i += 1 {
		screen.SetContent(x+i, y, ' ', nil, tcell.StyleDefault.Background(menuBar.Box.GetBackgroundColor()))
	}

	menuItemOffset := 1
	for _, mi := range menuBar.MenuItems {
		itemLen := len([]rune(mi.Title))
		mi.SetRect(menuItemOffset, y, itemLen, 1)
		mi.Draw(screen)
		menuItemOffset += itemLen + 1
	}
}

func (menuBar *MenuBar) InputHandler() func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
	return menuBar.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
		switch event.Key() {
		case tcell.KeyLeft:
			menuBar.currentOption--
			if menuBar.currentOption < 0 {
				menuBar.currentOption = -1
			}
		case tcell.KeyRight:
			menuBar.currentOption++
			if menuBar.currentOption >= len(menuBar.MenuItems) {
				menuBar.currentOption = len(menuBar.MenuItems) - 1
			}
		}
	})
}

func (menuBar *MenuBar) MouseHandler() func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
	return menuBar.WrapMouseHandler(func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
		if menuBar.subMenu != nil {
			consumed, capture = menuBar.subMenu.MouseHandler()(action, event, setFocus)
			if consumed {
				//p.subMenu = nil
				return
			}
		}
		if !menuBar.InRect(event.Position()) {
			return false, nil
		}
		// Pass mouse events down.
		for _, item := range menuBar.MenuItems {
			consumed, capture = item.MouseHandler()(action, event, setFocus)
			if consumed {
				menuBar.subMenu = NewSubMenu(menuBar, item.SubItems)
				x, y, _, _ := item.GetRect()
				if menuBar.GetBorder() {
					x++
				}
				menuBar.subMenu.Box.SetRect(x, y+1, 15, 10)
				return
			}
		}

		// ...handle mouse events not directed to the child primitive...
		return true, nil
	})
}

func (menuBar *MenuBar) Focus(delegate func(p cui.Primitive)) {
	//if menuBar.subMenu != nil {
	//	delegate(menuBar.subMenu)
	//} else {
	menuBar.Box.Focus(delegate)
	menuBar.subMenu = nil
	//}
}
