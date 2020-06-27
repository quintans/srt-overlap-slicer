package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Subtitle struct {
	idx      int
	fromTime time.Duration
	toTime   time.Duration
	text     string
}

var timeFramePattern, _ = regexp.Compile(`(\d+):(\d+):(\d+),(\d+) --> (\d+):(\d+):(\d+),(\d+)`)

func getDuration(parts []string) time.Duration {
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	second, _ := strconv.Atoi(parts[2])
	millisecond, _ := strconv.Atoi(parts[3])
	return time.Millisecond*time.Duration(millisecond) +
		time.Second*time.Duration(second) +
		time.Minute*time.Duration(minute) +
		time.Hour*time.Duration(hour)
}

func printDuration(duration time.Duration) string {
	hour := duration / time.Hour
	duration -= hour * time.Hour
	minute := duration / time.Minute
	duration -= minute * time.Minute
	second := duration / time.Second
	duration -= second * time.Second
	millisecond := duration / time.Millisecond
	return fmt.Sprintf(`%02d:%02d:%02d,%03d`, hour, minute, second, millisecond)
}

func readOneSubtitle(scanner *bufio.Scanner) (*Subtitle, error) {
	// read idx
	if !scanner.Scan() {
		return nil, nil
	}
	idxRaw := scanner.Text()
	idx, err := strconv.Atoi(idxRaw)
	if err != nil {
		return nil, errors.New("invalid subtitle index")
	}
	// read timing
	if !scanner.Scan() {
		return nil, errors.New("could not find subtitle timing")
	}
	timing := timeFramePattern.FindStringSubmatch(scanner.Text())
	if timing == nil {
		return nil, errors.New("invalid subtitle timing")
	}
	fromTime := getDuration(timing[1:5])
	toTime := getDuration(timing[5:9])
	// read content
	if !scanner.Scan() {
		return nil, errors.New("could not find subtitle text")
	}
	content := scanner.Text()
	for scanner.Scan() && scanner.Text() != "" {
		content += "\n"
		content += scanner.Text()
	}
	subtitle := &Subtitle{idx, fromTime, toTime, content}
	return subtitle, nil
}

func writeOneSubtitle(w io.Writer, subtitle *Subtitle, idx int) error {
	_, err := fmt.Fprint(w,
		idx, "\n",
		printDuration(subtitle.fromTime), " --> ", printDuration(subtitle.toTime), "\n",
		subtitle.text, "\n\n")
	return err
}

func slicer(r io.Reader, w io.Writer) {
	subtitles := readAllSubtittles(r)
	// find all intervals
	intervals := make([]time.Duration, len(subtitles)*2)
	for k, sub := range subtitles {
		intervals[k*2] = sub.fromTime
		intervals[k*2+1] = sub.toTime
	}
	// sort intervals
	sort.SliceStable(intervals, func(i, j int) bool {
		return intervals[i].Milliseconds() < intervals[j].Milliseconds()
	})

	// split by intervals
	newSubs := []*Subtitle{}
	size := len(intervals)
	for i := 1; i < size; i++ {
		low := intervals[i-1]
		high := intervals[i]
		// ignore lesser than 200ms
		if high.Milliseconds()-low.Milliseconds() < 100 {
			continue
		}
		txt := &bytes.Buffer{}
		newSub := &Subtitle{
			idx:      i,
			fromTime: low,
			toTime:   high,
		}
		// concat subtitles that fall inside the interval
		for _, sub := range subtitles {
			if low >= sub.fromTime && high <= sub.toTime {
				if txt.Len() > 0 {
					txt.WriteString("\n- ")
				}
				txt.WriteString(sub.text)
			}
		}
		if txt.Len() > 0 {
			newSub.text = txt.String()
			newSubs = append(newSubs, newSub)
		}
	}

	size = len(newSubs)
	for i := 0; i < size; i++ {
		writeOneSubtitle(w, newSubs[i], i+1)
	}
}

func readAllSubtittles(r io.Reader) []*Subtitle {
	scanner := bufio.NewScanner(r)
	subtitles := []*Subtitle{}

	for {
		subtitle, err := readOneSubtitle(scanner)
		if err != nil {
			panic(err)
		}
		if subtitle == nil {
			break
		}
		subtitle.text = strings.Trim(subtitle.text, "\n ")
		if len(subtitle.text) == 0 { // skip over empty subtitles
			continue
		}
		subtitles = append(subtitles, subtitle)
	}
	return subtitles
}

func main() {
	if len(os.Args) < 2 {
		println("Provide a subtitle file to fix.\ne.g. srt-overlap-fixer mysubtitle.srt")
		return
	}

	original := os.Args[1]
	fixed := original + ".fixed"

	file, _ := os.Open(original)
	newFile, _ := os.Create(fixed)
	defer file.Close()
	defer newFile.Close()

	slicer(file, newFile)

	os.Rename(original, original+".bak")
	os.Rename(fixed, original)
}
