import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostsService } from '../services/posts.service';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { DocumentSnapshot, Timestamp } from 'firebase/firestore';
import { PostPreviewComponent } from '../post-preview/post-preview.component';

@Component({
  selector: 'app-posts',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent],
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.scss']
})
export class PostsComponent implements OnInit {
  posts: any[] = [];

  lastPublished?: Date;

  constructor(private postsService: PostsService, private breadcrumbService: BreadcrumbService) { }
  
  ngOnInit(): void {
    this.breadcrumbService.setBreadcrumbs([
      { text: 'Home', routerLink: '/' },
      { text: 'Posts' },
    ]);
    this.loadMorePosts();
  }

  loadMorePosts() {
    this.postsService.getTopTen(this.lastPublished).subscribe((posts) => {
      this.posts = [...this.posts, ...posts];
      console.warn('posts', posts);
      if (posts.length === 0) {
        return;
      }
      this.lastPublished = posts[posts.length - 1].publish_date.toDate();
    });
  }

  onLoadMore() {
    this.loadMorePosts();
  }
}
