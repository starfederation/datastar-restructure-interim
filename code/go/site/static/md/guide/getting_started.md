# Getting Started

Datastar brings the functionality provided by libraries like [AlpineJs](https://alpinejs.dev/) (frontend reactivity) and [HTMX](https://htmx.org/) (backend reactivity) together, into one cohesive solution. It's a tiny, extensible framework that allows you to:

1. Manage state and build reactivity into your frontend using HTML attributes.
2. Modify the DOM and state by sending events from your backend.

With Datastar, you can build any UI that a full-stack framework like React, Vue.js or Svelte can, using a much simpler, hypermedia-driven approach.

<div class="alert alert-info">
    <iconify-icon icon="simple-icons:rocket"></iconify-icon>
    <div>
        We're so confident that Datastar can be used as a JavaScript framework replacement that we challenge anyone to find a use-case for a web app that Datastar <em>cannot</em> realistically be used to build!
    </div>
</div>

## Installation

### Using a Script Tag

The quickest way to use Datastar is to include it in your HTML using a script tag hosted on a CDN.

```html
<script type="module" defer src="https://cdn.jsdelivr.net/npm/@sudodevnull/datastar"></script>
```

If you prefer to host the file yourself, [download](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.min.js) it or create your own [custom bundle](/bundler), then include it from the appropriate path.
    
```html
<script type="module" defer src="/path/to/datastar.min.js"></script>
```

If you want a version with source maps, download and include the [unminified file](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.js) and the [source map](https://cdn.jsdelivr.net/npm/@sudodevnull/datastar/dist/datastar.js.map).

### Using NPM

You can alternatively install Datastar via [npm](https://www.npmjs.com/package/@sudodevnull/datastar) and then use `node_modules/@sudodevnull/datastar/dist/datastar.js` (or `datastar.min.js`).

```bash
npm install @sudodevnull/datastar
```

## Handling State

Let's take a look at how Datastar allows you to handle state using the [`data-store`](/reference/plugins_core#store) attribute.

```html
<div data-store="{variable: ''}"></div>
```

This is a global store. If you add `data-store` to multiple elements, the values will be merged into the global store (values defined later in the DOM tree override those defined earlier). 

Store values are nestable, which can be useful for namespacing values. The values must be written as a JavaScript object literal _or_ using JSON syntax.

```html
<div data-store="{primary: {variable: ''}, secondary: {variable: '' }}"></div>
```

## Adding Reactivity

Datastar provides us with a way to set up two-way data binding on an element using the [`data-model`](/reference/plugins_attributes#model) attribute, which can be placed on `input`, `textarea`, `select`, `checkbox` and `radio` elements.

```html
<input data-model="variable" type="text" placeholder="Type here!">
```

This binds the input field's value to the store value of the same name (`variable`). If either is changed, the other will automatically update. 

To see this in action, we can use the [`data-text`](/reference/plugins_attributes#text) attribute.

```html
<div data-text="$variable"></div>
```

<div data-store="{variable1: ''}" class="alert flex flex-col items-start p-8">
    <div class="flex gap-2">
        Input:
        <input data-model="variable1" class="input input-bordered">
    </div>
    <div class="flex gap-2">
        Output:
        <div data-text="$variable1"></div>
    </div>
</div>

This sets the text content of an element to the store value with the name `variable`. The `$` in `data-text="$variable"` is required because `$variable` is a store value.

The value of the `data-text` attribute is an expression that is evaluated, meaning that we can include JavaScript in it.

```html
<div data-text="$variable.toUpperCase()"></div>
```

<div data-store="{variable2: ''}" class="alert flex flex-col items-start p-8">
    <div class="flex gap-2">
        Input:
        <input data-model="variable2" class="input input-bordered">
    </div>
    <div class="flex gap-2">
        Output:
        <div data-text="$variable2.toUpperCase()"></div>
    </div>
</div>

Another useful attribute is `data-show`, which can be used to show or hide an element based on whether a JavaScript expression evaluates to `true` or `false`.

```html
<button data-show="$variable != ''">Save</button>
```

This results in the button being visible only when the variable is _not_ empty.

<div data-store="{variable3: ''}" class="alert flex flex-col items-start p-8">
    <div class="flex gap-2">
        Input:
        <input data-model="variable3" class="input input-bordered">
    </div>
    <div class="flex gap-2">
        Output:
        <div data-text="$variable3"></div>
    </div>
    <button data-show="$variable3 != ''" class="btn btn-primary">
        Save
    </button>
</div>

The `data-bind-*` attribute can be used to bind a JavaScript expression to any valid HTML attribute.

```html
<button data-bind-disabled="$variable == ''">Save</button>
```

This results in the button being given the `disabled` attribute whenever the variable _is_ empty.

<div data-store="{variable4: ''}" class="alert flex flex-col items-start p-8">
    <div class="flex gap-2">
        Input:
        <input data-model="variable4" class="input input-bordered">
    </div>
    <div class="flex gap-2">
        Output:
        <div data-text="$variable4"></div>
    </div>
    <button data-bind-disabled="$variable4 == ''" class="btn btn-primary">
        Save
    </button>
</div>

## Events

The [`data-on-*](/reference/plugins_attributes#on) attribute can be used to execute a JavaScript expression whenever an event is triggered on an element. 

```html
<button data-on-click="$variable = ''">
    Reset
</button>
```

This results in the `variable` store value being set to an empty string when the button element is clicked. If the `variable` store value is used elsewhere, its value will automatically update.

<div data-store="{variable5: 'Some input'}" class="alert flex flex-col items-start p-8">
    <div class="flex gap-2">
        Input:
        <input data-model="variable5" class="input input-bordered">
    </div>
    <div class="flex gap-2">
        Output:
        <div data-text="$variable5"></div>
    </div>
    <button data-on-click="$variable5 = ''" class="btn btn-secondary">
        Reset
    </button>
</div>

So what else can we do with these expressions? Anything we want, really. 

Here's how we might set a store value using a confirmation dialog.

```html
<div data-store="{confirmed: false}">
    <button data-on-click="$confirmed = confirm('Are you sure?)">
        Confirm an action
    </button>
</div>
```

Now you try! See if you can figure out how you'd code the demo below _before_ you reveal the solution.

<div data-store="{answer: ''}" class="alert flex items-center gap-4 p-8">
    <button data-on-click="$answer = prompt()" class="btn btn-primary">
        Prompt for some input
    </button>
    <div data-show="$answer != ''">
        <span data-text="$answer"></span>
    </div>
</div>

<div data-store="{reveal: false}">
    <button data-on-click="$reveal = true" class="btn btn-primary">
        Reveal solution
    </button>
    <div data-show="$reveal">
        ```html
        <div data-store="{answer: ''}">
            <button data-on-click="$answer = prompt()">
                Click me
            </button>
            <div data-text="$answer"></div>
        </div>
        ```
    </div>
</div>

We've just scratched the surface of frontend reactivity. Now let's take a look at how we can bring the backend into play.

## Backend Setup

Datastar uses [Server-Sent Events](https://en.wikipedia.org/wiki/Server-sent_events) or SSE. There's no special backend plumbing required to use SSE, just some special syntax. Fortunately, SSE is straightforward and [provides us with some advantages](/essays/event_streams_all_the_way_down).

First, set up your backend in the language of your choice. Using one of the helper SDKs (available for Go, PHP, TypeScript and .NET) will help you get up and running faster. We're going to use the SDKs in the examples below, which set the appropriate headers and format the events for us, but this is optional.

The following code would exist in a controller action endpoint in your backend.

```php
use starfederation\datastar\ServerSentEventGenerator;

// Creates a new `ServerSentEventGenerator` instance.
$sseGenerator = new ServerSentEventGenerator();

// Updates the `title` store value.
$sseGenerator->mergeStore(['title' => 'Greetings']);

// Merges a fragment into the DOM.
$sseGenerator->mergeFragment('<div id="greeting">Hello, world!</div>');
```

The `mergeStore()` method merges one or more store values in the frontend store, or creates them if they don't already exist.

The `mergeFragment()` method merges the HTML fragment into the DOM, replacing the element with the ID `greeting`. An element with the ID `greeting` must already exist in the DOM.

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

Now when the button is clicked, the server will respond with a new greeting, updating the `title` store value and the `greeting` element in the DOM. We're driving state from the backend!

We're not limited to just `GET` requests. We can also send `GET`, `POST`, `PUT`, `PATCH` and `DELETE` requests, using `$get()`, `$post()`, `$put()`, `$patch()` and `$delete()` respectively.

```html
<button data-on-click="$post('/actions/greeting')">
    Send a greeting
</button>
```    

One of the advantages of using SSE is that we can send multiple events (HTML fragments, store value updates, etc.) in a single response.

```php
$sseGenerator->mergeStore(['title' => 'Greetings']);
$sseGenerator->mergeFragment('<div id="greeting-world">Hello, world!</div>');
$sseGenerator->mergeStore(['subtitle' => 'Earthlings']);
$sseGenerator->mergeFragment('<div id="greeting-universe">Hello, universe!</div>');
```

## An Overview of Datastar

Using `data-*` attributes (hence the name), you can introduce reactive state to your frontend and access it anywhere in the DOM and in your backend. You can set up events that trigger requests to backed endpoints that respond with HTML fragments and store updates.

- Declare global state: `data-store="{foo: ''}"`
- Bind element values to store values: `data-model="foo"`
- Set the text content of an element to an expression.: `data-text="$foo"`
- Show or hide an element using an expression: `data-show="$foo"`
- Modify the classes on an element: `data-class="{'font-bold': $foo}"`
- Bind an expression to an HTML attribute: `data-bind-disabled="$foo == ''"`
- Execute an expression on an event: `data-on-click="$get(/endpoint)"`
- Persist all store values in local storage: `data-persist`
- Create a computed store value: `data-computed-foo="$bar + 1"`
- Create a reference to an element: `data-ref="alert"`
- Send a header with a request: `data-header-foo="{'x-powered-by': $foo}"`
- Replaces the URL: `data-replace-url="'/page1'"`
