import { Injectable } from '@angular/core';
import { Firestore, addDoc, collection, collectionData, orderBy, query, where } from '@angular/fire/firestore/lite';
import { DocumentReference, Timestamp, limit, startAfter } from 'firebase/firestore/lite';
import { Observable, first, from, map } from 'rxjs';

import { Marked } from "marked";
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js/lib/core';
import typescript from 'highlight.js/lib/languages/typescript';
import go from 'highlight.js/lib/languages/go';
import plaintext from 'highlight.js/lib/languages/plaintext';
import shell from 'highlight.js/lib/languages/shell';

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
  upvotes: number;
}

@Injectable({
  providedIn: 'root'
})
export class PostsService {
  parser: Marked;

  constructor(private firestore: Firestore) {
    hljs.registerLanguage('typescript', typescript);
    hljs.registerLanguage('go', go);
    hljs.registerLanguage('plaintext', plaintext);
    hljs.registerLanguage('shell', shell);
    this.parser = new Marked(markedHighlight({
      langPrefix: 'hljs language-',
      highlight(code, lang, info) {
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        return hljs.highlight(code, { language }).value;
      }
    }))
  }

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

    return collectionData(query(...postsQuery as [any, ...any[]])).pipe(first()) as Observable<Post[]>;
  }

  getPostBySlug(slug: string): Observable<Post> {
    const postCollection = query(
      collection(this.firestore, 'blog'),
      where('slug', '==', slug),
      limit(1),
    );
    return collectionData(postCollection, { idField: 'id' }).pipe(map(posts => posts[0])).pipe(first()) as Observable<Post>;
  }

  getPostComments(id: string): Observable<Comment[]> {
    return collectionData(collection(this.firestore, `blog/${id}/comments`)).pipe(first()) as Observable<Comment[]>;
  }

  addPostComment(id: string, comment: Comment) {
    return from(addDoc(collection(this.firestore, `blog/${id}/comments`), comment)).pipe(first());
  }
}
