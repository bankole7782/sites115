# sites115

a static sites generator with search functionality


## How to use
1.  Install with `sudo snap install sites115`
1.  To create a site called `sitename` run `sites115 cs sitename`.
1.  To render a site run `sites115 rs sitename` to render a site called `sitename`
1.  To view your site run `sites115.sites115d /path/to/sitename/out`. Replace /path/to/sitename/out with the actual dir.
1.  To deploy use docker. The dockerfile needed is part of the created site.

## Folder Structure
- `out` contains the generated site
- `static` contains the assets that would not be rendered but returned as is.
- `stuffs` contains the html and markdown files that would be rendered.
- `templates` contains the templates that is used to render html and markdown files contained in `stuffs`


## Conventions used in folder `stuffs`

1.  Every folder must contain a 'index.html' for the path /
1.  Every sub-folder of `stuffs` must contain a toc.txt.
    Sample of a toc.txt is as follows:
    ```
    how_to_blog.md
    on_css_frameworks.md
    long_running_tasks.md
    ```
1.  Every html or markdown file must begin with a variables part. It must begin with `---` and end with `---`. Sample contents include
    ```
    template: base.html
    title: Saenuma - A beautiful programs website
    meta: Saenuma homepage. Saenuma delivers beautiful programs.
    keywords: programs, database, program, forms, git, stories
    ```
    it must contain the variables: template, title, meta and keywords. You can include your own variables.

1.  Include a '404.html' in your `stuffs` folder for not found pages.
1.  Include a 'search_results.html' in your `stuffs` folder for results of a search
1.  A markdown file must have the `template` and the `md_template`  variables.


## Conventions used in the folder `templates`

1.  All templates files must not begin with a variables section.


## Structs that would be passed to your templates

### Page
Which is of type `map[string]string` that would contain your page variables.

This is passed to non-index html pages and markdown pages.

### Paginator
Which is declared as
```go
type PaginatorStruct struct {
  Page int
  PaginationCount int
  Pages []map[string]string
  TotalPages int
  PreviousPage int
  PreviousPagePath string
  NextPage int
  NextPagePath string
  TotalPagesArr []int
}
```
Where `Pages` is a list of `Page` variables declared above. This `Page` variables would also contain a `url` field

Paginator object is passed to the 'index.html' and the 'search_results.html' pages.

### SearchStr
Contains the query for example. `s=bank+account`

This is passed to the 'search_results.html' pages.

### HTML
Contains a HTML generated from markdown.

This is passed to the markdown pages.


## About searches

1.  The search page should send queries using javascript to `/search_results`
1.  The end result should look like `/search_results?s=bank+account`


## Don't use Nginx or Apache

The projects ships with its own alternative in the folder `sites115d`.
This is the only way to enable search functionality.

The `sites115d` expects the path to a folder. This folder must be the the `out` folder described above.
