package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createSellTab() *fyne.Container {
	var selectedTrade *Trade

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

	// 리스트의 최소 크기 설정
	tradesList.Resize(fyne.NewSize(350, 150))

	// 스크롤 가능한 컨테이너로 감싸기
	tradesScroll := container.NewVScroll(tradesList)
	tradesScroll.SetMinSize(fyne.NewSize(350, 150))

	sellPriceEntry := widget.NewEntry()
	sellPriceEntry.SetPlaceHolder("매도가격 (USD)")
	sellQuantityEntry := widget.NewEntry()
	sellQuantityEntry.SetPlaceHolder("매도수량")

	profitLabel := widget.NewLabel("계산할 매수내역을 선택해 주세요")
	profitRateLabel := widget.NewLabel("")
	sellCommissionLabel := widget.NewLabel("")

	tradesList.OnSelected = func(id widget.ListItemID) {
		selectedTrade = &trades[id]
		sellQuantityEntry.SetText(fmt.Sprintf("%d", selectedTrade.RemainingQty))

		// 선택 시 레이블 초기화 (updateCalculations가 곧바로 새 값을 설정할 것임)
		profitLabel.SetText("")
		profitRateLabel.SetText("")
		sellCommissionLabel.SetText("")
	}

	tradesList.OnUnselected = func(id widget.ListItemID) {
		selectedTrade = nil
		sellQuantityEntry.SetText("")
		sellPriceEntry.SetText("")

		// 선택 해제 시 안내 메시지로 되돌리기
		profitLabel.SetText("계산할 매수내역을 선택해 주세요")
		profitRateLabel.SetText("")
		sellCommissionLabel.SetText("")
	}

	updateCalculations := func() {
		if selectedTrade == nil {
			return
		}

		sellPrice, err1 := strconv.ParseFloat(sellPriceEntry.Text, 64)
		sellQty, err2 := strconv.Atoi(sellQuantityEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		totalSellAmount := sellPrice * float64(sellQty)
		sellCommission := totalSellAmount * settings.CommissionRate / 100.0

		buyAmount := selectedTrade.BuyPrice * float64(sellQty)
		profit := totalSellAmount - buyAmount - sellCommission -
			(selectedTrade.Commission * float64(sellQty) / float64(selectedTrade.Quantity))

		profitRate := (profit / buyAmount) * 100

		profitLabel.SetText(fmt.Sprintf("예상 수익: $%.4f", profit))
		profitRateLabel.SetText(fmt.Sprintf("수익률: %.2f%%", profitRate))
		sellCommissionLabel.SetText(fmt.Sprintf("매도수수료: $%.4f", sellCommission))
	}

	sellPriceEntry.OnChanged = func(string) { updateCalculations() }
	sellQuantityEntry.OnChanged = func(string) { updateCalculations() }

	// 매도 버튼 추가 및 처리 함수
	sellButton := widget.NewButton("매도", func() {
		if selectedTrade == nil {
			return
		}

		_, err1 := strconv.ParseFloat(sellPriceEntry.Text, 64)
		sellQty, err2 := strconv.Atoi(sellQuantityEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		// 매도 수량만큼 차감
		newRemainingQty := selectedTrade.RemainingQty - sellQty

		// trades 슬라이스에서 해당 거래의 인덱스 찾기
		var tradeIndex int
		for i, trade := range trades {
			if trade.ID == selectedTrade.ID {
				tradeIndex = i
				break
			}
		}

		// updateTrade 함수 사용
		updateTrade(tradeIndex, newRemainingQty)

		// UI 초기화
		sellPriceEntry.SetText("")
		sellQuantityEntry.SetText("")
		selectedTrade = nil
		tradesList.UnselectAll()
	})

	// 매도 정보 컨테이너 생성
	sellContent := container.NewVBox(
		widget.NewLabel("매도가격 (USD)"),
		sellPriceEntry,
		widget.NewLabel("매도수량"),
		sellQuantityEntry,
		sellButton, // 매도 버튼 추가
		profitLabel,
		profitRateLabel,
		sellCommissionLabel,
	)
	sellContent.Resize(fyne.NewSize(400, sellContent.MinSize().Height))

	sellCard := widget.NewCard("매도 정보", "", sellContent)
	sellContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, sellCard.MinSize().Height)), sellCard)

	// 매수 내역 컨테이너 생성
	historyCard := widget.NewCard("매수 내역", "", tradesScroll)
	historyCard.Resize(fyne.NewSize(400, historyCard.MinSize().Height))
	historyContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, historyCard.MinSize().Height)), historyCard)

	// 전체 컨테이너를 수평 배치로 변경
	mainContainer := container.NewHBox(sellContainer, historyContainer)

	registerTradeUpdateCallback(func() {
		tradesList.Refresh()
	})

	return mainContainer
}
