import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Post, PostsService } from '../services/posts.service';
import { BreadcrumbService } from '../services/breadcrumb.service';
import {  Subscription } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { SearchResultComponent } from '../search-result/search-result.component';
import { TagsComponent } from '../tags/tags.component';
import { TagsService } from '../services/tags.service';

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

  private _tagSubscription!: Subscription;

  lastPublished?: Date;

  showLoadMore = true;

  constructor(
    private postsService: PostsService,
    private tagsService: TagsService,
    private breadcrumbService: BreadcrumbService,
    private route: ActivatedRoute,
  ) { }
  
  ngOnInit(): void {
    this.tagsService.get().subscribe((tags) => {
      this.tags = tags.map(tag => tag.text);
    });
    this._tagSubscription = this.route.queryParams.subscribe(params => {
      this._tag = params['tag'];
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

  loadMorePosts() {
    this.postsService.getN(this.limit, this.lastPublished, this._tag).subscribe((posts) => {
      this.handlePosts(posts);
    });
  }

  newSearch() {
    this.postsService.getN(this.limit, undefined, this._tag).subscribe((posts) => {
      this.posts = [];
      this.handlePosts(posts);
    });
  }

  private handlePosts(posts: Post[]) {
    this.posts = [...this.posts, ...posts];
    if (posts.length > 0) {
      this.lastPublished = posts[posts.length - 1]?.publish_date.toDate();
    }
    this.showLoadMore = posts.length === this.limit;
  }

  ngOnDestroy(): void {
    this._tagSubscription.unsubscribe();
  }
}
