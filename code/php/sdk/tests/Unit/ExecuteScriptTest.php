<?php

use starfederation\datastar\events\ExecuteScript;

test('Event is correctly output', function() {
    $content = 'console.log("Hello, world!")';
    $event = new ExecuteScript($content);
    expect($event->getDataLines())
        ->toBe([
            'data: script ' . $content,
        ]);
});
