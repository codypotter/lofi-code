import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';

interface VideoResponse {
  id: string;
  title: string;
}

@Injectable({
  providedIn: 'root'
})
export class VideosService {

  constructor(private http: HttpClient) { }

  getFeaturedVideo() {
    return this.http.get<VideoResponse>(environment.baseUrl + '/api/featured-video');
  }
}
