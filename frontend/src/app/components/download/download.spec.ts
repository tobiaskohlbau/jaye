import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { Download } from './download';

describe('Download', () => {
  let component: Download;
  let fixture: ComponentFixture<Download>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [Download]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(Download);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
