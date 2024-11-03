package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var selectedTradeIndex = -1 // 추가된 전역 변수

func createAverageDownTab() *fyne.Container {
	// 입력 필드 생성
	newPriceEntry := widget.NewEntry()
	newPriceEntry.SetPlaceHolder("추가 매수가격 (USD)")
	newQuantityEntry := widget.NewEntry()
	newQuantityEntry.SetPlaceHolder("추가 수량")

	// 결과 표시 레이블
	avgPriceLabel := widget.NewLabel("계산할 매수내역을 선택해 주세요")
	totalQuantityLabel := widget.NewLabel("")
	totalAmountLabel := widget.NewLabel("")
	commissionLabel := widget.NewLabel("")

	// 매수 내역 리스트
	tradesList := widget.NewList(
		func() int { return len(trades) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			trade := trades[id]
			label := item.(*widget.Label)
			label.SetText(fmt.Sprintf("매수가: $%.2f, 수량: %d, 총액: $%.2f",
				trade.BuyPrice, trade.RemainingQty, trade.TotalAmount))
		},
	)
	registerTradeUpdateCallback(func() {
		tradesList.Refresh()
	})

	tradesScroll := container.NewVScroll(tradesList)
	tradesScroll.SetMinSize(fyne.NewSize(350, 150))

	// 계산 함수
	updateCalculations := func() {
		if selectedTradeIndex == -1 {
			return
		}

		selectedTrade := trades[selectedTradeIndex]
		newPrice, err1 := strconv.ParseFloat(newPriceEntry.Text, 64)
		newQuantity, err2 := strconv.Atoi(newQuantityEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		// 평균 매수가 계산
		totalQuantity := selectedTrade.RemainingQty + newQuantity
		totalAmount := (selectedTrade.BuyPrice * float64(selectedTrade.RemainingQty)) +
			(newPrice * float64(newQuantity))
		averagePrice := totalAmount / float64(totalQuantity)

		// 수수료 계산
		newCommission := newPrice * float64(newQuantity) * (settings.CommissionRate / 100.0)
		totalCommission := selectedTrade.Commission + newCommission

		// 결과 업데이트
		avgPriceLabel.SetText(fmt.Sprintf("평균 매수가: $%.2f", averagePrice))
		totalQuantityLabel.SetText(fmt.Sprintf("총 수량: %d", totalQuantity))
		totalAmountLabel.SetText(fmt.Sprintf("총 매수금액: $%.2f", totalAmount))
		commissionLabel.SetText(fmt.Sprintf("총 수수료: $%.2f", totalCommission))
	}

	// 이벤트 핸들러 설정
	newPriceEntry.OnChanged = func(string) { updateCalculations() }
	newQuantityEntry.OnChanged = func(string) { updateCalculations() }

	// 매수 버튼 추가
	averageDownButton := widget.NewButton("매수", func() {
		if selectedTradeIndex == -1 {
			return
		}

		newPrice, err1 := strconv.ParseFloat(newPriceEntry.Text, 64)
		newQuantity, err2 := strconv.Atoi(newQuantityEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		selectedTrade := &trades[selectedTradeIndex]

		// 새로운 평균가 계산
		totalQuantity := selectedTrade.RemainingQty + newQuantity
		totalAmount := (selectedTrade.BuyPrice * float64(selectedTrade.RemainingQty)) +
			(newPrice * float64(newQuantity))
		averagePrice := totalAmount / float64(totalQuantity)

		// 새로운 수수료 계산
		newCommission := newPrice * float64(newQuantity) * (settings.CommissionRate / 100.0)

		// 매수내역 업데이트
		selectedTrade.BuyPrice = averagePrice
		selectedTrade.RemainingQty = totalQuantity
		selectedTrade.Quantity = totalQuantity
		selectedTrade.TotalAmount = totalAmount
		selectedTrade.Commission += newCommission

		// updateTrade 함수 사용
		updateTrade(selectedTradeIndex, totalQuantity)

		// UI 초기화
		newPriceEntry.SetText("")
		newQuantityEntry.SetText("")
		selectedTradeIndex = -1
		tradesList.UnselectAll()
	})

	tradesList.OnSelected = func(id widget.ListItemID) {
		selectedTradeIndex = int(id)
		newPriceEntry.SetText("")
		newQuantityEntry.SetText("")

		// 선택 시 레이블 초기화
		avgPriceLabel.SetText("")
		totalQuantityLabel.SetText("")
		totalAmountLabel.SetText("")
		commissionLabel.SetText("")

		updateCalculations()
	}

	tradesList.OnUnselected = func(id widget.ListItemID) {
		selectedTradeIndex = -1
		newPriceEntry.SetText("")
		newQuantityEntry.SetText("")

		// 선택 해제 시 안내 메시지로 되돌리기
		avgPriceLabel.SetText("계산할 매수내역을 선택해 주세요")
		totalQuantityLabel.SetText("")
		totalAmountLabel.SetText("")
		commissionLabel.SetText("")
	}

	// 물타기 정보 컨테이너
	averageDownContent := container.NewVBox(
		widget.NewLabel("추가 매수가격 (USD)"),
		newPriceEntry,
		widget.NewLabel("추가 수량"),
		newQuantityEntry,
		averageDownButton,
		widget.NewSeparator(),
		avgPriceLabel,
		totalQuantityLabel,
		totalAmountLabel,
		commissionLabel,
	)
	averageDownContent.Resize(fyne.NewSize(400, averageDownContent.MinSize().Height))

	averageDownCard := widget.NewCard("물타기 정보", "", averageDownContent)
	averageDownContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, averageDownCard.MinSize().Height)), averageDownCard)

	// 매수 내역 컨테이너
	historyCard := widget.NewCard("매수 내역", "", tradesScroll)
	historyCard.Resize(fyne.NewSize(400, historyCard.MinSize().Height))
	historyContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, historyCard.MinSize().Height)), historyCard)

	// 전체 컨테이너를 수평 배치
	mainContainer := container.NewHBox(averageDownContainer, historyContainer)

	return mainContainer
}
