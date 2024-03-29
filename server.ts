import 'zone.js/node';

import { APP_BASE_HREF } from '@angular/common';
import { CommonEngine } from '@angular/ssr';
import * as express from 'express';
import { existsSync } from 'node:fs';
import { join } from 'node:path';
import bootstrap from './src/main.server';
import 'dotenv/config';
import axios from 'axios';

// The Express app is exported so that it can be used by serverless Functions.
export function app(): express.Express {
  const server = express();
  const distFolder = join(process.cwd(), 'dist/lofi-code/browser');
  const indexHtml = existsSync(join(distFolder, 'index.original.html'))
    ? join(distFolder, 'index.original.html')
    : join(distFolder, 'index.html');

  const commonEngine = new CommonEngine();

  server.set('view engine', 'html');
  server.set('views', distFolder);

  let cache = {
    time: Date.now(),
    data: null as any,
  };

  server.get('/api/featured-video', async (req, res) => {
    const apiKey = process.env['YOUTUBE_API_KEY'];
    const cacheValid = Date.now() - cache.time < 1000 * 60 * 60 * 24; // Cache valid for 24 hours

    if (cacheValid && cache.data) {
      res.json(cache.data);
      return
    }

    try {
      const response = await axios.get('https://www.googleapis.com/youtube/v3/search', {
        params: {
          part: 'snippet',
          channelId: 'UCsPXgrtO5bTfdVdNLLB_Erw',
          maxResults: 1,
          order: 'date',
          key: apiKey,
        },
        headers: {
          referer: 'https://lofi-code.com',
        }
      });
  
      const video = response.data.items[0];
      const videoId = video.id.videoId;
      const videoTitle = video.snippet.title;

      const videoData = { title: videoTitle, id: videoId };

      cache = {
        time: Date.now(),
        data: videoData
      };
  
      res.json(videoData);
    } catch (error: any) {
      console.error(error.message);
      res.status(500).json({ error: 'An error occurred while fetching the featured video' });
    }
  });

  // Example Express Rest API endpoints
  // server.get('/api/**', (req, res) => { });
  // Serve static files from /browser
  server.get('*.*', express.static(distFolder, {
    maxAge: '1y',
    setHeaders: (res, path) => {
      if (path.endsWith('.ttf')) {
        res.removeHeader('Content-Encoding');
      }
    }
  }));

  // All regular routes use the Angular engine
  server.get('*', (req, res, next) => {
    const { protocol, originalUrl, baseUrl, headers } = req;

    commonEngine
      .render({
        bootstrap,
        documentFilePath: indexHtml,
        url: `${protocol}://${headers.host}${originalUrl}`,
        publicPath: distFolder,
        providers: [{ provide: APP_BASE_HREF, useValue: baseUrl }],
      })
      .then((html) => res.send(html))
      .catch((err) => next(err));
  });

  return server;
}

function run(): void {
  const port = process.env['PORT'] || 4000;

  // Start up the Node server
  const server = app();
  server.listen(port, () => {
    console.log(`Node Express server listening on http://localhost:${port}`);
  });
}

// Webpack will replace 'require' with '__webpack_require__'
// '__non_webpack_require__' is a proxy to Node 'require'
// The below code is to ensure that the server is run only when not requiring the bundle.
declare const __non_webpack_require__: NodeRequire;
const mainModule = __non_webpack_require__.main;
const moduleFilename = mainModule && mainModule.filename || '';
if (moduleFilename === __filename || moduleFilename.includes('iisnode')) {
  run();
}

export default bootstrap;
