package vtt_parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type VTTFile struct {
	Cues []*Cue
}

type Cue struct {
	Speaker string
	Start   string
	End     string
	Text    string
}

func ParseVTT(r io.Reader) (*VTTFile, error) {
	vttFile := &VTTFile{}

	scanner := bufio.NewScanner(r)

	// Skip header and empty line that suceeds it
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	// Parse all cues
	for scanner.Scan() {
		cue, err := parseCue(scanner)
		if err != nil {
			return nil, err
		}

		// If the speaker is the same as the previous one, consolidate them
		if len(vttFile.Cues) > 0 && vttFile.Cues[len(vttFile.Cues)-1].Speaker == cue.Speaker {
			prevCue := vttFile.Cues[len(vttFile.Cues)-1]
			prevCue.Text += "\n" + cue.Text
			prevCue.End = cue.End
		} else {
			vttFile.Cues = append(vttFile.Cues, cue)
		}

		// Parse empty line that suceeds cue
		if !scanner.Scan() && scanner.Err() != nil {
			return nil, scanner.Err()
		}
	}

	return vttFile, nil
}

func parseCue(scanner *bufio.Scanner) (*Cue, error) {
	c := &Cue{}

	// Ignore index (first line)

	// Parse timestamps (second line)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	timestampsLine := strings.TrimSpace(scanner.Text())
	timestamps := strings.Split(timestampsLine, " --> ")
	if len(timestamps) != 2 {
		return nil, fmt.Errorf("invalid timestamp format: %s", timestampsLine)
	}
	c.Start = timestamps[0]
	c.End = timestamps[1]

	// Parse speaker and text (third line)
	// The third line looks like this: <SPEAKER>: <TEXT>
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	speakerAndTextLine := strings.TrimSpace(scanner.Text())
	speakerAndTextParts := strings.SplitN(speakerAndTextLine, ":", 2)
	if len(speakerAndTextParts) != 2 {
		return nil, fmt.Errorf("invalid cue format: %s", scanner.Text())
	}

	c.Speaker = strings.TrimSpace(speakerAndTextParts[0])
	c.Text = strings.TrimSpace(speakerAndTextParts[1])

	return c, nil
}
