package bibliography

import (
	"io"
	"reflect"
	"testing"

	"github.com/tkw1536/gotexml/utils"
)

func Test_readTag(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantTag BibTag
		wantEOF bool
	}{
		// value only
		{
			"empty tag",
			``,
			BibTag{},
			true,
		},

		{
			"literal value",
			`value`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringLiteral,
							value: `value`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 4, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 5, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 5, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 4, EOF: false},
				},
			},
			false,
		},

		{
			"quoted value",
			`"value"`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringQuote,
							value: `value`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
				},
			},
			false,
		},
		{
			"braced value",
			`{value}`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringBracket,
							value: `value`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
				},
			},
			false,
		},
		{
			"concated literals",
			`value1 # value2`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringLiteral,
							value: `value1`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 5, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ` # `,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 8, EOF: false},
							},
						},
					},
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringLiteral,
							value: `value2`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 9, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 14, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 15, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 15, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 14, EOF: false},
				},
			},
			false,
		},
		{
			"concated quote and literal",
			`"value1" # value2`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringQuote,
							value: `value1`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ` # `,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 8, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 10, EOF: false},
							},
						},
					},
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringLiteral,
							value: `value2`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 16, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 17, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 17, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 16, EOF: false},
				},
			},
			false,
		},

		// key = value
		{
			"simple name",
			`name = value`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: true,
						name: BibString{
							kind:  BibStringLiteral,
							value: `name`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 3, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ` = `,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 4, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
							},
						},
					},
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringLiteral,
							value: `value`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 12, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 12, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 11, EOF: false},
				},
			},
			false,
		},
		{
			"simple name (compact)",
			`name=value`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: true,
						name: BibString{
							kind:  BibStringLiteral,
							value: `name`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 3, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: `=`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 4, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 4, EOF: false},
							},
						},
					},
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringLiteral,
							value: `value`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 5, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 9, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 10, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 10, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 9, EOF: false},
				},
			},
			false,
		},
		{
			"name + compact value",
			`name=a#"b"`,
			BibTag{
				prefix: BibString{
					kind:  BibStringOther,
					value: ``,
					source: utils.ReaderRange{
						Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
						End:   utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					},
				},
				elements: []BibTagElement{
					BibTagElement{
						isKeyElement: true,
						name: BibString{
							kind:  BibStringLiteral,
							value: `name`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 3, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: `=`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 4, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 4, EOF: false},
							},
						},
					},
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringLiteral,
							value: `a`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 5, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 5, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: `#`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 6, EOF: false},
							},
						},
					},
					BibTagElement{
						isKeyElement: false,
						name: BibString{
							kind:  BibStringQuote,
							value: `b`,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 7, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 9, EOF: false},
							},
						},
						suffix: BibString{
							kind:  BibStringOther,
							value: ``,
							source: utils.ReaderRange{
								Start: utils.ReaderPosition{Line: 0, Column: 10, EOF: false},
								End:   utils.ReaderPosition{Line: 0, Column: 10, EOF: false},
							},
						},
					},
				},
				source: utils.ReaderRange{
					Start: utils.ReaderPosition{Line: 0, Column: 0, EOF: false},
					End:   utils.ReaderPosition{Line: 0, Column: 9, EOF: false},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTag, err := readTag(utils.NewRuneReaderFromString(tt.input + ", "))

			// if we want eof, only test for EOF
			if tt.wantEOF {
				if err != io.EOF {
					t.Errorf("BibTag.readTag() error = %v, wantErr %v", err, io.EOF)
					return
				}
				return
			}

			if (err != nil) != false {
				t.Errorf("BibTag.readTag() error = %v, wantErr %v", err, false)
				return
			}

			if !reflect.DeepEqual(gotTag, tt.wantTag) {
				t.Errorf("BibTag.readTag() = %v, want %v", gotTag, tt.wantTag)
			}
		})
	}
}
