package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	if err := loadSettings(); err != nil {
		settings = TradeSettings{
			StopLossRate:      -5.0,
			CommissionRate:    0.015,
			TargetProfitRates: []float64{5.0, 10.0, 15.0},
			IsDarkMode:        true,
		}
	}

	myApp := app.New()
	initialTheme := &myTheme{dark: settings.IsDarkMode}
	myApp.Settings().SetTheme(initialTheme)

	myWindow := myApp.NewWindow("주식 거래 계산기")

	buyTab := createBuyTab()
	sellTab := createSellTab()
	averageDownTab := createAverageDownTab()
	compoundTab := createCompoundTab()
	settingsTab := createSettingsTab(myApp)

	tabs := container.NewAppTabs(
		container.NewTabItem("매수", buyTab),
		container.NewTabItem("매도", sellTab),
		container.NewTabItem("물타기", averageDownTab),
		container.NewTabItem("복리계산", compoundTab),
		container.NewTabItem("설정", settingsTab),
	)

	// 탭 변경 시 매수내역 리스트 새로고침
	tabs.OnSelected = func(tab *container.TabItem) {
		switch tab.Text {
		case "매수":
			if list := findTradesList(buyTab); list != nil {
				list.Refresh()
			}
		case "매도":
			if list := findTradesList(sellTab); list != nil {
				list.Refresh()
			}
		case "물타기":
			if list := findTradesList(averageDownTab); list != nil {
				list.Refresh()
			}
		}
	}

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(800, 500))
	myWindow.ShowAndRun()
}

// 컨테이너 내에서 trades 리스트 위젯을 찾는 헬퍼 함수
func findTradesList(c *fyne.Container) *widget.List {
	for _, obj := range c.Objects {
		if card, ok := obj.(*widget.Card); ok {
			if scroll, ok := card.Content.(*container.Scroll); ok {
				if list, ok := scroll.Content.(*widget.List); ok {
					return list
				}
			}
		}
	}
	return nil
}
