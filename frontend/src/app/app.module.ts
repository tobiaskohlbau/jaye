import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { Location, LocationStrategy, PathLocationStrategy } from '@angular/common';

import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import { AppComponent } from './app.component';

import 'hammerjs';

import { JAYE_ROUTES } from './routes';

import { VideoService, VideoInterceptor } from './shared/video';

import { NavBarModule } from './shared/navbar';
import { DashboardModule } from './components/dashboard';
import { VideoListModule } from './components/video-list';
import { DownloadModule } from './components/download';

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    RouterModule.forRoot(JAYE_ROUTES),

    NavBarModule,
    DashboardModule,
    VideoListModule,
    DownloadModule
  ],
  declarations: [
    AppComponent
  ],
  providers: [
    Location,
    {
      provide: LocationStrategy,
      useClass: PathLocationStrategy
    },
    VideoService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: VideoInterceptor,
      multi: true,
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
