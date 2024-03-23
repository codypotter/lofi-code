import { Injectable } from '@angular/core';
import { DocumentReference, Firestore, doc, getDoc, setDoc } from '@angular/fire/firestore/lite';
import { first, from } from 'rxjs';

export interface Email {
  email?: string;
  mailingList?: boolean;
  user?: DocumentReference;
}

@Injectable({
  providedIn: 'root'
})
export class EmailsService {

  constructor(private firestore: Firestore) { }

  set(uid: string, email: Email) {
    return from(setDoc(doc(this.firestore, 'emails', uid), email, { merge: true })).pipe(first())
  }

  get(uid: string) {
    return from(getDoc(doc(this.firestore, 'emails', uid))).pipe(first());
  }
}
