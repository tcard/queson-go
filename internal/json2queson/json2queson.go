package json2queson
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "JSON_text",
			pos:  position{line: 43, col: 1, offset: 1866},
			expr: &actionExpr{
				pos: position{line: 44, col: 5, offset: 1880},
				run: (*parser).callonJSON_text1,
				expr: &seqExpr{
					pos: position{line: 44, col: 5, offset: 1880},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 44, col: 5, offset: 1880},
							name: "ws",
						},
						&labeledExpr{
							pos:   position{line: 44, col: 8, offset: 1883},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 44, col: 14, offset: 1889},
								name: "value",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 44, col: 20, offset: 1895},
							name: "ws",
						},
					},
				},
			},
		},
		{
			name: "begin_array",
			pos:  position{line: 46, col: 1, offset: 1921},
			expr: &seqExpr{
				pos: position{line: 46, col: 19, offset: 1939},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 46, col: 19, offset: 1939},
						name: "ws",
					},
					&litMatcher{
						pos:        position{line: 46, col: 22, offset: 1942},
						val:        "[",
						ignoreCase: false,
						want:       "\"[\"",
					},
					&ruleRefExpr{
						pos:  position{line: 46, col: 26, offset: 1946},
						name: "ws",
					},
				},
			},
		},
		{
			name: "begin_object",
			pos:  position{line: 47, col: 1, offset: 1949},
			expr: &seqExpr{
				pos: position{line: 47, col: 19, offset: 1967},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 47, col: 19, offset: 1967},
						name: "ws",
					},
					&litMatcher{
						pos:        position{line: 47, col: 22, offset: 1970},
						val:        "{",
						ignoreCase: false,
						want:       "\"{\"",
					},
					&ruleRefExpr{
						pos:  position{line: 47, col: 26, offset: 1974},
						name: "ws",
					},
				},
			},
		},
		{
			name: "end_array",
			pos:  position{line: 48, col: 1, offset: 1977},
			expr: &seqExpr{
				pos: position{line: 48, col: 19, offset: 1995},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 48, col: 19, offset: 1995},
						name: "ws",
					},
					&litMatcher{
						pos:        position{line: 48, col: 22, offset: 1998},
						val:        "]",
						ignoreCase: false,
						want:       "\"]\"",
					},
					&ruleRefExpr{
						pos:  position{line: 48, col: 26, offset: 2002},
						name: "ws",
					},
				},
			},
		},
		{
			name: "end_object",
			pos:  position{line: 49, col: 1, offset: 2005},
			expr: &seqExpr{
				pos: position{line: 49, col: 19, offset: 2023},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 49, col: 19, offset: 2023},
						name: "ws",
					},
					&litMatcher{
						pos:        position{line: 49, col: 22, offset: 2026},
						val:        "}",
						ignoreCase: false,
						want:       "\"}\"",
					},
					&ruleRefExpr{
						pos:  position{line: 49, col: 26, offset: 2030},
						name: "ws",
					},
				},
			},
		},
		{
			name: "name_separator",
			pos:  position{line: 50, col: 1, offset: 2033},
			expr: &seqExpr{
				pos: position{line: 50, col: 19, offset: 2051},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 50, col: 19, offset: 2051},
						name: "ws",
					},
					&litMatcher{
						pos:        position{line: 50, col: 22, offset: 2054},
						val:        ":",
						ignoreCase: false,
						want:       "\":\"",
					},
					&ruleRefExpr{
						pos:  position{line: 50, col: 26, offset: 2058},
						name: "ws",
					},
				},
			},
		},
		{
			name: "value_separator",
			pos:  position{line: 51, col: 1, offset: 2061},
			expr: &seqExpr{
				pos: position{line: 51, col: 19, offset: 2079},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 51, col: 19, offset: 2079},
						name: "ws",
					},
					&litMatcher{
						pos:        position{line: 51, col: 22, offset: 2082},
						val:        ",",
						ignoreCase: false,
						want:       "\",\"",
					},
					&ruleRefExpr{
						pos:  position{line: 51, col: 26, offset: 2086},
						name: "ws",
					},
				},
			},
		},
		{
			name:        "ws",
			displayName: "\"whitespace\"",
			pos:         position{line: 53, col: 1, offset: 2090},
			expr: &zeroOrMoreExpr{
				pos: position{line: 53, col: 19, offset: 2108},
				expr: &charClassMatcher{
					pos:        position{line: 53, col: 19, offset: 2108},
					val:        "[ \\t\\n\\r]",
					chars:      []rune{' ', '\t', '\n', '\r'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "value",
			pos:  position{line: 57, col: 1, offset: 2146},
			expr: &choiceExpr{
				pos: position{line: 58, col: 5, offset: 2156},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 58, col: 5, offset: 2156},
						val:        "false",
						ignoreCase: false,
						want:       "\"false\"",
					},
					&litMatcher{
						pos:        position{line: 59, col: 5, offset: 2168},
						val:        "null",
						ignoreCase: false,
						want:       "\"null\"",
					},
					&litMatcher{
						pos:        position{line: 60, col: 5, offset: 2179},
						val:        "true",
						ignoreCase: false,
						want:       "\"true\"",
					},
					&ruleRefExpr{
						pos:  position{line: 61, col: 5, offset: 2190},
						name: "object",
					},
					&ruleRefExpr{
						pos:  position{line: 62, col: 5, offset: 2201},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 5, offset: 2211},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 64, col: 5, offset: 2222},
						name: "str",
					},
				},
			},
		},
		{
			name: "object",
			pos:  position{line: 68, col: 1, offset: 2254},
			expr: &actionExpr{
				pos: position{line: 69, col: 5, offset: 2265},
				run: (*parser).callonobject1,
				expr: &seqExpr{
					pos: position{line: 69, col: 5, offset: 2265},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 69, col: 5, offset: 2265},
							name: "begin_object",
						},
						&labeledExpr{
							pos:   position{line: 70, col: 5, offset: 2282},
							label: "members",
							expr: &zeroOrOneExpr{
								pos: position{line: 70, col: 13, offset: 2290},
								expr: &actionExpr{
									pos: position{line: 71, col: 7, offset: 2298},
									run: (*parser).callonobject6,
									expr: &seqExpr{
										pos: position{line: 71, col: 7, offset: 2298},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 71, col: 7, offset: 2298},
												label: "head",
												expr: &ruleRefExpr{
													pos:  position{line: 71, col: 12, offset: 2303},
													name: "member",
												},
											},
											&labeledExpr{
												pos:   position{line: 72, col: 7, offset: 2316},
												label: "tail",
												expr: &zeroOrMoreExpr{
													pos: position{line: 72, col: 12, offset: 2321},
													expr: &actionExpr{
														pos: position{line: 72, col: 13, offset: 2322},
														run: (*parser).callonobject12,
														expr: &seqExpr{
															pos: position{line: 72, col: 13, offset: 2322},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 72, col: 13, offset: 2322},
																	name: "value_separator",
																},
																&labeledExpr{
																	pos:   position{line: 72, col: 29, offset: 2338},
																	label: "m",
																	expr: &ruleRefExpr{
																		pos:  position{line: 72, col: 31, offset: 2340},
																		name: "member",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 5, offset: 2588},
							name: "end_object",
						},
					},
				},
			},
		},
		{
			name: "member",
			pos:  position{line: 92, col: 1, offset: 2779},
			expr: &actionExpr{
				pos: position{line: 93, col: 5, offset: 2790},
				run: (*parser).callonmember1,
				expr: &seqExpr{
					pos: position{line: 93, col: 5, offset: 2790},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 93, col: 5, offset: 2790},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 93, col: 10, offset: 2795},
								name: "str",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 93, col: 14, offset: 2799},
							name: "name_separator",
						},
						&labeledExpr{
							pos:   position{line: 93, col: 29, offset: 2814},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 93, col: 35, offset: 2820},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "array",
			pos:  position{line: 99, col: 1, offset: 2933},
			expr: &actionExpr{
				pos: position{line: 100, col: 5, offset: 2943},
				run: (*parser).callonarray1,
				expr: &seqExpr{
					pos: position{line: 100, col: 5, offset: 2943},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 100, col: 5, offset: 2943},
							name: "begin_array",
						},
						&labeledExpr{
							pos:   position{line: 101, col: 5, offset: 2959},
							label: "values",
							expr: &zeroOrOneExpr{
								pos: position{line: 101, col: 12, offset: 2966},
								expr: &actionExpr{
									pos: position{line: 102, col: 7, offset: 2974},
									run: (*parser).callonarray6,
									expr: &seqExpr{
										pos: position{line: 102, col: 7, offset: 2974},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 102, col: 7, offset: 2974},
												label: "head",
												expr: &ruleRefExpr{
													pos:  position{line: 102, col: 12, offset: 2979},
													name: "value",
												},
											},
											&labeledExpr{
												pos:   position{line: 103, col: 7, offset: 2991},
												label: "tail",
												expr: &zeroOrMoreExpr{
													pos: position{line: 103, col: 12, offset: 2996},
													expr: &actionExpr{
														pos: position{line: 103, col: 13, offset: 2997},
														run: (*parser).callonarray12,
														expr: &seqExpr{
															pos: position{line: 103, col: 13, offset: 2997},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 103, col: 13, offset: 2997},
																	name: "value_separator",
																},
																&labeledExpr{
																	pos:   position{line: 103, col: 29, offset: 3013},
																	label: "v",
																	expr: &ruleRefExpr{
																		pos:  position{line: 103, col: 31, offset: 3015},
																		name: "value",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 113, col: 5, offset: 3262},
							name: "end_array",
						},
					},
				},
			},
		},
		{
			name:        "number",
			displayName: "\"number\"",
			pos:         position{line: 125, col: 1, offset: 3477},
			expr: &actionExpr{
				pos: position{line: 126, col: 5, offset: 3497},
				run: (*parser).callonnumber1,
				expr: &seqExpr{
					pos: position{line: 126, col: 5, offset: 3497},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 126, col: 5, offset: 3497},
							label: "sign",
							expr: &zeroOrOneExpr{
								pos: position{line: 126, col: 10, offset: 3502},
								expr: &ruleRefExpr{
									pos:  position{line: 126, col: 10, offset: 3502},
									name: "minus",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 126, col: 17, offset: 3509},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 126, col: 19, offset: 3511},
								name: "int",
							},
						},
						&labeledExpr{
							pos:   position{line: 126, col: 23, offset: 3515},
							label: "f",
							expr: &zeroOrOneExpr{
								pos: position{line: 126, col: 25, offset: 3517},
								expr: &ruleRefExpr{
									pos:  position{line: 126, col: 25, offset: 3517},
									name: "frac",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 126, col: 31, offset: 3523},
							label: "e",
							expr: &zeroOrOneExpr{
								pos: position{line: 126, col: 33, offset: 3525},
								expr: &ruleRefExpr{
									pos:  position{line: 126, col: 33, offset: 3525},
									name: "exp",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "decimal_point",
			pos:  position{line: 137, col: 1, offset: 3729},
			expr: &litMatcher{
				pos:        position{line: 138, col: 5, offset: 3747},
				val:        ".",
				ignoreCase: false,
				want:       "\".\"",
			},
		},
		{
			name: "digit1_9",
			pos:  position{line: 140, col: 1, offset: 3752},
			expr: &charClassMatcher{
				pos:        position{line: 141, col: 5, offset: 3765},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "e",
			pos:  position{line: 143, col: 1, offset: 3772},
			expr: &charClassMatcher{
				pos:        position{line: 144, col: 5, offset: 3778},
				val:        "[eE]",
				chars:      []rune{'e', 'E'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "exp",
			pos:  position{line: 146, col: 1, offset: 3784},
			expr: &actionExpr{
				pos: position{line: 147, col: 5, offset: 3792},
				run: (*parser).callonexp1,
				expr: &seqExpr{
					pos: position{line: 147, col: 5, offset: 3792},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 147, col: 5, offset: 3792},
							name: "e",
						},
						&labeledExpr{
							pos:   position{line: 147, col: 7, offset: 3794},
							label: "sign",
							expr: &zeroOrOneExpr{
								pos: position{line: 147, col: 12, offset: 3799},
								expr: &choiceExpr{
									pos: position{line: 147, col: 13, offset: 3800},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 147, col: 13, offset: 3800},
											name: "minus",
										},
										&ruleRefExpr{
											pos:  position{line: 147, col: 21, offset: 3808},
											name: "plus",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 147, col: 28, offset: 3815},
							label: "i",
							expr: &oneOrMoreExpr{
								pos: position{line: 147, col: 30, offset: 3817},
								expr: &ruleRefExpr{
									pos:  position{line: 147, col: 30, offset: 3817},
									name: "DIGIT",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "frac",
			pos:  position{line: 159, col: 1, offset: 4051},
			expr: &actionExpr{
				pos: position{line: 160, col: 5, offset: 4060},
				run: (*parser).callonfrac1,
				expr: &seqExpr{
					pos: position{line: 160, col: 5, offset: 4060},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 160, col: 5, offset: 4060},
							name: "decimal_point",
						},
						&oneOrMoreExpr{
							pos: position{line: 160, col: 19, offset: 4074},
							expr: &ruleRefExpr{
								pos:  position{line: 160, col: 19, offset: 4074},
								name: "DIGIT",
							},
						},
					},
				},
			},
		},
		{
			name: "int",
			pos:  position{line: 162, col: 1, offset: 4105},
			expr: &choiceExpr{
				pos: position{line: 163, col: 5, offset: 4113},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 163, col: 5, offset: 4113},
						name: "zero",
					},
					&actionExpr{
						pos: position{line: 163, col: 12, offset: 4120},
						run: (*parser).callonint3,
						expr: &seqExpr{
							pos: position{line: 163, col: 13, offset: 4121},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 163, col: 13, offset: 4121},
									name: "digit1_9",
								},
								&zeroOrMoreExpr{
									pos: position{line: 163, col: 22, offset: 4130},
									expr: &ruleRefExpr{
										pos:  position{line: 163, col: 22, offset: 4130},
										name: "DIGIT",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "minus",
			pos:  position{line: 165, col: 1, offset: 4162},
			expr: &litMatcher{
				pos:        position{line: 166, col: 5, offset: 4172},
				val:        "-",
				ignoreCase: false,
				want:       "\"-\"",
			},
		},
		{
			name: "plus",
			pos:  position{line: 168, col: 1, offset: 4177},
			expr: &actionExpr{
				pos: position{line: 169, col: 5, offset: 4186},
				run: (*parser).callonplus1,
				expr: &litMatcher{
					pos:        position{line: 169, col: 5, offset: 4186},
					val:        "+",
					ignoreCase: false,
					want:       "\"+\"",
				},
			},
		},
		{
			name: "zero",
			pos:  position{line: 171, col: 1, offset: 4219},
			expr: &litMatcher{
				pos:        position{line: 172, col: 5, offset: 4228},
				val:        "0",
				ignoreCase: false,
				want:       "\"0\"",
			},
		},
		{
			name:        "str",
			displayName: "\"string\"",
			pos:         position{line: 176, col: 1, offset: 4260},
			expr: &actionExpr{
				pos: position{line: 177, col: 5, offset: 4277},
				run: (*parser).callonstr1,
				expr: &seqExpr{
					pos: position{line: 177, col: 5, offset: 4277},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 177, col: 5, offset: 4277},
							name: "quotation_mark",
						},
						&labeledExpr{
							pos:   position{line: 177, col: 20, offset: 4292},
							label: "chars",
							expr: &zeroOrMoreExpr{
								pos: position{line: 177, col: 26, offset: 4298},
								expr: &ruleRefExpr{
									pos:  position{line: 177, col: 26, offset: 4298},
									name: "char",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 177, col: 32, offset: 4304},
							name: "quotation_mark",
						},
					},
				},
			},
		},
		{
			name: "char",
			pos:  position{line: 187, col: 1, offset: 4515},
			expr: &choiceExpr{
				pos: position{line: 188, col: 5, offset: 4524},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 188, col: 5, offset: 4524},
						run: (*parser).callonchar2,
						expr: &litMatcher{
							pos:        position{line: 188, col: 5, offset: 4524},
							val:        ".",
							ignoreCase: false,
							want:       "\".\"",
						},
					},
					&actionExpr{
						pos: position{line: 189, col: 5, offset: 4561},
						run: (*parser).callonchar4,
						expr: &litMatcher{
							pos:        position{line: 189, col: 5, offset: 4561},
							val:        "_",
							ignoreCase: false,
							want:       "\"_\"",
						},
					},
					&actionExpr{
						pos: position{line: 190, col: 5, offset: 4598},
						run: (*parser).callonchar6,
						expr: &litMatcher{
							pos:        position{line: 190, col: 5, offset: 4598},
							val:        " ",
							ignoreCase: false,
							want:       "\" \"",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 191, col: 5, offset: 4634},
						name: "unescaped",
					},
					&actionExpr{
						pos: position{line: 192, col: 5, offset: 4648},
						run: (*parser).callonchar9,
						expr: &seqExpr{
							pos: position{line: 192, col: 5, offset: 4648},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 192, col: 5, offset: 4648},
									name: "escape",
								},
								&labeledExpr{
									pos:   position{line: 193, col: 5, offset: 4660},
									label: "sequence",
									expr: &choiceExpr{
										pos: position{line: 194, col: 9, offset: 4679},
										alternatives: []interface{}{
											&actionExpr{
												pos: position{line: 194, col: 9, offset: 4679},
												run: (*parser).callonchar14,
												expr: &litMatcher{
													pos:        position{line: 194, col: 9, offset: 4679},
													val:        "\"",
													ignoreCase: false,
													want:       "\"\\\"\"",
												},
											},
											&actionExpr{
												pos: position{line: 195, col: 9, offset: 4720},
												run: (*parser).callonchar16,
												expr: &litMatcher{
													pos:        position{line: 195, col: 9, offset: 4720},
													val:        "\\",
													ignoreCase: false,
													want:       "\"\\\\\"",
												},
											},
											&actionExpr{
												pos: position{line: 196, col: 9, offset: 4762},
												run: (*parser).callonchar18,
												expr: &litMatcher{
													pos:        position{line: 196, col: 9, offset: 4762},
													val:        "/",
													ignoreCase: false,
													want:       "\"/\"",
												},
											},
											&actionExpr{
												pos: position{line: 197, col: 9, offset: 4803},
												run: (*parser).callonchar20,
												expr: &litMatcher{
													pos:        position{line: 197, col: 9, offset: 4803},
													val:        "b",
													ignoreCase: false,
													want:       "\"b\"",
												},
											},
											&actionExpr{
												pos: position{line: 198, col: 9, offset: 4845},
												run: (*parser).callonchar22,
												expr: &litMatcher{
													pos:        position{line: 198, col: 9, offset: 4845},
													val:        "f",
													ignoreCase: false,
													want:       "\"f\"",
												},
											},
											&actionExpr{
												pos: position{line: 199, col: 9, offset: 4887},
												run: (*parser).callonchar24,
												expr: &litMatcher{
													pos:        position{line: 199, col: 9, offset: 4887},
													val:        "n",
													ignoreCase: false,
													want:       "\"n\"",
												},
											},
											&actionExpr{
												pos: position{line: 200, col: 9, offset: 4929},
												run: (*parser).callonchar26,
												expr: &litMatcher{
													pos:        position{line: 200, col: 9, offset: 4929},
													val:        "r",
													ignoreCase: false,
													want:       "\"r\"",
												},
											},
											&actionExpr{
												pos: position{line: 201, col: 9, offset: 4971},
												run: (*parser).callonchar28,
												expr: &litMatcher{
													pos:        position{line: 201, col: 9, offset: 4971},
													val:        "t",
													ignoreCase: false,
													want:       "\"t\"",
												},
											},
											&actionExpr{
												pos: position{line: 202, col: 9, offset: 5013},
												run: (*parser).callonchar30,
												expr: &seqExpr{
													pos: position{line: 202, col: 9, offset: 5013},
													exprs: []interface{}{
														&litMatcher{
															pos:        position{line: 202, col: 9, offset: 5013},
															val:        "u",
															ignoreCase: false,
															want:       "\"u\"",
														},
														&labeledExpr{
															pos:   position{line: 202, col: 13, offset: 5017},
															label: "digits",
															expr: &seqExpr{
																pos: position{line: 202, col: 21, offset: 5025},
																exprs: []interface{}{
																	&ruleRefExpr{
																		pos:  position{line: 202, col: 21, offset: 5025},
																		name: "HEXDIG",
																	},
																	&ruleRefExpr{
																		pos:  position{line: 202, col: 28, offset: 5032},
																		name: "HEXDIG",
																	},
																	&ruleRefExpr{
																		pos:  position{line: 202, col: 35, offset: 5039},
																		name: "HEXDIG",
																	},
																	&ruleRefExpr{
																		pos:  position{line: 202, col: 42, offset: 5046},
																		name: "HEXDIG",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "escape",
			pos:  position{line: 212, col: 1, offset: 5278},
			expr: &litMatcher{
				pos:        position{line: 213, col: 5, offset: 5289},
				val:        "\\",
				ignoreCase: false,
				want:       "\"\\\\\"",
			},
		},
		{
			name: "quotation_mark",
			pos:  position{line: 215, col: 1, offset: 5295},
			expr: &litMatcher{
				pos:        position{line: 216, col: 5, offset: 5314},
				val:        "\"",
				ignoreCase: false,
				want:       "\"\\\"\"",
			},
		},
		{
			name: "unescaped",
			pos:  position{line: 218, col: 1, offset: 5319},
			expr: &charClassMatcher{
				pos:        position{line: 219, col: 5, offset: 5333},
				val:        "[^\\\\\"]",
				chars:      []rune{'\\', '"'},
				ignoreCase: false,
				inverted:   true,
			},
		},
		{
			name: "DIGIT",
			pos:  position{line: 224, col: 1, offset: 5439},
			expr: &charClassMatcher{
				pos:        position{line: 224, col: 10, offset: 5448},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HEXDIG",
			pos:  position{line: 225, col: 1, offset: 5454},
			expr: &charClassMatcher{
				pos:        position{line: 225, col: 10, offset: 5463},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
	},
}

func (c *current) onJSON_text1(value interface{}) (interface{}, error) {
	return value, nil
}

func (p *parser) callonJSON_text1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onJSON_text1(stack["value"])
}

func (c *current) onobject12(m interface{}) (interface{}, error) {
	return m, nil
}

func (p *parser) callonobject12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onobject12(stack["m"])
}

func (c *current) onobject6(head, tail interface{}) (interface{}, error) {
	s := bytes.NewBuffer(head.([]byte))
	for _, v := range tail.([]interface{}) {
		s.WriteByte('_')
		s.Write(v.([]byte))
	}
	return s.Bytes(), nil

}

func (p *parser) callonobject6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onobject6(stack["head"], stack["tail"])
}

func (c *current) onobject1(members interface{}) (interface{}, error) {
	s := bytes.NewBuffer([]byte(`X.`))
	if members != nil {
		s.Write(members.([]byte))
	}
	s.Write([]byte(`.X`))
	return s.Bytes(), nil

}

func (p *parser) callonobject1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onobject1(stack["members"])
}

func (c *current) onmember1(name, value interface{}) (interface{}, error) {
	return append(append(name.([]byte), '-'), value.([]byte)...), nil

}

func (p *parser) callonmember1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onmember1(stack["name"], stack["value"])
}

func (c *current) onarray12(v interface{}) (interface{}, error) {
	return v, nil
}

func (p *parser) callonarray12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray12(stack["v"])
}

func (c *current) onarray6(head, tail interface{}) (interface{}, error) {
	s := bytes.NewBuffer(head.([]byte))
	for _, v := range tail.([]interface{}) {
		s.WriteByte('_')
		s.Write(v.([]byte))
	}
	return s.Bytes(), nil

}

func (p *parser) callonarray6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray6(stack["head"], stack["tail"])
}

func (c *current) onarray1(values interface{}) (interface{}, error) {
	s := bytes.NewBuffer([]byte(`I.`))
	if values != nil {
		s.Write(values.([]byte))
	}
	s.Write([]byte(`.I`))
	return s.Bytes(), nil

}

func (p *parser) callonarray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray1(stack["values"])
}

func (c *current) onnumber1(sign, i, f, e interface{}) (interface{}, error) {
	var s bytes.Buffer
	for _, v := range []interface{}{sign, i, f, e} {
		if v != nil {
			s.Write(v.([]byte))
		}
	}
	return s.Bytes(), nil

}

func (p *parser) callonnumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnumber1(stack["sign"], stack["i"], stack["f"], stack["e"])
}

func (c *current) onexp1(sign, i interface{}) (interface{}, error) {
	s := bytes.NewBuffer([]byte{'e'})
	if sign != nil {
		s.Write(sign.([]byte))
	}
	for _, v := range i.([]interface{}) {
		s.Write(v.([]byte))
	}
	return s.Bytes(), nil

}

func (p *parser) callonexp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onexp1(stack["sign"], stack["i"])
}

func (c *current) onfrac1() (interface{}, error) {
	return c.text, nil
}

func (p *parser) callonfrac1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onfrac1()
}

func (c *current) onint3() (interface{}, error) {
	return c.text, nil
}

func (p *parser) callonint3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onint3()
}

func (c *current) onplus1() (interface{}, error) {
	return []byte(nil), nil
}

func (p *parser) callonplus1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplus1()
}

func (c *current) onstr1(chars interface{}) (interface{}, error) {
	s := bytes.NewBuffer([]byte(`w.`))
	for _, v := range chars.([]interface{}) {
		s.Write(v.([]byte))
	}
	s.Write([]byte(`.w`))
	return s.Bytes(), nil

}

func (p *parser) callonstr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstr1(stack["chars"])
}

func (c *current) onchar2() (interface{}, error) {
	return []byte(".."), nil
}

func (p *parser) callonchar2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar2()
}

func (c *current) onchar4() (interface{}, error) {
	return []byte("._"), nil
}

func (p *parser) callonchar4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar4()
}

func (c *current) onchar6() (interface{}, error) {
	return []byte("_"), nil
}

func (p *parser) callonchar6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar6()
}

func (c *current) onchar14() (interface{}, error) {
	return []byte(`"`), nil
}

func (p *parser) callonchar14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar14()
}

func (c *current) onchar16() (interface{}, error) {
	return []byte("\\"), nil
}

func (p *parser) callonchar16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar16()
}

func (c *current) onchar18() (interface{}, error) {
	return []byte("/"), nil
}

func (p *parser) callonchar18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar18()
}

func (c *current) onchar20() (interface{}, error) {
	return []byte(".b"), nil
}

func (p *parser) callonchar20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar20()
}

func (c *current) onchar22() (interface{}, error) {
	return []byte(".f"), nil
}

func (p *parser) callonchar22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar22()
}

func (c *current) onchar24() (interface{}, error) {
	return []byte(".n"), nil
}

func (p *parser) callonchar24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar24()
}

func (c *current) onchar26() (interface{}, error) {
	return []byte(".r"), nil
}

func (p *parser) callonchar26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar26()
}

func (c *current) onchar28() (interface{}, error) {
	return []byte(".t"), nil
}

func (p *parser) callonchar28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar28()
}

func (c *current) onchar30(digits interface{}) (interface{}, error) {
	s := bytes.NewBuffer([]byte(`.u`))
	for _, v := range digits.([]interface{}) {
		s.Write(v.([]byte))
	}
	return s.Bytes(), nil

}

func (p *parser) callonchar30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar30(stack["digits"])
}

func (c *current) onchar9(sequence interface{}) (interface{}, error) {
	return sequence, nil
}

func (p *parser) callonchar9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchar9(stack["sequence"])
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return strconv.Itoa(p.line) + ":" + strconv.Itoa(p.col) + " [" + strconv.Itoa(p.offset) + "]"
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
	want       string
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

var statePool = &sync.Pool{
	New: func() interface{} { return make(storeDict) },
}

func (sd storeDict) Discard() {
	for k := range sd {
		delete(sd, k)
	}
	statePool.Put(sd)
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	state := statePool.Get().(storeDict)
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state.Discard()
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return strings.Join(list[:len(list)-1], sep) + " " + lastSep + " " + list[len(list)-1]
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, lit.want)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, lit.want)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.cloneState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restoreState(state)
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}



