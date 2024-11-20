# Getting Started

Datastar brings the functionality provided by libraries like [AlpineJs](https://alpinejs.dev/) (frontend reactivity) and [HTMX](https://htmx.org/) (backend reactivity) together, into one cohesive solution. It's a lightweight, extensible framework that allows you to:

1. Manage state and build reactivity into your frontend using HTML attributes.
2. Modify the DOM and state by sending events from your backend.

With Datastar, you can build any UI that a full-stack framework like React, Vue.js or Svelte can, but with a much simpler, hypermedia-driven approach.

<div class="alert alert-info">
    <iconify-icon icon="simple-icons:rocket"></iconify-icon>
    <div>
        We're so confident that Datastar can be used as a JavaScript framework replacement that we challenge anyone to find a use-case for a web app that Datastar <em>cannot</em> realistically be used to build!
    </div>
</div>

## Installation

### Using a Script Tag

The quickest way to use Datastar is to include it in your HTML using a script tag hosted on a CDN.

### BEN this is out of date start: import maps are the way to go
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
### BEN this is out of date end

## Handling State

### BEN data-store isn't needed at this stage to teach now!

Let's take a look at how Datastar allows you to handle state using the [`data-store`](/reference/plugins_core#store) attribute.

```html
<div data-store="{input: ''}"></div>
```

The “store” is a global collection of reactive signals. You can think of signals as variables that automatically track and propagate changes to expressions. Don't worry if this sounds complicated; it will become clearer as we look at some examples.

If you add `data-store` to multiple elements, the signals provided will be _merged_ into the store (values defined later in the DOM tree override those defined earlier).

Signals are nestable, which can be useful for namespacing. The `data-store` value must be written as a JavaScript object literal _or_ using JSON syntax.

```html
<div data-store="{primary: {input: ''}, secondary: {input: '' }}"></div>
```

## Adding Reactivity

Datastar provides us with a way to set up two-way data binding on an element using the [`data-model`](/reference/plugins_attributes#model) attribute, which can be placed on `input`, `textarea`, `select`, `checkbox` and `radio` elements.

```html
<input data-model="input" type="text">
```

This binds the element's value to the signal of the same name. If either is changed, the other will automatically update.

Note how when using `data-model`, the value is the bare name of the signal (without the `$`).

To see this in action, we can use the [`data-text`](/reference/plugins_attributes#text) attribute.

```html
<div data-text="$input"></div>
```

<div class="alert flex justify-between items-start p-8">
    <div class="flex flex-col gap-4">
        <div class="flex items-center">
            <div class="w-20">Input:</div>
            <input data-model="input1" class="input input-bordered">
        </div>
        <div class="flex items-center">
            <div class="w-20">Output:</div>
            <div data-text="$input1" class="output"></div>
        </div>
    </div>
</div>

This sets the text content of an element to the value of the signal `$input`. The `$` is required to denote a signal.

The value of the `data-text` attribute is an expression that is evaluated, meaning that we can include JavaScript in it.

```html
<div data-text="$input.toUpperCase()"></div>
```

<div class="alert flex justify-between items-start p-8">
    <div class="flex flex-col gap-4">
        <div class="flex items-center">
            <div class="w-20">Input:</div>
            <input data-model="input2" class="input input-bordered">
        </div>
        <div class="flex items-center">
            <div class="w-20">Output:</div>
            <div data-text="$input2.toUpperCase()" class="output"></div>
        </div>
    </div>
</div>

The `data-computed-*` attribute creates a new signal that is computed based on an expression. The computed signal is read-only, and its value is automatically updated when any signals in the expression are updated.

```html
<div data-computed-repeated="$input.repeat(2)">
    <input data-model="input" type="text">
    <div data-text="$repeated"></div>
</div>
```

<div data-computed-repeated="$input3.repeat(2)" class="alert flex justify-between items-start p-8">
    <div class="flex flex-col gap-4">
        <div class="flex items-center">
            <div class="w-20">Input:</div>
            <input data-model="input3" class="input input-bordered">
        </div>
        <div class="flex items-center">
            <div class="w-20">Output:</div>
            <div data-text="$repeated" class="output"></div>
        </div>
    </div>
</div>

Another useful attribute is `data-show`, which can be used to show or hide an element based on whether a JavaScript expression evaluates to `true` or `false`.

```html
<button data-show="$input != ''">Save</button>
```

This results in the button being visible only when the input is _not_ empty.

<div data-store="{input4: ''}" class="alert flex justify-between items-start p-8">
    <div class="flex flex-col gap-4">
        <div class="flex items-center">
            <div class="w-20">Input:</div>
            <input data-model="input4" class="input input-bordered">
        </div>
        <div class="flex items-center">
            <div class="w-20">Output:</div>
            <div data-text="$input4" class="output"></div>
        </div>
    </div>
    <button data-show="$input4 != ''" class="btn btn-primary">
        Save
    </button>
</div>

The `data-bind-*` attribute can be used to bind a JavaScript expression to any valid HTML attribute.

```html
<button data-bind-disabled="$input == ''">Save</button>
```

This results in the button being given the `disabled` attribute whenever the input is empty.

<div data-store="{input5: ''}" class="alert flex justify-between items-start p-8">
    <div class="flex flex-col gap-4">
        <div class="flex items-center">
            <div class="w-20">Input:</div>
            <input data-model="input5" class="input input-bordered">
        </div>
        <div class="flex items-center">
            <div class="w-20">Output:</div>
            <div data-text="$input5" class="output"></div>
        </div>
    </div>
    <button data-bind-disabled="$input5 == ''" class="btn btn-primary">
        Save
    </button>
</div>

## Events

The [`data-on-*`](/reference/plugins_attributes#on) attribute can be used to execute a JavaScript expression whenever an event is triggered on an element.

```html
<button data-on-click="$input = ''">
    Reset
</button>
```

This results in the `$input` signal being set to an empty string when the button element is clicked. If the `$input` signal is used elsewhere, its value will automatically update.

<div data-store="{input6: 'Some input'}" class="alert flex justify-between items-start p-8">
    <div class="flex flex-col gap-4">
        <div class="flex items-center">
            <div class="w-20">Input:</div>
            <input data-model="input6" class="input input-bordered">
        </div>
        <div class="flex items-center">
            <div class="w-20">Output:</div>
            <div data-text="$input6" class="output"></div>
        </div>
    </div>
    <button data-on-click="$input6 = ''" class="btn btn-secondary">
        Reset
    </button>
</div>

So what else can we do with these expressions? Anything we want, really.

See if you can follow the code below _before_ trying the demo.

```html
<div data-store="{response: '', answer: 'bread'}"
     data-computed-correct="$response.toLowerCase() == $answer"
>
    <div id="question">
        What do you put in a toaster?
    </div>
    <button data-on-click="$response = prompt('Answer:')">
        BUZZ
    </button>
    <div data-show="$response != ''">
        You answered “<span data-text="$response"></span>”.
        <span data-show="$correct2">That is correct ✅</span>
        <span data-show="!$correct2">
            The correct answer is “<span data-text="$answer2"></span>” 🤷
        </span>
    </div>
</div>
```

<div data-store="{response1: '', answer1: 'bread'}" data-computed-correct1="$response1.toLowerCase() == $answer1" class="alert flex justify-between items-start gap-4 p-8">
    <div class="space-y-3">
        <div id="question1">
            What do you put in a toaster?
        </div>
        <div data-show="$response1 != ''">
            You answered “<span data-text="$response1"></span>”.
            <span data-show="$correct1">That is correct ✅</span>
            <span data-show="!$correct1">
                The correct answer is “<span data-text="$answer1"></span>” 🤷
            </span>
        </div>
    </div>
    <button data-on-click="$response1 = prompt('Answer:')" class="btn btn-primary">
        BUZZ
    </button>
</div>

We've just scratched the surface of frontend reactivity. Now let's take a look at how we can bring the backend into play.

## Backend Setup

Datastar uses [Server-Sent Events](https://en.wikipedia.org/wiki/Server-sent_events) or SSE. There's no special backend plumbing required to use SSE, just some special syntax. Fortunately, SSE is straightforward and [provides us with some advantages](/essays/event_streams_all_the_way_down).

First, set up your backend in the language of your choice. Using one of the helper SDKs (available for Go, PHP, TypeScript and .NET) will help you get up and running faster. We're going to use the SDKs in the examples below, which set the appropriate headers and format the events for us, but this is optional.

The following code would exist in a controller action endpoint in your backend.

!!!CODE_SNIPPET:getting_started/setup!!!

The `mergeFragments()` method merges the provided HTML fragment into the DOM, replacing the element with `id="question"`. An element with the ID `question` must already exist in the DOM.

The `mergeSignals()` method merges the `response` and `answer` signals into the frontend store.

With our backend in place, we can now use the `data-on-click` attribute to send a `GET` request to the `/actions/quiz` endpoint on the server when a button is clicked.

```html
<div data-store="{response: '', answer: ''}"
     data-computed-correct="$response.toLowerCase() == $answer"
>
    <div id="question"></div>
    <button data-on-click="$get('/actions/quiz')">
        Fetch a question
    </button>
    <button data-show="$answer != ''"
            data-on-click="$response = prompt('Answer:') ?? ''"
    >
        BUZZ
    </button>
    <div data-show="$response != ''">
        You answered “<span data-text="$response"></span>”.
        <span data-show="$correct2">That is correct ✅</span>
        <span data-show="!$correct2">
            The correct answer is “<span data-text="$answer2"></span>” 🤷
        </span>
    </div>
</div>
```

Now when the `Fetch a question` button is clicked, the server will respond with an event to modify the `question` element in the DOM and an event to modify the `response` and `answer` signals. We're driving state from the backend!

<div data-store="{response2: '', answer2: '', lastQuestionId: ''}" data-computed-correct2="$response2.toLowerCase() == $answer2" class="alert flex justify-between items-start gap-4 p-8">
    <div class="space-y-3 pb-3">
        <div id="question2"></div>
        <div data-show="$response2 != ''">
            You answered “<span data-text="$response2"></span>”.
            <span data-show="$correct2">That is correct ✅</span>
            <span data-show="!$correct2">
                The correct answer is “<span data-text="$answer2"></span>” 🤷
            </span>
        </div>
        <button data-on-click="$get('/examples/quiz/data')" class="btn btn-secondary">
            Fetch a question
        </button>
    </div>
    <button data-show="$answer2 != ''" data-on-click="$response2 = prompt('Answer:') ?? ''" class="btn btn-primary">
        BUZZ
    </button>
</div>

We're not limited to just `GET` requests. We can also send `GET`, `POST`, `PUT`, `PATCH` and `DELETE` requests, using `$get()`, `$post()`, `$put()`, `$patch()` and `$delete()` respectively.

Here's how we could send an answer to the server for processing, using a `POST` request.

```html
<button data-on-click="$post('/actions/quiz')">
    Submit answer
</button>
```

One of the benefits of using SSE is that we can send multiple events (HTML fragments, signal updates, etc.) in a single response.

!!!CODE_SNIPPET:getting_started/multiple_events!!!

## A Quick Overview

Using `data-*` attributes (hence the name), you can introduce reactive state to your frontend and access it anywhere in the DOM and in your backend. You can set up events that trigger requests to backend endpoints that respond with HTML fragment and signal updates.

- Merge signals into the store: `data-store="{foo: ''}"`
- Bind element values to signals: `data-model="foo"`
- Set the text content of an element to an expression.: `data-text="$foo"`
- Show or hide an element using an expression: `data-show="$foo"`
- Modify the classes on an element: `data-class="{'font-bold': $foo}"`
- Bind an expression to an HTML attribute: `data-bind-disabled="$foo == ''"`
- Execute an expression on an event: `data-on-click="$get(/endpoint)"`
- Persist all signals in local storage: `data-persist`
- Create a computed signal: `data-computed-foo="$bar + 1"`
- Create a reference to an element: `data-ref="alert"`
- Send a header with a request: `data-header-foo="{'x-powered-by': $foo}"`
- Replace the URL: `data-replace-url="'/page1'"`
