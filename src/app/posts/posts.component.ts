import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostsService } from '../services/posts.service';
import { BreadcrumbService } from '../services/breadcrumb.service';
import {  Subscription } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { SearchResultComponent } from '../search-result/search-result.component';
import { TagsComponent } from '../tags/tags.component';

@Component({
  selector: 'app-posts',
  standalone: true,
  imports: [CommonModule, SearchResultComponent, TagsComponent],
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.scss']
})
export class PostsComponent implements OnInit, OnDestroy {
  limit = 10;

  posts: any[] = [];

  tags: string[] = ['architecture', 'clean-code', 'go', 'programming'];

  private _tag?: string;

  private _tagSubscription!: Subscription;

  lastPublished?: Date;

  showLoadMore = true;

  constructor(private postsService: PostsService, private breadcrumbService: BreadcrumbService, private route: ActivatedRoute) { }
  
  ngOnInit(): void {
    this._tagSubscription = this.route.queryParams.subscribe(params => {
      console.warn('params', params);
      this._tag = params['tag'];
      this.breadcrumbService.setBreadcrumbs([
        { text: 'home', routerLink: '/' },
        { text: 'posts', routerLink: 'posts' },
        ...this._tag ? [{ text: this._tag }] : [],
      ]);
      this.loadMorePosts();
    });
  }

  loadMorePosts() {
    console.warn('loadMorePosts', this._tag, this.lastPublished, this.limit, this.posts.length)
    this.postsService.getN(this.limit, this.lastPublished, this._tag).subscribe((posts) => {
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
    this._tagSubscription.unsubscribe();
  }
}
