package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type TradeCalculation struct {
	buyPrice       float64
	sellPrice      float64
	quantity       int
	totalBuyAmount float64
	buyCommission  float64
	sellCommission float64
}

func calculateTrade(buyPrice float64, sellPrice float64, quantity int) TradeCalculation {
	totalBuyAmount := buyPrice * float64(quantity)
	buyCommission := totalBuyAmount * (settings.CommissionRate / 100.0)

	totalSellAmount := sellPrice * float64(quantity)
	sellCommission := totalSellAmount * (settings.CommissionRate / 100.0)

	return TradeCalculation{
		buyPrice:       buyPrice,
		sellPrice:      sellPrice,
		quantity:       quantity,
		totalBuyAmount: totalBuyAmount,
		buyCommission:  buyCommission,
		sellCommission: sellCommission,
	}
}

func createBuyTab() *fyne.Container {
	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("매수가격 (USD)")
	quantityEntry := widget.NewEntry()
	quantityEntry.SetPlaceHolder("수량")

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

	buyButton := widget.NewButton("매수", func() {
		price, err1 := strconv.ParseFloat(priceEntry.Text, 64)
		quantity, err2 := strconv.Atoi(quantityEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		addTrade(price, quantity)
		tradesList.Refresh()

		// 입력 필드 초기화
		priceEntry.SetText("")
		quantityEntry.SetText("")
	})

	totalAmountLabel := widget.NewLabel("")
	buyCommissionLabel := widget.NewLabel("")
	stopLossLabel := widget.NewLabel("")
	profitTargetsContainer := container.NewVBox()

	updateCalculations := func() {
		buyPrice, err1 := strconv.ParseFloat(priceEntry.Text, 64)
		quantity, err2 := strconv.Atoi(quantityEntry.Text)

		if err1 != nil || err2 != nil {
			return
		}

		// 매수 정보 계산
		initialTrade := calculateTrade(buyPrice, buyPrice, quantity)
		totalAmountLabel.SetText(fmt.Sprintf("총 매수금액: $%.2f", initialTrade.totalBuyAmount))
		buyCommissionLabel.SetText(fmt.Sprintf("매수수수료: $%.4f", initialTrade.buyCommission))

		// 손소 매도가격 계산
		minSellPrice := buyPrice * (1 + settings.CommissionRate/100) / (1 - settings.CommissionRate/100)
		minSellLabel := widget.NewLabel(fmt.Sprintf("최소 매도가격: $%.4f (수수료 손익분기점)", minSellPrice))

		// 손절가 계산
		stopLossPrice := buyPrice * (1 + settings.StopLossRate/100)
		stopLossTrade := calculateTrade(buyPrice, stopLossPrice, quantity)
		stopLossProfit := (stopLossPrice * float64(quantity)) -
			initialTrade.totalBuyAmount - initialTrade.buyCommission - stopLossTrade.sellCommission

		stopLossLabel.SetText(fmt.Sprintf("손절가: $%.2f (손익: $%.4f, 매도수수료: $%.4f)",
			stopLossPrice, stopLossProfit, stopLossTrade.sellCommission))

		// 목표수익 계산
		profitTargetsContainer.Objects = nil
		profitTargetsContainer.Add(minSellLabel) // 최소 매도가격 표시 추가
		for _, rate := range settings.TargetProfitRates {
			targetPrice := buyPrice * (1 + rate/100)
			targetTrade := calculateTrade(buyPrice, targetPrice, quantity)

			profit := (targetPrice * float64(quantity)) -
				initialTrade.totalBuyAmount - initialTrade.buyCommission - targetTrade.sellCommission

			targetLabel := widget.NewLabel(fmt.Sprintf(
				"목표가(%.1f%%): $%.2f (수익: $%.4f, 매도수수료: $%.4f)",
				rate, targetPrice, profit, targetTrade.sellCommission))

			profitTargetsContainer.Add(targetLabel)
		}
		profitTargetsContainer.Refresh()
	}

	priceEntry.OnChanged = func(string) { updateCalculations() }
	quantityEntry.OnChanged = func(string) { updateCalculations() }

	// 리스트 선택 이벤트 추가
	tradesList.OnSelected = func(id widget.ListItemID) {
		selectedTrade := trades[id]

		// 선택된 거래 정보로 입력 필드 업데이트
		priceEntry.SetText(fmt.Sprintf("%.2f", selectedTrade.BuyPrice))
		quantityEntry.SetText(fmt.Sprintf("%d", selectedTrade.RemainingQty))

		// 계산 함수 호출
		updateCalculations()
	}

	// 선택 해제 이벤트 추가 (선택사항)
	tradesList.OnUnselected = func(id widget.ListItemID) {
		priceEntry.SetText("")
		quantityEntry.SetText("")
	}

	// 입력 필드들의 최소 넓이 설정
	priceEntry.Resize(fyne.NewSize(200, priceEntry.MinSize().Height))
	quantityEntry.Resize(fyne.NewSize(200, quantityEntry.MinSize().Height))

	contentContainer := container.NewVBox(
		widget.NewLabel("매수가격 (USD)"),
		priceEntry,
		widget.NewLabel("수량"),
		quantityEntry,
		buyButton,
	)
	contentContainer.Resize(fyne.NewSize(400, contentContainer.MinSize().Height))

	inputCard := widget.NewCard("매매정보", "", contentContainer)
	inputContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, inputCard.MinSize().Height)), inputCard)

	targetContent := container.NewVBox(
		totalAmountLabel,
		buyCommissionLabel,
		stopLossLabel,
		widget.NewSeparator(),
		profitTargetsContainer,
	)
	targetContent.Resize(fyne.NewSize(400, targetContent.MinSize().Height))

	targetCard := widget.NewCard("목표가 정보", "", targetContent)
	targetContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, targetCard.MinSize().Height)), targetCard)

	historyCard := widget.NewCard("매수 내역", "",
		tradesScroll,
	)
	historyCard.Resize(fyne.NewSize(400, historyCard.MinSize().Height))

	historyContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, historyCard.MinSize().Height)), historyCard)

	// 왼쪽 컨테이너 (매매정보와 목표가 정보)
	leftContainer := container.NewVBox(inputContainer, targetContainer)

	// 전체 컨테이너를 수평 배치로 변경
	mainContainer := container.NewHBox(leftContainer, historyContainer)

	registerTradeUpdateCallback(func() {
		tradesList.Refresh()
	})

	return mainContainer
}
