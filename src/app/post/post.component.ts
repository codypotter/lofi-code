import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { Post, PostsService } from '../services/posts.service';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { marked } from 'marked';
import { environment } from 'src/environments/environment';
import { NgxSpinnerService } from 'ngx-spinner';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { Subject, takeUntil } from 'rxjs';

@Component({
  selector: 'app-post',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent],
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.scss']
})
export class PostComponent implements OnInit, OnDestroy {
  post?: Post;

  relatedPosts: Post[] = [];

  destroy$ = new Subject<void>();

  constructor(
    private postsService: PostsService,
    private route: ActivatedRoute,
    private spinner: NgxSpinnerService,
    private breadcrumbService: BreadcrumbService,
  ) {}

  ngOnInit() {
    this.route.params.pipe(takeUntil(this.destroy$)).subscribe((params) => {
      const slug = params['slug'];
      this.loadPost(slug);
      this.setBreadcrumbs(slug);
    });
  }

  private getRelatedPosts() {
    const [firstTag] = this.post?.tags ?? [];
    if (firstTag) {
      this.postsService.getN(4, undefined, firstTag).subscribe((relatedPosts) => {
        this.relatedPosts = relatedPosts.filter(p => p.slug !== this.post?.slug).slice(0, 3);
      });
    }
  }

  loadPost(slug: string) {
    this.spinner.show();
    this.postsService.getPostBySlug(slug).pipe(takeUntil(this.destroy$)).subscribe((post) => {
      this.post = post;
      this.getRelatedPosts();
      this.spinner.hide();
    });
  }

  setBreadcrumbs(slug: string) {
    this.breadcrumbService.setBreadcrumbs([
      { text: 'home', routerLink: '/' },
      { text: 'posts', routerLink: '/posts' },
      { text: slug },
    ]);
  }

  getContent() {
    return marked(this.post?.content[0].value as string ?? '', { async: false }) as string;
  }

  getHeaderImg() {
    return `${environment.storageUrl}${encodeURIComponent(this.post?.header_image ?? '')}?alt=media`;
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
