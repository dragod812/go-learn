package main

import (
	"reflect"
	"testing"
)

func TestJustifyText(t *testing.T) {
	type args struct {
		wordSequence     []string
		targetLineLength int
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{
			"OneLine",
			args{
				[]string{"The", "quick"},
				11,
			},
			[]string{
				"The-quick--",
			},
		},
		{
			"MultiLineSimple",
			args{
				[]string{"the", "quick", "brown", "fox"},
				11,
			},
			[]string{
				"the---quick",
				"brown-fox--",
			},
		},
		{
			"OneWordEachLine",
			args{
				[]string{
					"aaaaaaaaaa",
					"aaaaaaaaaaa",
					"aaa",
					"aaaaaaaaaa",
				},
				11,
			},
			[]string{
				"aaaaaaaaaa-",
				"aaaaaaaaaaa",
				"aaa--------",
				"aaaaaaaaaa-",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := JustifyText(tt.args.wordSequence, tt.args.targetLineLength); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("JustifyText() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
