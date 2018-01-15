import { Component, NgModule, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormControl, NgForm, FormGroupDirective } from '@angular/forms';
import { MdCardModule, MdInputModule, MdButtonModule, MdIconModule, MdListModule, MdRadioModule } from '@angular/material';

import { HttpEventType } from '@angular/common/http';

import { VideoService, VideoInfo } from '../../shared/video';

import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-download',
  templateUrl: './download.html',
  styleUrls: ['./download.scss']
})
export class Download implements OnInit {
  @Input()
  selectedVideo: string;

  @Output()
  updated: EventEmitter<VideoInfo> = new EventEmitter();

  searchForm = new FormControl();
  videos: Array<VideoInfo>;
  error: string = "Invalid YouTube Video or Query";

  constructor(private videoService: VideoService) { }

  ngOnInit() {
    this.searchForm.valueChanges
      .debounceTime(200)
      .switchMap(query => {
        return this.videoService.search(query)
      })
      .flatMap((ids: Array<String>) => {
        return Observable.forkJoin(
          ids.map((id: string) => this.videoService.info(id)));
      })
      .subscribe(videos => {
        if (videos.length > 0)
          this.selectedVideo = videos[0].id;

        this.videos = videos;
      });
  }

  download(id: string): void {
    this.videos = null;
    this.selectedVideo = null;
    this.searchForm.setValue("");
    this.videoService.video(id)
      .subscribe(
      video => {
        this.updated.emit(video);
      });
  }
}

@NgModule({
  imports: [CommonModule, FormsModule, ReactiveFormsModule, MdCardModule, MdInputModule, MdButtonModule, MdIconModule, MdListModule, MdRadioModule],
  exports: [Download],
  declarations: [Download],
  providers: [VideoService],
})
export class DownloadModule { }
