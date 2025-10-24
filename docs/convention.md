# Coding Conventions

## Variable Names

For variable names, all export variables should be written in 
[PascalCase](https://en.wikipedia.org/wiki/Camel_case#Variations_and_synonyms).

Local variables should be written in camelCase. Acronyms are to be written with
only the first letter capitalized in all cases except where said acronym is
at the beginning of a local variable. I.E. between `HTMLParser` and `HtmlParser`

For an export variable, the latter would be preferred, and between `parsedHTML`
and `parsedHtml` for a local variable, the latter would also be preferred.[^1]
This rule extends to shorter acronyms such as `Id` or `Io`.

## Indentation

For all languages, tabs are preferred over spaces, as is Go convention.[^2]

## File and Folder Names

In general, choose short, lowercase folder and filenames.

In Go, use short filenames and try not to include hyphens or underscores,
unless required for certain behavior (`_test.go`).

In Angular, lowercase kebab-case is preferred for both folder names and file names.[^3]

For any other programming languages, the most standard convention
for that language should be followed.

Most plaintext and markdown documents should have a lowercase filename,
however certain special documents may be in all caps
(excluding the extension which remains lowercase).

Examples include but are not limited to README.md, LICENSE, COPYING, etc.[^4]

## Git Hooks

In each repository, a `.githooks` directory is present.

Follow the instructions in `.githooks/README.md` to configure the
relevant code formatter and corresponding pre-commit git hooks.

[^1]: https://en.wikipedia.org/wiki/Camel_case#Programming_and_coding

[^2]: https://go.dev/doc/effective_go#formatting

[^3]: https://angular.dev/style-guide#separate-words-in-file-names-with-hyphens

[^4]: https://en.wikipedia.org/wiki/README#As_a_generic_term
