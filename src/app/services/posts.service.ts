import { Injectable } from '@angular/core';
import { Firestore, collection, collectionData, orderBy, query, where } from '@angular/fire/firestore';
import { Timestamp, limit, startAfter } from 'firebase/firestore';
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
    let postsCollection;
    if (lastPublishDate && tag) {
      postsCollection = query(
        collection(this.firestore, 'blog'),
        orderBy('publish_date', 'desc'),
        where('status', '==', 'published'),
        where('tags', 'array-contains-any', [tag]),
        limit(n),
        startAfter(lastPublishDate),
      );
    } else if (lastPublishDate) {
      postsCollection = query(
        collection(this.firestore, 'blog'),
        orderBy('publish_date', 'desc'),
        where('status', '==', 'published'),
        limit(n),
        startAfter(lastPublishDate),
      );
    } else if (tag) {
      postsCollection = query(
        collection(this.firestore, 'blog'),
        orderBy('publish_date', 'desc'),
        where('status', '==', 'published'),
        where('tags', 'array-contains-any', [tag]),
        limit(n),
      );
    } else {
      postsCollection = query(
        collection(this.firestore, 'blog'),
        orderBy('publish_date', 'desc'),
        where('status', '==', 'published'),
        limit(n),
      );
    }
    return collectionData(postsCollection) as Observable<Post[]>;
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
