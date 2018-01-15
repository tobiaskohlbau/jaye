import { Pipe, PipeTransform, Component, OnInit, OnDestroy, NgModule, Directive, ElementRef, AfterViewChecked, Input, HostListener } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdCardModule, MdButtonModule, MdGridListModule } from '@angular/material';

import { FlexLayoutModule } from '@angular/flex-layout';

import { VideoService, VideoInfo } from '../../shared/video';

import { Observable } from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';


// @Pipe({
//   name: "sort"
// })
// export class ArraySortPipe implements PipeTransform {
//   transform(array: Array<string>, args: string): Array<VideoItem> {
//     array.sort((a: VideoItem, b: VideoItem) => {
//       if (a.id < b.id) {
//         return -1;
//       } else if (a.id > b.id) {
//         return 1;
//       } else {
//         return 0;
//       }
//     });
//     return array;
//   }
// }

@Directive({
  selector: '[matchHeight]'
})
export class MatchHeight implements AfterViewChecked {
  @Input()
  matchHeight: string;

  constructor(private el: ElementRef) {
  }

  ngAfterViewChecked() {
    this.match(this.el.nativeElement, this.matchHeight);
  }

  match(parent: HTMLElement, className: string) {
    const children = parent.getElementsByClassName(className);

    if (!children) return;

    Array.from(children).forEach((x: HTMLElement) => {
      x.style.height = 'initial';
    });

    const itemHeights = Array.from(children)
      .map(x => x.getBoundingClientRect().height);

    const maxHeight = itemHeights.reduce((prev, curr) => {
      return curr > prev ? curr : prev;
    }, 0);

    Array.from(children)
      .forEach((x: HTMLElement) => x.style.height = `${maxHeight}px`);
  }

  @HostListener('window:resize')
  onResize() {
    this.match(this.el.nativeElement, this.matchHeight);
  }
}

@Component({
  selector: 'app-video-list',
  templateUrl: './video-list.html',
  styleUrls: ['./video-list.scss'],
})
export class VideoList implements OnInit, OnDestroy {
  videos: VideoInfo[];

  constructor(private videoService: VideoService) { }

  ngOnInit() {
    this.getVideos();
  }

  ngOnDestroy() {
  }

  getVideos() {
    this.videoService.list()
      .subscribe(
      videos => {
        if (this.videos == null || this.videos.length != videos.length)
          this.videos = videos;
      });
  }

  addVideo(video: VideoInfo) {
    this.videos.push(video);
  }
}

@NgModule({
  imports: [FlexLayoutModule, CommonModule, MdCardModule, MdButtonModule, MdGridListModule],
  exports: [VideoList],
  declarations: [VideoList, MatchHeight],
  providers: [VideoService],
})
export class VideoListModule { }
