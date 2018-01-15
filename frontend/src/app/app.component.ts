import { Component, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';

import { environment } from '../environments/environment';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class AppComponent {
  showShadow = false;

  constructor(router: Router) { }
}
