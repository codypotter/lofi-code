import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Post } from '../services/posts.service';
import { TagsComponent } from '../tags/tags.component';
import { RouterLink } from '@angular/router';
import { environment } from 'src/environments/environment';

@Component({
    selector: 'app-search-result',
    imports: [CommonModule, TagsComponent, RouterLink],
    templateUrl: './search-result.component.html',
    styleUrls: ['./search-result.component.scss']
})
export class SearchResultComponent {
  @Input() post!: Post;

  buildImageUrl() {
    return `${environment.storageUrl}${encodeURIComponent(this.post.open_graph_image ?? '')}?alt=media`;
  }
}
