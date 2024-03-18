import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostsService } from '../services/posts.service';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { RouterModule } from '@angular/router';
import { TagsComponent } from '../tags/tags.component';
import { YoutubePlayerComponent } from '../youtube-player/youtube-player.component';
import { VideosService } from '../services/videos.service';
import { Subject, takeUntil } from 'rxjs';
import { NgxSpinnerService } from 'ngx-spinner';
import { NgxSkeletonLoaderModule } from 'ngx-skeleton-loader';

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

  constructor(
    private postsService: PostsService,
    private breadcrumbService: BreadcrumbService,
    private videosService: VideosService,
    private spinner: NgxSpinnerService,
  ) { }
  
  ngOnInit(): void {
    this.spinner.show();
    this.breadcrumbService.setBreadcrumbs([]);
    this.postsService.getN(3).pipe(takeUntil(this.destroy$)).subscribe((posts) => { this.posts = posts;});
    this.videosService.getFeaturedVideo().pipe(takeUntil(this.destroy$)).subscribe({
      next: (res) => {
        const featuredVideo = res.items[0];
        this.featuredVideoId = featuredVideo?.id?.videoId;
        this.spinner.hide();
      },
      error: (err) => {
        this.spinner.hide();
        console.warn(err);
        this.featuredVideoId = 'Rk2SBoBwtRU';
      }
    })
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
