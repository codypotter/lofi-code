import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export interface Breadcrumb {
  text: string;
  routerLink?: string | any[];
}

@Injectable({
  providedIn: 'root'
})
export class BreadcrumbService {
  private _breadcrumbs = new BehaviorSubject<Array<Breadcrumb>>([]);
  
  public readonly breadcrumbs$ = this._breadcrumbs.asObservable();

  setBreadcrumbs(breadcrumbs: Array<Breadcrumb>) {
    this._breadcrumbs.next(breadcrumbs);
  }
}
