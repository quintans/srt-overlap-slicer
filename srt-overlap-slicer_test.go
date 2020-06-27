package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOneOverlap(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name: "no overlap",
			input: `1
00:00:00,000 --> 00:00:01,000
This is the first line

2
00:00:03,000 --> 00:00:04,000
And this is the second line
`,
			expect: `1
00:00:00,000 --> 00:00:01,000
This is the first line

2
00:00:03,000 --> 00:00:04,000
And this is the second line

`,
		},
		{
			name: "right edge overlap",
			input: `1
00:00:00,000 --> 00:00:02,000
This is the first line

2
00:00:01,000 --> 00:00:04,000
And this is the second line
`,
			expect: `1
00:00:00,000 --> 00:00:01,000
This is the first line

2
00:00:01,000 --> 00:00:02,000
This is the first line
- And this is the second line

3
00:00:02,000 --> 00:00:04,000
And this is the second line

`,
		},
		{
			name: "multiple overlap",
			input: `1
00:00:00,000 --> 00:00:02,000
FIRST

2
00:00:01,000 --> 00:00:08,000
SECOND

3
00:00:03,000 --> 00:00:04,000
THIRD

4
00:00:05,000 --> 00:00:06,000
FOURTH
`,
			expect: `1
00:00:00,000 --> 00:00:01,000
FIRST

2
00:00:01,000 --> 00:00:02,000
FIRST
- SECOND

3
00:00:02,000 --> 00:00:03,000
SECOND

4
00:00:03,000 --> 00:00:04,000
SECOND
- THIRD

5
00:00:04,000 --> 00:00:05,000
SECOND

6
00:00:05,000 --> 00:00:06,000
SECOND
- FOURTH

7
00:00:06,000 --> 00:00:08,000
SECOND

`,
		},
		{
			name: "complete overlap",
			input: `1
00:00:00,000 --> 00:00:04,000
This is the first line

2
00:00:01,000 --> 00:00:03,000
And this is the second line
`,
			expect: `1
00:00:00,000 --> 00:00:01,000
This is the first line

2
00:00:01,000 --> 00:00:03,000
This is the first line
- And this is the second line

3
00:00:03,000 --> 00:00:04,000
This is the first line

`,
		},
		{
			name: "same overlap",
			input: `1
00:00:00,000 --> 00:00:04,000
This is the first line

2
00:00:00,000 --> 00:00:04,000
And this is the second line
`,
			expect: `1
00:00:00,000 --> 00:00:04,000
This is the first line
- And this is the second line

`,
		},
		{
			name: "200ms diff overlap",
			input: `1
00:00:00,000 --> 00:00:04,000
This is the first line

2
00:00:00,050 --> 00:00:03,950
And this is the second line
`,
			expect: `1
00:00:00,050 --> 00:00:03,950
This is the first line
- And this is the second line

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			w := &bytes.Buffer{}
			slicer(r, w)

			result := w.String()
			fmt.Println(result)
			require.Equal(t, tt.expect, result)
		})
	}
}
