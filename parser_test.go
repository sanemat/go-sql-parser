package sqlparser

import (
	"errors"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		input   []Token
		want    Node
		wantErr error
	}{
		{
			name: "simple select",
			input: []Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenIdentifier, Literal: "column1"},
				{Type: TokenKeyword, Literal: "from"},
				{Type: TokenIdentifier, Literal: "tablename"},
				{Type: TokenSymbol, Literal: ";"},
				{Type: TokenEOF, Literal: ""},
			},
			want: &SelectStatement{
				Columns: []string{"column1"},
				Table:   "table1",
				Where:   nil,
			},
		},
		{
			name: "select multiple columns",
			input: []Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenIdentifier, Literal: "id"},
				{Type: TokenSymbol, Literal: ","},
				{Type: TokenIdentifier, Literal: "title"},
				{Type: TokenKeyword, Literal: "from"},
				{Type: TokenIdentifier, Literal: "table1"},
				{Type: TokenSymbol, Literal: ";"},
				{Type: TokenEOF, Literal: ""},
			},
			want: &SelectStatement{
				Columns: []string{"id", "title"},
				Table:   "table1",
				Where:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			got, err := p.Parse()
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Praser.Parse() error = %+v, want %+v", err, tt.want)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parser.Parse() = %+v, want %+v", got, tt.want)
				}
			} else {
				if err == nil {
					t.Errorf("Parser.Parse() = %+v, wantErr %+v", got, tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Parser.Parse() error = %+v, wantErr %+v", err, tt.wantErr)
				}
			}
		})
	}
}
