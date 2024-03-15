import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { Post, PostsService } from '../services/posts.service';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { marked } from 'marked';
import { environment } from 'src/environments/environment';
import { NgxSpinnerService } from 'ngx-spinner';
import { BreadcrumbService } from '../services/breadcrumb.service';

@Component({
  selector: 'app-post',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent],
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.scss']
})
export class PostComponent implements OnInit {
  post!: Post;

  constructor(
    private postsService: PostsService,
    private route: ActivatedRoute,
    private spinner: NgxSpinnerService,
    private breadcrumbService: BreadcrumbService,
  ) {}

  ngOnInit() {
    this.breadcrumbService.setBreadcrumbs([
      { text: 'Home', routerLink: '/' },
      { text: 'Posts', routerLink: '/posts' },
      { text: this.route.snapshot.params['slug'] },
    ]);
    this.spinner.show();
    this.postsService.getPostBySlug(this.route.snapshot.params['slug']).subscribe((post) => {
      this.post = post;
      this.spinner.hide();
    });
  }

  getContent() {
    return marked(this.post.content[0].value as string, { async: false }) as string;
  }

  getHeaderImg() {
    return `${environment.storageUrl}${encodeURIComponent(this.post.header_image ?? '')}?alt=media`;

  }
}
