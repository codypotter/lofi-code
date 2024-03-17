import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostsService } from '../services/posts.service';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { RouterModule } from '@angular/router';
import { TagsComponent } from '../tags/tags.component';
import { YoutubePlayerComponent } from '../youtube-player/youtube-player.component';
import { VideosService } from '../services/videos.service';
import { Subject, forkJoin, takeUntil } from 'rxjs';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent, RouterModule, TagsComponent, YoutubePlayerComponent],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  posts: any[] = [];

  destroy$ = new Subject<void>();

  constructor(private postsService: PostsService, private breadcrumbService: BreadcrumbService, private videosService: VideosService) { }
  
  ngOnInit(): void {
    this.breadcrumbService.setBreadcrumbs([]);
    this.postsService.getN(3).pipe(takeUntil(this.destroy$)).subscribe((posts) => { this.posts = posts;});
    this.videosService.getVideos().pipe(takeUntil(this.destroy$)).subscribe((videos) => { console.error('videos', videos); })
    this.videosService.getFeaturedVideo().pipe(takeUntil(this.destroy$)).subscribe((video) => { console.error('video', video); })
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
