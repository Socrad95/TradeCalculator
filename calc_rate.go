package main

import (
	"fmt"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CalculateDailyRate(targetReturn float64, days int) float64 {
	// (1 + x)^days = 1 + (targetReturn/100)
	// x = (1 + targetReturn/100)^(1/days) - 1
	return ((math.Pow(1+(targetReturn/100), 1/float64(days)) - 1) * 100)
}

func createCompoundTab() *fyne.Container {
	// 일일 수익률 계산 섹션
	targetReturnEntry := widget.NewEntry()
	targetReturnEntry.SetPlaceHolder("목표 수익률 (%)")
	daysEntry := widget.NewEntry()
	daysEntry.SetPlaceHolder("투자 기간 (일)")
	dailyRateResult := widget.NewLabel("")

	updateDailyRate := func() {
		targetReturn, err1 := strconv.ParseFloat(targetReturnEntry.Text, 64)
		days, err2 := strconv.Atoi(daysEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		dailyRate := CalculateDailyRate(targetReturn, days)
		dailyRateResult.SetText(fmt.Sprintf("필요 일일 수익률: %.4f%%", dailyRate))
	}

	targetReturnEntry.OnChanged = func(string) { updateDailyRate() }
	daysEntry.OnChanged = func(string) { updateDailyRate() }

	// 최종 수익률 계산 섹션
	dailyRateEntry := widget.NewEntry()
	dailyRateEntry.SetPlaceHolder("일일 수익률 (%)")
	periodEntry := widget.NewEntry()
	periodEntry.SetPlaceHolder("투자 기간 (일)")
	totalReturnResult := widget.NewLabel("")

	updateTotalReturn := func() {
		dailyRate, err1 := strconv.ParseFloat(dailyRateEntry.Text, 64)
		period, err2 := strconv.Atoi(periodEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		// (1 + 일일수익률)^기간 - 1 = 총 수익률
		totalReturn := (math.Pow(1+(dailyRate/100), float64(period)) - 1) * 100
		totalReturnResult.SetText(fmt.Sprintf("최종 수익률: %.2f%%", totalReturn))
	}

	dailyRateEntry.OnChanged = func(string) { updateTotalReturn() }
	periodEntry.OnChanged = func(string) { updateTotalReturn() }

	// 직접 계산 섹션 추가
	directInvestmentEntry := widget.NewEntry()
	directInvestmentEntry.SetPlaceHolder("투자금")
	directReturnEntry := widget.NewEntry()
	directReturnEntry.SetPlaceHolder("목표 수익률 (%)")
	directProfitResult := widget.NewLabel("")

	calculateDirectProfit := func() {
		investment, err1 := strconv.ParseFloat(directInvestmentEntry.Text, 64)
		returnRate, err2 := strconv.ParseFloat(directReturnEntry.Text, 64)

		if err1 != nil || err2 != nil {
			return
		}

		profit := investment * (returnRate / 100)
		totalAmount := investment + profit

		directProfitResult.SetText(fmt.Sprintf("예상 수익금: %.0f원\n최종 금액: %.0f원", profit, totalAmount))
	}

	// 자동 계산을 위한 이벤트 핸들러
	directInvestmentEntry.OnChanged = func(string) { calculateDirectProfit() }
	directReturnEntry.OnChanged = func(string) { calculateDirectProfit() }

	// UI 배치 수정
	content := container.NewVBox(
		targetReturnEntry,
		daysEntry,
		dailyRateResult,
		dailyRateEntry,
		periodEntry,
		totalReturnResult,
		widget.NewLabel("\n직접 계산:"),
		directInvestmentEntry,
		directReturnEntry,
		directProfitResult,
	)

	return content
}
