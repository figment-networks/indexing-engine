package structs_test

import (
	"math/big"
	"testing"

	"github.com/figment-networks/indexing-engine/structs"
)

func TestTransactionAmount_Add(t *testing.T) {
	tests := []struct {
		name string
		args []structs.TransactionAmount
		want structs.TransactionAmount
	}{
		{
			name: "test with same exp",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12344)},
				{Numeric: big.NewInt(1)},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(12345), Exp: 0},
		},
		{
			name: "test with different exp",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(6), Exp: 1},
				{Numeric: big.NewInt(12345)},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(123456), Exp: 1},
		},
		{
			name: "test with different exp reverse order",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12345)},
				{Numeric: big.NewInt(6), Exp: 1},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(123456), Exp: 1},
		},
		{
			name: "test with multiple additions in a row",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12345)},
				{Numeric: big.NewInt(6), Exp: 1},
				{Numeric: big.NewInt(7), Exp: 2},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(1234567), Exp: 2},
		},
		{
			name: "test with different exps",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(123), Exp: 3},
				{Numeric: big.NewInt(456), Exp: 6},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(123456), Exp: 6},
		},
		{
			name: "test with different exps reverse order",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(456), Exp: 6},
				{Numeric: big.NewInt(123), Exp: 3},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(123456), Exp: 6},
		},
		{
			name: "test with overflow",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(9), Exp: 1},
				{Numeric: big.NewInt(9), Exp: 1},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(18), Exp: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answer := structs.Clone(tt.args[0])
			var err error

			for _, tr := range tt.args[1:] {
				answer, err = structs.Add(answer, tr)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
			}

			if tt.want.Numeric.Cmp(answer.Numeric) != 0 {
				t.Errorf("want: %+v, got: %+v", tt.want.Numeric, answer.Numeric)
				return
			}

			if tt.want.Exp != answer.Exp {
				t.Errorf("want: %+v, got: %+v", tt.want.Exp, answer.Exp)
				return
			}
		})
	}
}

func TestTransactionAmount_Sub(t *testing.T) {

	tests := []struct {
		name string
		args []structs.TransactionAmount
		want structs.TransactionAmount
	}{
		{
			name: "test with same exp",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12346)},
				{Numeric: big.NewInt(1)},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(12345), Exp: 0},
		},
		{
			name: "test with same exp",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12346)},
				{Numeric: big.NewInt(1)},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(12345), Exp: 0},
		},
		{
			name: "test with negative answer",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(1)},
				{Numeric: big.NewInt(12346)},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(-12345), Exp: 0},
		},
		{
			name: "test with subtracting negative",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12346)},
				{Numeric: big.NewInt(-1)},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(12347), Exp: 0},
		},
		{
			name: "test with different exp",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12346)},
				{Numeric: big.NewInt(4), Exp: 1},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(123456), Exp: 1},
		},
		{
			name: "test with multiple subtractions in a row",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(200)},
				{Numeric: big.NewInt(2), Exp: 1},
				{Numeric: big.NewInt(5), Exp: 2},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(19975), Exp: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answer := structs.Clone(tt.args[0])
			var err error

			for _, tr := range tt.args[1:] {
				answer, err = structs.Sub(answer, tr)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
			}

			if tt.want.Numeric.Cmp(answer.Numeric) != 0 {
				t.Errorf("want: %+v, got: %+v", tt.want.Numeric, answer.Numeric)
				return
			}

			if tt.want.Exp != answer.Exp {
				t.Errorf("want: %+v, got: %+v", tt.want.Exp, answer.Exp)
				return
			}
		})
	}
}

func TestTransactionAmount_Div(t *testing.T) {
	bigNum := new(big.Int)
	bigNum, _ = bigNum.SetString("27238008064810922235187", 10)
	tests := []struct {
		name   string
		args   []structs.TransactionAmount
		want   structs.TransactionAmount
		errExp bool
	}{
		{
			name: "test dividing by zero",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(12346)},
				{Numeric: new(big.Int)},
			},
			errExp: true,
		},
		{
			name: "test different currencies",
			args: []structs.TransactionAmount{
				{Currency: "c1", Numeric: big.NewInt(12346)},
				{Currency: "c2", Numeric: new(big.Int)},
			},
			errExp: true,
		},
		{
			name: "test with same exp, bigger exp",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(202), Exp: 2},
				{Numeric: big.NewInt(2), Exp: 1},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(101), Exp: 1},
		},
		{
			name: "test with same exp, smaller exp",
			args: []structs.TransactionAmount{
				{Numeric: big.NewInt(2), Exp: 1},
				{Numeric: big.NewInt(202), Exp: 2},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(9900990099), Exp: 11},
		},
		{
			name: "test with big number",
			args: []structs.TransactionAmount{
				{Numeric: bigNum, Exp: 18},
				{Numeric: big.NewInt(120000000), Exp: 0},
			},
			want: structs.TransactionAmount{Numeric: big.NewInt(226983400540091), Exp: 18},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answer := structs.Clone(tt.args[0])
			answer, err := structs.Div(answer, tt.args[1])

			if !tt.errExp {
				if tt.want.Numeric.Cmp(answer.Numeric) != 0 {
					t.Errorf("Numeric want: %+v, got: %+v", tt.want.Numeric, answer.Numeric)
				}

				if tt.want.Exp != answer.Exp {
					t.Errorf("Exp want: %+v, got: %+v", tt.want.Exp, answer.Exp)
				}
			} else {
				if err == nil {
					t.Errorf("err expected")
				}
			}
		})
	}
}
