import { Component, NgModule, OnInit } from '@angular/core';
import { YouTubePlayerModule } from '@angular/youtube-player';

@Component({
  standalone: true,
  imports: [YouTubePlayerModule],
  templateUrl: './youtube-player.component.html',
  selector: 'app-youtube-player',
})
export class YoutubePlayerComponent implements OnInit {
  ngOnInit(): void {
    const tag = document.createElement('script');
    tag.src = 'https://www.youtube.com/iframe_api';
    document.body.appendChild(tag);
  }
}
