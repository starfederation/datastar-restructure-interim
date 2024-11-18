<?php
use starfederation\datastar\ServerSentEventGenerator;

$sse = new ServerSentEventGenerator();
$sse->mergeFragments('<div id="question">...</div>');
$sse->mergeFragments('<div id="instructions">...</div>');
$sse->mergeSignals(['answer' => '...']);
$sse->mergeSignals(['prize' => '...']);
