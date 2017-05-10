package progress

import (
	"io"
	"os"
	"sync"
	"strconv"
	"bytes"
	"fmt"
)

var (
	LeftEdge = "["
	RightEdge = "]"
	Fill = "\u2588"
	Empty = " "
	ESC = 27
)

type Bar struct {
	progress int
	total int
	writer io.Writer
	mutex sync.Mutex
	bar_width int
	prefix string
	left_bar_edge string
	right_bar_edge string
	fill_bar string
	empty_bar string
}

func NewBar(total int, name string) *Bar {
	return &Bar{
		progress: 0,
		total: total,
		writer: os.Stdout,
		prefix: name + " |",
		bar_width: 10,
		left_bar_edge: LeftEdge,
		right_bar_edge: RightEdge,
		fill_bar: Fill,
		empty_bar: Empty,
	}
}

func ManualBar(total int, writer io.Writer, prefix string) *Bar {
	return &Bar{
		progress: 0,
		total: total,
		writer: writer,
		prefix: prefix,
		bar_width: 10,
		left_bar_edge: LeftEdge,
		right_bar_edge: RightEdge,
		fill_bar: Fill,
		empty_bar: Empty,
	}
}

func (bar *Bar) display_bar() {
	fmt.Fprintf(bar.writer, "%c[%dA", ESC, 0)
	fmt.Fprintf(bar.writer, "%c[2K\r", ESC)
	percentage := int(float64(bar.progress) / float64(bar.total) * 100)
	partially_fill := percentage % 10
	full_fill := percentage / 10
	var buf bytes.Buffer
	buf.WriteString(bar.prefix + " | ")
	buf.WriteString(strconv.Itoa(bar.progress) + "/" + strconv.Itoa(bar.total) + " " + bar.left_bar_edge)
	for i := 0; i < full_fill; i++ {
		buf.WriteString(bar.fill_bar)
	}
	buf.WriteString(strconv.Itoa(partially_fill))
	empty_cells := bar.bar_width - full_fill - 1
	for i := 0; i < empty_cells; i++ {
		buf.WriteString(bar.empty_bar)
	}
	buf.WriteString(bar.right_bar_edge)
	buf.WriteString(" " + strconv.Itoa(percentage) + "%")
	buf.WriteString("\n")
	bar.writer.Write(buf.Bytes())
}

func (bar *Bar) Current() int {
	return bar.total
}

func (bar *Bar) Set(num int) bool {
	bar.progress = num
	if bar.progress >= bar.total || bar.progress < 0 {
		return false
	} else {
		bar.display_bar()
		return true
	}
}

func (bar *Bar) Add(num int) bool {
	bar.progress += num
	if bar.progress >= bar.total || bar.progress < 0 {
		return false
	} else {
		bar.display_bar()
		return true
	}
}

func (bar *Bar) Increase() bool {
	bar.progress += 1
	if bar.progress >= bar.total {
		return false
	} else {
		bar.display_bar()
		return true
	}
}