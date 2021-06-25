package structs

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

var (
	tenInt = big.NewInt(10)
)

func (t *TransactionAmount) GetFloat() *big.Float {
	divider := new(big.Int).Exp(tenInt, big.NewInt(int64(t.Exp)), nil)
	dvdr := new(big.Float).SetInt(divider)
	value := new(big.Float).SetInt(t.Numeric)
	value.Quo(value, dvdr)
	return value
}

func Add(t, o TransactionAmount) (TransactionAmount, error) {
	result := Clone(t)

	if t.Currency != o.Currency {
		return result, fmt.Errorf("coin currency different: %v %v\n", t.Currency, o.Currency)
	}
	expDiff := (t.Exp - o.Exp)

	if expDiff < 0 {
		zerosToAdd := math.Abs(float64(expDiff))
		multiplier := new(big.Int).Exp(tenInt, big.NewInt(int64(zerosToAdd)), nil)
		result.Numeric.Mul(result.Numeric, multiplier)
		result.Numeric.Add(result.Numeric, o.Numeric)
		result.Exp = result.Exp - expDiff
		result.Text = ""
		return result, nil
	}

	zerosToAdd := int64(expDiff)
	multiplier := new(big.Int).Exp(tenInt, big.NewInt(zerosToAdd), nil)
	tmp := Clone(o)
	tmp.Numeric.Mul(tmp.Numeric, multiplier)
	result.Numeric.Add(tmp.Numeric, result.Numeric)
	result.Text = ""
	return result, nil
}

func Sub(t, o TransactionAmount) (TransactionAmount, error) {
	result := Clone(t)

	if t.Currency != o.Currency {
		return result, fmt.Errorf("coin currency different: %v %v\n", t.Currency, o.Currency)
	}

	expDiff := (t.Exp - o.Exp)

	if expDiff < 0 {
		zerosToAdd := math.Abs(float64(expDiff))
		multiplier := new(big.Int).Exp(tenInt, big.NewInt(int64(zerosToAdd)), nil)
		result.Numeric.Mul(result.Numeric, multiplier)
		result.Numeric.Sub(result.Numeric, o.Numeric)
		result.Exp = result.Exp - expDiff
		result.Text = ""
		return result, nil
	}

	zerosToAdd := int64(expDiff)
	multiplier := new(big.Int).Exp(tenInt, big.NewInt(zerosToAdd), nil)
	tmp := Clone(o)
	tmp.Numeric.Mul(tmp.Numeric, multiplier)
	result.Numeric.Sub(result.Numeric, tmp.Numeric)
	result.Text = ""
	return result, nil
}

func Div(t, o TransactionAmount) (TransactionAmount, error) {
	result := Clone(t)

	if t.Currency != o.Currency {
		return result, fmt.Errorf("coin currency different: %v %v\n", t.Currency, o.Currency)
	}

	if o.Numeric.Int64() == 0 {
		return result, fmt.Errorf("division by zero")
	}

	expDiff := result.Exp - o.Exp

	res := new(big.Float)
	tValue := new(big.Float).SetInt(result.Numeric)
	oValue := new(big.Float).SetInt(o.Numeric)
	res.Quo(tValue, oValue)

	var dec int32
	if res.Cmp(big.NewFloat(1)) < 0 {
		s := strings.Split(res.Text('g', 10), ".")
		dec = int32(len(s[1]))
		multiplier := new(big.Int).Exp(tenInt, big.NewInt(int64(dec)), nil)
		mltp := new(big.Float).SetInt(multiplier)
		res.Mul(res, mltp)
	}

	result.Numeric = new(big.Int)
	res.Int(result.Numeric)
	result.Exp = dec + expDiff

	return result, nil
}

func Clone(t TransactionAmount) TransactionAmount {
	tmp := TransactionAmount{
		Text:     t.Text,
		Currency: t.Currency,
		Numeric:  &big.Int{},
		Exp:      t.Exp,
	}
	if t.Numeric == nil {
		tmp.Numeric.Set(big.NewInt(0))
		return tmp
	}
	tmp.Numeric.Set(t.Numeric)
	return tmp
}
