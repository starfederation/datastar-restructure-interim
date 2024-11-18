use starfederation\datastar\ServerSentEventGenerator;

// Get a random question and answer from somewhere in your code.
[$question, $answer] = getRandomQuestionAnswer();

// Creates a new `ServerSentEventGenerator` instance.
$sse = new ServerSentEventGenerator();

// Merges the HTML fragment into the DOM.
$sse->mergeFragments(
    '<div id="question">' . $question . '</div>'
);

// Merges the `answer` signal into the store.
$sse->mergeSignals(['answer' => $answer]);
