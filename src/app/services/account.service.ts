import { Injectable } from '@angular/core';
import { Auth, GithubAuthProvider, User, authState, createUserWithEmailAndPassword, signInWithEmailAndPassword, signInWithPopup, updateCurrentUser, updateProfile } from '@angular/fire/auth';
import { updateEmail, updatePassword } from 'firebase/auth';
import { NGXLogger } from 'ngx-logger';
import { Observable, catchError, defer, first, from, switchMap, tap } from 'rxjs';
import { UsersService } from './users.service';

@Injectable({
  providedIn: 'root'
})
export class AccountService {
  currentUser: Observable<User | null>;

  constructor(
    private logger: NGXLogger,
    private auth: Auth,
    private usersService: UsersService,
  ) {
    this.currentUser = authState(auth)
  }

  login(email: string, password: string) {
    this.logger.trace('logging in', email);
    return defer(() => signInWithEmailAndPassword(this.auth, email, password)).pipe(first());
  }

  loginWithGitHub() {
    this.logger.trace('logging in with GitHub');
    return defer(() => signInWithPopup(this.auth, new GithubAuthProvider().addScope('read:user'))).pipe(
      first(),
      tap(() => {
        const { uid, providerData } = this.getCurrentUserInfo() ?? {}
        const { email, displayName, photoURL } = providerData?.[0] ?? {};
        if (!uid) {
          this.logger.error('No user ID somehow?');
          return;
        }
        this.usersService.set(uid, {
          email: email ?? '',
          displayName: displayName ?? 'anonymous',
          mailingList: false,
          photoURL: photoURL ?? '',
        }).subscribe({
          next: () => this.logger.debug('User created'),
          error: (err) => this.logger.error('Error creating user:', err),
        });
      }),
    );
  }

  createAccount(email: string, password: string, displayName: string, mailingList: boolean) {
    this.logger.trace('creating account', email, displayName);
    return defer(() => createUserWithEmailAndPassword(this.auth, email, password)).pipe(
      switchMap(userCredential => {
        const {uid, providerData} = this.getCurrentUserInfo() ?? {};
        if (!uid) {
          this.logger.error('No user ID somehow?');
          return from(updateProfile(userCredential.user, { displayName }));;
        }
        this.usersService.set(uid, {
          email,
          displayName,
          mailingList,
          photoURL: providerData?.[0]?.photoURL ?? '',
        }).subscribe({
          next: () => this.logger.debug('User created'),
          error: (err) => this.logger.error('Error creating user:', err),
        });
        return from(updateProfile(userCredential.user, { displayName }));
      }),
      first(),
    );
  }

  getCurrentUserInfo() {
    this.logger.trace('getting current user info');
    return this.auth.currentUser;
  }

  updateUsername(displayName: string) {
    this.logger.trace('updating username', displayName);
    return from(updateProfile(this.auth.currentUser!, { displayName })).pipe(first());
  }

  logout() {
    this.logger.trace('logging out');
    return from(this.auth.signOut()).pipe(first());
  }

  updatePassword(password: string) {
    this.logger.trace('updating password');
    return this.currentUser.pipe(
      first(),
      switchMap(user => from(updatePassword(user!, password))),
    );
  }

  updateEmail(password: string) {
    this.logger.trace('updating password');
    return this.currentUser.pipe(
      first(),
      switchMap(user => from(updateEmail(user!, password))),
    );
  }
}
