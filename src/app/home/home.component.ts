import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostsService } from '../services/posts.service';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { BreadcrumbService } from '../services/breadcrumb.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  posts: any[] = [];

  constructor(private postsService: PostsService, private breadcrumbService: BreadcrumbService) { }
  
  ngOnInit(): void {
    this.breadcrumbService.setBreadcrumbs([]);
    this.postsService.getTopTen().subscribe((posts) => {
      this.posts = posts.slice(0, 3);
    });
  }
}
