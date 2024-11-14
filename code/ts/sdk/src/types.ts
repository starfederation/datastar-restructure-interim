const eventTypes = [
    "datastar-fragment"
    , "datastar-signal"
    , "datastar-remove"
    , "datastar-redirect"
    , "datastar-console"
] as consts;

export type EventType = typeof eventTypes[number];

export interface DatastarEventOptions {
    eventId: string;
    retryDuration: number;
};

export interface DatastarEvent {
    type: EventType;
};
