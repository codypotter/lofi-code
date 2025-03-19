import { Component, OnInit } from '@angular/core';

import { PostsService } from '../services/posts.service';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { RouterModule } from '@angular/router';
import { TagsComponent } from '../tags/tags.component';
import { YoutubePlayerComponent } from '../youtube-player/youtube-player.component';
import { VideosService } from '../services/videos.service';
import { Subject, takeUntil } from 'rxjs';
import { NgxSkeletonLoaderModule } from 'ngx-skeleton-loader';
import { NGXLogger } from 'ngx-logger';
import { Title } from '@angular/platform-browser';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent, RouterModule, TagsComponent, YoutubePlayerComponent, NgxSkeletonLoaderModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  posts: any[] = [];

  destroy$ = new Subject<void>();

  featuredVideoId = '';

  featuredVideoTitle = '';

  constructor(
    private logger: NGXLogger,
    private postsService: PostsService,
    private breadcrumbService: BreadcrumbService,
    private videosService: VideosService,
    private title: Title,
  ) { }
  
  ngOnInit(): void {
    this.title.setTitle('lofi code');
    this.breadcrumbService.setBreadcrumbs([]);
    this.postsService.getN(3).subscribe((posts) => { this.posts = posts;});
    this.logger.debug('skipping featured video');
    this.videosService.getFeaturedVideo().pipe(takeUntil(this.destroy$)).subscribe({
      next: (res) => {
        this.featuredVideoId = res.id;
        this.featuredVideoTitle = res.title;
      },
      error: (err) => {
        this.logger.warn(err);
        this.featuredVideoId = 'Rk2SBoBwtRU';
      }
    })
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
