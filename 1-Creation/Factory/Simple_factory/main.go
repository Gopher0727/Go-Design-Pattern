package main

import "fmt"

type Payment interface {
	pay(int64) (string, error)
}

type AliPay struct{}

func (a *AliPay) pay(money int64) (string, error) {
	return "使用支付宝支付了" + fmt.Sprint(money) + "元", nil
}

type WeChatPay struct{}

func (w *WeChatPay) pay(money int64) (string, error) {
	return "使用微信支付了" + fmt.Sprint(money) + "元", nil
}

type PayType int8

const (
	AliPayType    PayType = 1
	WeChatPayType PayType = 2
)

func NewPayment(payType PayType) Payment {
	switch payType {
	case AliPayType:
		return &AliPay{}
	case WeChatPayType:
		return &WeChatPay{}
	default:
		return nil
	}
}

func main() {
	p := NewPayment(AliPayType)
	result, err := p.pay(100)
	if err != nil {
		fmt.Println("支付失败:", err)
		return
	}
	fmt.Println(result)
}
