import { EventType, DatastarEventOptions } from "./types.ts";
import { PatchStoreOptions } from "./signals.ts";
import { RenderFragmentOptions, FragmentEventOptions} from "./fragments.ts";
import { requireThatString } from "@cowwoc/requirements";

export abstract class ServerSentEventGenerator {
    // runtimes should override this method to create an sse stream
    abstract constructor()

    // runtimes should override this method and use it's output to send an event
    protected send(
         eventType: EventType,
         dataLines: string[],
         options: DatastarEventOptions
    ): string[] {
        const { eventId, retryDuration } = options || {};

        const typeLine = [`event: ${eventType}\n`];
        const idLine = eventId ? [`id: ${eventId}\n`] : [];
        const retryLine = retryDuration ? [`retry: ${retryDuration}\n`] : [];

        return typeLine.concat(
            idLine,
            retryLine,
            dataLines.map((data) => {
                return `data: ${data}\n`;
            }),
            ['\n\n']
        );
    }

    private eachNewlineIsADataLine(prefix: string, data: string) {
        const [ head, ...tail] = data.split('\n');

        return [ `${prefix} ${head}`].concat(tail);
    }

    private eachOptionIsADataLine<T extends Record<string, any>>(options?: T): string[] {
        return Object.keys(options ?? {}).map((key) => {
            return `${key} ${options[key]}`;
        });
    }

    public renderFragment(data: string, options?: RenderFragmentOptions): void {
        const { eventId, retryDuration, ...renderOptions } = options || {};
        const dataLines = this.eachOptionIsADataLine(renderOptions)
            .concat(this.eachNewlineIsADataLine('fragment', data));

        return this.send('datastar-fragment', dataLines, { eventId, retryDuration });
    }

    public removeFragments(selector: string, options?: FragmentEventOptions): void {
        const { eventId, retryDuration, ...eventOptions } = options || {};
        const dataLines = this.eachOptionIsADataLine(eventOptions)
            .concat([`selector ${selector}`]);

        return this.send('datastar-remove', dataLines, { eventId, retryDuration });
    }

    public patchStore(data: Record<string, any>, options?: PatchStoreOptions): void {
        const { eventId, retryDuration, ...eventOptions } = options || {};
        const dataLines = this.eachOptionIsADataLine(eventOptions)
            .concat(this.eachNewlineIsADataLine('signal', JSON.stringify(data)));

        return this.send('datastar-signal', dataLines, { eventId, retryDuration });
    }

    public removeFromStore(paths: string[], options?: DatastarEventOptions): void {
        const { eventId, retryDuration } = options || {};
        const dataLines = [`paths ${paths.join(' ')}`];

        return this.send('datastar-remove', dataLines, { eventId, retryDuration });
    }

    // runtimes should validate that the url is valid with their preferred method
    protected redirect(url: string, options?: DatastarEventOptions): void {
        const { eventId, retryDuration } = options || {};
        const dataLines = [`url ${url}`];

        return this.send('datastar-redirect', dataLines, { eventId, retryDuration });
    }

    public console(mode: "groupEnd", options?: DatastarEventOptions): void
    public console(
        mode: Omit<ConsoleMode, "groupEnd">,
        message: string,
        options?: DatastarEventOptions
    ): void
    public console(
        mode: ConsoleMode,
        message?: string,
        options?: DatastarEventOptions
    ): void {
        const msg = typeof message === 'string' ? message : "";
        const opts = typeof message === 'string' ? options : message;

        const { eventId, retryDuration } = opts || {};

        const msgDataLines = this.eachNewlineIsADataLine('message', msg);
        const dataLines = [`mode ${mode}`].concat(msgDataLines);

        return this.send('datastar-console', dataLines, { eventId, retryDuration });
    }

    abstract parseIncoming()
}
