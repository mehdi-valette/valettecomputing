# JSDoc

JSDoc is a documentation standard for annotating JavaScript code. To add a JSDoc comment, start with `/**` (notice the double asterisk), and add metadata using `@tag info` syntax.

With JSDoc, you can provide information to other developers

```javascript
/** @license Apache-2.0 */

/**
* @author Alice
* @description This function does nothing
*/
function test() {
  return null;
}
```

and to your IDE (you may need to adjust your IDE's configuration).

```javascript
/** @type {null | string} */
let myVar = null;

myVar = 45; // IDE warning: Type 'number' is not assignable to type 'string'
```

With JSDoc, you can extract your code's documentation into HTML files, show type-related warnings in your IDE, or even type-check your files before a Git push (for example, as part of a pre-commit or pre-push hook).

## Documentation

JSDoc is primarily built for documentation. You can annotate files, classes, functions, and variables. Compared to simple comments, JSDoc allows your IDE to warn about deprecated functions or indicate which events are emitted by a class. Tools such as `jsdoc/jsdoc` and `jsdoc-to-markdown` look for JSDoc tags and turn them into HTML or Markdown documentation.

Here are a few examples of what you can document with JSDoc:

- Legal information: author, license, and copyright.

- Versioning: version, deprecated, and since (the date a feature was added).

- Code explanation: description, events emitted or listened to, and dependencies.

## Type information

JSDoc supports type annotations that your IDE uses to type-check your code. Adding types is arguably *the* biggest help for advanced scripting in JavaScript. They enable autocompletion, warnings, and errors, and catch many bugs during development.

JSDoc can naturally be used with TypeScript, offering an easy integration between \*.js, \*.ts, and \*.d.ts files.

## Advantages over TypeScript

JSDoc declarations are just comments, you don't need transpilation (i.e. transforming *.ts => *.js). No need for transpilers, source map files, or separate source and output directories in many cases. Enjoy the power of types with the simplicity of JavaScript!

## Limitations compared to TypeScript

JSDoc is more verbose than TypeScript, especially for advanced types. Also, JSDoc cannot add keywords to JavaScript; there are no enums, decorators, or non-null assertion operators.

Additionally, some frameworks depend on TypeScript features. For example, NestJS uses TypeScript's experimental decorators, making TypeScript mandatory.

## Conclusion

JSDoc allows you to type your code and add documentation simply by writing comments. While its syntax isn't as elegant as TypeScript's, it's a simpler alternative that gets you pretty far. Besides, JSDoc works in TypeScript files too.

Have you ever used JSDoc to type your project?
