<?php

use starfederation\datastar\events\ExecuteJS;

test('Event is correctly output', function() {
    $content = 'console.log("Hello, world!")';
    $event = new ExecuteJS($content);
    expect($event->getDataLines())
        ->toBe([
            'data: script ' . $content,
        ]);
});
