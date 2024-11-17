package app

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2"
	"image/color"
)

type myTheme struct{}

func (m myTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameDisabled:
		return color.RGBA{R: 120, G: 120, B: 120, A: 255}
	case theme.ColorNameForeground:
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	case theme.ColorNameBackground:
		return color.RGBA{R: 250, G: 250, B: 250, A: 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255}
	case theme.ColorNameButton:
		return color.RGBA{R: 220, G: 220, B: 220, A: 255}
	case theme.ColorNameHover:
		return color.RGBA{R: 180, G: 180, B: 180, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
