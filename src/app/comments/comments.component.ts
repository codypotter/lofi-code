import { Component, Input, OnInit } from '@angular/core';
import { Comment, Post, PostsService } from '../services/posts.service';
import { NGXLogger } from 'ngx-logger';
import { CommentComponent } from '../comment/comment.component';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-comments',
  standalone: true,
  imports: [CommonModule, CommentComponent],
  templateUrl: './comments.component.html',
  styleUrl: './comments.component.scss'
})
export class CommentsComponent implements OnInit {
  @Input() post!: Post;

  comments: Comment[] = [];

  constructor(private postsService: PostsService, private logger: NGXLogger) { }

  ngOnInit(): void {
    this.postsService.getPostComments(this.post.id!).subscribe((comments) => {
      this.logger.trace('comments', comments);
      this.comments = comments;
    });
  }
}
