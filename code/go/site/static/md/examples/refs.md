## Refs

## Demo

<div>
     <div data-ref="foo">I'm a div that is getting referenced</div>
     <pre
          class=""
          data-text="JSON.stringify(ctx.store(),null,2)"
     >
          Stuff in store
     </pre>
     <div data-text="`I'm using content of '${$foo}.innerHTML}'`"></div>
</div>

## Explanation

```html
<div>
  <div data-ref="foo">I'm a div that is getting referenced</div>
  <div data-text="`I'm using content of '${$foo}.innerHTML}'`"></div>
</div>
```

Adding `data-ref="foo"` to an element creates a signal called `$foo` that points to that element.
