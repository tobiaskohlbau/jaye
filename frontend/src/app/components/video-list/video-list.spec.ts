import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VideoList } from './video-list';

describe('VideoList', () => {
  let component: VideoList;
  let fixture: ComponentFixture<VideoList>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [VideoList]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VideoList);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
