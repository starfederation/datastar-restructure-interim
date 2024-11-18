<?php

use starfederation\datastar\events\ExecuteJs;

test('Event is correctly output', function() {
    $content = 'console.log("Hello, world!")';
    $event = new ExecuteJs($content);
    expect($event->getDataLines())
        ->toBe([
            'data: script ' . $content,
        ]);
});
