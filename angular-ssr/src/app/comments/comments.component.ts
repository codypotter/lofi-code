import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Comment, Post, PostsService } from '../services/posts.service';
import { NGXLogger } from 'ngx-logger';
import { CommentComponent } from '../comment/comment.component';
import { CommonModule } from '@angular/common';
import { User } from '@angular/fire/auth';
import { AccountService } from '../services/account.service';
import { FormsModule } from '@angular/forms';
import { Timestamp } from '@angular/fire/firestore/lite';
import { Subject, takeUntil } from 'rxjs';

@Component({
    selector: 'app-comments',
    imports: [CommonModule, CommentComponent, FormsModule],
    templateUrl: './comments.component.html',
    styleUrl: './comments.component.scss'
})
export class CommentsComponent implements OnInit, OnDestroy {
  @Input() post!: Post;

  comments: Comment[] = [];

  currentUser?: User;

  commentText = '';

  destroy$ = new Subject<void>();

  constructor(
    private postsService: PostsService,
    private accountService: AccountService,
    private logger: NGXLogger,
  ) { }

  ngOnInit(): void {
    this.getComments();

    this.accountService.currentUser.pipe(takeUntil(this.destroy$)).subscribe((user) => {
      this.currentUser = user ?? undefined;
    });
  }

  private getComments() {
    this.postsService.getPostComments(this.post.id!).subscribe((comments) => {
      this.logger.trace('comments', comments);
      this.comments = comments;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onSubmit() {
    if (!this.currentUser) {
      this.logger.error('No current user');
      return;
    }
    this.accountService.getCurrentUserRef().subscribe((user) => {
      this.logger.debug('user', user);
      this.logger.debug('post', this.post);
      this.logger.debug('currentUser', this.currentUser)
      this.postsService.addPostComment(this.post.id!, {
        text: this.commentText,
        timestamp: Timestamp.fromDate(new Date()),
        user: user.ref,
      }).subscribe(() => {
        this.logger.debug('Comment added');
        this.commentText = '';
        this.getComments();
      });
    })
  }

  captureUnauthorized() {
    if (!this.currentUser) {
      this.accountService.setShowCreateAccount(true);
    }
  }
}
