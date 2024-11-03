package main

import "time"

type Trade struct {
	ID           string    // 거래 식별자
	BuyPrice     float64   // 매수가
	Quantity     int       // 수량
	BuyDate      time.Time // 매수일시
	TotalAmount  float64   // 총 매수금액
	Commission   float64   // 매수수수료
	RemainingQty int       // 남은 수량 (부분 매도 지원)
}

var trades = make([]Trade, 0)

// 매수내역 업데이트를 위한 콜백 함수들을 저장
var tradeUpdateCallbacks []func()

// 매수내역 업데이트 콜백 등록
func registerTradeUpdateCallback(callback func()) {
	tradeUpdateCallbacks = append(tradeUpdateCallbacks, callback)
}

// 모든 매수내역 리스트 업데이트
func updateAllTradeLists() {
	for _, callback := range tradeUpdateCallbacks {
		callback()
	}
}

func addTrade(price float64, quantity int) Trade {
	trade := Trade{
		ID:           time.Now().Format("20060102150405"), // YYYYMMDDhhmmss 형식
		BuyPrice:     price,
		Quantity:     quantity,
		BuyDate:      time.Now(),
		TotalAmount:  price * float64(quantity),
		Commission:   price * float64(quantity) * settings.CommissionRate / 100.0,
		RemainingQty: quantity,
	}
	trades = append(trades, trade)
	updateAllTradeLists() // 모든 리스트 업데이트
	return trade
}

// 매수내역 수정 함수 추가
func updateTrade(index int, remainingQty int) {
	if index >= 0 && index < len(trades) {
		trades[index].RemainingQty = remainingQty
		if remainingQty <= 0 {
			trades = append(trades[:index], trades[index+1:]...)
		}
		updateAllTradeLists() // 모든 리스트 업데이트
	}
}
