package blif

import "testing"

func TestBlifParser(t *testing.T) {
	parser, err := NewBlifParser()
	if err != nil {
		t.Errorf("Failed to build parser: %s", err)
	}

	tests := []struct {
		blif string
	}{
		{`# Benchmark "pla/clip_124" written by ABC on Thu Jan 25 13:45:29 2024
.model pla/clip_124
.inputs x0 x1 x2 x3 x4 x5 x6 x7 x8
.outputs f0 f1 f2 f3 f4
.gate NOT      A=x0 Y=new_n15_
.gate NOT      A=x2 Y=new_n16_
.gate NOR      A=x7 B=x4 Y=new_n17_
.gate NOR      A=new_n17_ B=new_n16_ Y=new_n18_
.gate NOT      A=new_n17_ Y=new_n19_
.end`,
		},
		{`# Benchmark "pla/parity_188" written by ABC on Thu Jan 25 13:45:30 2024
.model pla/parity_188
.inputs x0 x1 x2 x3 x4 x5 x6 x7 x8 x9 x10 x11 x12 x13 x14 x15
.outputs f0
.gate NOT      A=x15 Y=new_n18_
.gate NOR      A=new_n18_ B=x14 Y=new_n19_
.gate NOT      A=x14 Y=new_n20_
.end`},
	}

	for _, test := range tests {
		_, err := parser.ParseString("", test.blif)
		if err != nil {
			t.Errorf("Failed to parse BLIF file: %s", err)
		}
	}
}
