use starfederation\datastar\ServerSentEventGenerator;

// Creates a new `ServerSentEventGenerator` instance.
$sseGenerator = new ServerSentEventGenerator();

// Merges the HTML fragment into the DOM.
$sseGenerator->mergeFragments(
    '<div id="question">What do you put in a toaster?</div>'
);

// Merges the `answer` value into the store.
$sseGenerator->mergeSignals(['answer' => 'bread']);