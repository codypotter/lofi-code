import { Injectable } from '@angular/core';
import { Firestore, collection, collectionData, doc, getDoc, setDoc } from '@angular/fire/firestore';
import { first, from } from 'rxjs';

interface User {
  displayName?: string;
  email?: string;
  photoURL?: string;
  mailingList?: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  constructor(private firestore: Firestore) { }

  set(uid: string, user: User) {
    return from(setDoc(doc(this.firestore, 'users', uid), user, { merge: true })).pipe(first())
  }

  get(uid: string) {
    return from(getDoc(doc(this.firestore, 'users', uid))).pipe(first());
  }
}
