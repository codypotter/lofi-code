import { Injectable } from '@angular/core';
import { Auth, GithubAuthProvider, createUserWithEmailAndPassword, signInWithEmailAndPassword, signInWithPopup, updateProfile } from '@angular/fire/auth';
import { defer, first, from, switchMap } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AccountService {

  constructor(private auth: Auth) { }

  login(email: string, password: string) {
    return defer(() => signInWithEmailAndPassword(this.auth, email, password)).pipe(first());
  }

  loginWithGitHub() {
    return defer(() => signInWithPopup(this.auth, new GithubAuthProvider().addScope('read:user'))).pipe(first());
  }

  createAccount(email: string, password: string, displayName: string) {
    return defer(() => createUserWithEmailAndPassword(this.auth, email, password)).pipe(
      switchMap(userCredential => {
        return from(updateProfile(userCredential.user, { displayName }));
      }),
      first(),
    );
  }

  getCurrentUserInfo() {
    return this.auth.currentUser;
  }
}
