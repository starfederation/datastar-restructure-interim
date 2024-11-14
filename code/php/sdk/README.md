# Datastar PHP SDK

This package provides a PHP SDK for working with [Datastar](https://data-star.dev/).

## License

This package is licensed for free under the MIT License.

## Requirements

This package requires PHP 8.1 or later.

## Installation

Install using composer.

```shell
composer require starfederation/datastar-php
```

## Usage

```php
use starfederation\datastar\enums\ConsoleMode;
use starfederation\datastar\enums\EventType;
use starfederation\datastar\enums\FragmentMergeMode;
use starfederation\datastar\ServerSentEventGenerator;

// Creates a new `ServerSentEventGenerator` instance.
$sseGenerator = new ServerSentEventGenerator();

// Renders a fragment in the DOM.
$sseGenerator->renderFragment('<div></div>', [
    'selector' => '#my-div',
    'mergeMode' => FragmentMergeMode::Append,
    'settleDuration' => 1000,
    'useViewTransition' => true,
]);

// Removes one or more fragments from the DOM.
$sseGenerator->removeFragments('#my-div');

// Updates values in the store.
$sseGenerator->patchStore(['foo' => 123], ['onlyIfMissing' => true]);

// Removes one or more paths from the store.
$sseGenerator->removeFromStore(['foo', 'bar']);

// Redirects the browser.
$sseGenerator->redirect('/success');

// Sends a message to the browser console.
$sseGenerator->console(ConsoleMode::Log, 'Hello, world!');
```

```php
use starfederation\datastar\ParseIncoming;

$store = ParseIncoming::getStore();
```

---

Created by [PutYourLightsOn](https://putyourlightson.com/).
