import * as http from "http";
import { serverSentEventGenerator as BaseSseGenerator } from "@starfederation/datastar-sdk/src/serverSentGenerator.ts";
import { sseHeaders } from "@starfederation/datastar-sdk/src/sse.ts";
import { requireThatString } from "@cowwoc/requirements";
import { URL } from 'url';

export class ServerSentGenerator extends BaseSseGenerator {
    res: http.Response;
    req: http.Request;

    public constructor(req: http.Request, res: http.Response) {
        super();
        this.res = res;
        this.req = req;

        this.res.writeHead(200, sseHeaders);

        // When client closes connection, stop sending events
        this.req.on('close', () => {
            this.res.end();
        });
    }

    private send(eventType: EventType, dataLines: string[], options: DatastarEventOptions) {
        const eventLines = super.send(eventType, dataLines, options);

        eventLines.forEach((line) => {
            this.res.write(line);
        });
    }

    public redirect(url: URL, options?: DatastarEventOptions): void {
        requireThatString(url.protocol, "url protocol must be http(s)").matches(\^https?$\);

        super(url.href, options);
    }
}
