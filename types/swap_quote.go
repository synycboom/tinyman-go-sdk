package types

import (
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
)

// SwapQuote represents a swap quote
type SwapQuote struct {
	// SwapType is a swap type
	SwapType string

	// AmountIn is an input asset amount
	AmountIn AssetAmount

	// AmountOut is an output asset amount
	AmountOut AssetAmount

	// SwapFee is a swap fee
	SwapFee AssetAmount

	// Slippage is a slippage
	Slippage float64
}

// AmountOutWithSlippage calculates the output asset amount after applying the slippage
func (s *SwapQuote) AmountOutWithSlippage() (*AssetAmount, error) {
	if s.SwapType == constants.SwapFixedOutput {
		return &s.AmountOut, nil
	}

	amountOutWithSlippage, err := s.AmountOut.Mul(nil, &s.Slippage)
	if err != nil {
		return nil, err
	}
	amount, err := s.AmountOut.Sub(amountOutWithSlippage, nil)
	if err != nil {
		return nil, err
	}

	return amount, nil
}

// AmountInWithSlippage calculates the input asset amount after applying the slippage
func (s *SwapQuote) AmountInWithSlippage() (*AssetAmount, error) {
	if s.SwapType == constants.SwapFixedInput {
		return &s.AmountIn, nil
	}

	amountInWithSlippage, err := s.AmountIn.Mul(nil, &s.Slippage)
	if err != nil {
		return nil, err
	}
	amount, err := s.AmountIn.Add(amountInWithSlippage, nil)
	if err != nil {
		return nil, err
	}

	return amount, nil
}

// Price returns the price
func (s *SwapQuote) Price() float64 {
	return float64(s.AmountOut.Amount) / float64(s.AmountIn.Amount)
}

// PriceWithSlippage returns the price after applying the slippage
func (s *SwapQuote) PriceWithSlippage() (float64, error) {
	in, err := s.AmountInWithSlippage()
	if err != nil {
		return 0, err
	}
	out, err := s.AmountOutWithSlippage()
	if err != nil {
		return 0, err
	}

	return float64(out.Amount) / float64(in.Amount), nil
}
