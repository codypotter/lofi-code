import { NgxLoggerLevel } from "ngx-logger";

export const environment = {
    storageUrl: 'https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/',
    youtubeApiKey: '',
    baseUrl: '',
    logLevel: NgxLoggerLevel.DEBUG,
    firebaseOptions: {
        apiKey: "",
        authDomain: "",
        projectId: "",
        storageBucket: "",
        messagingSenderId: "",
        appId: "",
        measurementId: ""
    },
}