import { Component, OnDestroy, OnInit } from '@angular/core';

import { Post, PostsService } from '../services/posts.service';
import { BreadcrumbService } from '../services/breadcrumb.service';
import {  Subject, take, takeUntil } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { SearchResultComponent } from '../search-result/search-result.component';
import { TagsComponent } from '../tags/tags.component';
import { Title } from '@angular/platform-browser';

@Component({
    selector: 'app-posts',
    imports: [SearchResultComponent, TagsComponent],
    templateUrl: './posts.component.html',
    styleUrls: ['./posts.component.scss']
})
export class PostsComponent implements OnInit, OnDestroy {
  limit = 10;

  posts: Post[] = [];

  private _tag?: string;

  lastPublished?: Date;

  showLoadMore = true;

  destroy$ = new Subject<void>();

  constructor(
    private postsService: PostsService,
    private breadcrumbService: BreadcrumbService,
    private route: ActivatedRoute,
    private title: Title,
  ) { }
  
  ngOnInit(): void {
    this.title.setTitle('lofi code - posts');
    this.posts = [];
    this._tag = this.route.snapshot.queryParams['tag'];
    this.setBreadcrumbs();
    this.newSearch();

    this.route.queryParams.pipe(takeUntil(this.destroy$)).subscribe(params => {
      const newTag = params['tag'];
      const normalizedTag = newTag === 'all' ? undefined : newTag;
      if (normalizedTag !== this._tag) {
        this._tag = normalizedTag;
        this.newSearch();
        this.setBreadcrumbs();
      }
    });
  }

  private setBreadcrumbs() {
    this.breadcrumbService.setBreadcrumbs([
      { text: 'home', routerLink: '/' },
      { text: 'posts', routerLink: 'posts' },
      ...this._tag ? [{ text: this._tag }] : [],
    ]);
  }

  loadPosts() {
    this.postsService.getN(this.limit, this.lastPublished, this._tag).pipe(take(1)).subscribe((posts) => {
      this.handlePosts(posts);
    });
  }

  newSearch() {
    this.lastPublished = undefined;
    this.posts = [];
    this.loadPosts();
  }

  private handlePosts(posts: Post[]) {
    this.posts = [...this.posts, ...posts];
    if (posts.length > 0) {
      this.lastPublished = posts[posts.length - 1]?.publish_date.toDate();
    }
    this.showLoadMore = posts.length === this.limit;
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}