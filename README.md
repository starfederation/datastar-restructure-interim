![Version](https://img.shields.io/npm/v/@sudodevnull/datastar)
![npm bundle size](https://img.shields.io/bundlephobia/minzip/%40sudodevnull%2Fdatastar)

<p align="center"><img width="200" src="https://media.githubusercontent.com/media/starfederation/datastar/refs/heads/main/code/go/site/static/images/rocket.gif"></p>

# Datastar

### A real-time hypermedia framework.

Datastar helps you build real-time web applications with the simplicity of server-side rendering and the power of a full-stack SPA framework.

Here’s what frontend reactivity looks like using Datastar:

```html
<div data-store="{input: ''}">
    <input data-model="input" type="text">
    <div data-text="$input.toUpperCase()"></div>
    <button data-on-click="$post('/endpoint')">Save</button>
</div> 
```

## Getting Started

Include Datastar with a single 14 KiB file and start adding reactivity to your frontend immediately. Write your backend in the language of your choice, and use the helper SDKs (available for Go, PHP, TypeScript and .NET) to get up and running even faster.

```html
<script type="module" defer src="https://cdn.jsdelivr.net/npm/@sudodevnull/datastar"></script>
```

Read the [Getting Started guide »](https://data-star.dev/guide/getting_started)  
Visit the [Datastar website »](https://data-star.dev/)

## Contributing

Read the [Contribution Guidelines »](CONTRIBUTING.md)

## Development

Read the [Development Guidelines »](DEVELOPMENT.md)
