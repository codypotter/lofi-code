import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostsService } from '../services/posts.service';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { DocumentSnapshot, Timestamp } from 'firebase/firestore';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { Subscribable, Subscription } from 'rxjs';
import { ActivatedRoute, Route } from '@angular/router';

@Component({
  selector: 'app-posts',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent],
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.scss']
})
export class PostsComponent implements OnInit, OnDestroy {
  limit = 10;

  posts: any[] = [];

  tags: string[] = ['architecture', 'clean-code', 'go', 'programming',];

  private _tags = [];

  private _tagsSubscription!: Subscription;

  lastPublished?: Date;

  showLoadMore = true;

  constructor(private postsService: PostsService, private breadcrumbService: BreadcrumbService, private route: ActivatedRoute) { }
  
  ngOnInit(): void {
    this.breadcrumbService.setBreadcrumbs([
      { text: 'home', routerLink: '/' },
      { text: 'posts' },
    ]);
    this._tagsSubscription = this.route.queryParams.subscribe(params => {
      console.warn('params', params);
      this._tags = params['tags'];
      this.loadMorePosts();
    });
  }

  loadMorePosts() {
    this.postsService.getN(this.limit, this.lastPublished).subscribe((posts) => {
      this.posts = [...this.posts, ...posts];
      if (posts.length < this.limit) {
        this.showLoadMore = false;
      }
      this.lastPublished = posts[posts.length - 1].publish_date.toDate();
    });
  }

  onLoadMore() {
    this.loadMorePosts();
  }

  ngOnDestroy(): void {
    this._tagsSubscription.unsubscribe();
  }
}
