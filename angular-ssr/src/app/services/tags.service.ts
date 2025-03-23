import { Injectable, Input } from '@angular/core';
import { Firestore, collection, collectionData, query } from '@angular/fire/firestore/lite';
import { first, Observable } from 'rxjs';

export interface Tag {
  text: string;
}

@Injectable({
  providedIn: 'root'
})
export class TagsService {
  constructor(private firestore: Firestore) { }

  get(): Observable<Tag[]> {
    return collectionData(query(collection(this.firestore, 'tags'))).pipe(first()) as Observable<Tag[]>;
  }
}
