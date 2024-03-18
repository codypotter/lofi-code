import { Injectable } from '@angular/core';
import { Firestore, collection, collectionData, orderBy, query, where } from '@angular/fire/firestore/lite';
import { Timestamp, limit, startAfter } from 'firebase/firestore/lite';
import { Observable, map } from 'rxjs';

export interface Post {
  name: string;
  slug: string;
  tags: string[];
  publish_date: Timestamp;
  created_on: Timestamp;
  status: string;
  reviewed: boolean;
  header_image: string | null;
  content: Array<{ type: string; value: string | string[] }>;
}

@Injectable({
  providedIn: 'root'
})
export class PostsService {
  constructor(private firestore: Firestore) { }

  getN(n: number, lastPublishDate?: Date, tag?: string): Observable<Post[]> {
    const postsQuery: any = [
      collection(this.firestore, 'blog'),
      orderBy('publish_date', 'desc'),
      where('status', '==', 'published'),
    ];
    if (tag) {
      postsQuery.push(where('tags', 'array-contains-any', [tag]));
    }
    postsQuery.push(limit(n));
    if (lastPublishDate) {
      postsQuery.push(startAfter(lastPublishDate));
    }

    return collectionData(query(...postsQuery as [any, ...any[]])) as Observable<Post[]>;
  }

  getPostBySlug(slug: string): Observable<Post> {
    const postCollection = query(
      collection(this.firestore, 'blog'),
      where('slug', '==', slug),
      limit(1),
    );
    return collectionData(postCollection).pipe(map(posts => posts[0])) as Observable<Post>;
  }
}
