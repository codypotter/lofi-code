import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, RouterOutlet } from '@angular/router';
import { NgxSpinnerModule } from 'ngx-spinner';
import { Breadcrumb, BreadcrumbService } from './services/breadcrumb.service';
import { CreateAccountComponent } from './create-account/create-account.component';
import { AccountService } from './services/account.service';
import { User } from 'firebase/auth';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, NgxSpinnerModule, RouterModule, CreateAccountComponent],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'lofi code';

  breadcrumbs: Breadcrumb[] = [];

  user: User | null = null;

  showCreateAccount = false;

  constructor(
    private breadcrumbService: BreadcrumbService,
    private cdr: ChangeDetectorRef,
    private accountService: AccountService,
  ) { }

  ngOnInit() {
    this.user = this.accountService.getCurrentUserInfo();
    console.warn('user', this.user);
    this.breadcrumbService.breadcrumbs$.subscribe((breadcrumbs) => {
      this.breadcrumbs = breadcrumbs;
      this.cdr.detectChanges();
    });
  }
}
