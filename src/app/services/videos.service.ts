import { Injectable } from '@angular/core';
import { Firestore, collection, collectionData, query } from '@angular/fire/firestore';
import { Observable } from 'rxjs';

interface Video {
  videoId: string;
}

@Injectable({
  providedIn: 'root'
})
export class VideosService {

  constructor(private firestore: Firestore) { }

  getVideos(): Observable<Video[]> {
    return collectionData(query(collection(this.firestore, 'videos'))) as Observable<Video[]>;
  }
}
