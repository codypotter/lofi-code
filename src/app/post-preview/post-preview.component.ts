import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Post } from '../services/posts.service';
import { environment } from 'src/environments/environment';
import { RouterModule } from '@angular/router';
import { TagsComponent } from '../tags/tags.component';
import { NGXLogger } from 'ngx-logger';

@Component({
    selector: 'app-post-preview',
    imports: [CommonModule, RouterModule, TagsComponent],
    templateUrl: './post-preview.component.html',
    styleUrls: ['./post-preview.component.scss'],
    host: { ngSkipHydration: 'true' }
})
export class PostPreviewComponent implements OnInit {
  @Input() post!: Post;

  headerImg = '';

  constructor(private logger: NGXLogger) {}

  ngOnInit(): void {
    this.headerImg = `${environment.storageUrl}${encodeURIComponent(this.post.header_image ?? '')}?alt=media`;
  }
}
