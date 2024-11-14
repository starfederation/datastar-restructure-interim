# Getting Started

Datastar brings libraries like [AlpineJs](https://alpinejs.dev/) (frontend reactivity) and [HTMX](https://htmx.org/) (backend reactivity) together, in one cohesive solution. It's a tiny (extensible) framework that allows you to:

1. Manage state and build reactivity into your frontend using HTML attributes.
2. Modify the DOM and state by sending HTML fragments from your backend.

With Datastar, you can build any UI that a full-stack framework like React, Vue.js or Svelte can, using a much simpler, hypermedia-driven approach.

<div class="alert alert-info">
    <p>
        We're so confident that Datastar can be used as a JavaScript framework replacement that we challenge anyone to find a use-case for a web app that Datastar _cannot_ realistically be used to build!
    </p> 
</div>

## Installation

### Script Tag

The quickest way to use Datastar is to include it in your HTML using a script tag hosted on a CDN.

```html
<script type="module" defer src="https://cdn.jsdelivr.net/npm/@sudodevnull/datastar"></script>
```

If you prefer to host the file yourself, [download](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.min.js) the file or create your own custom bundle using the [bundler](/bundler), then include it from the appropriate path:
    
```html
<script type="module" defer src="/path/to/datastar.min.js"></script>
```

If you want a version with source maps, download and include the [unminified file](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.js) and the [source map](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.js.map).

### NPM

You can alternatively install Datastar via [npm](https://www.npmjs.com/package/@sudodevnull/datastar) and then use `node_modules/@sudodevnull/datastar/dist/datastar.js` (or `datastar.min.js`).

```bash
npm install @sudodevnull/datastar
```

## Handling State

Let's start with how Datastar allows you to handle state using the [`data-store`](/reference/plugins_core#store) attribute.

```html
<div data-store="{title: ''}"></div>
```

This is a global store. If you add `data-store` to multiple elements, the values will be merged into the global store (values defined later in the DOM tree override those defined earlier). The value must be written as a JavaScript object literal _or_ using JSON syntax.

Store values are nestable, which can be useful for namespacing values.

```html
<div data-store="{primary: {title: ''}, secondary: {title: '' }}"</div>
```

## Adding Reactivity

Datastar provides us with a way to set up two-way data binding on an element using the [`data-model`](/reference/plugins_attributes#model) attribute, which can be place on any form field (`input`, `textarea`, `select`, `checkbox` and `radio` elements).

```html
<input data-model="title" type="text" placeholder="Type here!">
```

This binds the input field's value to the store value of the same name (`title`). If either is changed, the other will automatically update. 

To see this in action, we can use the [`data-text`](/reference/plugins_attributes#text) attribute.

```html
<div data-text="$title"></div>
```

This sets the text content of an element to the store value with the name `title`. The `$` indicates that `$title` is a store value.

The value of the `data-text` attribute is an expression that is evaluated, meaning that we can include JavaScript in it.

```html
<div data-text="$title.toUpperCase()"></div>
```

<div data-store="{title1: ''}" class="alert flex flex-col items-start p-8">
    <div>
        Title:
        <span data-text="$title1.toUpperCase()"></span>
    </div>
    <input data-model="title1" placeholder="Enter a title" class="input input-bordered">
</div>

Another common attribute is `data-show`, which can be used to show or hide an element based on whether a JavaScript expression evaluates to `true` or `false`.

```html
<input data-show="$title != ''" type="submit" value="Save">
```

This results in the submit button being visible only when the title is _not_ an empty string.

<div data-store="{title2: ''}" class="alert flex flex-col items-start p-8">
    <div>
        Title:
        <span data-text="$title2.toUpperCase()"></span>
    </div>
    <input data-model="title2" placeholder="Enter a title" class="input input-bordered">
    <button data-show="$title2 != ''">Save</button>
</div>

## Events

The [`data-on-*](/reference/plugins_attributes#on) attribute can be used to execute a JavaScript expression whenever an event is triggered on an element. 

```html
<button data-on-click="$title = 'New title'">
    Reset
</button>
```

This results in the `title` store value being set to `New title` when the button element is clicked. If the `title` store value is used elsewhere, its value will automatically update.

<div data-store="{title3: ''}" class="alert flex flex-col items-start p-8">
    <div>
        Title:
        <span data-text="$title3.toUpperCase()"></span>
    </div>
    <input data-model="title3" placeholder="Enter a title" class="input input-bordered">
    <button data-on-click="$title3 = 'New title'">Reset</button>
</div>

So what else can we do with these expressions? Well anything we want, really:. 

```html
<button data-on-click="$prompt = prompt('Enter a value', $prompt); confirm('Are you sure?') && console.log($prompt)">
    Click me to log a prompt value
</button>
```

<div data-store="{prompt: ''}" class="alert flex flex-col items-start p-8">
    <button data-on-click="$prompt = prompt('Enter a value', $prompt); confirm('Are you sure?') && console.log($prompt)">
        Click me to log a prompt value
    </button>
</div>

We've only scratched the surface of frontend reactivity, but let's now take a look at how the backend can come into play.

## Backend Setup

Datastar uses [Server-Sent Events](https://en.wikipedia.org/wiki/Server-sent_events) or SSE. There's no special backend plumbing required to use SSE, just some special syntax. Fortunately, SSE is straightforward and [provides us with many advantages](/essays/event_streams_all_the_way_down).

First, set up your backend in the language of your choice. Using one of the helper SDKs (Go, PHP, TypeScript, .NET) will help you get up and running faster. We're going to use the SDKs in the examples below, which set the appropriate headers and format the events for us, but this is optional.

The following code would exist in a controller action endpoint in your backend.

```php
use starfederation\datastar\ServerSentEventGenerator;

// Creates a new `ServerSentEventGenerator` instance.
$sseGenerator = new ServerSentEventGenerator();

// Updates the `title` store value.
$sseGenerator->patchStore(['title' => 'Greetings']);

// Swaps out an existing fragment in the DOM.
$sseGenerator->renderFragment('<div id="greeting">Hello, world!</div>');
```

The `patchStore()` method updates one or more store values in the frontend, or creates them if they don't already exist.

The `renderFragment()` method sends an event to the client with the HTML fragment, which replaces the target element with the ID `greeting`. An element with the ID `greeting` must already exist in the DOM.

With our backend in place, we can now use a `data-on-click` on a button to send a `GET` request to the `/actions/greeting` endpoint on the server.

```html
<div data-store="{title: ''}"></div>
    <h1 data-text="$title"></h1>
    <div id="greeting"></div>
    <button data-on-click="$get('/actions/greeting')">
        Request a greeting
    </button>
</div>
```

Now when the button is clicked, the server will respond with a new greeting, updating the `title` store value and the `greeting` element in the DOM. We're driving state from the backend – neato!

We're not limited to just `GET` requests. We can also send `POST`, `PUT`, `PATCH` and `DELETE` requests, using `$post()`, `$put()`, `$patch()` and `$delete()` respectively.

```html
<button data-on-click="$post('/actions/greeting')">
    Send a greeting
</button>
```    

One of the advantages of using SSE is that we can send multiple events (HTML fragments, store value updates, etc.) in a single response.

```php
$sseGenerator->patchStore(['title' => 'Greetings']);
$sseGenerator->patchStore(['subtitle' => 'Earthlings']);
$sseGenerator->renderFragment('<div id="greeting-world">Hello, world!</div>');
$sseGenerator->renderFragment('<div id="greeting-universe">Hello, universe!</div>');
```

## An Overview of What's Possible

You can think of Datastar as an extension to HTML's [data attributes](https://developer.mozilla.org/en-US/docs/Learn/HTML/Howto/Use_data_attributes). Using `data-*` attributes (hence the name), you can introduce state to your frontend and access it anywhere in the DOM or from your backend. You can set up events that trigger requests to endpoints that respond with HTML fragments and store updates.

- Declare global state: `data-store="{foo: ''}"`
- Bind element values to store values: `data-model="foo"`
- Set the text content of an element to an expression.: `data-text="$foo"`
- Show or hide an element using an expression: `data-show="$foo"`
- Modify the classes on an element: `data-class="{'font-bold': $foo}"`
- Bind an expression to an HTML attribute: `data-bind-disabled="$foo == ''"`
- Execute an expression whenever an event is triggered on an element: `data-on-click="$get(/endpoint)"`
- Persist all store values in local storage: `data-persist`
- Create a new computed store value from an expression: `data-computed-foo="'Hello, ' + $name"`
- Create a reference to an element that can be referenced: `data-ref="alert"`
- Send a header along with a request: `data-header-foo="{'x-powered-by': $foo}"`
- Replaces the URL: `data-replace-url="'/page1'"`
