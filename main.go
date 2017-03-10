package main

import (
	"flag"
	"fmt"
	"math"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Mortgage struct {
	payment       float64
	extra         float64
	original      float64
	balance       float64
	interest      float64
	principal     float64
	totalInterest float64
	equity        float64
	rate          float64
	length        float64
}

func main() {
	loan := flag.Float64("loan", 100000.00, "Loan amount")
	rate := flag.Float64("rate", 6.0, "Mortgage rate in percentage")
	term := flag.Int("term", 15, "Mortgage duration in years")
	extra := flag.Float64("extra", 0.0, "Mortgage duration in years")
	flag.Parse()

	mortgage := Mortgage{
		payment:  payment(*loan, getRate(*rate), getMonths(*term)),
		original: *loan,
		balance:  *loan,
		extra:    *extra,
		rate:     getRate(*rate),
		length:   getMonths(*term),
	}
	log.WithFields(log.Fields{"Monthly Payment": mortgage.payment}).
		Info("Monthly Payment")

	mortgage.amortize(mortgage.length)
}

func (m *Mortgage) amortize(length float64) {
	if length == 0 || m.balance <= 0 {
		return
	}
	//calculate interest amount
	m.interest = m.rate * m.balance
	//calculate principle amount
	m.principal = m.payment - m.interest + m.extra
	if m.principal > m.balance {
		m.principal = m.balance
	}
	//add to equity and to total interest
	m.equity += m.principal
	m.totalInterest += m.interest
	//subtract balance
	m.balance -= m.principal

	length--
	l := m.length - length
	log.WithFields(log.Fields{
		"#":      int(l),
		"B":      toCurrency(m.balance),
		"P":      toCurrency(m.principal),
		"I":      toCurrency(m.interest),
		"Waste":  toCurrency(m.totalInterest),
		"Equity": toCurrency(m.equity),
	}).Info()
	if int(l)%12 == 0 {
		log.Warn(strings.Repeat("-", 120))
	}
	m.amortize(length)
}

func payment(loan float64, rate float64, months float64) float64 {
	return loan * (rate * math.Pow((rate+1), months)) / (math.Pow((rate+1), months) - 1)
}

func getRate(rate float64) float64 {
	return rate / 100.0 / 12
}

func getMonths(years int) float64 {
	return float64(years * 12)
}

func toCurrency(n float64) string {
	return fmt.Sprintf("$%.2f", n)
}
