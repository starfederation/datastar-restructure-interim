import { expect, describe, test } from 'vitest'
import fc from 'fast-check';

import { ServerSentEventGenerator }  from "./serverSentEventGenerator.ts";
import { fragmentMergeModes } from "./fragments.ts";
import { consoleModes } from "./console.ts";

class MockedServerSentEventGenerator extends ServerSentEventGenerator {
    public constructor() {
        super();
        // no setup needed
    }

    private send(eventType: EventType, dataLines: string[], options: DatastarEventOptions): string[] {
        return super.send(eventType, dataLines, options);
    }

    public redirect(url: string, options?: DatastarEventOptions) {
        return super.redirect(url, options);
    }
}

test('ServerSentEventGenerator can be instantiated', () => {
    const sseGenerator = new MockedServerSentEventGenerator();

    expect(sseGenerator).toBeInstanceOf(MockedServerSentEventGenerator);
});

expect.extend({
  anyItemStartsWith(actual: string[], required: string) {
      const pass = actual.some((line) => {
          return line.startsWith(required);
      });
      const message = () => `Expected an item of [${actual}] to start with ${required}`;

      return { pass, message };
  }
});

function containsAllKeys(
    eventType: DatastarEventType,
    event: Record<string, any>,
    res: string[]
) {
    const { eventId, retryDuration, ...eventOptions } = event;

    expect(res).anyItemStartsWith(`event: ${eventType}`);

    if (eventId) {
        expect(res).anyItemStartsWith(`id: ${eventId}`);
    }

    if (retryDuration) {
        expect(res).anyItemStartsWith(`retry: ${retryDuration}`);
    }

    Object.keys(eventOptions).forEach((key) => {
        expect(res).anyItemStartsWith(`data: ${key}`);
    });
}

const renderFragmentRecord = fc.record({
    eventId: fc.string(),
    retryDuration: fc.nat(),
    fragment: fc.string(),
    selector: fc.string(),
    useViewTransitions: fc.boolean(),
    settleDuration: fc.nat(),
    mergeMode: fc.constantFrom(...fragmentMergeModes)
}, {
    withRequiredKeys: ['fragment']
});

test('renderFragment has valid line endings', () => fc.assert(
    fc.property(fc.array(renderFragmentRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { fragment, ...options } = event;

            const res = sseGenerator.renderFragment(fragment, options);

            res.forEach((line) => {
                expect(line).toMatch(/\n$/);
            });

            expect(res.at(-1)).toEqual("\n\n");
        });
    })
));

test('renderFragment contains all provided keys', () => fc.assert(
    fc.property(fc.array(renderFragmentRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { fragment, ...options } = event;

            const res = sseGenerator.renderFragment(fragment, options);

            containsAllKeys('datastar-fragment', event, res);
        });
    })
));

const removeFragmentsRecord = fc.record({
    eventId: fc.string(),
    retryDuration: fc.nat(),
    selector: fc.string(),
    useViewTransitions: fc.boolean(),
    settleDuration: fc.nat()
}, {
    withRequiredKeys: ['selector']
});

test('removeFragments has valid line endings', () => fc.assert(
    fc.property(fc.array(removeFragmentsRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { selector, ...options } = event;

            const res = sseGenerator.removeFragments(selector, options);

            res.forEach((line) => {
                expect(line).toMatch(/\n$/);
            });

            expect(res.at(-1)).toEqual("\n\n");
        });
    })
));

test('removeFragments contains all provided keys', () => fc.assert(
    fc.property(fc.array(removeFragmentsRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { selector, ...options } = event;

            const res = sseGenerator.removeFragments(selector, options);

            containsAllKeys('datastar-remove', event, res);
        });
    })
));

const patchStoreRecord = fc.record({
    eventId: fc.string(),
    retryDuration: fc.nat(),
    data: fc.object(),
    onlyIfMissing: fc.boolean(),
}, {
    withRequiredKeys: ['selector']
});

test('patchStore has valid line endings', () => fc.assert(
    fc.property(fc.array(patchStoreRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { data, ...options } = event;

            const res = sseGenerator.patchStore(data, options);

            res.forEach((line) => {
                expect(line).toMatch(/\n$/);
            });

            expect(res.at(-1)).toEqual("\n\n");
        });
    })
));

test('patchStore contains all provided keys', () => fc.assert(
    fc.property(fc.array(patchStoreRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { data, ...options } = event;

            const res = sseGenerator.patchStore(data, options);

            containsAllKeys('datastar-signal', { signal: data, ...options }, res);
        });
    })
));

const removeFromStoreRecord = fc.record({
    eventId: fc.string(),
    retryDuration: fc.nat(),
    paths: fc.array(fc.string()),
}, {
    withRequiredKeys: ['selector']
});

test('removeFromStore has valid line endings', () => fc.assert(
    fc.property(fc.array(removeFromStoreRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { paths, ...options } = event;

            const res = sseGenerator.removeFromStore(paths, options);

            res.forEach((line) => {
                expect(line).toMatch(/\n$/);
            });

            expect(res.at(-1)).toEqual("\n\n");
        });
    })
));

test('removeFromStore contains all provided keys', () => fc.assert(
    fc.property(fc.array(removeFromStoreRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { paths, ...options } = event;

            const res = sseGenerator.removeFromStore(paths, options);

            containsAllKeys('datastar-remove', event, res);
        });
    })
));

const redirectRecord = fc.record({
    eventId: fc.string(),
    retryDuration: fc.nat(),
    url: fc.webUrl(),
}, {
    withRequiredKeys: ['selector']
});

test('redirect has valid line endings', () => fc.assert(
    fc.property(fc.array(redirectRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { url, ...options } = event;

            const res = sseGenerator.redirect(url, options);

            res.forEach((line) => {
                expect(line).toMatch(/\n$/);
            });

            expect(res.at(-1)).toEqual("\n\n");
        });
    })
));

test('redirect contains all provided keys', () => fc.assert(
    fc.property(fc.array(redirectRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { url, ...options } = event;

            const res = sseGenerator.redirect(url, options);

            containsAllKeys('datastar-redirect', event, res);
        });
    })
));

const consoleRecord = fc.record({
    eventId: fc.string(),
    retryDuration: fc.nat(),
    mode: fc.constantFrom(...consoleModes),
    message: fc.string()
}, {
    withRequiredKeys: ['mode', 'message']
});

test('console has valid line endings', () => fc.assert(
    fc.property(fc.array(consoleRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { mode, message, ...options } = event;

            const res = mode === "groupEnd" ?
                 sseGenerator.console(mode, options) :
                 sseGenerator.console(mode, message, options);

            res.forEach((line) => {
                expect(line).toMatch(/\n$/);
            });

            expect(res.at(-1)).toEqual("\n\n");
        });
    })
));

test('console contains all provided keys', () => fc.assert(
    fc.property(fc.array(consoleRecord), (data) => {
        const sseGenerator = new MockedServerSentEventGenerator();
        data.forEach((event) => {
            const { mode, message, ...options } = event;

            const res = mode === "groupEnd" ?
                 sseGenerator.console(mode, options) :
                 sseGenerator.console(mode, message, options);

            containsAllKeys('datastar-console', event, res);
        });
    })
));
