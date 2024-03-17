import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PostsService } from '../services/posts.service';
import { PostPreviewComponent } from '../post-preview/post-preview.component';
import { BreadcrumbService } from '../services/breadcrumb.service';
import { RouterModule } from '@angular/router';
import { TagsComponent } from '../tags/tags.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, PostPreviewComponent, RouterModule, TagsComponent],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  posts: any[] = [];

  constructor(private postsService: PostsService, private breadcrumbService: BreadcrumbService) { }
  
  ngOnInit(): void {
    this.breadcrumbService.setBreadcrumbs([]);
    this.postsService.getN(10).subscribe((posts) => {
      this.posts = posts.slice(0, 3);
    });
  }
}
