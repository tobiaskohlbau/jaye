import { Component, NgModule } from '@angular/core';

import { DownloadModule } from '../download';
import { VideoListModule } from '../video-list';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.html',
  styleUrls: ['./dashboard.scss']
})
export class Dashboard { }

@NgModule({
  imports: [DownloadModule, VideoListModule],
  exports: [Dashboard],
  declarations: [Dashboard],
})
export class DashboardModule { }
