import { Injectable } from '@angular/core';
import { Firestore, addDoc, arrayUnion, collection, collectionData, orderBy, query, setDoc, updateDoc, where } from '@angular/fire/firestore/lite';
import { DocumentReference, Timestamp, doc, limit, startAfter } from 'firebase/firestore/lite';
import { Observable, from, map, of, switchMap } from 'rxjs';

export interface Comment {
  text: string;
  timestamp: Timestamp;
  user: DocumentReference;
}

export interface Post {
  id?: string;
  name: string;
  slug: string;
  description: string;
  tags: string[];
  publish_date: Timestamp;
  created_on: Timestamp;
  status: string;
  reviewed: boolean;
  header_image: string | null;
  open_graph_image: string | null;
  content: Array<{ type: string; value: string | string[] }>;
  comments: Array<Comment>;
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
    return collectionData(postCollection, { idField: 'id' }).pipe(map(posts => posts[0])) as Observable<Post>;
  }

  getPostComments(id: string): Observable<Comment[]> {
    return collectionData(collection(this.firestore, `blog/${id}/comments`)) as Observable<Comment[]>;
  }

  addPostComment(id: string, comment: Comment) {
    return from(addDoc(collection(this.firestore, `blog/${id}/comments`), comment));
  }
}
