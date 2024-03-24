import { Component, Inject, OnDestroy, OnInit, PLATFORM_ID } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { Post, PostsService } from '../services/posts.service';
import { ActivatedRoute, Router } from '@angular/router';
import { marked } from 'marked';
import { environment } from 'src/environments/environment';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { Subject, takeUntil } from 'rxjs';
import { NgxSkeletonLoaderModule } from 'ngx-skeleton-loader';
import { CommentsComponent } from '../comments/comments.component';
import { NGXLogger } from 'ngx-logger';
import { ShareButtonsModule } from 'ngx-sharebuttons/buttons';
import { ShareIconsModule } from 'ngx-sharebuttons/icons';
import { TagsComponent } from '../tags/tags.component';
import { Meta, Title } from '@angular/platform-browser';
import { Location } from '@angular/common';

@Component({
  selector: 'app-post',
  standalone: true,
  imports: [
    CommonModule,
    PostPreviewComponent,
    NgxSkeletonLoaderModule,
    CommentsComponent,
    ShareButtonsModule,
    ShareIconsModule,
    TagsComponent,
  ],
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
    private breadcrumbService: BreadcrumbService,
    private logger: NGXLogger,
    private meta: Meta,
    private location: Location,
    private title: Title,
    private router: Router,
  ) {}

  ngOnInit() {
    this.post = this.route.snapshot.data['post'];
    this.setMetaTags();
    this.getRelatedPosts();
    this.setBreadcrumbs(this.post!.slug);
  }

  private getRelatedPosts() {
    const [firstTag] = this.post?.tags ?? [];
    if (firstTag) {
      this.postsService.getN(4, undefined, firstTag).subscribe((relatedPosts) => {
        this.relatedPosts = relatedPosts.filter(p => p.slug !== this.post?.slug).slice(0, 3);
      });
    }
  }

  setMetaTags() {
    const title = `lofi code - ${this.post?.name}`;
    const description = this.post?.description ?? '';

    this.title.setTitle(title);
    this.meta.updateTag({ name: 'description', content: description });
    this.meta.updateTag({ name: 'og:title', content: title });
    this.meta.updateTag({ name: 'og:description', content: description });
    this.meta.updateTag({ name: 'og:image', content: this.getOGImage() });
    this.meta.updateTag({ name: 'og:url', content: `https://lofi-code.com${this.location.path()}` });
    this.meta.updateTag({ name: 'site_name', content: 'lofi code'});
    this.meta.updateTag({ name: 'twitter:title', content: title });
    this.meta.updateTag({ name: 'twitter:description', content: description });
    this.meta.updateTag({ name: 'twitter:image', content: this.getOGImage() });
    this.meta.updateTag({ name: 'twitter:card', content: 'summary_large_image' });
    this.meta.updateTag({ name: 'type', content: 'article' });
    this.meta.updateTag({ name: 'article:published_time', content: this.post?.publish_date.toDate().toISOString() ?? '' });
    this.meta.updateTag({ name: 'article:modified_time', content: this.post?.created_on.toDate().toISOString() ?? '' });
    this.meta.updateTag({ name: 'article:author', content: 'Cody Potter' });
    this.meta.updateTag({ name: 'article:section', content: 'Technology' });
    this.meta.updateTag({ name: 'article:tag', content: this.post?.tags.join(',') ?? '' });
  }

  setBreadcrumbs(slug: string) {
    this.breadcrumbService.setBreadcrumbs([
      { text: 'home', routerLink: '/' },
      { text: 'posts', routerLink: '/posts' },
      { text: slug },
    ]);
  }

  getContent() {
    const allContents = this.post?.content.map(item => {
      if (item.type === 'images') {
        return `<img class="centered-image" src="${this.imageUrl(item.value as string)}" alt="content image"/>`;
      } else {
        return item.value;
      }
    }).join('\n') ?? '';
    return marked(allContents, { async: false }) as string;
  }

  getHeaderImg() {
    return this.imageUrl(this.post?.header_image ?? '');
  }

  getOGImage() {
    return this.imageUrl(this.post?.open_graph_image ?? '');
  }

  imageUrl(url: string) {
    return `${environment.storageUrl}${encodeURIComponent(url)}?alt=media`;
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
