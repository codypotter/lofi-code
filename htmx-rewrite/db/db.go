package db

import (
	"bytes"
	"context"
	"loficode/config"
	"loficode/templates/components"
	"time"

	firebase "firebase.google.com/go"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// getPostBySlug(slug: string): Observable<Post> {
//     const postCollection = query(
//       collection(this.firestore, 'blog'),
//       where('slug', '==', slug),
//       limit(1),
//     );
//     return collectionData(postCollection, { idField: 'id' }).pipe(map(posts => posts[0])).pipe(first()) as Observable<Post>;
//   }

func New(cfg *config.Config) (*firebase.App, error) {
	app, err := firebase.NewApp(context.Background(), &firebase.Config{})
	if err != nil {
		return nil, err
	}

	return app, nil
}

// sample post
// {
//   "document": {
//     "name": "projects/lofi-code/databases/(default)/documents/blog/rYCtl6x1M4C5IAepn9c4",
//     "fields": {
//       "header_image": {
//         "stringValue": "images/klmah_elmo-hell-3x1.png"
//       },
//       "upvotes": {
//         "nullValue": null
//       },
//       "open_graph_image": {
//         "stringValue": "images/2wbwg_elmo-hell-16x9.png"
//       },
//       "created_on": {
//         "timestampValue": "2024-04-01T16:49:05.854Z"
//       },
//       "name": {
//         "stringValue": "The Hellscape of Multi-Language Angular"
//       },
//       "description": {
//         "stringValue": "Let's talk about the multi-lingual hellscape that is localizing your Angular app. Simply put, there are no good options and no matter what you will need to compromise to provide multi-lingual features in your Angular app."
//       },
//       "reviewed": {
//         "booleanValue": true
//       },
//       "publish_date": {
//         "timestampValue": "2024-03-31T04:00:00Z"
//       },
//       "slug": {
//         "stringValue": "angular-language-hell"
//       },
//       "tags": {
//         "arrayValue": {
//           "values": [
//             {
//               "stringValue": "angular"
//             },
//             {
//               "stringValue": "typescript"
//             }
//           ]
//         }
//       },
//       "status": {
//         "stringValue": "published"
//       },
//       "content": {
//         "arrayValue": {
//           "values": [
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "text"
//                   },
//                   "value": {
//                     "stringValue": "Ok, so you want to add multi-lingual features to your angular app -- an admirable task. I have recently gone down this path, and I'd like to share my experience with you.\n\nSimply put, there are no good options. No matter what you do, you will need to make a compromise in order to provide multi-lingual features in your Angular app.\n\nIf you want to check all of these boxes, you're going to be disappointed:\n\n- [ ] Cheap, static web hosting\n- [ ] Simple, fast build system\n- [ ] Easy to read/organize/update translations\n- [ ] Dynamic (runtime) translations\n- [ ] Clean, focused, performant application code\n\nFirst of all, none of these requirements are crazy or outlandish. But here's the harsh reality -- achieving all these goals simultaneously is not easily attainable by any means I have found."
//                   }
//                 }
//               }
//             },
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "images"
//                   },
//                   "value": {
//                     "arrayValue": {
//                       "values": [
//                         {
//                           "stringValue": "images/9mdcw_elmo-hell.gif"
//                         }
//                       ]
//                     }
//                   }
//                 }
//               }
//             },
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "text"
//                   },
//                   "value": {
//                     "stringValue": "Aside: If you have a super slick way to do this and it checks all the boxes, please let me know what it is! You can share your solutions and ideas in the comments or [contact me](https://codypotter.com).\n\n## Why? A decision tree"
//                   }
//                 }
//               }
//             },
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "images"
//                   },
//                   "value": {
//                     "arrayValue": {
//                       "values": [
//                         {
//                           "stringValue": "images/2ds5m_decision-tree.png"
//                         }
//                       ]
//                     }
//                   }
//                 }
//               }
//             },
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "text"
//                   },
//                   "value": {
//                     "stringValue": "Ok, lets talk about the current options. At the bottom of the flowchart are the two big players in the localization game, Angular Internationalization (i18n) and ngx-translate, are the only options I'm discussing today. \n\n## Your options\n\nLet's discuss both, including how they work and what they look like in your project.\n\n### Angular Internationalization (i18n)\n\nFirst of all, I could (and should) have a whole blog post about abbreviations that use this pattern. i18n is unclear and confusing because the abbreviation itself is complex, un-pronouncable, and easy to mis-type for newbies. Can you spot the difference -\u003e i18n vs il8n.\n\nAngular i18n, is the Angular-prescribed method for adding multi-language features to your application. \n\n#### How it appears in your code\n\n__Markup Translation__: Angular i18n primarily involves marking up text in your templates and components for translation using special Angular syntax. For example, you can use the i18n attribute on HTML elements to mark them for translation.\n\n```html\n\u003cp i18n\u003eHello, world!\u003c/p\u003e\n```\n\n__Translation Files__: Angular generates translation files (.xlf or .xlf2 format) during the build process, which contain the translatable strings marked in your code. These files can be provided to translators for localization into different languages.\n\n__Translation Placeholder__: Angular i18n uses placeholder elements (`\u003cng-container\u003e` and `\u003cng-template\u003e`) to represent translatable content in the generated translation files. Translators can then replace these placeholders with the corresponding translations.\n\n#### How it is Hosted/Served\n\n__Build Process__: During the Angular build process (ng build), Angular CLI compiles your application and generates separate bundles for each supported language. It also creates translation files for each language, which contain the translated text.\n\n__Serving Translations__: When serving your Angular application, you can __configure your server__ to serve the appropriate language bundle and translation files based on the user's locale. This can be done using server-side routing or by embedding the language in the URL.\n\n__Dynamic Loading__: Alternatively, you can dynamically load translations at runtime using Angular's TranslateService. This allows you to fetch translations asynchronously from a server or storage service, providing flexibility in managing translations. Note: this service isn't like google translate, your content must already have translations provided somewhere, or fetchable from a service.\n\n#### What's the problem?\n\n- It takes a lot of work to set everything up, and working with the translation files can be a pain without the software you need.\n- Say you need to add ONE translation, a very common use case in iterative development, the cli tools make you regenerate your entire translation file, and then its up to you to manually add the xlf translations you need. Its a big headache, and it is confusing for newbies.\n- Its a pain to iterate locally. You need to stop your dev server, extract translations, merge everything together.\n- You can't translate libraries unless they add support for i18n.\n\n### ngx-translate (the third-party option)\n\nngx-translate is a popular third-party library for internationalization (i18n) in Angular applications. It provides methods to manage translations and localize content dynamically at runtime. \n\n🚨🚨🚨\n\nBig asterisk! Long-term support for ngx-translate is questionable. It posted a maintance mode notice, and [this GitHub issue](https://github.com/ngx-translate/core/issues/783) regarding the future of the library is a little bleak. As of Feb 26, 2024, it is still open and there are discussions regarding a lack of worthy maintainers. As of this writing, there is no support for the most recent Angular version, 17, as there should be.\n\n🚨🚨🚨\n\n### How it appears in your code\n\n__Setup__: To use ngx-translate in your Angular application, you first need to install the library via npm and import the necessary modules into your Angular application module.\n\n__Translation Files__: ngx-translate typically uses JSON files to store translations for different languages. These translation files contain key-value pairs, where keys represent the original text in the application, and values represent the translated text.\n\n__Translation Service__: ngx-translate provides a TranslateService that allows you to load translation files, switch between languages, and translate text dynamically in your components and services.\n\n__Markup Translation__: In your Angular templates, you can use the translate pipe provided by ngx-translate to translate text dynamically based on the current language.\n\n```html\n\u003cp\u003e{{ 'HELLO_WORLD' | translate }}\u003c/p\u003e\n```\n\n### How it is Hosted/Served\n\nHosting and serving an Angular application that uses ngx-translate is no different from hosting any other Angular application. You can serve it statically or deploy it to a server, CDN, or cloud platform like AWS, Firebase, or Azure.\n\n__Static Hosting__: You can build your Angular application with ngx-translate (ng build) and deploy the generated static files to any static hosting provider, such as GitHub Pages, Netlify, Vercel, or AWS S3. The translation files (JSON) are included in the build output and served alongside other static assets.\n\n__Server-side Rendering (SSR)__: If you're using Angular Universal for server-side rendering (SSR), ngx-translate can still be used in a similar manner. Translation files are preloaded on the server, and translated content is rendered dynamically before being sent to the client.\n\n__Dynamic Loading__: ngx-translate also supports dynamic loading of translation files at runtime, allowing you to fetch translations asynchronously from a server or storage service. This can be useful for applications with a large number of translations or when translations are updated frequently. Again, these translations have to live somewhere -- its not Google translate.\n\n#### What's the problem?\n\nYou might be thinking, \"What's the problem Cody?\". This looks fine. Here's an example where ngx-translate gets a little out of hand.\n\nLet's look at this template:\n\n```html\n\u003cp\u003e{{ 'HELLO_WORLD' | translate }}\u003c/p\u003e\n```\n\nPretty simple right? In ngx-translate, the syntax for template interpolation is the same as with Angular i18n. However, there are additional complexities introduced in managing the translation service and handling translation logic in the application code.\n\n__Managing the TranslateService__: In _EVERY_ Angular component that uses translations, you need to inject and use the TranslateService to load translations and handle translation logic. This involves additional setup and management compared to Angular i18n, where translation handling is built-in.\n\n__Translation File Loading__: With ngx-translate, you need to load translation files dynamically at runtime using the TranslateService. This involves additional asynchronous logic to fetch translation files from a server or storage service.\n\n__Handling Translation Logic__: In your application code, you may need to handle translation logic manually, such as switching between languages, updating translations dynamically, and handling fallbacks for missing translations. This requires additional code compared to Angular i18n, where translation handling is handled by the framework.\n\n__Testing Considerations__: When testing components that use the TranslateService, additional setup is required to provide mocks or stubs for the service. This can add complexity to unit tests and may require mocking HTTP requests or translation methods.\n\n__Bindings__: Since ngx-translate uses bindings (either a pipe or directive), the translations are found at runtime, which results in a flassh of content effect (FOC) while the translations are loaded. There is a performance overhead as well.\n\n__Splitting translations__: This is buggy and just straight up doesn't work right.\n\n__JSON Format__: If you want to offload your translation work to a third party, you will run into a problem where third parties do not use the same format as ngx-translate. JSON is non-standard.\n\n__No text shows up__: If it breaks, you get no text.\n\n__AOT__: No ahead-of-time compilation optimizations are available. "
//                   }
//                 }
//               }
//             },
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "images"
//                   },
//                   "value": {
//                     "arrayValue": {
//                       "values": [
//                         {
//                           "stringValue": "images/z2r98_scared-dog.gif"
//                         }
//                       ]
//                     }
//                   }
//                 }
//               }
//             },
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "text"
//                   },
//                   "value": {
//                     "stringValue": "## Ok so now what?\n\nGenerally speaking, ngx-translate offers more flexibility and additional features compared to Angular i18n. However, it also introduces additional complexity in managing the translation service and handling translation logic in your application code. This is a trade-off you'll have to weigh for your project.\n\nI don't recommend using Angular i18n in its current state. It is an immediate no-go if you're using a hosting service like AWS S3, Firebase on a free/spark plan, Vercel, GitHub Pages, or Netlify. The build system makes bespoke builds, per translation, that can make your deployment process much more complicated.\n\n### What's on the horizon?\n\nNothing really exciting.\n\nAngular is continually evolving. In fact there have been some recent enhancements, providing much better support for runtime i18n. However, it is still [undocumented](https://github.com/angular/angular/issues/46851) (laugh out loud).\n\nThe Angular team has proven that i18n is clearly not a priority for the framework. The maintainer of ngx-translate has actually been working on Angular i18n, improving it and hoping folks will stop using ngx-translate. His alternative suggestion is to put ngx-translate behind a paywall. To be fair, if ngx-translate could be improved to solve the litany of problems it has, it woud be extremely valuable.\n\nIf you're interested in helping out `ngx-translate`, they are open-source and they are in need qualified maintainers. The future of the project is unclear.\n\n### This part sucks to hear\n\nIf you're hosting your site statically, you're SOL, out of luck. `ngx-translate` has questionable support, and the official angular i18n tools do not document any support. Your best bet is to try to explor a path toward `@angular/ssr` or an `nginx` or `apache` web server solution."
//                   }
//                 }
//               }
//             },
//             {
//               "mapValue": {
//                 "fields": {
//                   "type": {
//                     "stringValue": "images"
//                   },
//                   "value": {
//                     "arrayValue": {
//                       "values": [
//                         {
//                           "stringValue": "images/o4xcn_that-sucks.gif"
//                         }
//                       ]
//                     }
//                   }
//                 }
//               }
//             }
//           ]
//         }
//       }
//     },
//     "createTime": "2024-04-01T16:49:05.880755Z",
//     "updateTime": "2024-04-02T11:50:11.293410Z"
//   },
//   "readTime": "2025-03-25T01:23:09.316456Z"
// }

func newParser() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(highlighting.Highlighting),
	)
}

func GetPostBySlug(app *firebase.App, slug string) (*components.Post, error) {
	cl, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}
	query := cl.Collection("blog").Where("slug", "==", slug).Limit(1)
	postIt := query.Documents(context.Background())
	defer postIt.Stop()
	post, err := postIt.Next()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	// err = newParser().Convert(post.Data()["content"], &buf)
	if err != nil {
		return nil, err
	}
	// i need to use something like github.com/yuin/goldmark to parse the content
	// and convert it to html
	return &components.Post{
		Id:          &post.Ref.ID,
		Name:        post.Data()["name"].(string),
		Slug:        post.Data()["slug"].(string),
		Description: post.Data()["description"].(string),
		// HeaderImage:    &post.Data()["header_image"].(string),
		// OpenGraphImage: &post.Data()["open_graph_image"].(string),
		// CreatedOn:      post.Data()["created_on"].(string),
		PublishDate: post.Data()["publish_date"].(time.Time),
		Content:     buf.String(),
		Comments:    nil,
		Tags:        nil,
		Upvotes:     0,
	}, nil
}
