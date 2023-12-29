import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Post } from '../services/posts.service';

@Component({
  selector: 'app-post-preview',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './post-preview.component.html',
  styleUrls: ['./post-preview.component.scss']
})
export class PostPreviewComponent {
  @Input() post!: Post;
}
