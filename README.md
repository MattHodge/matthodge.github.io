# Making changes to the website
Altering the website is fairly easy. All content can be found in the `_pages` folder. This folders holds all Markdown files.

The website is build via the [Jekyll](https://jekyllrb.com/) static content generator. We've chosen to write our content in Markdown. If you create a new page make sure you add the [FrontMatter](https://jekyllrb.com/docs/frontmatter/) header to the file.

# Generating the website locally

Running the website locally is fairly easy. Please see the steps below for your operating system.

## On a Mac

Make sure you have Ruby and [Bundler](http://bundler.io/) installed.

Checkout this repository and from the root execute the following commands to install Jekyll and its dependencies.

```bash
bundle install --path _vendor/bundle
```

In the root of the folder execute the following command:

```bash
bundle exec jekyll serve --watch
```

The `--watch` will watch the files in the folder and rebuild the website after any change of the files.
