# Coding Conventions

## Variable Names

For variable names, all export variables should be written in 
[PascalCase](https://en.wikipedia.org/wiki/Camel_case#Variations_and_synonyms).
Local variables should be written in camelCase. Acronyms are to be written with
only the first letter capitalized in all cases except where said acronym is
at the beginning of a local variable. I.E. between `HTMLParser` and `HtmlParser`
for an export variable, the latter would be preferred, and between `parsedHTML`
and `parsedHtml` for a local variable, the latter would also be preferred.[^1]
This rule extends to shorter acronyms such as `Id` or `Io`.


[^1]: https://en.wikipedia.org/wiki/Camel_case#Programming_and_coding