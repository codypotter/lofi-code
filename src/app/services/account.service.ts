import { Injectable } from '@angular/core';
import { Auth, GithubAuthProvider, User, authState, createUserWithEmailAndPassword, signInWithEmailAndPassword, signInWithPopup, updateProfile } from '@angular/fire/auth';
import { Observable, catchError, defer, first, from, switchMap } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AccountService {
  currentUser: Observable<User | null>;

  constructor(private auth: Auth) {
    this.currentUser = authState(auth)
  }

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

  logout() {
    return this.auth.signOut();
  }
}
