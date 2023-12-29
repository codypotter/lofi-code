import { Injectable } from '@angular/core';
import { Firestore, collection, collectionData, orderBy, query, where } from '@angular/fire/firestore';
import { limit } from 'firebase/firestore';
import { Observable } from 'rxjs';

export interface Post {
  name: string;
  tags: string[];
  publish_date: { seconds: number; nanoseconds: number };
  created_on: { seconds: number; nanoseconds: number };
  status: string;
  reviewed: boolean;
  header_image: string | null;
  content: Array<{ type: string; value: string | string[] }>;
}

@Injectable({
  providedIn: 'root'
})
export class PostsService {
  constructor(private firestore: Firestore) {}

  getTopTen(): Observable<Post[]> {
    const postsCollection = query(
      collection(this.firestore, 'blog'),
      orderBy('publish_date', 'desc'),
      where('status', '==', 'published'),
      limit(10),
    );
    return collectionData(postsCollection) as Observable<Post[]>;
  }
}
