import { Injectable } from '@angular/core';
import { NGXLogger } from 'ngx-logger';
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

  constructor(private logger: NGXLogger) { }

  setBreadcrumbs(breadcrumbs: Array<Breadcrumb>) {
    this.logger.trace('setting breadcrumbs', breadcrumbs);
    this._breadcrumbs.next(breadcrumbs);
  }
}
