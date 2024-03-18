
import { Component, ElementRef, HostListener, Input, NgModule, OnInit, ViewChild } from '@angular/core';
import { YouTubePlayerModule } from '@angular/youtube-player';

@Component({
  standalone: true,
  imports: [YouTubePlayerModule],
  templateUrl: './youtube-player.component.html',
  selector: 'app-youtube-player',
})
export class YoutubePlayerComponent implements OnInit {
  @Input() videoId = '';

  playerWidth = 0;

  @HostListener('window:resize', ['$event']) onResize() {
    this.playerWidth = this.container.nativeElement.offsetWidth;
  }

  @ViewChild('container', { static: true }) container!: ElementRef;

  ngOnInit(): void {
    const tag = document.createElement('script');
    tag.src = 'https://www.youtube.com/iframe_api';
    document.body.appendChild(tag);
    this.playerWidth = this.container.nativeElement.offsetWidth;
  }
}
