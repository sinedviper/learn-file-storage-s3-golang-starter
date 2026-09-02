package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"os/exec"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var result ffprobeOutput
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return "", err
	}

	if len(result.Streams) == 0 {
		return "", errors.New("no video streams found")
	}

	width := result.Streams[0].Width
	height := result.Streams[0].Height

	ratio := float64(width) / float64(height)

	const tolerance = 0.1

	switch {
	case math.Abs(ratio-16.0/9.0) < tolerance:
		return "16:9", nil
	case math.Abs(ratio-9.0/16.0) < tolerance:
		return "9:16", nil
	default:
		return "other", nil
	}
}

func processVideoForFastStart(filePath string) (string, error) {
	outputFilePath := filePath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFilePath)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return outputFilePath, nil
}
