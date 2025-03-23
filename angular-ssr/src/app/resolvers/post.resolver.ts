import { inject } from '@angular/core';
import { ResolveFn } from '@angular/router';
import { Post, PostsService } from '../services/posts.service';

export const postResolver: ResolveFn<Post> = (route, state) => {
  const postsService = inject(PostsService);
  const slug = route.paramMap.get('slug');
  return postsService.getPostBySlug(slug!);
};
