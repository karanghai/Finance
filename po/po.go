//TODO: Detailed Description

package po

import (
	"math"
	"math/rand"
)

type Security struct {
	Name       string
	Symbol     string
	SecurityId int64
}

type OptionType uint64

const (
	CALL OptionType = 0
	PUT  OptionType = 1
)

// TODO: Add greeks
type Option struct {
	Ticker            Security
	StockPrice        float64
	StrikePrice       float64
	ImpliedVolatility float64
	RiskFreeRate      float64
	TimeToExpiration  float64
	Type              OptionType
	OptionPrice       float64
}

// TODO: Comments
func (o *Option) UnderlyingPriceAtExpiration() float64 {
	power := (o.RiskFreeRate-0.5*math.Pow(o.ImpliedVolatility, 2))*(o.TimeToExpiration) + (o.ImpliedVolatility * math.Sqrt(o.TimeToExpiration) * rand.NormFloat64())
	return o.StockPrice * math.Exp(power)
}

// TODO: Comments
func (o *Option) CallPayOff(stockPrice float64) float64 {
	return math.Max(stockPrice-o.StrikePrice, 0)
}

// TODO: Comments
func (o *Option) PutPayOff(stockPrice float64) float64 {
	return math.Max(o.StrikePrice-stockPrice, 0)
}

// TODO: Comments
type MonteCarlo struct {
	Subject     Option
	Simulations int64
}

// TODO: Comments
func (m *MonteCarlo) Simulate() Option {
	if m.Subject.Type == CALL {
		m.CalculatePrice(m.Subject.CallPayOff)
	} else if m.Subject.Type == PUT {
		m.CalculatePrice(m.Subject.PutPayOff)
	}
	return m.Subject
}

// TODO: Comments
func (m *MonteCarlo) CalculatePrice(payOff func(stockPrice float64) float64) {
	var simulatedPayoffsSum float64
	for i := int64(0); i < m.Simulations; i++ {
		St := m.Subject.UnderlyingPriceAtExpiration()
		simulatedPayoffsSum += payOff(St)
	}
	discount := math.Exp(-m.Subject.RiskFreeRate * m.Subject.TimeToExpiration)
	m.Subject.OptionPrice = discount * simulatedPayoffsSum / float64(m.Simulations)
}