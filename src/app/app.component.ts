import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, RouterOutlet } from '@angular/router';
import { Breadcrumb, BreadcrumbService } from './services/breadcrumb.service';
import { NavbarComponent } from './navbar/navbar.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterModule, NavbarComponent],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'lofi code';

  breadcrumbs: Breadcrumb[] = [];

  constructor(
    private breadcrumbService: BreadcrumbService,
    private cdr: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.breadcrumbService.breadcrumbs$.subscribe((breadcrumbs) => {
      this.breadcrumbs = breadcrumbs;
      this.cdr.detectChanges();
    });
  }
}
