package multimedia

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type ffmpegConverter struct {
}

func NewFFMPEG() Converter {
	return ffmpegConverter{}
}

func (c ffmpegConverter) Convert(ctx context.Context, src io.Reader, dst io.Writer) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", "-", "-f", "mp3", "-")

	cmd.Stdout = dst
	cmd.Stdin = src
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg failed to convert video: %v", err)
	}
	return nil
}

func (c ffmpegConverter) Merge(ctx context.Context, video, audio io.Reader, dst io.Writer) error {
	vf, err := ioutil.TempFile("", "ytdl-video")
	if err != nil {
		return fmt.Errorf("failed to create tmp video file")
	}
	defer os.Remove(vf.Name())

	af, err := ioutil.TempFile("", "ytdl-audio")
	if err != nil {
		return fmt.Errorf("failed to create tmp video file")
	}
	defer os.Remove(af.Name())

	of, err := ioutil.TempFile("", "ytdl-merged")
	if err != nil {
		return fmt.Errorf("failed to create tmp merge file")
	}
	defer os.Remove(of.Name())

	if _, err := io.Copy(vf, video); err != nil {
		return fmt.Errorf("failed to copy video input to temp file: %v", err)
	}
	vf.Close()

	if _, err := io.Copy(af, audio); err != nil {
		return fmt.Errorf("failed to copy audio input to temp file: %v", err)
	}
	af.Close()

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", vf.Name(), "-i", af.Name(), "-f", "mp4", "-y", of.Name())
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to merge audio and video: %v", err)
	}

	if _, err := of.Seek(0, 0); err != nil {
		return fmt.Errorf("failed seeking merged file: %v", err)
	}
	if _, err := io.Copy(dst, of); err != nil {
		return fmt.Errorf("failed writing merged filo to dst: %v", err)
	}
	of.Close()

	return nil
}
