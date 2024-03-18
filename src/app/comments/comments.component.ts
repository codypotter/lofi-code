import { Component, Input } from '@angular/core';
import { Post } from '../services/posts.service';

@Component({
  selector: 'app-comments',
  standalone: true,
  imports: [],
  templateUrl: './comments.component.html',
  styleUrl: './comments.component.scss'
})
export class CommentsComponent {
  @Input() post!: Post;
}
