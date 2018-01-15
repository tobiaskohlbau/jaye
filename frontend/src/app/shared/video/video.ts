import { Injectable } from '@angular/core';

import { Location } from '@angular/common';
import { HttpClient, HttpInterceptor, HttpHandler, HttpEvent, HttpRequest } from '@angular/common/http';

import { environment } from '../../../environments/environment';

import { Observable } from 'rxjs/Observable';
import 'rxjs/Rx';

export interface Response {
    success: string;
    response: Object;
}

export interface VideoInfo {
    id: string;
    title: string;
    url: string;
    thumbnail: string;
    service: string;
}

Injectable()
export class VideoInterceptor implements HttpInterceptor {
    constructor() { }

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        let url = environment.apiEndpoint + req.url;
        const nReq = req.clone({ headers: req.headers.set('Accept', 'application/json'), url: url });
        return next.handle(nReq);
    }
}

@Injectable()
export class VideoService {
    constructor(private http: HttpClient,
        private location: Location) {
    }

    videoURL(id: string): string {
        return environment.apiEndpoint + "/video?service=youtube&id=" + id;
    }

    audioURL(id: string): string {
        return environment.apiEndpoint + "/audio?service=youtube&id=" + id;
    }

    video(id: string): Observable<VideoInfo> {
        return this.http.get<Response>('/video?service=youtube&id=' + id).map(res => <VideoInfo>res.response);
    }

    list(): Observable<Array<VideoInfo>> {
        return this.http.get<Response>("/list?service=youtube")
            .map(res => {
                return <Array<VideoInfo>>res.response;
            });
    }

    info(id: string): Observable<VideoInfo> {
        return this.http.get<Response>("/info?service=youtube&id=" + id).map(res => <VideoInfo>res.response);
    }

    search(query: string): Observable<Array<String>> {
        if (query == "")
            return Observable.create(obs => obs.next(Array<String>()));

        let id = this.idFromURL(query);
        if (id != "") {
            return Observable.create(obs => obs.next(Array<String>(id)));
        }

        return this.http.get<Response>("/search?service=youtube&q=" + encodeURI(query))
        .map(data => {
            return <Array<string>>data.response;
        });
    }

    idFromURL(url: string): string {
        let m = url.match(/(^|=|\/)([0-9A-Za-z_-]{11})(\/|&|$|\?|#)/)
        if (m == null) {
            return "";
        }
        return m[2];
    }
}
