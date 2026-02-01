# Stop stopPropagation!

A clickable card, a nested button, and stopPropagation to clean the mess. This setup holds a spot of honor in the long and sad history of bad practices.

## TL;DR:

Nested interactive elements often lead to accidental event propagation and introduce serious accessibility issues for keyboard and screen reader users. Stopping propagation in inner handlers hides useful DOM behavior and can break tooling. A better approach is to avoid nested interactivity altogether. Either by flattening the DOM and using CSS for visual stacking, or, ideally, by designing cards with clearly separated clickable regions that make each action explicit and accessible.

## The problem

Clicking the inner element also triggers the outer element’s event listener, because events propagate up the DOM tree. Additionaly, nesting interactive elements confuses plugins and assistive technologies.

## Worst solution

`evt.stopPropagation();`

Many developers prevent propagation in the inner listener. But event propagation is a useful DOM feature. You may later want to listen to all document events for auditing, and some plugins or accessibility tools rely on propagation as well. And, of course, it doesn't address accessibility issues.

## So-so solution

`if(evt.target === inner) return;`

Ignore inner-element events in the outer listener. This avoids interfering with the events’ default behavior. Unfortunately, accessibility still suffers.

The W3C warns against nesting interactive elements. This practice requires event-handler workarounds and challenges accessibility tools like screen readers.

During keyboard navigation, should focus go to the inner element or the outer one first?

When a screen reader is on the outer element, should it read the inner element’s text? If it doesn’t, information is lost. If it does, users may think the outer element triggers the inner action.

## Better solution

Flatten the DOM and use CSS to visually stack elements. Follow accessibility best practices, including correct semantics and focus order.

Because the elements are DOM siblings, clicks on the inner element don’t propagate to the outer one. Accessibility improves, since assistive technologies prioritize DOM and focus order over visual presentation.

## Best solution

Whenever possible, avoid nesting interactive elements in either the DOM or the presentation. Instead, split the card into clearly labeled clickable zones. This helps users and tools alike to understand what each click does.

How do you handle nested interactive elements in cards?