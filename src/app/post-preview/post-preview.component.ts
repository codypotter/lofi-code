import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Post } from '../services/posts.service';
import { marked } from 'marked';
import { environment } from 'src/environments/environment';
import { RouterModule } from '@angular/router';
import { TagsComponent } from '../tags/tags.component';

@Component({
  selector: 'app-post-preview',
  standalone: true,
  imports: [CommonModule, RouterModule, TagsComponent],
  templateUrl: './post-preview.component.html',
  styleUrls: ['./post-preview.component.scss']
})
export class PostPreviewComponent implements OnInit {
  @Input() post!: Post;

  content = '';

  headerImg = '';

  ngOnInit(): void {
    this.content = marked(this.post.content[0].value as string, { async: false }) as string;
    this.headerImg = `${environment.storageUrl}${encodeURIComponent(this.post.header_image ?? '')}?alt=media`;
  }
}
