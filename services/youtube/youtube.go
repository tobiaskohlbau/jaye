// Copyright (c) 2017 Tobias Kohlbau
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/rylio/ytdl"
	"kohlbau.de/x/jaye/multimedia"
	"kohlbau.de/x/jaye/services"
)

type youtubeService struct {
	m            sync.Mutex
	youtubeURL   string
	youtubeToken string
	videoPath    string
	cl           http.Client
	converter    multimedia.Converter
}

func New(youtubeURL, youtubeToken, videoPath string) services.Service {
	return youtubeService{
		youtubeURL:   youtubeURL,
		youtubeToken: youtubeToken,
		videoPath:    videoPath,
		cl:           http.Client{},
		converter:    multimedia.NewFFMPEG(),
	}
}

func (s youtubeService) Search(ctx context.Context, query string) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/search/?q=%s&part=snippet&type=video&key=%s", s.youtubeURL, query, s.youtubeToken), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req = req.WithContext(ctx)
	resp, err := s.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query youtube api: %v", err)
	}
	defer resp.Body.Close()

	var search search
	if err := json.NewDecoder(resp.Body).Decode(&search); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %v", err)
	}

	var vids []string

	for _, vid := range search.Items {
		vids = append(vids, vid.ID.VideoID)
	}

	return vids, nil
}

func (s youtubeService) Info(ctx context.Context, id string) (services.VideoInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/videos/?id=%s&part=snippet&key=%s", s.youtubeURL, id, s.youtubeToken), nil)
	if err != nil {
		return services.VideoInfo{}, fmt.Errorf("failed to create request: %v", err)
	}

	req = req.WithContext(ctx)
	resp, err := s.cl.Do(req)
	if err != nil {
		return services.VideoInfo{}, fmt.Errorf("failed to query youtube api: %v", err)
	}
	defer resp.Body.Close()

	var vid video
	if err := json.NewDecoder(resp.Body).Decode(&vid); err != nil {
		return services.VideoInfo{}, fmt.Errorf("failed to decode video response: %v", err)
	}

	if len(vid.Items) == 0 {
		return services.VideoInfo{}, fmt.Errorf("failed to find video for id: %s", id)
	}

	return services.VideoInfo{
		ID:        vid.Items[0].ID,
		Title:     vid.Items[0].Snippet.Title,
		URL:       "https://youtube.com/watch?v=" + vid.Items[0].ID,
		Thumbnail: vid.Items[0].Snippet.Thumbnails.High.URL,
		Service:   "youtube",
	}, nil
}

func (s youtubeService) download(vid *ytdl.VideoInfo, fm ytdl.Format, id, name string) (io.ReadCloser, error) {
	dir := path.Join(s.videoPath, id)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, errors.New("failed to create video directory")
	}

	vidp := path.Join(dir, fmt.Sprintf("%s.%s", name, fm.Extension))
	if _, err := os.Stat(vidp); err == nil {
		log.Printf("video already exists: %v", id)
		vidf, err := os.OpenFile(vidp, os.O_RDONLY, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to open pre downloaded file: %v", err)
		}
		return vidf, nil
	}

	file, err := os.OpenFile(vidp, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for download: %v", err)
	}

	log.Printf("downloading %s: %v", name, id)

	if err := vid.Download(fm, file); err != nil {
		file.Close()
		if err := os.Remove(vidp); err != nil {
			log.Printf("failed to delete video file: %v", err)
		}
		return nil, fmt.Errorf("failed to download video file: %v", err)
	}
	file.Sync()
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek video file: %v", err)
	}

	log.Printf("finished downloading %s: %v", name, id)

	return file, nil
}

func (s youtubeService) VideoFile(ctx context.Context, id string) (io.ReadCloser, error) {
	s.m.Lock()
	defer s.m.Unlock()

	// fetch video info
	vid, err := ytdl.GetVideoInfoFromID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find video by id: %v", err)
	}

	fm := vid.Formats.Copy()
	fm.Sort(ytdl.FormatResolutionKey, true)
	if len(fm) == 0 {
		return nil, errors.New("failed to retrieve video format")
	}

	vrc, err := s.download(vid, fm[0], id, "video")
	if err != nil {
		return nil, err
	}
	defer vrc.Close()

	fm = vid.Formats.Copy()
	fm.Sort(ytdl.FormatAudioBitrateKey, true)
	if len(fm) == 0 {
		return nil, errors.New("failed to retrieve video format")
	}

	arc, err := s.download(vid, fm[0], id, "audio")
	if err != nil {
		return nil, err
	}
	defer arc.Close()

	vidp := path.Join(s.videoPath, id, "combined.mp4")
	if _, err := os.Stat(vidp); err == nil {
		log.Printf("video already exists: %v", id)
		vidf, err := os.OpenFile(vidp, os.O_RDONLY, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to open pre downloaded file: %v", err)
		}
		return vidf, nil
	}

	file, err := os.OpenFile(vidp, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for merging: %v", err)
	}

	if err := s.converter.Merge(ctx, vrc, arc, file); err != nil {
		return nil, fmt.Errorf("failed to merge video and audio files: %v", err)
	}
	file.Sync()
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek combined file: %v", err)
	}

	return file, nil
}

func (s youtubeService) AudioFile(ctx context.Context, id string) (io.ReadCloser, error) {
	s.m.Lock()
	defer s.m.Unlock()

	// fetch video info
	vid, err := ytdl.GetVideoInfoFromID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find video by id: %v", err)
	}

	fm := vid.Formats.Copy()
	fm.Sort(ytdl.FormatAudioEncodingKey, true)
	if len(fm) == 0 {
		return nil, errors.New("failed to retrieve video format")
	}

	rc, err := s.download(vid, fm[0], id, "audio")
	if err != nil {
		return nil, err
	}

	audp := path.Join(s.videoPath, id, "audio.mp3")
	if _, err := os.Stat(audp); err == nil {
		log.Printf("audio already exists: %v", id)
		audf, err := os.OpenFile(audp, os.O_RDONLY, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to open pre converted file: %v", err)
		}
		return audf, nil
	}

	file, err := os.Create(audp)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio file: %v", err)
	}

	log.Printf("converting video: %v", id)

	if err := s.converter.Convert(context.Background(), rc, file); err != nil {
		file.Close()
		if err := os.Remove(audp); err != nil {
			log.Printf("failed to delete audio file: %v", err)
		}
		return nil, fmt.Errorf("failed to convert video: %v", err)
	}
	file.Sync()
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek audio file: %v", err)
	}

	log.Printf("finished converting video: %v", id)

	return file, nil
}

func (s youtubeService) List(ctx context.Context) ([]services.VideoInfo, error) {
	folders, err := ioutil.ReadDir(s.videoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve video list: %v", err)
	}

	var vids []services.VideoInfo
	var tims []time.Time
	for _, f := range folders {
		vi, err := s.Info(ctx, f.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve video info: %v", err)
		}
		vids = append(vids, vi)
		tims = append(tims, f.ModTime())
	}

	sort.Slice(vids, func(i, j int) bool { return tims[i].UnixNano() < tims[j].UnixNano() })

	return vids, nil
}

type search struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string `json:"channelTitle"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

type video struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
				Standard struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"standard"`
				Maxres struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle         string   `json:"channelTitle"`
			Tags                 []string `json:"tags"`
			CategoryID           string   `json:"categoryId"`
			LiveBroadcastContent string   `json:"liveBroadcastContent"`
			Localized            struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"localized"`
			DefaultAudioLanguage string `json:"defaultAudioLanguage"`
		} `json:"snippet"`
	} `json:"items"`
}
