import { Store } from "../plugins/attributes/core/store";
import { ActionsProcessor } from "../plugins/preprocessors/core/actions";
import { SignalsProcessor } from "../plugins/preprocessors/core/signals";
import { Engine } from "./engine";

export * from "./const";
export type * from "./types";

const ds = new Engine();
ds.load(
    ActionsProcessor,
    SignalsProcessor,
    Store,
);

export const Datastar = ds;
