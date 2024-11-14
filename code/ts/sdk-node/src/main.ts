import * as http from "http";
import { serverSentEventGenerator as BaseSseGenerator } from "@starfederation/datastar-sdk/src/serverSentGenerator.ts";
import { sseHeaders } from "@starfederation/datastar-sdk/src/sse.ts";
import { requireThatString } from "@cowwoc/requirements";
import { URL } from 'url';
import  * as querystring from "querystring"
;
export class ServerSentGenerator extends BaseSseGenerator {
    res: http.Response;
    req: http.Request;

    public constructor(req: Request, res: Response) {
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

    public async parseIncoming() {
        if (this.req.method === "GET") {
            const query = queryString.parse(this.req.url.search);
            return JSON.parse(query.datastar);
        }

        await json()(this.req, this.res, (err) => void err && this.res.end(err));

        return this.req.body.datastar;
    }
}
