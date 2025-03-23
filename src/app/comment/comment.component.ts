import { Component, Input, OnInit } from '@angular/core';
import { Comment } from '../services/posts.service';
import { UsersService } from '../services/users.service';
import { NGXLogger } from 'ngx-logger';
import { TimeagoModule } from 'ngx-timeago';

@Component({
    selector: 'app-comment',
    imports: [TimeagoModule],
    templateUrl: './comment.component.html',
    styleUrl: './comment.component.scss'
})
export class CommentComponent implements OnInit{
  @Input() comment!: Comment;

  author = '';

  photoUrl = '';

  constructor(private usersService: UsersService, private logger: NGXLogger) {}

  ngOnInit(): void {
    this.usersService.get(this.comment.user.id).subscribe(user => {
      this.author = user.get('displayName') ?? 'Anonymous';
      this.photoUrl = user.get('photoURL') ?? 'assets/anonymous.png';
    });
  }
}
