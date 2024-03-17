import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TagsService } from '../services/tags.service';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-tags',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './tags.component.html',
  styleUrls: ['./tags.component.scss']
})
export class TagsComponent implements OnInit {
  @Input() tags: string[] = [];

  @Input() size = 'is-normal';

  constructor(private tagsService: TagsService) { }

  ngOnInit(): void {
    if (this.tags.length === 0) {
      this.tagsService.get().subscribe((tags) => this.tags = ['all', ...tags.map(tag => tag.text)]);
    }
  }
}
