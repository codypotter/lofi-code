import { inject } from '@angular/core';
import { ResolveFn } from '@angular/router';
import { Post, PostsService } from './posts.service';

export const postsResolver: ResolveFn<Post[]> = (route, state) => {
  const tag = route.queryParams['tag'];
  return inject(PostsService).getN(10, undefined, tag === 'all' ? undefined : tag);
};
