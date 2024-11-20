## SSE Indicator

## Demo

<div class="flex flex-col gap-4">
  <div class="flex gap-2">
    <div data-class="{loading: $fetching}" class="indicator flex items-center gap-2">
      <iconify-icon icon="svg-spinners:blocks-wave"></iconify-icon>
      <span>Loading</span>
    </div>
    <button class="flex-1 btn btn-primary" data-on-click="$get('/examples/fetch_indicator/greet')" data-sse-indicator="fetching" data-testid="greeting_button" data-bind-disabled="$fetching" >
      Click me for a greeting
    </button>
  </div>
  <div id="greeting"></div>
</div>

## Explanation

```html
<style>
    .indicator {
        opacity:0;
        transition: opacity 300ms ease-out;
    }
    .indicator.loading {
        opacity:1;
        transition: opacity 300ms ease-in;
    }
</style>
<div
  class="indicator"
  data-class="{loading: $fetching}"
>
    Loading Indicator
</div>
<button
  data-on-click="$get('/examples/fetch_indicator/greet')"
  data-sse-indicator="fetching"
  data-bind-disabled="$fetching" >
  Click me for a greeting
</button>
<div id="greeting"></div>
```

The `data-sse-indicator` attribute accepts the name of a signal whose value is set to `true` when a fetch request initiated from the same element is in progress, otherwise `false`. If the signal does not exist in the store, it will be added.
