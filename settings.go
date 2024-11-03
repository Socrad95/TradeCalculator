package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TradeSettings struct {
	StopLossRate      float64
	CommissionRate    float64
	TargetProfitRates []float64
	IsDarkMode        bool
}

var settings = TradeSettings{
	StopLossRate:      -5.0,
	CommissionRate:    0.015,
	TargetProfitRates: []float64{5.0, 10.0, 15.0},
	IsDarkMode:        true,
}

func saveSettings() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".trade_calculator")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "settings.json")
	file, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, file, 0644)
}

func loadSettings() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(homeDir, ".trade_calculator", "settings.json")
	file, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &settings)
}

func createSettingsTab(myApp fyne.App) *fyne.Container {
	stopLossEntry := widget.NewEntry()
	stopLossEntry.SetText(strconv.FormatFloat(settings.StopLossRate, 'f', 1, 64))
	stopLossEntry.OnChanged = func(value string) {
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			settings.StopLossRate = v
			saveSettings()
		}
	}

	commissionEntry := widget.NewEntry()
	commissionEntry.SetText(strconv.FormatFloat(settings.CommissionRate, 'f', 3, 64))
	commissionEntry.OnChanged = func(value string) {
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			settings.CommissionRate = v
			saveSettings()
		}
	}

	profitEntry := widget.NewEntry()
	profitEntry.SetPlaceHolder("수익률을 입력하세요 (예: 5.0)")

	profitTagsContainer := container.NewGridWrap(fyne.NewSize(80, 35))

	var updateProfitTags func()
	updateProfitTags = func() {
		profitTagsContainer.Objects = nil
		for i, rate := range settings.TargetProfitRates {
			rateStr := strconv.FormatFloat(rate, 'f', 1, 64)
			index := i

			tagButton := &widget.Button{
				Text:          rateStr + "%",
				Icon:          theme.CancelIcon(),
				IconPlacement: widget.ButtonIconTrailingText,
				OnTapped: func() {
					settings.TargetProfitRates = append(
						settings.TargetProfitRates[:index],
						settings.TargetProfitRates[index+1:]...,
					)
					saveSettings()
					updateProfitTags()
				},
			}
			tagButton.Importance = widget.LowImportance
			tagButton.Resize(fyne.NewSize(4, 4))

			profitTagsContainer.Add(container.NewPadded(tagButton))
		}
		profitTagsContainer.Refresh()
	}

	updateProfitTags()

	profitEntry.OnSubmitted = func(value string) {
		if rate, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			settings.TargetProfitRates = append(settings.TargetProfitRates, rate)
			sort.Float64s(settings.TargetProfitRates)
			saveSettings()
			updateProfitTags()
			profitEntry.SetText("")
		}
	}

	darkModeToggle := widget.NewCheck("다크 모드", func(checked bool) {
		settings.IsDarkMode = checked
		if checked {
			myApp.Settings().SetTheme(&myTheme{dark: true})
		} else {
			myApp.Settings().SetTheme(&myTheme{dark: false})
		}
		saveSettings()
	})

	currentTheme, ok := myApp.Settings().Theme().(*myTheme)
	if ok {
		settings.IsDarkMode = currentTheme.dark
	}
	darkModeToggle.Checked = settings.IsDarkMode
	darkModeToggle.Refresh()

	return container.NewVBox(
		widget.NewCard("손절 설정", "",
			container.NewVBox(
				widget.NewLabel("손절 손실율 (%)"),
				stopLossEntry,
			),
		),
		widget.NewCard("수수료 설정", "",
			container.NewVBox(
				widget.NewLabel("거래수수료율 (%)"),
				commissionEntry,
			),
		),
		widget.NewCard("목표수익 설정", "",
			container.NewVBox(
				widget.NewLabel("목표수익율 (%, 쉼표로 구분)"),
				profitEntry,
				profitTagsContainer,
			),
		),
		widget.NewCard("테마 설정", "",
			container.NewVBox(
				darkModeToggle,
			),
		),
	)
}
