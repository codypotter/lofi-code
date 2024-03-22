import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';

interface YouTubeResponse {
  items: {
    id: {
      videoId: string;
    };
  }[];
}

@Injectable({
  providedIn: 'root'
})
export class VideosService {

  constructor(private http: HttpClient) { }

  getFeaturedVideo() {
    return this.http.get<YouTubeResponse>('https://www.googleapis.com/youtube/v3/search', {
      params: {
        part: 'snippet',
        channelId: 'UCsPXgrtO5bTfdVdNLLB_Erw',
        maxResults: 1,
        order: 'date',
        key: environment.youtubeApiKey,
      }
    });
  }
}
