import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Post } from '../services/posts.service';
import { TagsComponent } from '../tags/tags.component';
import { marked } from 'marked';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-search-result',
  standalone: true,
  imports: [CommonModule, TagsComponent, RouterLink],
  templateUrl: './search-result.component.html',
  styleUrls: ['./search-result.component.scss']
})
export class SearchResultComponent {
  @Input() post!: Post;

  getContent() {
    return marked(this.post.content[0].value as string, { async: false }) as string;
  }
}
