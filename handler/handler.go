// Copyright (c) 2017 Tobias Kohlbau
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"kohlbau.de/x/jaye/services"
)

func New(yt services.Service) http.Handler {
	mux := http.NewServeMux()
	h := handler{yt: yt}
	mux.HandleFunc("/search", h.serviceHandler(search))
	mux.HandleFunc("/info", h.serviceHandler(info))
	mux.HandleFunc("/video", h.serviceHandler(video))
	mux.HandleFunc("/audio", h.serviceHandler(audio))
	mux.HandleFunc("/list", h.serviceHandler(list))
	return mux
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

type handler struct {
	yt services.Service
}

type writer interface {
	io.Writer
	Header() http.Header
}

func (h handler) serviceHandler(fn func(writer, *http.Request, services.Service) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		service := r.FormValue("service")
		if service == "" {
			http.Error(w, "service id not supplied", http.StatusBadRequest)
			return
		}

		var data interface{}
		var status int
		var err error

		switch service {
		case "youtube":
			data, status, err = fn(w, r, h.yt)
		default:
			data = nil
			status = http.StatusBadRequest
			err = errors.New("service not found")
		}

		if err != nil {
			data = err.Error()
		}

		// used if handler does not use json
		if data == nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(response{Data: data, Success: err == nil})
		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}
	}
}

func info(w writer, r *http.Request, s services.Service) (interface{}, int, error) {
	id := r.FormValue("id")
	if id == "" {
		return nil, http.StatusBadRequest, errors.New("no id supplied")
	}

	vid, err := s.Info(r.Context(), url.QueryEscape(id))
	if err != nil {
		log.Printf("failed to retrieve video info: %v", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to retrieve video info: %s", url.QueryEscape(id))
	}

	return vid, http.StatusOK, nil
}

func search(w writer, r *http.Request, s services.Service) (interface{}, int, error) {
	q := r.FormValue("q")
	if q == "" {
		return nil, http.StatusBadRequest, errors.New("missing query parameter")
	}

	vids, err := s.Search(r.Context(), url.QueryEscape(q))
	if err != nil {
		log.Printf("failed to find youtube video: %v", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to find youtube video for query: %s", url.QueryEscape(q))
	}

	return vids, http.StatusOK, nil
}

func video(w writer, r *http.Request, s services.Service) (interface{}, int, error) {
	id := r.FormValue("id")
	if id == "" {
		return nil, http.StatusBadRequest, errors.New("no id supplied")
	}

	rc, err := s.VideoFile(r.Context(), id)
	if err != nil {
		log.Printf("failed to retrieve video file: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve video file")
	}
	defer rc.Close()

	vi, err := s.Info(r.Context(), id)
	if err != nil {
		log.Printf("failed to retrieve video info: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve video info")
	}

	if r.Header.Get("Accept") == "application/json" {
		return vi, http.StatusOK, nil
	}

	w.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.mp4\"", vi.Title))
	io.Copy(w, rc)
	return nil, http.StatusOK, nil
}

func audio(w writer, r *http.Request, s services.Service) (interface{}, int, error) {
	id := r.FormValue("id")
	if id == "" {
		return nil, http.StatusBadRequest, errors.New("no id supplied")
	}

	rc, err := s.AudioFile(r.Context(), id)
	if err != nil {
		log.Printf("failed to retrieve audio file: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve audio file")
	}
	defer rc.Close()

	vi, err := s.Info(r.Context(), id)
	if err != nil {
		log.Printf("failed to retrieve video info: %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve video info")
	}

	if r.Header.Get("Accept") == "application/json" {
		return vi, http.StatusOK, nil
	}

	w.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.mp3\"", vi.Title))
	io.Copy(w, rc)
	return nil, http.StatusOK, nil
}

func list(w writer, r *http.Request, s services.Service) (interface{}, int, error) {
	videos, err := s.List(r.Context())
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to list videos: %v", err)
	}
	return videos, http.StatusOK, nil
}
