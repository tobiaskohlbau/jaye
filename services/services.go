// Copyright (c) 2017 Tobias Kohlbau
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package services

import (
	"context"
	"io"
)

// Service describes an interface for interacting with a video service.
type Service interface {
	Search(ctx context.Context, query string) ([]string, error)
	Info(ctx context.Context, id string) (VideoInfo, error)
	AudioFile(ctx context.Context, id string) (io.ReadCloser, error)
	VideoFile(ctx context.Context, id string) (io.ReadCloser, error)
	List(ctx context.Context) ([]VideoInfo, error)
}

type VideoInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Service   string `json:"service"`
}
