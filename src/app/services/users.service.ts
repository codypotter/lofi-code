import { Injectable } from '@angular/core';
import { Firestore, doc, getDoc, setDoc } from '@angular/fire/firestore/lite';
import { first, from } from 'rxjs';
import { Storage, ref, uploadBytesResumable } from '@angular/fire/storage';
import { getDownloadURL } from 'firebase/storage';

interface User {
  displayName?: string;
  photoURL?: string;
}

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  constructor(private firestore: Firestore, private storage: Storage) { }

  set(uid: string, user: User) {
    return from(setDoc(doc(this.firestore, 'users', uid), user, { merge: true })).pipe(first())
  }

  get(uid: string) {
    return from(getDoc(doc(this.firestore, 'users', uid))).pipe(first());
  }

  uploadAvatar(blob: Blob, uid: string) {
    return from(uploadBytesResumable(ref(this.storage, `avatars/${uid}`), blob).then((snapshot) => {
      return getDownloadURL(snapshot.ref);
    })).pipe(first());
  }
}
