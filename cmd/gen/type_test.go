package main

import (
	"reflect"
	"testing"
)

func TestParseTypeExpr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *TypeExpr
		wantErr bool
	}{
		{
			name:  "simple",
			input: "PATH",
			want: &TypeExpr{
				Base: "PATH",
			},
		},
		{
			name:  "simple list",
			input: "PATH [...]",
			want: &TypeExpr{
				Base:     "PATH",
				Repeated: true,
			},
		},
		{
			name:  "nested list",
			input: "PATH [ARGUMENT [...]]",
			want: &TypeExpr{
				Base: "PATH",
				Inner: &TypeExpr{
					Base:     "ARGUMENT",
					Repeated: true,
				},
			},
		},
		{
			name:  "deep nesting",
			input: "PATH [ARGUMENT [UNIT [STRING [...]]]]",
			want: &TypeExpr{
				Base: "PATH",
				Inner: &TypeExpr{
					Base: "ARGUMENT",
					Inner: &TypeExpr{
						Base: "UNIT",
						Inner: &TypeExpr{
							Base:     "STRING",
							Repeated: true,
						},
					},
				},
			},
		},
		{
			name:    "missing closing bracket",
			input:   "PATH [ARGUMENT [...]",
			wantErr: true,
		},
		{
			name:    "unexpected tokens",
			input:   "PATH PATH",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTypeExpr(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf(
					"ParseTypeExpr(%q)\n got  %#v\n want %#v",
					tt.input,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestTypeExpr_toGoType_Functional(t *testing.T) {
	tests := []struct {
		name      string
		expr      *TypeExpr
		wantType  string
		wantPanic bool
	}{
		{
			name:     "nil expression",
			expr:     nil,
			wantType: "any",
		},
		{
			name: "single base string: STRING",
			expr: &TypeExpr{
				Base: TypeString,
			},
			wantType: "string",
		},
		{
			name: "repeated base string: STRING[...]",
			expr: &TypeExpr{
				Base:     TypeString,
				Repeated: true,
			},
			wantType: "[]string",
		},
		{
			name: "two-layer nested string: STRING[STRING]",
			expr: &TypeExpr{
				Base: TypeString,
				Inner: &TypeExpr{
					Base:     TypeString,
					Repeated: true,
				},
			},
			wantType: "[][]string",
		},
		{
			name: "three-layer nested int: INT[INT[INT[...]]]",
			expr: &TypeExpr{
				Base: TypeInteger,
				Inner: &TypeExpr{
					Base: TypeLevel,
					Inner: &TypeExpr{
						Base:     TypeServiceExitType,
						Repeated: true,
					},
				},
			},
			wantType: "[][][]int",
		},
		{
			name: "nested mismatch panics",
			expr: &TypeExpr{
				Base: TypeString,
				Inner: &TypeExpr{
					Base: TypeBoolean,
				},
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				} else if tt.wantPanic {
					t.Errorf("expected panic but did not occur")
				}
			}()

			gotType, _ := tt.expr.toGoType()
			if !tt.wantPanic && gotType != tt.wantType {
				t.Errorf("toGoType() = %q, want %q", gotType, tt.wantType)
			}
		})
	}
}
