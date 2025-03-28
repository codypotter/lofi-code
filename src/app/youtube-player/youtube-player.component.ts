
import { ChangeDetectorRef, Component, ElementRef, HostListener, Input, ViewChild, afterRender } from '@angular/core';
import { YouTubePlayerModule } from '@angular/youtube-player';
import { NGXLogger } from 'ngx-logger';

@Component({
    imports: [YouTubePlayerModule],
    templateUrl: './youtube-player.component.html',
    selector: 'app-youtube-player'
})
export class YoutubePlayerComponent {
  @Input() videoId = '';

  playerWidth = 0;

  @HostListener('window:resize', ['$event']) onResize() {
    this.playerWidth = this.container.nativeElement.offsetWidth;
    this.cd.detectChanges()
  }

  @ViewChild('container', { static: true }) container!: ElementRef;

  constructor(private logger: NGXLogger, private cd: ChangeDetectorRef) {
    afterRender(() => {
      this.logger.trace('youtube-player: afterRender');
      const tag = document.createElement('script');
      tag.src = 'https://www.youtube.com/iframe_api';
      document.body.appendChild(tag);
      this.playerWidth = this.container.nativeElement.offsetWidth;
      this.cd.detectChanges();
      this.logger.trace('youtube-player: initialized');
    });
  }
}
