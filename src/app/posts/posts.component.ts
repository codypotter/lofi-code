import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Post, PostsService } from '../services/posts.service';
import { BreadcrumbService } from '../services/breadcrumb.service';
import {  Subject, take, takeUntil } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { SearchResultComponent } from '../search-result/search-result.component';
import { TagsComponent } from '../tags/tags.component';
import { TagsService } from '../services/tags.service';
import { NgxSpinnerService } from 'ngx-spinner';

@Component({
  selector: 'app-posts',
  standalone: true,
  imports: [CommonModule, SearchResultComponent, TagsComponent],
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.scss']
})
export class PostsComponent implements OnInit, OnDestroy {
  limit = 10;

  posts: Post[] = [];

  tags: string[] = [];

  private _tag?: string;

  lastPublished?: Date;

  showLoadMore = true;

  destroy$ = new Subject<void>();

  constructor(
    private postsService: PostsService,
    private tagsService: TagsService,
    private breadcrumbService: BreadcrumbService,
    private route: ActivatedRoute,
    private spinnerService: NgxSpinnerService,
  ) { }
  
  ngOnInit(): void {
    this.tagsService.get().pipe(takeUntil(this.destroy$)).subscribe((tags) => {
      this.tags = tags.map(tag => tag.text);
    });
    this.route.queryParams.pipe(takeUntil(this.destroy$)).subscribe(params => {
      const tag = params['tag'];
      this._tag = tag === 'all' ? undefined : tag;
      this.setBreadcrumbs();
      this.newSearch();
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
    this.spinnerService.show();
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
    this.spinnerService.hide();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}