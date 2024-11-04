# Datastar

Real-time Hypermedia first Library and Framework

# How to include on page
```html
<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar/ts/library/dist/datastar-allinone.js"></script>
```

If you just want the core and add your own plugins


but how do I add my own plugins?

```html
<script type="importmap">
{
    "imports": {
      "datastar": "https://cdn.jsdelivr.net/gh/starfederation/datastar/ts/library/dist/datastar-core.js"
    }
}
</script>
<script type="module">
import {Datastar} from 'datastar'

Datastar.load(
    // all my preprocessor, action and attribute plugins!
)
</script>
```
