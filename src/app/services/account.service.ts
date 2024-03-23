import { Injectable } from '@angular/core';
import { Auth, GithubAuthProvider, User, authState, createUserWithEmailAndPassword, signInWithEmailAndPassword, signInWithPopup, updateCurrentUser, updateProfile } from '@angular/fire/auth';
import { updateEmail, updatePassword } from 'firebase/auth';
import { NGXLogger } from 'ngx-logger';
import { BehaviorSubject, Observable, defer, first, from, switchMap, tap } from 'rxjs';
import { UsersService } from './users.service';
import { EmailsService } from './emails.service';

@Injectable({
  providedIn: 'root'
})
export class AccountService {
  currentUser: Observable<User | null>;

  private _showLogin = new BehaviorSubject<boolean>(false);
  
  public readonly showLogin = this._showLogin.asObservable();

  private _showCreateAccount = new BehaviorSubject<boolean>(false);
  
  public readonly showCreateAccount = this._showCreateAccount.asObservable();

  constructor(
    private logger: NGXLogger,
    private auth: Auth,
    private usersService: UsersService,
    private emailsService: EmailsService,
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
          displayName: displayName ?? 'anonymous',
          photoURL: photoURL ?? '',
        }).subscribe({
          next: () => {
            this.logger.debug('User created');
            this.getCurrentUserRef().subscribe((user) => {
              this.emailsService.set(user.id, {
                email: email ?? '',
                mailingList: false,
                user: user.ref,
              }).subscribe({
                next: () => this.logger.debug('Email created'),
                error: (err) => this.logger.error('Error creating email:', err),
              });
            });
          },
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
          displayName,
          photoURL: providerData?.[0]?.photoURL ?? '',
        }).subscribe({
          next: () => {
            this.logger.debug('User created')
            this.getCurrentUserRef().subscribe((user) => {
              this.emailsService.set(user.id, {
                email,
                mailingList,
                user: user.ref,
              }).subscribe({
                next: () => this.logger.debug('Email created'),
                error: (err) => this.logger.error('Error creating email:', err),
              });
            });
          },
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

  getCurrentUserRef() {
    this.logger.trace('getting current user ref');
    return this.usersService.get(this.getCurrentUserInfo()?.uid!);
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

  setShowLogin(value: boolean) {
    this._showLogin.next(value);
  }

  setShowCreateAccount(value: boolean) {
    console.debug('setShowCreateAccount', value)
    this._showCreateAccount.next(value);
  }
}
