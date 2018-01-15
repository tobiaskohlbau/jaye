import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { Dashboard } from './dashboard';

describe('DashboardComponent', () => {
  let component: Dashboard;
  let fixture: ComponentFixture<Dashboard>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [Dashboard]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(Dashboard);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
