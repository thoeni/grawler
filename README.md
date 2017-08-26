# Grawler - a sitemap generator
Given a website url (and an optional depth level) this program can start from the url provided and extract all the links, static assets and images referenced by the page itself, and continue following all the links that belong to the same domain.

The result of the processing is returned as `json` (to the standard output or to a file), for example, running the following command

`go build && ./grawler -url=https://thoeni.io`

will produce the following:

<details><summary>Json output</summary><p>

```
{
   "https://thoeni.io": {
      "Url": "https://thoeni.io",
      "Links": [
         "https://thoeni.io/tags/java",
         "https://thoeni.io/tags/conferences",
         "https://thoeni.io/",
         "https://thoeni.io/page/resources",
         "https://thoeni.io/tags",
         "https://thoeni.io/tags/engineering",
         "https://thoeni.io/tags/philosophy",
         "https://thoeni.io/post/macos-sierra-java/",
         "https://thoeni.io/tags/meetups",
         "https://thoeni.io/tags/sierra",
         "https://thoeni.io/post/welcome/",
         "https://thoeni.io/post/hacker-news-london-october-2016/",
         "https://thoeni.io/tags/london",
         "https://thoeni.io/tags/macos",
         "https://thoeni.io/post/socratesuk-2016/",
         "https://thoeni.io/tags/software-craftsmanship",
         "https://thoeni.io/index.xml",
         "https://thoeni.io/project/slack-tube-service/",
         "https://thoeni.io/page/about",
         "https://thoeni.io/post/the-art-of-coding/",
         "https://thoeni.io/tags/dev",
         "https://thoeni.io/tags/agile",
         "https://thoeni.io/tags/tech"
      ],
      "Assets": [
         "https://thoeni.io/css/main.css",
         "https://thoeni.io/css/pygment_highlights.css",
         "https://thoeni.io/css/highlight.min.css",
         "https://thoeni.io/",
         "https://thoeni.io/index.xml"
      ],
      "Images": [
         "https://thoeni.io/images/avatar.jpg"
      ]
   }
}
```

</p></details>



The command instructs the program to crawl the website at `https://thoeni.io` up to depth `0` which means parsing just the first page and extracting a set of `Links`, `Assets` and `Images`. A depth of `1` would cause the whole array of `Links` to be browsed as well and parsed in the same way, and so on and so forth.

A depth of `-1` assumes that the navigation will stop when no more links are found. Depending on the website this might take a long time and/or fill the memory.

By default the depth is set to 0, therefore when no flag is passed only one page will be parsed.

The starting `-url` doesn't have to be the home/root of a website but it could be a subpath (e.g. `https://thoeni.io/blog`)

## How to build it/run it

This program has been developed using Go 1.9.0 (1.8.x should work as well)

```
go build
./grawler -url=https://thoeni.io
```
or

```
go install
grawler -url=https://thoeni.io
```

## Flags
It is possible to call the program passing three flags:

- `-url`: is mandatory, and can be a root URL (e.g. https://thoeni.io) or a sub path (e.g. https://thoeni.io/blog)
- `-depth`: depth of the search. A value of `-1` means "unbounded". It is defaulted to `0` which can be useful to visualise a specific subset of the graph, for example: `./grawler -depth=0 -url=https://thoeni.io/blog/` will show all the blog pages referenced by the blog section removing unnecessary noise
- `-output`: output file, if specified will generate a file with named as indicated containing the minified json output (e.g. `./grawler -url=https://thoeni.io -depth=1 -output=out.json`)
- `-graph`: graph option, if set to true will generate a js file named `sitemap.js` that can be used by the Javascript visualisation tool (e.g. `./grawler -url=https://thoeni.io -depth=0 -graph=true`)
   Note: in order to use the Javascript visualisation page there's a huge limitation on the number of nodes (pages elements) that can be handled (which means pages and links), which ideally should be less than 50, therefore a `-depth=0` is recommended when using this option.
  
## Visual output

The tool allows to specify a `-graph=true` flag for `Javascript` formatting, which means creating a `.js` file with a `var sitemap = {}` in it that can be easily imported and parsed by an open source library called [Dracula](https://www.graphdracula.net/) which can print it like a directed graph. The limitation of the library is quite strong as with the increasing number of nodes the browser could just give up (especially Internet Explorer).

It's still a nice way to represent a page with a depth level of `0` to see all the first-level children pages.

## Testing
Tests depend on a utility library for assertions, therefore prior to running tests either one of the two following commands is needed:

 - `go get -t`
 - `dep ensure` if the go `dep` tool is installed

This being done, a `go test -v` will run the tests in the folder.

### Improvements

 - a `sync.Map` could be used in place of the current `map` guarded by `sync.Mutex`
 - a library for parsing html files could be used
 - on the visualisation side, anything deeper than `-depth=0` becomes unreadable: a proper frontend with the ability to dive deeper into each node one level at a time could be useful