# Getting Started

Datastar brings libraries like [AlpineJs](https://alpinejs.dev/) (frontend reactivity) and [HTMX](https://htmx.org/) (backend reactivity) together, in one cohesive solution. It's a tiny (extensible) framework that allows you to:

1. Manage state and build reactivity into your frontend using HTML attributes.
2. Modify the DOM and state by sending HTML fragments from your backend.

With Datastar, you can build any UI that a full-stack framework like React, Vue.js or Svelte can, using a much simpler, hypermedia-driven approach.

<div class="alert alert-info">
    <p>
        We’re so confident that Datastar can be used as a JavaScript framework replacement that we challenge anyone to find a use-case that Datastar _cannot_ handle!
    </p> 
</div>

## Installation

### Script Tag

Include Datastar in your HTML using a script tag:

```html
<script type="module" src="https://cdn.jsdelivr.net/npm/@sudodevnull/datastar@PACKAGE_VERSION/dist/datastar.min.js" defer></script>
```

If you prefer to host the file yourself, [download](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar) and include it from the appropriate path:
    
```html
<script type="module" src="/path/to/datastar.min.js" defer></script>
```

If you want a version with source maps, download and include the [unminified file](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.js) and the [source map](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.js.map).

### NPM

You can install Datastar via [npm](https://www.npmjs.com/package/@sudodevnull/datastar) and then use `node_modules/@sudodevnull/datastar/dist/datastar.js` (or `datastar.min.js`).

```bash
npm install @sudodevnull/datastar
```

## Handling State

Let's start with how Datastar handles state. Enter the [`data-store`](/reference/plugins_core#store) attribute.

```html
<div data-store="{ title: '' }"></div>
```

This is a global store. If you add `data-store` to multiple elements, the values will be merged into one single store (values defined later in the DOM tree override those defined earlier). The value must be written as a JavaScript object literal _or_ using JSON syntax.

Store values can be nested, which is useful for namespacing values with similar names.

```html
<div data-store="{ primary: { title: '' }, secondary: { title: '' }}"></div>
```

## Adding Reactivity

Datastar provides us with a way to set up two-way data binding on an element. In our case this `input` element. Say hello to the [`data-model`](/reference/plugins_attributes#model) attribute.

```html
<input data-model="title" type="text" placeholder="Type here!">
```

This binds the input field's value to the store value of the same name (`title`). If either is changed, the other will automatically update. 

Good stuff so far. So how can we see this? We can check the changes locally using the [`data-text`](/reference/plugins_attributes#text) attribute.

```html
<div data-text="$title"></div>
```

This sets the text content of an element to the store value with the name `title`. The `$` indicates that `$title` is a store value.

The value of the `data-text` attribute is an expression that is evaluated, meaning that we can include JavaScript in it.

```html
<div data-text="$title.toUpperCase()"></div>
```

Another common attribute is `data-show`, which can be used to show or hide an element based on whether a JavaScript expression evaluates to `true` or `false`.

```html
<input data-show="$title != ''" type="submit" value="Save">
```

This results in the submit button being visible only when the title is _not_ an empty string.

## Events

The [`data-on-*](/reference/plugins_attributes#on) attribute can be used to execute a JavaScript expression whenever an event is triggered on an element. 

```html
<button data-on-click="$title = 'New title'">
    Reset
</button>
```

This results in the `title` store value being set to `New title` when the button element is clicked. If the `title` store value is used elsewhere, its value will automatically update.

So what else can we do with these expressions? Anything we want, really:. 

```html
<div data-on-click="$prompt = prompt('Enter a value', $prompt); confirm('Are you sure?') && console.log($prompt)">
    Click me to log a prompt value
</div>
```

We've only scratched the surface of frontend reactivity, but let's take a look at backend reactivity.

## Backend Setup

Datastar uses [Server-Sent Events](https://en.wikipedia.org/wiki/Server-sent_events) or SSE. There's no special backend plumbing required to use SSE, just some special syntax. Fortunately, SSE is straightforward and [provides us with many advantages](/essays/event_streams_all_the_way_down).

First, set up your backend in the language of your choice. Using one of the official SDKs (Go, PHP, TypeScript, .NET) will help you get up and running faster. We're going to use the SDKs in the examples below, which set the appropriate headers and format the events for us, but this is optional.

```php
use starfederation\datastar\ServerSentEventGenerator;

// Creates a new `ServerSentEventGenerator` instance.
$sseGenerator = new ServerSentEventGenerator();

// Swaps out an existing fragment in the DOM.
$sseGenerator->renderFragment('<div id="main">Hello, world!</div>');
```

The `renderFragment` method sends an event to the client with the HTML fragment to replace the target element with the ID `main`.

```html
<div data-text="$input"></div>
```

To this:

```html
<div id="output"></div>
```

Give ourselves a button to perform this action.

Add this to your `<main>` element:

```html
<button data-on-click="$$put('/put')">Send State</button>
```

...and give ourselves a place to show our new state on the client.

Voilà! Now if you check out what you've done, you'll find you're able to send data to your `/put` endpoint and respond with HTML updating the output `div`. Neato!

Let's retrieve the backend data we're now storing.

Add this to your server code:

```js
app.get("/get", (req, res) => {
  setHeaders(res);

  const output = `Backend State: ${JSON.stringify(backendData)}.`;
  let frag = `<div id="output2">${output}</div>`;

  sendSSE({
    res,
    frag,
    end: true,
  });
});
```

And this to your HTML:

```html
<button data-on-click="$$get('/get')">Get Backend State</button>
<div id="output2"></div>
```

We're now fetching state that's stored on the backend.

Let's try something for fun. In your `/get` route, change your call to `sendSSE` so that we do not immediately end the request connection.

Change your `sendSSE` function call in your `\get` route.

```js
sendSSE({
  ...
  end: false,
});
```

Add this to your `sendSSE` function below the first call:

```js
frag = `<div id="output3">Check this out!</div>;`;
sendSSE({
  res,
  frag,
  selector: "#main",
  mergeType: "prepend",
  end: true,
});
```

Now you'll notice you're sending two events in one call. That's because Datastar uses SSE. So using `prepend` we're able to prepend what we want to a target element. We do this using a `selector` and in our case this is the `<main>` element. Good stuff! You can check out all of Datastar's event types [here](/reference/plugins_backend).

There's one last thing we're going to do. Let's add a simple data feed upon loading the page.

Copy this to your server code:

```js
app.get("/feed", async (req, res) => {
  setHeaders(res);
  while (res.writable) {
    const rand = randomBytes(8).toString("hex");
    const frag = `<span id="feed">${rand}</span>`;
    sendSSE({
      res,
      frag,
      end: false,
    });
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  res.end();
});
```

Add this inside your `<main>` element:

```html
<div>
  <span>Feed from server: </span>
  <span id="feed" data-on-load="$$get('/feed')"></span>
</div>
```

I told you we would use another `data-on` action earlier and here it is. `data-on-load` will perform this request when the page loads. If you check things out now you should see a feed that updates using SSE upon loading. Cool!

Datastar supports all the verbs without requiring a `<form>` element: `GET, POST, PUT, PATCH, DELETE`.

So that concludes our primer! Check out the full code for our Node example [here](/examples/node).

If you're still here I imagine you want to know more. Let's define things a little better.

## A Better View

To be more precise, think of Datastar as an extension to HTML's [data attributes](https://developer.mozilla.org/en-US/docs/Learn/HTML/Howto/Use_data_attributes). Using attributes, you can introduce state to your frontend, then access it anywhere in your DOM, or a backend of your choice. You can also set up events that trigger endpoints, then respond with HTML that targets fragments of your DOM.

- Declare global state: `data-store="{foo: ''}"`
- Link-up HTML elements to state slots: `data-model="foo"`
- Adjust HTML elements text content: `data-text="$foo"`
- Hookup other effects on your DOM to the state: `data-show="$foo"`
- Setup events using `data-on-click="$$get(/endpoint)"`
- Respond in HTML wrapped in SSE with a target element ID to update

It's that simple. To dive deeper check out some of the other links or just click below.
